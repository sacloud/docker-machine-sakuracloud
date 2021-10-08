package driver

import (
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/mcnflag"
)

var (
	defaultRegion          = "is1b"   // 石狩第2ゾーン
	defaultOSType          = "coreos" // OSタイプ
	defaultCore            = 1        // デフォルトコア数
	defaultMemorySize      = 1        // デフォルトメモリサイズ
	defaultDiskPlan        = "ssd"    // ディスクプラン(ssd/hdd)
	defaultDiskSize        = 20       // 20GB
	defaultDiskConnection  = "virtio" // ディスク接続ドライバ
	defaultInterfaceDriver = "virtio" // NIC接続ドライバ
	defaultPacketFilter    = ""
	defaultEnablePWAuth    = false
)

var (
	allowOSTypes          = []string{"rancheros", "centos", "ubuntu", "coreos"}
	allowDiskPlans        = []string{"hdd", "ssd"}
	allowSSDSizes         = []int{20, 40, 100, 250, 500, 1024, 2048, 4096}
	allowHDDSizes         = []int{40, 60, 80, 100, 250, 500, 750, 1024, 2048, 4096}
	allowDiskConnections  = []string{"virtio", "ide"}
	allowInterfaceDrivers = []string{"virtio", "e1000"}
)

type sakuraServerConfig struct {
	HostName        string
	OSType          string
	Core            int
	Memory          int
	GPU int
	DiskPlan        string
	DiskSize        int
	DiskConnection  string
	InterfaceDriver string
	Password        string
	PacketFilter    string
	EnablePWAuth    bool
	EnginePort      int
}

var defaultServerConfig = &sakuraServerConfig{
	Core:         defaultCore,
	Memory:       defaultMemorySize,
	DiskPlan:     defaultDiskPlan,
	DiskSize:     defaultDiskSize,
	PacketFilter: defaultPacketFilter,
	EnablePWAuth: defaultEnablePWAuth,
	EnginePort:   engine.DefaultPort,
}

func (c *sakuraServerConfig) SSHUserName() string {
	switch c.OSType {
	case "ubuntu":
		return "ubuntu"
	case "rancheros":
		return "rancher"
	case "coreos":
		return "core"
	default:
		return "root"
	}
}

func (c *sakuraServerConfig) IsUbuntu() bool {
	return c.OSType == "ubuntu"
}
func (c *sakuraServerConfig) IsCentOS() bool {
	return c.OSType == "centos"
}

func (c *sakuraServerConfig) IsNeedWaitingRestart() bool {
	return c.IsUbuntu() || c.IsCentOS()
}

func (c *sakuraServerConfig) Validate() error {
	// os-type
	if !c.isStrInValue(c.OSType, allowOSTypes...) {
		return fmt.Errorf("%q must be set to one of [%s]", "--sakuracloud-os-type", strings.Join(allowOSTypes, "/"))
	}

	// disk-plan
	if !c.isStrInValue(c.DiskPlan, allowDiskPlans...) {
		return fmt.Errorf("%q must be set to one of [%s]", "--sakuracloud-disk-plan", strings.Join(allowDiskPlans, "/"))
	}

	// disk-size(per disk-plan)
	var allowDiskSizes []int
	switch c.DiskPlan {
	case "ssd":
		allowDiskSizes = allowSSDSizes
	case "hdd":
		allowDiskSizes = allowHDDSizes
	}
	if !c.isIntInValue(c.DiskSize, allowDiskSizes...) {
		return fmt.Errorf("%q must be set to one of [20(SSD)/40/60(HDD)/80(HDD)/100/250/500/750(HDD)/1024/2048/4096]", "--sakuracloud-disk-size")
	}

	// disk-connection
	if !c.isStrInValue(c.DiskConnection, allowDiskConnections...) {
		return fmt.Errorf("%q must be set to one of [%s]", "--sakuracloud-disk-connection", strings.Join(allowDiskConnections, "/"))
	}

	// interface-driver
	if !c.isStrInValue(c.InterfaceDriver, allowInterfaceDrivers...) {
		return fmt.Errorf("%q must be set to one of [%s]", "--sakuracloud-interface-driver", strings.Join(allowInterfaceDrivers, "/"))
	}

	return nil
}

