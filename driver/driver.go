package driver

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/api"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cli"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/spec"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

// Driver sakuracloud driver
type Driver struct {
	*drivers.BaseDriver
	serverConfig *spec.SakuraServerConfig
	Client       *api.APIClient
	ID           string
	DiskID       string
	EnginePort   int
	SSHKey       string
	DNSZone      string
	GSLB         string
}

// GetCreateFlags create flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return spec.Options.McnFlags()
}

// NewDriver create driver instance
func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Client:       &api.APIClient{},
		serverConfig: spec.DefaultServerConfig,
		EnginePort:   spec.DefaultServerConfig.EnginePort,
		DNSZone:      spec.DefaultServerConfig.DNSZone,
		GSLB:         spec.DefaultServerConfig.GSLB,
	}
}

func validateSakuraServerConfig(c *api.APIClient, config *spec.SakuraServerConfig) error {
	//さくら用設定のバリデーション

	//ex. プランの存在確認や矛盾した設定の検出など
	if config.ConnectedSwitch != "" && config.PrivateIP == "" {
		return fmt.Errorf("Missing Private IP --sakuracloud-private-ip")
	}

	if config.PrivateIPOnly && config.PrivateIP == "" {
		return fmt.Errorf("Missing Private IP --sakuracloud-private-ip")
	}

	if config.PrivatePacketFilter != "" && config.PrivateIP == "" {
		return fmt.Errorf("Missing Private IP --sakuracloud-private-ip")
	}

	res, err := c.IsValidPlan(config.Core, config.MemorySize)
	if !res || err != nil {
		return fmt.Errorf("Invalid Parameter: core or memory is invalid : %v", err)
	}

	if config.PrivateIPOnly && config.GSLB != "" {
		return fmt.Errorf("GSLB Needs Global IP. Please unset --sakuracloud-private-ip-only.")
	}

	return nil
}

// SetConfigFromFlags create config values from flags
func (d *Driver) SetConfigFromFlags(srcFlags drivers.DriverOptions) error {
	cliClient := cli.NewClient()
	flags := cliClient.GetDriverOptions(srcFlags)

	d.Client = api.NewAPIClient(
		flags.String("sakuracloud-access-token"),
		flags.String("sakuracloud-access-token-secret"),
		flags.String("sakuracloud-region"),
	)

	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SSHUser = "ubuntu"
	d.SSHPort = 22
	d.SSHKey = flags.String("sakuracloud-ssh-key")

	if err := d.getClient().ValidateClientConfig(); err != nil {
		return err
	}

	d.serverConfig = &spec.SakuraServerConfig{
		HostName:            "",
		Core:                flags.Int("sakuracloud-core"),
		MemorySize:          flags.Int("sakuracloud-memory-size"),
		ConnectedSwitch:     flags.String("sakuracloud-connected-switch"),
		PrivateIP:           flags.String("sakuracloud-private-ip"),
		PrivateIPSubnetMask: flags.String("sakuracloud-private-ip-subnet-mask"),
		PrivateIPOnly:       flags.Bool("sakuracloud-private-ip-only"),
		Gateway:             flags.String("sakuracloud-gateway"),
		DiskPlan:            flags.String("sakuracloud-disk-plan"),
		DiskSize:            flags.Int("sakuracloud-disk-size"),
		DiskName:            flags.String("sakuracloud-disk-name"),
		DiskConnection:      flags.String("sakuracloud-disk-connection"),
		Group:               flags.String("sakuracloud-group"),
		AutoReboot:          flags.Bool("sakuracloud-auto-reboot"),
		IgnoreVirtioNet:     flags.Bool("sakuracloud-ignore-virtio-net"),
		PacketFilter:        flags.String("sakuracloud-packet-filter"),
		PrivatePacketFilter: flags.String("sakuracloud-private-packet-filter"),
		EnablePWAuth:        flags.Bool("sakuracloud-enable-password-auth"),
	}

	if d.serverConfig.HostName == "" {
		d.serverConfig.HostName = d.GetMachineName()
	}

	if d.serverConfig.IsDiskNameDefault() {
		d.serverConfig.DiskName = d.GetMachineName()
	}

	if d.serverConfig.PrivateIPOnly {
		d.IPAddress = d.serverConfig.PrivateIP
	}

	if d.serverConfig.DiskSourceArchiveID == "" {
		archiveID, err := d.getClient().GetUbuntuArchiveID()
		if err != nil {
			return err
		}
		d.serverConfig.DiskSourceArchiveID = archiveID
	}

	d.EnginePort = flags.Int("sakuracloud-engine-port")
	d.DNSZone = flags.String("sakuracloud-dns-zone")
	d.GSLB = flags.String("sakuracloud-gslb")

	return validateSakuraServerConfig(d.Client, d.serverConfig)
}

