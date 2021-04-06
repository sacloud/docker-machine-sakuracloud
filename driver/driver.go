package driver

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/sacloud/docker-machine-sakuracloud/sakuracloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/server"
	"github.com/sacloud/libsacloud/v2/utils/server/ostype"
)

// Driver sakuracloud driver
type Driver struct {
	*drivers.BaseDriver
	serverConfig *sakuraServerConfig
	Client       *sakuracloud.APIClient
	ID           string
	DiskID       string
	EnginePort   int
	SSHKey       string
}

// GetCreateFlags create flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return mcnFlags
}

// NewDriver create driver instance
func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Client:       &sakuracloud.APIClient{},
		serverConfig: defaultServerConfig,
		EnginePort:   defaultServerConfig.EnginePort,
	}
}

func validateSakuraServerConfig(c *sakuracloud.APIClient, config *sakuraServerConfig) error {

	err := config.Validate()
	if err != nil {
		return fmt.Errorf("invalid parameter: %s", err)
	}

	res, err := c.IsValidPlan(config.Core, config.Memory)
	if !res || err != nil {
		return fmt.Errorf("invalid parameter: core or memory is invalid : %v", err)
	}

	if config.PacketFilter != "" {
		id := types.StringID(config.PacketFilter)
		if id.IsEmpty() {
			return fmt.Errorf("invalid parameter: invalid packet-filter-id")
		}
		exists, err := c.IsExistsPacketFilter(id)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("invalid parameter: packet-filter[id:%d] is not exists", id)
		}
	}

	return nil
}

// SetConfigFromFlags create config values from flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	// API Client
	d.Client = sakuracloud.NewAPIClient(
		flags.String("sakuracloud-access-token"),
		flags.String("sakuracloud-access-token-secret"),
		flags.String("sakuracloud-zone"),
		flags.String("sakuracloud-password"),
	)
	if err := d.getClient().ValidateClientConfig(); err != nil {
		return err
	}

	// Swarm(legacy swarm)
	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")

	d.serverConfig = &sakuraServerConfig{
		HostName:        "",
		OSType:          flags.String("sakuracloud-os-type"),
		Core:            flags.Int("sakuracloud-core"),
		Memory:          flags.Int("sakuracloud-memory"),
		DiskPlan:        flags.String("sakuracloud-disk-plan"),
		DiskSize:        flags.Int("sakuracloud-disk-size"),
		DiskConnection:  flags.String("sakuracloud-disk-connection"),
		InterfaceDriver: flags.String("sakuracloud-interface-driver"),
		Password:        flags.String("sakuracloud-password"),
		PacketFilter:    flags.String("sakuracloud-packet-filter"),
		EnablePWAuth:    flags.Bool("sakuracloud-enable-password-auth"),
	}

	if d.serverConfig.HostName == "" {
		d.serverConfig.HostName = d.GetMachineName()
	}

	// for SSH
	d.SSHUser = d.serverConfig.SSHUserName()
	d.SSHPort = 22
	d.SSHKey = flags.String("sakuracloud-ssh-key")

	// for docker engine port
	d.EnginePort = flags.Int("sakuracloud-engine-port")

	return validateSakuraServerConfig(d.Client, d.serverConfig)
}

func (d *Driver) getClient() *sakuracloud.APIClient {
	d.Client.Init()
	return d.Client
}

// DriverName return driver name
func (d *Driver) DriverName() string {
	return "sakuracloud"
}

// GetURL return docker url
func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, strconv.Itoa(d.EnginePort))), nil
}

// GetSSHHostname return ssh hostname
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetIP return public or private ip address
func (d *Driver) GetIP() (string, error) {
	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	return d.getClient().GetIP(d.ID)
}

// GetState get server power state
func (d *Driver) GetState() (state.State, error) {
	s, err := d.getClient().State(d.ID)
	if err != nil {
		return state.None, err
	}
	var vmState state.State
	switch s {
	case "up":
		vmState = state.Running
	case "cleaning":
		vmState = state.Stopping
	case "down":
		vmState = state.Stopped
	default:
		vmState = state.None
	}
	return vmState, nil
}