func (c *sakuraServerConfig) isStrInValue(value string, allows ...string) bool {
	for _, s := range allows {
		if value == s {
			return true
		}
	}
	return false
}

func (c *sakuraServerConfig) isIntInValue(value int, allows ...int) bool {
	for _, i := range allows {
		if value == i {
			return true
		}
	}
	return false
}

// mcnFlags OptionList
var mcnFlags = []mcnflag.Flag{
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
		EnvVar: "SAKURACLOUD_ZONE",
		Name:   "sakuracloud-zone",
		Usage:  "sakuracloud zone name[is1a/is1b/tk1a/tk1b]",
		Value:  defaultRegion,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_OS_TYPE",
		Name:   "sakuracloud-os-type",
		Usage:  fmt.Sprintf("sakuracloud os(public-archive) type[%s]", strings.Join(allowOSTypes, "/")),
		Value:  defaultOSType,
	},
	mcnflag.IntFlag{
		EnvVar: "SAKURACLOUD_CORE",
		Name:   "sakuracloud-core",
		Usage:  "sakuracloud number of CPU core",
		Value:  defaultCore,
	},
	mcnflag.IntFlag{
		EnvVar: "SAKURACLOUD_MEMORY",
		Name:   "sakuracloud-memory",
		Usage:  "sakuracloud memory size(GB)",
		Value:  defaultMemorySize,
	},
	mcnflag.IntFlag{
		EnvVar: "SAKURACLOUD_GPU",
		Name:   "sakuracloud-gpu",
		Usage:  "sakuracloud number of GPUs",
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_DISK_PLAN",
		Name:   "sakuracloud-disk-plan",
		Usage:  "sakuracloud disk plan[hdd/ssd]",
		Value:  defaultDiskPlan,
	},
	mcnflag.IntFlag{
		EnvVar: "SAKURACLOUD_DISK_SIZE",
		Name:   "sakuracloud-disk-size",
		Usage:  "sakuracloud disk size(GB)[20(SSD)/40/60(HDD)/80(HDD)/100/250/500/750(HDD)/1024/2048/4096]",
		Value:  defaultDiskSize,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_DISK_CONNECTION",
		Name:   "sakuracloud-disk-connection",
		Usage:  fmt.Sprintf("sakuracloud disk connection[%s]", strings.Join(allowDiskConnections, "/")),
		Value:  defaultDiskConnection,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_INTERFACE_DRIVER",
		Name:   "sakuracloud-interface-driver",
		Usage:  fmt.Sprintf("sakuracloud interface(NIC) driver[%s]", strings.Join(allowInterfaceDrivers, "/")),
		Value:  defaultInterfaceDriver,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_PASSWORD",
		Name:   "sakuracloud-password",
		Usage:  "sakuracloud user password",
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_PACKET_FILTER",
		Name:   "sakuracloud-packet-filter",
		Usage:  "sakuracloud packet-filter for eth0(shared)[filter ID]",
		Value:  defaultPacketFilter,
	},
	mcnflag.BoolFlag{
		EnvVar: "SAKURACLOUD_ENABLE_PASSWORD_AUTH",
		Name:   "sakuracloud-enable-password-auth",
		Usage:  "sakuracloud enable password auth flag",
	},
	mcnflag.IntFlag{
		EnvVar: "SAKURACLOUD_ENGINE_PORT",
		Name:   "sakuracloud-engine-port",
		Usage:  "Docker engine port",
		Value:  engine.DefaultPort,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_SSH_KEY",
		Name:   "sakuracloud-ssh-key",
		Usage:  "SSH Private Key Path",
		Value:  "",
	},
}
