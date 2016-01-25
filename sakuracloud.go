package sakuracloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	sakuraCloudAPIRoot                = "https://secure.sakura.ad.jp/cloud/zone"
	sakuraCloudAPIRootSuffix          = "api/cloud/1.1"
	sakuraCloudPublicImageSearchWords = "Ubuntu%20Server%2014%2064bit"
	sakuraAllowSudoScriptBody         = `#!/bin/bash

  # @sacloud-once
  # @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集します
  # @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
  # @sacloud-require-archive distro-debian
  # @sacloud-require-archive distro-ubuntu

  export DEBIAN_FRONTEND=noninteractive
	echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
  exit 0`

	sakuraAllowSudoWithKernelUpgradeScriptBody = `#!/bin/bash

  # @sacloud-once
  # @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集 + Kernelアップグレード(linux-generic-lts-vivid) + 再起動
  # @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
  # @sacloud-require-archive distro-debian
  # @sacloud-require-archive distro-ubuntu

  export DEBIAN_FRONTEND=noninteractive
	echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
  apt-get -y update || exit 1
	apt-get install -y --install-recommends linux-generic-lts-vivid || exit 1
	sh -c 'sleep 10; shutdown -h now' &
  exit 0`

	sakuraAddIPForEth1ScriptBodyFormat = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: setup ip address for eth1
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	echo "auto eth1" >> /etc/network/interfaces
	echo "iface eth1 inet static" >> /etc/network/interfaces
	echo "address %s" >> /etc/network/interfaces
	echo "netmask %s" >> /etc/network/interfaces
	ifdown eth1; ifup eth1
	exit 0`

	sakuraChangeDefaultGatewayScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: change default gateway
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/gateway/#gateway/g' /etc/network/interfaces
	echo "up route add default gw %s" >> /etc/network/interfaces
	exit 0`

	sakuraDisableEth0ScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: disable eth0
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/iface eth0 inet static/iface eth0 inet manual/g' /etc/network/interfaces
	ifdown eth0 || exit 0
	exit 0`
)

// Client type of sakuracloud api client config values
type Client struct {
	AccessToken       string
	AccessTokenSecret string
	Region            string
}

// NewClient create new sakuracloud api client
func NewClient(token, tokenSecret, region string) *Client {
	return &Client{AccessToken: token, AccessTokenSecret: tokenSecret, Region: region}
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

// Create create server
func (c *Client) Create(spec *sakura.Server, addIPAddress string) (*sakura.Response, error) {
	var (
		method = "POST"
		uri    = "server"
		body   = sakura.Request{Server: spec}
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	var res sakura.Response
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	if addIPAddress != "" && len(res.Server.Interfaces) > 1 {
		if err := c.updateIPAddress(spec, res, addIPAddress); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (c *Client) updateIPAddress(spec *sakura.Server, statusRes sakura.Response, ip string) error {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("interface/%s", statusRes.Server.Interfaces[1].ID)
		body   = sakura.Request{
			Interface: &sakura.Interface{UserIPAddress: ip},
		}
	)

	_, err := c.newRequest(method, uri, body)
	if err != nil {
		return err
	}

	return nil

}

// CreateDisk create disk
func (c *Client) CreateDisk(spec *sakura.Disk) (string, error) {
	var (
		method = "POST"
		uri    = "disk"
		body   = sakura.Request{Disk: spec}
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return "", err
	}

	//HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないため文字列で受ける
	var res struct {
		*sakura.Response
		Success string `json:",omitempty"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.Disk.ID, nil
}

// EditDisk edit disk
func (c *Client) EditDisk(diskID string, spec *sakura.DiskEditValue) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/config", "disk", diskID)
		body   = spec
	)

	data, err := c.newRequest(method, uri, body)
	if err != nil {
		return false, err
	}

	var res sakura.Response
	if err := json.Unmarshal(data, &res); err != nil {
		return false, err
	}
	return true, nil
}

// ConnectDisk connect disk
func (c *Client) ConnectDisk(diskID string, serverID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/server/%s", "disk", diskID, serverID)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return false, err
	}

	var res sakura.ResultFlagValue
	if err := json.Unmarshal(data, &res); err != nil {
		return false, err
	}
	return true, nil
}

// Delete delete server
func (c *Client) Delete(id string, disks []string) error {
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

// State get server state
func (c *Client) State(id string) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "server", id)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s sakura.Response
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}
	return s.Server.Instance.Status, nil
}

// DiskState get disk state
func (c *Client) DiskState(diskID string) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "disk", diskID)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s sakura.Response
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Disk.Availability, nil
}

// PowerOn power on
func (c *Client) PowerOn(id string) error {
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

// PowerOff power off
func (c *Client) PowerOff(id string) error {
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

// GetIP get public ip address
func (c *Client) GetIP(id string, privateIPOnly bool) (string, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s", "server", id)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var s sakura.Response
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	if privateIPOnly && len(s.Server.Interfaces) > 1 {
		return s.Server.Interfaces[1].UserIPAddress, nil
	}

	return s.Server.Interfaces[0].IPAddress, nil
}

// GetAllowSudoNoteID get ubuntu customize note id
// FIXME
// workaround for [Non root ssh create sudo can't get password](https://github.com/docker/machine/issues/1569)
// [PR #1586](https://github.com/docker/machine/pull/1586)がマージされるまで暫定
// スクリプト(Note)を使ってubuntuユーザがsudo可能にする
func (c *Client) GetAllowSudoNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_99_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraAllowSudoScriptBody)
}

// GetAllowSudoWithKernelUpgradeNoteID get ubuntu customize note id
func (c *Client) GetAllowSudoWithKernelUpgradeNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_99_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraAllowSudoWithKernelUpgradeScriptBody)
}

