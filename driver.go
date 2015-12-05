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
)

const (
	defaultRegion               = "is1a"          // 石狩第１ゾーン
	defaultPlan                 = "1001"          //TODO プラン名称から設定できるようにする? or コアとメモリを個別に指定できるようにする?
	defaultConnectedSwitch      = ""              // 追加で接続するSwitchのID
	defaultAdditionalIP         = ""              // 追加で接続するSwitch用NICのIP
	defaultAdditionalSubnetMask = "255.255.255.0" // 追加で接続するSwitch用NICのIP
	defaultDiskPlan             = "4"             // SSD
	defaultDiskSize             = 20480           // 20GB
	defaultDiskName             = "disk001"       // ディスク名
	defaultDiskConnection       = "virtio"        // virtio
)

type Driver struct {
	*drivers.BaseDriver
	serverConfig *sakuraServerConfig
	Client       *Client
	Id           string
	DiskId       string
}

type sakuraServerConfig struct {
	HostName             string
	Plan                 string
	ConnectedSwitch      string
	AdditionalIP         string
	AdditionalSubnetMask string
	DiskPlan             string
	DiskSize             int
	DiskName             string
	DiskConnection       string
	DiskSourceArchiveId  string
	Password             string
}

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
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_PLAN",
			Name:   "sakuracloud-plan",
			Usage:  "sakuracloud plan id [memory(GB) & core(NNN)]",
			Value:  defaultPlan,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_CONNECTED_SWITCH",
			Name:   "sakuracloud-connected-switch",
			Usage:  "sakuracloud connected switch['switch ID']",
			Value:  defaultConnectedSwitch,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_ADDITIONAL_IP",
			Name:   "sakuracloud-additional-ip",
			Usage:  "sakuracloud additional ip['xxx.xxx.xxx.xxx']",
			Value:  defaultAdditionalIP,
		},
		mcnflag.StringFlag{
			EnvVar: "SAKURACLOUD_ADDITIONAL_SUBNET_MASK",
			Name:   "sakuracloud-additional-subnet-mask",
			Usage:  "sakuracloud additional subnetmask['255.255.255.0']",
			Value:  defaultAdditionalIP,
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
	}
}

func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Client: &Client{},
		serverConfig: &sakuraServerConfig{
			Plan:                 defaultPlan,
			ConnectedSwitch:      defaultConnectedSwitch,
			AdditionalIP:         defaultAdditionalIP,
			AdditionalSubnetMask: defaultAdditionalSubnetMask,
			DiskPlan:             defaultDiskPlan,
			DiskSize:             defaultDiskSize,
			DiskName:             defaultDiskName,
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

func validateSakuraServerConfig(c *sakuraServerConfig) error {
	//TODO さくら用設定のバリデーション

	//ex. プランの存在確認や矛盾した設定の検出など
	if c.ConnectedSwitch != "" && (c.AdditionalIP == "" || c.AdditionalSubnetMask == "") {
		return fmt.Errorf("Missing Additional IP or subnet --sakuracloud-additional-ip/subnet-mask")
	}

	return nil
}

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
		HostName:             flags.String("sakuracloud-host-name"),
		Plan:                 flags.String("sakuracloud-plan"),
		ConnectedSwitch:      flags.String("sakuracloud-connected-switch"),
		AdditionalIP:         flags.String("sakuracloud-additional-ip"),
		AdditionalSubnetMask: flags.String("sakuracloud-additional-subnet-mask"),
		DiskPlan:             flags.String("sakuracloud-disk-plan"),
		DiskSize:             flags.Int("sakuracloud-disk-size"),
		DiskName:             flags.String("sakuracloud-disk-name"),
		DiskConnection:       flags.String("sakuracloud-disk-connection"),
		DiskSourceArchiveId:  flags.String("sakuracloud-disk-source-archive-id"),
	}

	if d.serverConfig.HostName == "" {
		d.serverConfig.HostName = d.GetMachineName()
	}

	if d.serverConfig.DiskSourceArchiveId == "" {
		archiveId, err := d.getClient().VirtualGuest().GetUbuntuArchiveId()
		if err != nil {
			return err
		}
		d.serverConfig.DiskSourceArchiveId = archiveId
	}

	return validateSakuraServerConfig(d.serverConfig)
}

func (d *Driver) getClient() *Client {
	return d.Client
}

func (d *Driver) DriverName() string {
	return "sakuracloud"
}

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

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	return d.getClient().VirtualGuest().GetIP(d.Id)
}

