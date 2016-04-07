package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
)

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
