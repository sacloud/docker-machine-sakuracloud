package api

import (
	"fmt"
)

// State get server state
func (c *APIClient) State(strId string) (string, error) {
	id, res := ToSakuraID(strId)
	if !res {
		return "", fmt.Errorf("ServerID is invalid: %s", strId)
	}

	server, err := c.Client.Server.Read(id)
	if err != nil {
		return "", err
	}
	return server.Instance.Status, nil
}

// PowerOn power on
func (c *APIClient) PowerOn(strId string) error {
	id, res := ToSakuraID(strId)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strId)
	}

	_, err := c.Client.Server.Boot(id)
	return err
}

// PowerOff power off
func (c *APIClient) PowerOff(strId string) error {
	id, res := ToSakuraID(strId)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strId)
	}
	_, err := c.Client.Server.Shutdown(id)
	return err
}

// GetIP get public ip address
func (c *APIClient) GetIP(strId string) (string, error) {
	id, res := ToSakuraID(strId)
	if !res {
		return "", fmt.Errorf("ServerID is invalid: %s", strId)
	}
	server, err := c.Client.Server.Read(id)
	if err != nil {
		return "", err
	}
	return server.Interfaces[0].IPAddress, nil
}

// Delete delete server
func (c *APIClient) Delete(strId string, strDisks []string) error {
	id, res := ToSakuraID(strId)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strId)
	}

	disks, res := ToSakuraIDAll(strDisks)
	if !res {
		return fmt.Errorf("DiskIDs are invalid: %#v", strDisks)
	}

	_, err := c.Client.Server.DeleteWithDisk(id, disks)
	return err
}