func (d *Driver) GetState() (state.State, error) {
	s, err := d.getClient().VirtualGuest().State(d.Id)
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
	id, err := d.getClient().VirtualGuest().Create(spec, d.serverConfig.AdditionalIP)
	if err != nil {
		return fmt.Errorf("Error creating host: %v", err)
	}
	log.Infof("Created Server ID: %s", id)
	d.Id = id

	// FIXME
	// workaround for [Non root ssh create sudo can't get password](https://github.com/docker/machine/issues/1569)
	// [PR #1586](https://github.com/docker/machine/pull/1586)がマージされるまで暫定
	// スクリプト(Note)を使ってubuntuユーザがsudo可能にする
	//setup note(script)
	noteId, err := d.getClient().VirtualGuest().GetUbuntuCustomizeNoteId()
	if err != nil || noteId == "" {
		return fmt.Errorf("Error creating custom note: %v", err)
	}

	var addIpNoteId = ""

	if d.serverConfig.ConnectedSwitch != "" {
		var err error
		addIpNoteId, err = d.getClient().VirtualGuest().GetAddIPCustomizeNoteId(d.serverConfig.AdditionalIP, d.serverConfig.AdditionalSubnetMask)
		if err != nil {
			return fmt.Errorf("Error creating custom note: %v", err)
		}
	}

	// create disk( from public archive 'Ubuntu')
	diskSpec := d.buildSakuraDiskSpec()
	diskId, err := d.getClient().VirtualGuest().CreateDisk(diskSpec)
	if err != nil {
		return fmt.Errorf("Error creating disk: %v", err)
	}
	log.Infof("Created Disk ID: %v", diskId)
	d.DiskId = diskId

	//wait for disk available
	d.waitForDiskAvailable()

	//connect disk for server
	connectSuccess, err := d.getClient().VirtualGuest().ConnectDisk(diskId, id)
	if err != nil || !connectSuccess {
		return fmt.Errorf("Error connecting disk to server: %v", err)
	}

	//edit disk
	editDiskSpec := d.buildSakuraDiskEditSpec(publicKey, noteId, addIpNoteId)
	editSuccess, err := d.getClient().VirtualGuest().EditDisk(diskId, editDiskSpec)
	if err != nil || !editSuccess {
		return fmt.Errorf("Error editting disk: %v", err)
	}
	log.Infof("Editted Disk Id: %v", diskId)
	d.waitForDiskAvailable()

	//start
	err = d.getClient().VirtualGuest().PowerOn(id)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	//wait for startup
	d.waitForServerByState(state.Running)

	time.Sleep(10 * time.Second)

	//wait for applay startup script and shutdown
	d.waitForServerByState(state.Stopped)

	//restart
	err = d.getClient().VirtualGuest().PowerOn(id)
	if err != nil {
		return fmt.Errorf("Error starting server: %v", err)
	}
	//wait for startup
	d.waitForServerByState(state.Running)

	if addIpNoteId != "" {
		err = d.getClient().VirtualGuest().DeleteNote(addIpNoteId)
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
		s, err := d.getClient().VirtualGuest().DiskState(d.DiskId)
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

func (d *Driver) buildSakuraServerSpec() *Server {

	var network []map[string]string
	if d.serverConfig.ConnectedSwitch == "" {
		network = []map[string]string{{"Scope": "shared"}}
	} else {
		network = []map[string]string{{"Scope": "shared"}, {"ID": d.serverConfig.ConnectedSwitch}}
	}

	spec := &Server{
		Name:        d.serverConfig.HostName,
		Description: "",
		ServerPlan: Resource{
			ID: d.serverConfig.Plan,
		},
		ConnectedSwitches: network,
	}

	log.Infof("Build host spec %#v", spec)
	return spec
}
func (d *Driver) buildSakuraDiskSpec() *Disk {
	spec := &Disk{
		Name: d.serverConfig.DiskName,
		Plan: Resource{
			ID: d.serverConfig.DiskPlan,
		},
		SizeMB:     d.serverConfig.DiskSize,
		Connection: d.serverConfig.DiskConnection,
		SourceArchive: Resource{
			ID: d.serverConfig.DiskSourceArchiveId,
		},
	}

	log.Infof("Build disk spec %#v", spec)
	return spec
}

func (d *Driver) buildSakuraDiskEditSpec(publicKey string, noteId string, addIpNoteId string) *DiskEditValue {

	var notes []Resource
	if addIpNoteId == "" {
		notes = []Resource{
			Resource{ID: noteId},
		}
	} else {
		notes = []Resource{
			Resource{ID: noteId},
			Resource{ID: addIpNoteId},
		}

	}

	spec := &DiskEditValue{
		Password: d.serverConfig.Password,
		SSHKey: SSHKey{
			PublicKey: publicKey,
		},
		Notes: notes,
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

func (d *Driver) Kill() error {
	return d.getClient().VirtualGuest().PowerOff(d.Id)
}

func (d *Driver) Remove() error {
	log.Infof("Removing sakura cloud server ...")

	err := d.Stop()
	if err != nil {
		log.Errorf("Error stopping server: %v", err)
	} else {
		d.waitForServerByState(state.Stopped)
	}

	err = d.getClient().VirtualGuest().Delete(d.Id, []string{d.DiskId})
	if err != nil {
		log.Errorf("Error deleting server: %v", err)
	} else {
		log.Infof("Removed sakura cloud server.")
	}

	return nil
}
func (d *Driver) Restart() error {
	// PowerOff
	d.getClient().VirtualGuest().PowerOff(d.Id)

	// wait
	d.waitForServerByState(state.Stopped)

	//poweron
	d.getClient().VirtualGuest().PowerOn(d.Id)

	//wait
	d.waitForServerByState(state.Running)

	//return
	return nil
}

func (d *Driver) Start() error {
	return d.getClient().VirtualGuest().PowerOn(d.Id)
}
func (d *Driver) Stop() error {
	return d.getClient().VirtualGuest().PowerOff(d.Id)
}
