package driver

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
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
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
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
		return fmt.Errorf("Invalid Parameter: %s", err)
	}

	res, err := c.IsValidPlan(config.Core, config.Memory)
	if !res || err != nil {
		return fmt.Errorf("Invalid Parameter: core or memory is invalid : %v", err)
	}

	if config.PacketFilter != "" {
		id, valid := sakuracloud.ToSakuraID(config.PacketFilter)
		if !valid {
			return fmt.Errorf("Invalid Parameter: invalid packet-filter-id")
		}
		exists, err := c.IsExistsPacketFilter(id)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Invalid Parameter: packet-filter[id:%d] is not exists", id)
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
	sb := d.buildSakuraServerSpec(publicKey)
	serverResponse, err := sb.Build()
	if err != nil {
		return fmt.Errorf("Error creating host: %v", err)
	}
	d.ID = serverResponse.Server.GetStrID()
	d.DiskID = serverResponse.Disks[0].Disk.GetStrID()
	d.IPAddress = serverResponse.Server.IPAddress()

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
	} else {
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

		pKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
		if err != nil {
			return "", err
		}

		return string(pKey), nil
	}

}

func (d *Driver) preparePassword() {
	if d.serverConfig.Password == "" {
		d.serverConfig.Password = generateRandomPassword()
		log.Infof("password is not set, generated.[password:%s]", d.serverConfig.Password)
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

func (d *Driver) buildSakuraServerSpec(publicKey string) builder.PublicArchiveUnixServerBuilder {

	name := d.serverConfig.HostName
	b := d.getClient().ServerBuilder(
		d.serverConfig.OSType,
		name,
		d.serverConfig.Password,
	)

	// set server spec
	b.SetServerName(name)
	b.SetHostName(name)
	b.SetCore(d.serverConfig.Core)
	b.SetMemory(d.serverConfig.Memory)
	b.AddPublicNWConnectedNIC()

	switch d.serverConfig.InterfaceDriver {
	case "virtio":
		b.SetInterfaceDriver(sacloud.InterfaceDriverVirtIO)
	case "e1000":
		b.SetInterfaceDriver(sacloud.InterfaceDriverE1000)
	}

	// set disk spec
	b.SetDiskSize(d.serverConfig.DiskSize)
	b.SetDiskPlan(d.serverConfig.DiskPlan)
	switch d.serverConfig.DiskConnection {
	case "virtio":
		b.SetDiskConnection(sacloud.DiskConnectionVirtio)
	case "ide":
		b.SetDiskConnection(sacloud.DiskConnectionIDE)
	}

	// edit disk params
	b.AddSSHKey(publicKey)
	b.SetDisablePWAuth(!d.serverConfig.EnablePWAuth)
	if d.serverConfig.IsUbuntu() {
		// add startup-script for allow sudo by ubuntu user
		b.AddNote(sakuraAllowSudoScriptBody)
	} else if d.serverConfig.IsCentOS() {
		b.AddNote(fmt.Sprintf(sakuraInstallNetToolsScriptBody, d.EnginePort))
	}

	b.SetNotesEphemeral(true)
	b.SetSSHKeysEphemeral(true)

	// event handlers(for logging)
	b.SetEventHandler(builder.ServerBuildOnCreateServerBefore, func(_ *builder.ServerBuildValue, _ *builder.ServerBuildResult) {
		log.Infof("Creating server...")
	})
	b.SetEventHandler(builder.ServerBuildOnBootBefore, func(_ *builder.ServerBuildValue, _ *builder.ServerBuildResult) {
		log.Infof("Booting server...")
	})
	b.SetDiskEventHandler(builder.DiskBuildOnCreateDiskBefore, func(_ *builder.DiskBuildValue, _ *builder.DiskBuildResult) {
		log.Infof("Creating disk...")
	})

	log.Debugf("Build host spec %#v", b)
	return b
}

func (d *Driver) createSSHKey() (string, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}

	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
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
