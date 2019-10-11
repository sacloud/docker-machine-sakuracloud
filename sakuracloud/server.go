package sakuracloud

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// State reads server state
func (c *APIClient) State(strID string) (string, error) {
	id := types.StringID(strID)
	if id.IsEmpty() {
		return "", fmt.Errorf("ServerID is invalid: %s", strID)
	}
	server, err := sacloud.NewServerOp(c.caller).Read(context.Background(), c.Zone, id)
	if err != nil {
		return "", err
	}
	return string(server.InstanceStatus), nil
}

// PowerOn power on
func (c *APIClient) PowerOn(strID string) error {
	id := types.StringID(strID)
	if id.IsEmpty() {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}

	return sacloud.NewServerOp(c.caller).Boot(context.Background(), c.Zone, id)
}

// PowerOff power off
func (c *APIClient) PowerOff(strID string) error {
	id := types.StringID(strID)
	if id.IsEmpty() {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}
	return sacloud.NewServerOp(c.caller).Shutdown(context.Background(), c.Zone, id, nil)
}

// GetIP get public ip address
func (c *APIClient) GetIP(strID string) (string, error) {
	id := types.StringID(strID)
	if id.IsEmpty() {
		return "", fmt.Errorf("ServerID is invalid: %s", strID)
	}
	server, err := sacloud.NewServerOp(c.caller).Read(context.Background(), c.Zone, id)
	if err != nil {
		return "", err
	}
	return server.Interfaces[0].IPAddress, nil
}

// Delete delete server
func (c *APIClient) Delete(strID string, strDisks []string) error {
	id := types.StringID(strID)
	if id.IsEmpty() {
		return fmt.Errorf("ServerID is invalid: %s", strID)
	}

	server, err := sacloud.NewServerOp(c.caller).Read(context.Background(), c.Zone, id)
	if err != nil {
		return fmt.Errorf("reading server is failed: %s", id)
	}
	var disks []types.ID
	for _, disk := range server.Disks {
		disks = append(disks, disk.ID)
	}

	return sacloud.NewServerOp(c.caller).DeleteWithDisks(context.Background(), c.Zone, id, &sacloud.ServerDeleteWithDisksRequest{
		IDs: disks,
	})
}

// ReadServer returns server info
func (c *APIClient) ReadServer(ctx context.Context, id types.ID) (*sacloud.Server, error) {
	sv, err := sacloud.NewServerOp(c.caller).Read(context.Background(), c.Zone, id)
	if err != nil {
		return nil, err
	}
	return sv, nil
}
