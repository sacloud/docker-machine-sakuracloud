package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
	"regexp"
)

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
