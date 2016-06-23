package api

import (
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

// CreateDisk create disk
func (c *APIClient) CreateDisk(spec *sacloud.Disk) (string, error) {
	disk, err := c.client.Disk.Create(spec)
	if err != nil {
		return "", err
	}

	return disk.ID, nil
}

// EditDisk edit disk
func (c *APIClient) EditDisk(diskID string, spec *sacloud.DiskEditValue) (bool, error) {
	_, err := c.client.Disk.Config(diskID, spec)
	if err != nil {
		return false, err
	}
	return true, err
}

// ConnectDisk connect disk
func (c *APIClient) ConnectDisk(diskID string, serverID string) (bool, error) {
	_, err := c.client.Disk.ConnectToServer(diskID, serverID)
	if err != nil {
		return false, err
	}
	return true, err
}

// DiskState get disk state
func (c *APIClient) GetDiskByID(diskID string) (*sacloud.Disk, error) {
	disk, err := c.client.Disk.Read(diskID)
	if err != nil {
		return nil, err
	}
	return disk, err

}
