package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"regexp"
)

// State get server state
func (c *APIClient) State(id string) (string, error) {
	server, err := c.client.Server.Read(id)
	if err != nil {
		return "", err
	}
	return server.Instance.Status, nil
}

// PowerOn power on
func (c *APIClient) PowerOn(id string) error {
	_, err := c.client.Server.Boot(id)
	return err
}

// PowerOff power off
func (c *APIClient) PowerOff(id string) error {
	_, err := c.client.Server.Shutdown(id)
	return err
}

// GetIP get public ip address
func (c *APIClient) GetIP(id string, privateIPOnly bool) (string, error) {

	server, err := c.client.Server.Read(id)
	if err != nil {
		return "", err
	}

	if privateIPOnly && len(server.Interfaces) > 1 {
		return server.Interfaces[1].UserIPAddress, nil
	}

	return server.Interfaces[0].IPAddress, nil
}

// Create create server
func (c *APIClient) Create(spec *sacloud.Server, addIPAddress string) (*sacloud.Server, error) {

	server, err := c.client.Server.Create(spec)
	if err != nil {
		return nil, err
	}

	if addIPAddress != "" && len(server.Interfaces) > 1 {
		if err := c.updateIPAddress(&server.Interfaces[1], addIPAddress); err != nil {
			return nil, err
		}
	}

	return server, nil
}

func (c *APIClient) updateIPAddress(nic *sacloud.Interface, ip string) error {
	nic.UserIPAddress = ip
	_, err := c.client.Interface.Update(nic.ID, nic)

	if err != nil {
		return err
	}

	return nil

}

// Delete delete server
func (c *APIClient) Delete(id string, disks []string) error {
	_, err := c.client.Server.DeleteWithDisk(id, disks)
	return err
}

// ConnectPacketFilterToSharedNIC connect packet filter to eth0(shared)
func (c *APIClient) ConnectPacketFilterToSharedNIC(server *sacloud.Server, idOrNameFilter string) error {
	if server.Interfaces != nil && len(server.Interfaces) > 0 {
		return c.connectPacketFilter(&server.Interfaces[0], idOrNameFilter)
	}
	return nil
}

// ConnectPacketFilterToPrivateNIC connect packet filter to eth1(private)
func (c *APIClient) ConnectPacketFilterToPrivateNIC(server *sacloud.Server, idOrNameFilter string) error {
	if server.Interfaces != nil && len(server.Interfaces) > 1 {
		return c.connectPacketFilter(&server.Interfaces[1], idOrNameFilter)
	}
	return nil
}

// ConnectPacketFilter connect filter to nic
func (c *APIClient) connectPacketFilter(nic *sacloud.Interface, idOrNameFilter string) error {
	if idOrNameFilter == "" {
		return nil
	}

	var id string
	//id or name ?
	if match, _ := regexp.MatchString(`^[0-9]+$`, idOrNameFilter); match {
		//IDでの検索
		p, err := c.client.PacketFilter.Read(idOrNameFilter)
		if err != nil {
			return err
		}
		id = p.ID
	}

	//search
	if id == "" {

		res, err := c.client.PacketFilter.WithNameLike(idOrNameFilter).Limit(1).Find()

		if err != nil {
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
	_, err := c.client.Interface.ConnectToPacketFilter(nic.ID, id)
	if err != nil {
		return err
	}
	return nil
}
