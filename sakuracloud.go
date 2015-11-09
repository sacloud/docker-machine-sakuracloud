package sakuracloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	sakuraCloudAPIRoot          = "https://secure.sakura.ad.jp/cloud/zone"
	sakuraCloudAPIRootSuffix    = "api/cloud/1.1"
	sakuraUbuntuSetupScriptName = "_allow-sudo-for-docker-machine_"
	sakuraUbuntuSetupScriptBody = `#!/bin/bash

  # @sacloud-once
  # @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集します
  # @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
  # @sacloud-require-archive distro-debian
  # @sacloud-require-archive distro-ubuntu

  export DEBIAN_FRONTEND=noninteractive
	echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
  sudo apt-get -y update || exit 1
  exit 0`
)

type Client struct {
	AccessToken       string
	AccessTokenSecret string
	Region            string
}

type virtualGuest struct {
	*Client
}

type Resource struct {
	ID string
}

type ResultFlagValue struct {
	IsOk bool `json:"is_ok"`
}

type ConnectedSwitch struct {
	//TODO 既存スイッチへの接続(ID指定),未接続(null)への対応
	Scope string
}

type Server struct {
	Name              string
	Description       string
	ServerPlan        Resource
	Tags              []string
	ConnectedSwitches []ConnectedSwitch
}

type Disk struct {
	Name          string
	Plan          Resource
	SizeMB        int
	Connection    string
	SourceArchive Resource
}

type SSHKey struct {
	PublicKey string
}

type DiskEditValue struct {
	Password string
	SSHKey   SSHKey
	Notes    []Resource
}

type resDisk struct {
	ReinstallCount int
	ID             string
	Availability   string
	SizeMB         int
	Storage        struct {
		ID         string
		Class      string
		MountIndex string
	}
}

type resInterface struct {
	ID         string
	IPAddress  string
	MACAddress string
}

type createServerRequest struct {
	Server Server
	Count  int
}
type ServerStatusResponse struct {
	Server struct {
		ID    string
		Icon  string
		Disks []resDisk

		HostName   string
		Interfaces []resInterface

		Instance struct {
			Status string
		}
	}

	ResultFlagValue
}

type createDiskRequest struct {
	Disk Disk
}

type DiskStatusResponse struct {
	Disk resDisk
	ResultFlagValue
}

func NewClient(token, tokenSecret, region string) *Client {
	return &Client{AccessToken: token, AccessTokenSecret: tokenSecret, Region: region}
}

func (c *Client) VirtualGuest() *virtualGuest {
	return &virtualGuest{c}
}

func (c *Client) getEndpoint() string {
	return fmt.Sprintf("%s/%s/%s", sakuraCloudAPIRoot, c.Region, sakuraCloudAPIRootSuffix)
}

func (c *Client) isOkStatus(code int) bool {
	codes := map[int]bool{
		200: true,
		201: true,
		202: true,
		204: true,
		305: false,
		400: false,
		401: false,
		403: false,
		404: false,
		405: false,
		406: false,
		408: false,
		409: false,
		411: false,
		413: false,
		415: false,
		500: false,
		503: false,
	}

	return codes[code]
}

func (c *Client) isReadyDiskStatus(state string) bool {
	return state == "available"
}

func (c *Client) newRequest(method, uri string, body interface{}) ([]byte, error) {
	var (
		client = &http.Client{}
		url    = fmt.Sprintf("%s/%s", c.getEndpoint(), uri)
		err    error
		req    *http.Request
	)

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyJSON))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Error with request: %v - %q", url, err)
	}

	req.SetBasicAuth(c.AccessToken, c.AccessTokenSecret)
	req.Method = method

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if !c.isOkStatus(resp.StatusCode) {
		return nil, fmt.Errorf("Error in response: %s", data)
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Create(spec *Server) (string, error) {
	var (
		method = "POST"
		uri    = "server"
		body   = createServerRequest{Server: *spec}
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return "", err
	}

	var res ServerStatusResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	return res.Server.ID, nil
}

func (c *Client) CreateDisk(spec *Disk) (string, error) {
	var (
		method = "POST"
		uri    = "disk"
		body   = createDiskRequest{Disk: *spec}
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return "", err
	}

	var res DiskStatusResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.Disk.ID, nil
}

func (c *Client) EditDisk(diskId string, spec *DiskEditValue) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/config", "disk", diskId)
		body   = spec
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return false, err
	}

	var res DiskStatusResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) ConnectDisk(diskId string, serverId string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/server/%s", "disk", diskId, serverId)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return false, err
	}

	var res ResultFlagValue
	if err := json.Unmarshal(data, &res); err != nil {
		return false, err
	}
	return true, nil
}

func (c *virtualGuest) Delete(id string, disks []string) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s", "server", id)
	)

	_, err := c.newRequest(method, uri, map[string]interface{}{"WithDisk": disks})
	if err != nil {
		return err
	}
	return nil
}

func (c *virtualGuest) State(id string) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "server", id)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s ServerStatusResponse
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}
	return s.Server.Instance.Status, nil
}

func (c *virtualGuest) DiskState(diskId string) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "disk", diskId)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s DiskStatusResponse
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Disk.Availability, nil
}

func (c *virtualGuest) PowerOn(id string) error {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/power", "server", id)
	)

	_, err := c.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *virtualGuest) PowerOff(id string) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", "server", id)
	)

	_, err := c.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *virtualGuest) GetIP(id string) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "server", id)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s ServerStatusResponse
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Server.Interfaces[0].IPAddress, nil
}

// FIXME
// workaround for [Non root ssh create sudo can't get password](https://github.com/docker/machine/issues/1569)
// [PR #1586](https://github.com/docker/machine/pull/1586)がマージされるまで暫定
// スクリプト(Note)を使ってubuntuユーザがsudo可能にする
func (c *virtualGuest) GetUbuntuCustomizeNoteId() (string, error) {
	//TODO ノートのバージョニング

	type filter struct {
		Name string
	}
	type noteRequest struct {
		Filter filter
	}

	type noteData struct {
		ID           string
		Name         string
		Content      string
		Availability string
	}
	type noteResponse struct {
		Count int
		Notes []noteData
	}
	type createNoteData struct {
		Name    string
		Content string
	}
	type createNodeWrap struct {
		Note createNoteData
	}
	type responseNoteWrap struct {
		Note noteData
	}
	var (
		method = "GET"
		uri    = "note"
		body   = noteRequest{
			Filter: filter{Name: sakuraUbuntuSetupScriptName},
		}
	)
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	uri = fmt.Sprintf("%s?%s", uri, bodyJSON)
	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var existsNote noteResponse
	if err := json.Unmarshal(data, &existsNote); err != nil {
		return "", err
	}

	//すでに登録されている場合
	if existsNote.Count > 0 {
		return existsNote.Notes[0].ID, nil
	}

	//ない場合はここで作成する
	method = "POST"
	uri = "note"
	n := createNodeWrap{
		Note: createNoteData{
			Name:    sakuraUbuntuSetupScriptName,
			Content: sakuraUbuntuSetupScriptBody,
		},
	}

	data, err = c.newRequest(method, uri, n)
	if err != nil {
		return "", err
	}
	var s responseNoteWrap
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Note.ID, nil
}
