package spec

import (
	"fmt"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/sacloud/libsacloud/sacloud"
)

var (
	defaultRegion          = "is1b"      // 石狩第2ゾーン
	defaultOSType          = "rancheros" // OSタイプ
	defaultCore            = 1           // デフォルトコア数
	defaultMemorySize      = 1           // デフォルトメモリサイズ
	defaultDiskPlan        = "ssd"       // ディスクプラン(ssd/hdd)
	defaultDiskSize        = 20480       // 20GB
	defaultDiskConnection  = "virtio"    // ディスク接続ドライバ
	defaultInterfaceDriver = "virtio"    // NIC接続ドライバ
	defaultPacketFilter    = "0"
	defaultEnablePWAuth    = false
)

type SakuraServerConfig struct {
	HostName        string
	OSType          string
	Core            int
	Memory          int
	DiskPlan        string
	DiskSize        int
	DiskConnection  string
	InterfaceDriver string
	Password        string
	PacketFilter    string
	EnablePWAuth    bool
	EnginePort      int
}

func (c *SakuraServerConfig) SSHUserName() string {
	switch c.OSType {
	case "ubuntu":
		return "ubuntu"
	case "rancheros":
		return "rancher"
	default:
		return "root"
	}
}

func (c *SakuraServerConfig) IsUbuntu() bool {
	return c.OSType == "ubuntu"
}

var DefaultServerConfig = &SakuraServerConfig{
	Core:         defaultCore,
	Memory:       defaultMemorySize,
	DiskPlan:     defaultDiskPlan,
	DiskSize:     defaultDiskSize,
	PacketFilter: defaultPacketFilter,
	EnablePWAuth: defaultEnablePWAuth,
	EnginePort:   engine.DefaultPort,
}

// McnFlags OptionList
var McnFlags = []mcnflag.Flag{
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
		Usage:  "sakuracloud zone name[is1b/tk1a/is1a]",
		Value:  defaultRegion,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_OS_TYPE",
		Name:   "sakuracloud-os-type",
		Usage:  "sakuracloud os(public-archive) type[centos/ubuntu/debian/rancheros]",
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
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_DISK_PLAN",
		Name:   "sakuracloud-disk-plan",
		Usage:  "sakuracloud disk plan[hdd/ssd]",
		Value:  string(defaultDiskPlan),
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
		Usage:  "sakuracloud disk connection[virtio/ide]",
		Value:  defaultDiskConnection,
	},
	mcnflag.StringFlag{
		EnvVar: "SAKURACLOUD_INTERFACE_DRIVER",
		Name:   "sakuracloud-interface-driver",
		Usage:  "sakuracloud interface(NIC) driver[virtio/e1000]",
		Value:  string(defaultInterfaceDriver),
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
		Value:  fmt.Sprintf("%d", defaultPacketFilter),
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
