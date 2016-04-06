package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
)

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
