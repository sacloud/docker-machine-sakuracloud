package sakuracloud

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
)

const (
	defaultRegion              = "is1a" // 石狩第１ゾーン
	defaultCore                = 1      //デフォルトコア数
	defaultMemorySize          = 1      // デフォルトメモリサイズ
	defaultConnectedSwitch     = ""     // 追加で接続するSwitchのID
	defaultPrivateIPOnly       = false
	defaultPrivateIP           = ""              // 追加で接続するSwitch用NICのIP
	defaultPrivateIPSubnetMask = "255.255.255.0" // 追加で接続するSwitch用NICのIP
	defaultGateway             = ""
	defaultDiskPlan            = "4"       // SSD
	defaultDiskSize            = 20480     // 20GB
	defaultDiskName            = "disk001" // ディスク名
	defaultDiskConnection      = "virtio"  // virtio
	defaultGroup               = ""        // グループタグ
	defaultAutoReboot          = false     // 自動再起動
	defaultIgnoreVirtioNet     = false     // virtioNICの無効化
	defaultPacketFilter        = ""
	defaultPrivatePacketFilter = ""
	defaultUpgradeKernel       = false
)

// Driver sakuracloud driver
type Driver struct {
	*drivers.BaseDriver
	serverConfig *sakuraServerConfig
	Client       *Client
	ID           string
	DiskID       string
}

type sakuraServerConfig struct {
	HostName            string
	Core                int
	MemorySize          int
	ConnectedSwitch     string
	PrivateIPOnly       bool
	PrivateIP           string
	PrivateIPSubnetMask string
	Gateway             string
	DiskPlan            string
	DiskSize            int
	DiskName            string
	DiskConnection      string
	DiskSourceArchiveID string
	Password            string
	Group               string
	AutoReboot          bool
	IgnoreVirtioNet     bool
	PacketFilter        string
	PrivatePacketFilter string
	UpgradeKernel       bool
}

func (c *sakuraServerConfig) getPlanID() int64 {
	planID, _ := strconv.ParseInt(fmt.Sprintf("%d%03d", c.MemorySize, c.Core), 10, 64)
	return planID
}

// GetCreateFlags create flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_ACCESS_TOKEN",
			Name:   "sakuracloud-access-token",
			Usage:  "sakuracloud access token",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_ACCESS_TOKEN_SECRET",
			Name:   "sakuracloud-access-token-secret",
			Usage:  "sakuracloud access token secret",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_REGION",
			Name:   "sakuracloud-region",
			Usage:  "sakuracloud region name[tk1a/is1a/is1b/tk1v]",
			Value:  defaultRegion,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_HOST_NAME",
			Name:   "sakuracloud-host-name",
			Usage:  "sakuracloud host name",
		},
		mcnflag.IntFlag{
			EnvVar: "SAKURACLOUD_CORE",
			Name:   "sakuracloud-core",
			Usage:  "sakuracloud number of CPU core",
			Value:  defaultCore,
		},
		mcnflag.IntFlag{
			EnvVar: "SAKURACLOUD_MEMORY_SIZE",
			Name:   "sakuracloud-memory-size",
			Usage:  "sakuracloud memory size(GB)",
			Value:  defaultMemorySize,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_CONNECTED_SWITCH",
			Name:   "sakuracloud-connected-switch",
			Usage:  "sakuracloud connected switch['switch ID']",
			Value:  defaultConnectedSwitch,
		},
		mcnflag.BoolFlag{
			EnvVar: "SAKURACLOUD_PRIVATE_IP_ONLY",
			Name:   "sakuracloud-private-ip-only",
			Usage:  "sakuracloud private ip only flag",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PRIVATE_IP",
			Name:   "sakuracloud-private-ip",
			Usage:  "sakuracloud private ip['xxx.xxx.xxx.xxx']",
			Value:  defaultPrivateIP,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PRIVATE_IP_SUBNET_MASK",
			Name:   "sakuracloud-private-ip-subnet-mask",
			Usage:  "sakuracloud private ip subnetmask['255.255.255.0']",
			Value:  defaultPrivateIPSubnetMask,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_GATEWAY",
			Name:   "sakuracloud-gateway",
			Usage:  "sakuracloud default gateway ip address",
			Value:  defaultGateway,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_DISK_PLAN",
			Name:   "sakuracloud-disk-plan",
			Usage:  "sakuracloud disk plan[HDD(2)/SSD(4)]",
			Value:  defaultDiskPlan,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_DISK_NAME",
			Name:   "sakuracloud-disk-name",
			Usage:  "sakuracloud disk name",
			Value:  defaultDiskName,
		},
		mcnflag.IntFlag{
			EnvVar: "SAKURACLOUD_DISK_SIZE",
			Name:   "sakuracloud-disk-size",
			Usage:  "sakuracloud disk size(MB)[20480,102400,256000,512000]",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_DISK_CONNECTION",
			Name:   "sakuracloud-disk-connection",
			Usage:  "sakuracloud disk connection[virtio/ide]",
			Value:  defaultDiskConnection,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_DISK_SOURCE_ARCHIVE_ID",
			Name:   "sakuracloud-disk-source-archive-id",
			Usage:  "sakuracloud disk source archive id",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PASSWORD",
			Name:   "sakuracloud-password",
			Usage:  "sakuracloud user password",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_GROUP",
			Name:   "sakuracloud-group",
			Usage:  "sakuracloud @group tag [a/b/c/d]",
			Value:  defaultGroup,
		},
		mcnflag.BoolFlag{
			EnvVar: "SAKURACLOUD_AUTO_REBOOT",
			Name:   "sakuracloud-auto-reboot",
			Usage:  "sakuracloud @auto-reboot tag flag",
		},
		mcnflag.BoolFlag{
			EnvVar: "SAKURACLOUD_IGNORE_VIRTIO_NET",
			Name:   "sakuracloud-ignore-virtio-net",
			Usage:  "sakuracloud ignore @virtio-net-pci tag flag",
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PACKET_FILTER",
			Name:   "sakuracloud-packet-filter",
			Usage:  "sakuracloud packet-filter for eth0(shared)[filter ID or NAME]",
			Value:  defaultPacketFilter,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PRIVATE_PACKET_FILTER",
			Name:   "sakuracloud-private-packet-filter",
			Usage:  "sakuracloud packet-filter for eth1(private)[filter ID or NAME]",
			Value:  defaultPacketFilter,
		},
		mcnflag.BoolFlag{
			EnvVar: "SAKURACLOUD_UPGRADE_KERNEL",
			Name:   "sakuracloud-upgrade-kernel",
			Usage:  "sakuracloud upgrade kernel flag",
		},
	}
}

