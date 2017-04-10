package api

import (
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
)

// CreateDisk create disk
func (c *APIClient) CreateDisk(spec *sacloud.Disk) (string, error) {
	disk, err := c.client.Disk.Create(spec)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", disk.ID), nil
}

// EditDisk edit disk
func (c *APIClient) EditDisk(diskID string, spec *sacloud.DiskEditValue) (bool, error) {
	id, res := ToSakuraID(diskID)
	if !res {
		return false, fmt.Errorf("DiskID is invalid: %v", diskID)
	}
	_, err := c.client.Disk.Config(id, spec)
	if err != nil {
		return false, err
	}
	return true, err
}

// ConnectDisk connect disk
func (c *APIClient) ConnectDisk(diskID string, serverID string) (bool, error) {

	dId, res := ToSakuraID(diskID)
	if !res {
		return false, fmt.Errorf("DiskID is invalid: %v", diskID)
	}
	sId, res := ToSakuraID(serverID)
	if !res {
		return false, fmt.Errorf("ServerID is invalid: %v", serverID)
	}

	_, err := c.client.Disk.ConnectToServer(dId, sId)
	if err != nil {
		return false, err
	}
	return true, err
}

// DiskState get disk state
func (c *APIClient) GetDiskByID(diskID string) (*sacloud.Disk, error) {
	id, res := ToSakuraID(diskID)
	if !res {
		return nil, fmt.Errorf("DiskID is invalid: %v", diskID)
	}
	disk, err := c.client.Disk.Read(id)
	if err != nil {
		return nil, err
	}
	return disk, err

}