func (d *Driver) getClient() *api.APIClient {
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

	if d.serverConfig.PrivateIPOnly {
		return d.serverConfig.PrivateIP, nil
	}

	return d.getClient().GetIP(d.ID, d.serverConfig.PrivateIPOnly)

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
	spec := d.buildSakuraServerSpec()

	var publicKey = ""
	if d.SSHKey == "" {
		log.Infof("Creating SSH public key...")
		pKey, err := d.createSSHKey()
		if err != nil {
			return err
		}
		publicKey = pKey
	} else {
		log.Info("Importing SSH key...")
		if err := mcnutils.CopyFile(d.SSHKey, d.GetSSHKeyPath()); err != nil {
			return fmt.Errorf("unable to copy ssh key: %s", err)
		}
		if err := os.Chmod(d.GetSSHKeyPath(), 0600); err != nil {
			return fmt.Errorf("unable to set permissions on the ssh key: %s", err)
		}
		if err := mcnutils.CopyFile(d.SSHKey+".pub", d.GetSSHKeyPath()+".pub"); err != nil {
			return fmt.Errorf("unable to copy ssh key: %s", err)
		}

		pKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
		if err != nil {
			return err
		}

		publicKey = string(pKey)
	}

	if d.serverConfig.Password == "" {
		d.serverConfig.Password = generateRandomPassword()
		log.Infof("password is not set, generated password:%s", d.serverConfig.Password)
	}

	//create server
	serverResponse, err := d.getClient().Create(spec, d.serverConfig.PrivateIP)
	if err != nil {
		return fmt.Errorf("Error creating host: %v", err)
	}
	id := serverResponse.ID
	log.Infof("Created Server ID: %s", id)
	d.ID = id

	var noteIDs []string

	noteID, err := d.getClient().GetAllowSudoNoteID(id)
	if err != nil || noteID == "" {
		return fmt.Errorf("Error creating custom note: %v", err)
	}
	noteIDs = append(noteIDs, noteID)

	var addIPNoteID = ""
	if d.serverConfig.PrivateIP != "" {
		var err error
		addIPNoteID, err = d.getClient().GetAddIPCustomizeNoteID(id, d.serverConfig.PrivateIP, d.serverConfig.PrivateIPSubnetMask)
		if err != nil {
			return fmt.Errorf("Error creating custom note: %v", err)
		}

		if addIPNoteID != "" {
			noteIDs = append(noteIDs, addIPNoteID)
		}
	}

	var changeGatewayNoteID = ""
	if d.serverConfig.Gateway != "" {
		var err error
		changeGatewayNoteID, err = d.getClient().GetChangeGatewayCustomizeNoteID(id, d.serverConfig.Gateway)
		if err != nil {
			return fmt.Errorf("Error creating custom note: %v", err)
		}

		if changeGatewayNoteID != "" {
			noteIDs = append(noteIDs, changeGatewayNoteID)
		}

	}

	var disableEth0NoteID = ""
	if d.serverConfig.PrivateIPOnly {
		var err error
		disableEth0NoteID, err = d.getClient().GetDisableEth0CustomizeNoteID(id)
		if err != nil {
			return fmt.Errorf("Error creating custom note: %v", err)
		}

		if disableEth0NoteID != "" {
			noteIDs = append(noteIDs, disableEth0NoteID)
		}

	}

	// create disk( from public archive 'Ubuntu')
	diskSpec := d.buildSakuraDiskSpec()
	diskID, err := d.getClient().CreateDisk(diskSpec)
	if err != nil {
		return fmt.Errorf("Error creating disk: %v", err)
	}
	log.Infof("Created Disk ID: %v", diskID)
	d.DiskID = diskID

	//wait for disk available
	d.waitForDiskAvailable()

	//connect disk for server
	connectSuccess, err := d.getClient().ConnectDisk(diskID, id)
	if err != nil || !connectSuccess {
		return fmt.Errorf("Error connecting disk to server: %v", err)
	}

	//edit disk
	editDiskSpec := d.buildSakuraDiskEditSpec(publicKey, noteIDs[:])
	editSuccess, err := d.getClient().EditDisk(diskID, editDiskSpec)
	if err != nil || !editSuccess {
		return fmt.Errorf("Error editting disk: %v", err)
	}
	log.Infof("Editted Disk Id: %v", diskID)
	d.waitForDiskAvailable()

	//connect packet filter
	if d.serverConfig.PacketFilter != "" {
		log.Infof("Connecting Packet Filter(shared): %v", d.serverConfig.PacketFilter)
		err := d.getClient().ConnectPacketFilterToSharedNIC(serverResponse, d.serverConfig.PacketFilter)
		if err != nil {
			return fmt.Errorf("Error connecting PacketFilter(shared): %v", err)
		}
	}

	if d.serverConfig.PrivatePacketFilter != "" {
		log.Infof("Connecting Packet Filter(private): %v", d.serverConfig.PrivatePacketFilter)
		err := d.getClient().ConnectPacketFilterToPrivateNIC(serverResponse, d.serverConfig.PrivatePacketFilter)
		if err != nil {
			return fmt.Errorf("Error connecting PacketFilter(prvate): %v", err)
		}

	}

	if d.DNSZone != "" {
		log.Infof("Setting SakuraCloud DNS: %v", d.DNSZone)
		ip, _ := d.GetIP()
		ns, err := d.getClient().SetupDnsRecord(d.DNSZone, d.GetMachineName(), ip)
		if err != nil {
			return fmt.Errorf("Error setting SakuraCloud DNS: %v", err)
		}

		if ns != nil {
			log.Infof("Added DNS Zone,Please Set Whois NameServer to [%s]", strings.Join(ns, ","))
		}

	}

	if d.GSLB != "" {
		log.Infof("Setting SakuraCloud GSLB: %v", d.GSLB)
		ip, _ := d.GetIP()
		fqdn, err := d.getClient().SetupGslbRecord(d.GSLB, ip)
		if err != nil {
			return fmt.Errorf("Error setting SakuraCloud GSLV: %v", err)
		}

		if fqdn != nil {
			log.Infof("Added GSLB,Please Set CNAME Record : ex. 'your-lb-hostname IN CNAME %s'", fqdn)
		}
	}

	//start
	err = d.getClient().PowerOn(id)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	//wait for startup
	d.waitForServerByState(state.Running)

	// wait for reboot (only upgradeKernel option is true)
	//if d.serverConfig.UpgradeKernel {
	//wait for applay startup script and shutdown
	d.waitForServerByState(state.Stopped)

	//restart
	err = d.getClient().PowerOn(id)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	//wait for startup
	d.waitForServerByState(state.Running)
	//}

	//cleanup notes
	for n := range noteIDs {
		err = d.getClient().DeleteNote(noteIDs[n])
		if err != nil {
			return fmt.Errorf("Error deleting note: %v", err)
		}

	}

	return nil
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

func (d *Driver) waitForDiskAvailable() {
	log.Infof("Waiting for disk to become available")
	for {
		s, err := d.getClient().GetDiskByID(d.DiskID)
		if err != nil {
			log.Debugf("Failed to get DiskState - %+v", err)
			continue
		}

		if s.Availability == "available" {
			break
		} else {
			log.Debugf("Still waiting - state is %s...", s)
		}
		time.Sleep(5 * time.Second)
	}
}

func (d *Driver) buildSakuraServerSpec() *sacloud.Server {

	var tags []string
	if !d.serverConfig.IgnoreVirtioNet {
		tags = append(tags, "@virtio-net-pci")
	}
	if d.serverConfig.Group != "" {
		tags = append(tags, fmt.Sprintf("@group=%s", d.serverConfig.Group))
	}
	if d.serverConfig.AutoReboot {
		tags = append(tags, "@auto-reboot")
	}

	spec := &sacloud.Server{
		Name:        d.serverConfig.HostName,
		Description: "",
		Tags:        tags[:],
	}
	spec.SetServerPlanByID(d.serverConfig.GetPlanID())
	spec.AddPublicNWConnectedParam()
	if d.serverConfig.ConnectedSwitch != "" {
		spec.AddExistsSwitchConnectedParam(d.serverConfig.ConnectedSwitch)
	} else if d.serverConfig.PrivateIP != "" {
		spec.AddEmptyConnectedParam()
	}

	log.Debugf("Build host spec %#v", spec)
	return spec
}
func (d *Driver) buildSakuraDiskSpec() *sacloud.Disk {
	spec := &sacloud.Disk{
		Name:       d.serverConfig.DiskName,
		SizeMB:     d.serverConfig.DiskSize,
		Connection: sacloud.EDiskConnection(d.serverConfig.DiskConnection),
	}

	spec.SetSourceArchive(d.serverConfig.DiskSourceArchiveID)
	if d.serverConfig.DiskPlan == "2" {
		spec.SetDiskPlanToHDD()
	} else {
		spec.SetDiskPlanToSSD()
	}

	log.Debugf("Build disk spec %#v", spec)
	return spec
}

func (d *Driver) buildSakuraDiskEditSpec(publicKey string, noteIDs []string) *sacloud.DiskEditValue {
	notes := make([]*sacloud.Resource, len(noteIDs))
	for n := range noteIDs {
		notes[n] = &sacloud.Resource{ID: noteIDs[n]}
	}

	pAuth := !d.serverConfig.EnablePWAuth

	spec := &sacloud.DiskEditValue{
		Password: &d.serverConfig.Password,
		SSHKey: &sacloud.SSHKey{
			PublicKey: publicKey,
		},
		DisablePWAuth: &pAuth,
		Notes:         notes[:],
	}
	log.Debugf("Build disk edit spec %#v", spec)
	return spec
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

	if d.DNSZone != "" {
		log.Infof("Removing DNS Record ...")
		ip, _ := d.GetIP()
		err := d.getClient().DeleteDnsRecord(d.DNSZone, d.GetMachineName(), ip)
		if err != nil {
			log.Errorf("Error deleting dns: %v", err)
		}
		log.Infof("Removed DNS Record.")
	}

	if d.GSLB != "" {
		log.Infof("Removing GSLB server ...")
		ip, _ := d.GetIP()
		err := d.getClient().DeleteGslbServer(d.GSLB, ip)
		if err != nil {
			log.Errorf("Error deleting GSLB: %v", err)
		}
		log.Infof("Removed GSLB. Wait 30s ...")
		time.Sleep(30 * time.Second)
		log.Infof("Done.")
	}

	err := d.Stop()
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