// NewDriver create driver instance
func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Client: &Client{},
		serverConfig: &sakuraServerConfig{
			Core:                defaultCore,
			MemorySize:          defaultMemorySize,
			PrivateIPOnly:       defaultPrivateIPOnly,
			ConnectedSwitch:     defaultConnectedSwitch,
			PrivateIP:           defaultPrivateIP,
			PrivateIPSubnetMask: defaultPrivateIPSubnetMask,
			Gateway:             defaultGateway,
			DiskPlan:            defaultDiskPlan,
			DiskSize:            defaultDiskSize,
			DiskName:            defaultDiskName,
			Group:               defaultGroup,
			AutoReboot:          defaultAutoReboot,
			IgnoreVirtioNet:     defaultIgnoreVirtioNet,
			PacketFilter:        defaultPacketFilter,
			PrivatePacketFilter: defaultPrivatePacketFilter,
			UpgradeKernel:       defaultUpgradeKernel,
		},
	}
}

func validateClientConfig(c *Client) error {
	if c.AccessToken == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token")
	}

	if c.AccessTokenSecret == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token-secret")
	}
	return nil
}

func validateSakuraServerConfig(c *Client, config *sakuraServerConfig) error {
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

	return nil
}

// SetConfigFromFlags create config values from flags
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {

	d.Client = &Client{
		Region:            flags.String("sakuracloud-region"),
		AccessToken:       flags.String("sakuracloud-access-token"),
		AccessTokenSecret: flags.String("sakuracloud-access-token-secret"),
	}

	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SSHUser = "ubuntu"
	d.SSHPort = 22

	if err := validateClientConfig(d.Client); err != nil {
		return err
	}

	d.serverConfig = &sakuraServerConfig{
		HostName:            flags.String("sakuracloud-host-name"),
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
		DiskSourceArchiveID: flags.String("sakuracloud-disk-source-archive-id"),
		Group:               flags.String("sakuracloud-group"),
		AutoReboot:          flags.Bool("sakuracloud-auto-reboot"),
		IgnoreVirtioNet:     flags.Bool("sakuracloud-ignore-virtio-net"),
		PacketFilter:        flags.String("sakuracloud-packet-filter"),
		PrivatePacketFilter: flags.String("sakuracloud-private-packet-filter"),
		UpgradeKernel:       flags.Bool("sakuracloud-upgrade-kernel"),
	}

	if d.serverConfig.HostName == "" {
		d.serverConfig.HostName = d.GetMachineName()
	}

	if d.serverConfig.DiskName == defaultDiskName {
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

	return validateSakuraServerConfig(d.Client, d.serverConfig)
}

func (d *Driver) getClient() *Client {
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
	return "tcp://" + ip + ":2376", nil
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

// Create create server on sakuracloud
func (d *Driver) Create() error {
	spec := d.buildSakuraServerSpec()

	log.Infof("Creating SSH key...")
	publicKey, err := d.createSSHKey()
	if err != nil {
		return err
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
	id := serverResponse.Server.ID
	log.Infof("Created Server ID: %s", id)
	d.ID = id

	var noteIDs []string

	if d.serverConfig.UpgradeKernel {
		noteID, err := d.getClient().GetAllowSudoWithKernelUpgradeNoteID(id)
		if err != nil || noteID == "" {
			return fmt.Errorf("Error creating custom note: %v", err)
		}
		noteIDs = append(noteIDs, noteID)

	} else {
		noteID, err := d.getClient().GetAllowSudoNoteID(id)
		if err != nil || noteID == "" {
			return fmt.Errorf("Error creating custom note: %v", err)
		}
		noteIDs = append(noteIDs, noteID)

	}

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
		err := d.getClient().ConnectPacketFilterToSharedNIC(serverResponse.Server, d.serverConfig.PacketFilter)
		if err != nil {
			return fmt.Errorf("Error connecting PacketFilter(shared): %v", err)
		}
	}

	if d.serverConfig.PrivatePacketFilter != "" {
		log.Infof("Connecting Packet Filter(private): %v", d.serverConfig.PrivatePacketFilter)
		err := d.getClient().ConnectPacketFilterToPrivateNIC(serverResponse.Server, d.serverConfig.PrivatePacketFilter)
		if err != nil {
			return fmt.Errorf("Error connecting PacketFilter(prvate): %v", err)
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
		s, err := d.getClient().DiskState(d.DiskID)
		if err != nil {
			log.Debugf("Failed to get DiskState - %+v", err)
			continue
		}

		if s == "available" {
			break
		} else {
			log.Debugf("Still waiting - state is %s...", s)
		}
		time.Sleep(5 * time.Second)
	}
}

func (d *Driver) buildSakuraServerSpec() *sakura.Server {

	var network []map[string]string
	if d.serverConfig.ConnectedSwitch != "" {
		network = []map[string]string{{"Scope": "shared"}, {"ID": d.serverConfig.ConnectedSwitch}}
	} else if d.serverConfig.PrivateIP != "" {
		network = []map[string]string{{"Scope": "shared"}, nil}
	} else {
		network = []map[string]string{{"Scope": "shared"}}
	}

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

	spec := &sakura.Server{
		Name:        d.serverConfig.HostName,
		Description: "",
		ServerPlan: sakura.NumberResource{
			ID: d.serverConfig.getPlanID(),
		},
		ConnectedSwitches: network,
		Tags:              tags[:],
	}

	log.Infof("Build host spec %#v", spec)
	return spec
}
func (d *Driver) buildSakuraDiskSpec() *sakura.Disk {
	diskPlan, _ := strconv.ParseInt(d.serverConfig.DiskPlan, 10, 64)
	spec := &sakura.Disk{
		Name: d.serverConfig.DiskName,
		Plan: sakura.NumberResource{
			ID: diskPlan,
		},
		SizeMB:     d.serverConfig.DiskSize,
		Connection: d.serverConfig.DiskConnection,
		SourceArchive: sakura.Resource{
			ID: d.serverConfig.DiskSourceArchiveID,
		},
	}

	log.Infof("Build disk spec %#v", spec)
	return spec
}

func (d *Driver) buildSakuraDiskEditSpec(publicKey string, noteIDs []string) *sakura.DiskEditValue {
	notes := make([]sakura.Resource, len(noteIDs))
	for n := range noteIDs {
		notes[n] = sakura.Resource{ID: noteIDs[n]}
	}

	spec := &sakura.DiskEditValue{
		Password: d.serverConfig.Password,
		SSHKey: sakura.SSHKey{
			PublicKey: publicKey,
		},
		Notes: notes[:],
	}
	log.Infof("Build disk edit spec %#v", spec)
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