// PreCreateCheck check before create
func (d *Driver) PreCreateCheck() error {
	if d.SSHKey != "" {
		if _, err := os.Stat(d.SSHKey); os.IsNotExist(err) {
			return fmt.Errorf("Ssh key does not exist: %q", d.SSHKey)
		}

		if _, err := os.Stat(d.SSHKey + ".pub"); os.IsNotExist(err) {
			return fmt.Errorf("Ssh public key does not exist: %q", d.SSHKey+".pub")
		}
	}
	return nil
}

// Create create server on sakuracloud
func (d *Driver) Create() error {

	publicKey, err := d.prepareSSHKey()
	if err != nil {
		return err
	}
	d.preparePassword()

	// build server
	ctx := context.Background()
	sb := d.buildSakuraServerSpec(publicKey)
	buildResult, err := sb.Build(ctx, d.Client.ServerBuilderClient(), d.Client.Zone)
	if err != nil {
		return fmt.Errorf("error creating host: %v", err)
	}

	// read server status
	sv, err := d.Client.ReadServer(ctx, buildResult.ServerID)
	if err != nil {
		return fmt.Errorf("error creating host: %v", err)
	}

	d.ID = sv.ID.String()
	d.DiskID = sv.Disks[0].ID.String()
	d.IPAddress = sv.Interfaces[0].IPAddress

	if d.serverConfig.IsNeedWaitingRestart() {
		// wait for shutdown
		d.waitForServerByState(state.Stopped)
		d.Start()
		d.waitForServerByState(state.Running)
	}

	return nil
}

func (d *Driver) prepareSSHKey() (string, error) {
	if d.SSHKey == "" {
		log.Infof("Creating SSH public key...")
		pKey, err := d.createSSHKey()
		if err != nil {
			return "", err
		}
		return pKey, nil
	}

	log.Info("Importing SSH key...")
	if err := mcnutils.CopyFile(d.SSHKey, d.GetSSHKeyPath()); err != nil {
		return "", fmt.Errorf("unable to copy ssh key: %s", err)
	}
	if err := os.Chmod(d.GetSSHKeyPath(), 0600); err != nil {
		return "", fmt.Errorf("unable to set permissions on the ssh key: %s", err)
	}
	if err := mcnutils.CopyFile(d.SSHKey+".pub", d.GetSSHKeyPath()+".pub"); err != nil {
		return "", fmt.Errorf("unable to copy ssh key: %s", err)
	}

	pKey, err := os.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}

	return string(pKey), nil
}

func (d *Driver) preparePassword() {
	if d.serverConfig.Password == "" {
		d.Client.Password = generateRandomPassword()
		log.Infof("password is not set, generated.[password:%s]", d.Client.Password)
		d.serverConfig.Password = d.Client.Password
	}
}

func generateRandomPassword() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func (d *Driver) waitForServerByState(waitForState state.State) {
	log.Infof("Waiting for server to become %v", waitForState)
	for {
		s, err := d.GetState()
		if err != nil {
			log.Debugf("Failed to get Server State - %+v", err)
			continue
		}

		if s == waitForState {
			break
		} else {
			log.Debugf("Still waiting - state is %s...", s)
		}
		time.Sleep(5 * time.Second)
	}
}

const sakuraAllowSudoScriptBody = `#!/bin/bash
# @sacloud-once
# @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集します
# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
# @sacloud-require-archive distro-debian
# @sacloud-require-archive distro-ubuntu
export DEBIAN_FRONTEND=noninteractive
echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
sh -c 'sleep 10; shutdown -h now' &
exit 0`

const sakuraInstallNetToolsScriptBody = `#!/bin/bash
# @sacloud-once
# @sacloud-desc docker-machine用のプロビジョニング準備を行います
# @sacloud-desc （このスクリプトは、CentOSでのみ動作します）
# @sacloud-require-archive distro-centos
yum install -y net-tools || exit 1
firewall-cmd --zone=public --add-port=%d/tcp --permanent || exit 1
sh -c 'sleep 10; shutdown -h now' &
exit 0`