// GetAddIPCustomizeNoteID get add ip customize note id
func (c *Client) GetAddIPCustomizeNoteID(serverID string, ip string, subnet string) (string, error) {
	noteName := fmt.Sprintf("_30_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraAddIPForEth1ScriptBodyFormat, ip, subnet)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetChangeGatewayCustomizeNoteID get change gateway address customize note id
func (c *Client) GetChangeGatewayCustomizeNoteID(serverID string, gateway string) (string, error) {
	noteName := fmt.Sprintf("_20_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraChangeDefaultGatewayScriptBody, gateway)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetDisableEth0CustomizeNoteID get disable eth0 customize note id
func (c *Client) GetDisableEth0CustomizeNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_10_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraDisableEth0ScriptBody)
}

func (c *Client) getCustomizeNoteID(noteName string, noteBody string) (string, error) {

	var (
		method = "GET"
		uri    = "note"
		body   = sakura.Request{
			Filter: map[string]interface{}{"Name": noteName},
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
	var existsNote sakura.SearchResponse
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
	n := sakura.Request{
		Note: &sakura.Note{
			Name:    noteName,
			Content: noteBody,
		},
	}

	data, err = c.newRequest(method, uri, n)
	if err != nil {
		return "", err
	}
	var s sakura.Response
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Note.ID, nil

}

// DeleteNote delete note
func (c *Client) DeleteNote(id string) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("note/%s", id)
	)

	_, err := c.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	return nil
}

// ConnectPacketFilterToSharedNIC connect packet filter to eth0(shared)
func (c *Client) ConnectPacketFilterToSharedNIC(server *sakura.Server, idOrNameFilter string) error {
	if server.Interfaces != nil && len(server.Interfaces) > 0 {
		return c.connectPacketFilter(&server.Interfaces[0], idOrNameFilter)
	}
	return nil
}

// ConnectPacketFilterToPrivateNIC connect packet filter to eth1(private)
func (c *Client) ConnectPacketFilterToPrivateNIC(server *sakura.Server, idOrNameFilter string) error {
	if server.Interfaces != nil && len(server.Interfaces) > 1 {
		return c.connectPacketFilter(&server.Interfaces[1], idOrNameFilter)
	}
	return nil
}

// ConnectPacketFilter connect filter to nic
func (c *Client) connectPacketFilter(nic *sakura.Interface, idOrNameFilter string) error {
	if idOrNameFilter == "" {
		return nil
	}

	var id string
	//id or name ?
	if match, _ := regexp.MatchString(`^[0-9]+$`, idOrNameFilter); match {
		//IDでの検索
		var (
			method = "GET"
			uri    = fmt.Sprintf("packetfilter/%s", idOrNameFilter)
		)
		data, _ := c.newRequest(method, uri, nil)

		var res sakura.Response
		if err := json.Unmarshal(data, &res); err != nil {
		} else {
			if res.IsOk {
				id = res.PacketFilter.ID
			}
		}
		// else {
		// 	return fmt.Errorf("PacketFilter [%s](id):Not Found", idOrNameFilter)
		// }
	}

	//search
	if id == "" {
		//名前での検索
		var (
			method = "GET"
			uri    = "packetfilter"
			body   = sakura.Request{
				Filter: map[string]interface{}{"Name": idOrNameFilter},
			}
		)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return err
		}
		uri = fmt.Sprintf("%s?%s", uri, bodyJSON)
		data, err := c.newRequest(method, uri, nil)
		if err != nil {
			return err
		}

		var res sakura.SearchResponse
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		if res.Count > 0 {
			id = res.PacketFilters[0].ID
		} else {
			return fmt.Errorf("PacketFilter [%s](name):Not Found", idOrNameFilter)
		}
	}

	// not found
	if id == "" {
		return nil
	}

	//connect
	var (
		method = "PUT"
		uri    = fmt.Sprintf("/interface/%s/to/packetfilter/%s", nic.ID, id)
	)

	_, err := c.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetUbuntuArchiveID get ubuntu archive id
func (c *Client) GetUbuntuArchiveID() (string, error) {

	var (
		method = "GET"
		uri    = "archive"
		body   = sakura.Request{
			Filter: map[string]interface{}{
				"Name":  sakuraCloudPublicImageSearchWords,
				"Scope": "shared",
			},
			Include: []string{"ID", "Name"},
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
	var ubuntu sakura.SearchResponse
	if err := json.Unmarshal(data, &ubuntu); err != nil {
		return "", err
	}

	//すでに登録されている場合
	if ubuntu.Count > 0 {
		return ubuntu.Archives[0].ID, nil
	}

	return "", errors.New("Archive'Ubuntu Server 14.04 LTS 64bit' not found.")
}

// IsValidPlan return valid plan
func (c *Client) IsValidPlan(core int, memGB int) (bool, error) {
	//assert args
	if core <= 0 {
		return false, fmt.Errorf("Invalid Parameter: CPU Core")
	}
	if memGB <= 0 {
		return false, fmt.Errorf("Invalid Parameter: Memory Size(GB)")
	}

	planID, _ := strconv.ParseInt(fmt.Sprintf("%d%03d", memGB, core), 10, 64)

	var (
		method = "GET"
		uri    = fmt.Sprintf("product/server/%d", planID)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return false, err
	}
	var plan sakura.Response
	if err := json.Unmarshal(data, &plan); err != nil {
		return false, err
	}

	if plan.ServerPlan != nil {
		return true, nil
	}

	return false, fmt.Errorf("Server Plan[%d] Not Found", planID)

}
