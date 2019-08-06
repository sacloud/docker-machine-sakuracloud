package sakuracloud

import (
	"fmt"
)

// State get server state
func (c *APIClient) State(strID string) (string, error) {
	id, res := ToSakuraID(strID)
	if !res {
		return "", fmt.Errorf("ServerID is invalid: %s", strID)
	}

	server, err := c.client.Server.Read(id)
	if err != nil {
		return "", err
	}
	return server.Instance.Status, nil
}

// PowerOn power on
func (c *APIClient) PowerOn(strID string) error {
	id, res := ToSakuraID(strID)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}

	_, err := c.client.Server.Boot(id)
	return err
}

// PowerOff power off
func (c *APIClient) PowerOff(strID string) error {
	id, res := ToSakuraID(strID)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}
	_, err := c.client.Server.Shutdown(id)
	return err
}

// GetIP get public ip address
func (c *APIClient) GetIP(strID string) (string, error) {
	id, res := ToSakuraID(strID)
	if !res {
		return "", fmt.Errorf("ServerID is invalid: %s", strID)
	}
	server, err := c.client.Server.Read(id)
	if err != nil {
		return "", err
	}
	return server.Interfaces[0].IPAddress, nil
}

// Delete delete server
func (c *APIClient) Delete(strID string, strDisks []string) error {
	id, res := ToSakuraID(strID)
	if !res {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}

	disks, res := ToSakuraIDAll(strDisks)
	if !res {
		return fmt.Errorf("DiskIDs are invalid: %#v", strDisks)
	}

	_, err := c.client.Server.DeleteWithDisk(id, disks)
	return err
}