func (d *Driver) buildSakuraServerSpec(publicKey string) *server.Builder {

	var interfaceDriver types.EInterfaceDriver
	switch d.serverConfig.InterfaceDriver {
	case "virtio":
		interfaceDriver = types.InterfaceDrivers.VirtIO
	case "e1000":
		interfaceDriver = types.InterfaceDrivers.E1000
	}

	var ost ostype.UnixPublicArchiveType
	switch d.serverConfig.OSType {
	case "rancheros":
		ost = ostype.RancherOS
	case "centos":
		ost = ostype.CentOS
	case "ubuntu":
		ost = ostype.Ubuntu
	case "coreos":
		ost = ostype.CoreOS
	}

	var diskPlan types.ID
	switch d.serverConfig.DiskPlan {
	case "ssd":
		diskPlan = types.DiskPlans.SSD
	case "hdd":
		diskPlan = types.DiskPlans.HDD
	}

	var diskConn types.EDiskConnection
	switch d.serverConfig.DiskConnection {
	case "virtio":
		diskConn = types.DiskConnections.VirtIO
	case "ide":
		diskConn = types.DiskConnections.IDE
	}

	var notes []string
	if d.serverConfig.IsUbuntu() {
		// add startup-script for allow sudo by ubuntu user
		notes = append(notes, sakuraAllowSudoScriptBody)
	} else if d.serverConfig.IsCentOS() {
		notes = append(notes, fmt.Sprintf(sakuraInstallNetToolsScriptBody, d.EnginePort))
	}

	diskBuilder := &server.FromUnixDiskBuilder{
		OSType: ost,
		Name:   d.serverConfig.HostName,
		SizeGB: d.serverConfig.DiskSize,
		//DistantFrom:   nil,
		PlanID:     diskPlan,
		Connection: diskConn,
		// Description:   "",
		// Tags:          nil,
		// IconID:        0,
		EditParameter: &server.UnixDiskEditRequest{
			HostName:            d.serverConfig.HostName,
			Password:            d.serverConfig.Password,
			DisablePWAuth:       !d.serverConfig.EnablePWAuth,
			EnableDHCP:          false,
			ChangePartitionUUID: true,
			// IPAddress:                 "",
			// NetworkMaskLen:            0,
			// DefaultRoute:              "",
			SSHKeys:            []string{publicKey},
			IsSSHKeysEphemeral: false,
			IsNotesEphemeral:   true,
			Notes:              notes,
		},
	}

	builder := &server.Builder{
		Name:            d.serverConfig.HostName,
		CPU:             d.serverConfig.Core,
		MemoryGB:        d.serverConfig.Memory,
		Commitment:      types.Commitments.Standard, // TODO パラメータ化
		Generation:      types.PlanGenerations.Default,
		InterfaceDriver: interfaceDriver,
		//Description:     "",
		//IconID:          0,
		//Tags:            nil,
		BootAfterCreate: true,
		//CDROMID:         0,
		//PrivateHostID:   0,
		NIC: &server.SharedNICSetting{
			PacketFilterID: types.StringID(d.serverConfig.PacketFilter),
		},
		//AdditionalNICs: nil,
		DiskBuilders: []server.DiskBuilder{diskBuilder},
	}

	log.Debugf("Build host spec %#v", builder)
	return builder
}

func (d *Driver) createSSHKey() (string, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}

	publicKey, err := os.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

// Kill force power off
func (d *Driver) Kill() error {
	return d.getClient().PowerOff(d.ID)
}

// Remove remove server
func (d *Driver) Remove() error {
	log.Infof("Removing sakura cloud server ...")

	err := d.Kill()
	if err != nil {
		log.Errorf("Error stopping server: %v", err)
	} else {
		d.waitForServerByState(state.Stopped)
	}

	err = d.getClient().Delete(d.ID, []string{d.DiskID})
	if err != nil {
		log.Errorf("Error deleting server: %v", err)
	} else {
		log.Infof("Removed sakura cloud server.")
	}

	return nil
}

// Restart restart server(call PowerOFf and PowerOn)
func (d *Driver) Restart() error {
	// PowerOff
	d.getClient().PowerOff(d.ID)

	// wait
	d.waitForServerByState(state.Stopped)

	//poweron
	d.getClient().PowerOn(d.ID)

	//wait
	d.waitForServerByState(state.Running)

	//return
	return nil
}

// Start power on server
func (d *Driver) Start() error {
	return d.getClient().PowerOn(d.ID)
}

// Stop power off server
func (d *Driver) Stop() error {
	return d.getClient().PowerOff(d.ID)
}
