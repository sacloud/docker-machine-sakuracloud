package spec

import (
	"fmt"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/api"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/option"
	"strconv"
)

const (
	defaultRegion              = "is1a" // 石狩第１ゾーン
	defaultCore                = 1      //デフォルトコア数
	defaultMemorySize          = 1      // デフォルトメモリサイズ
	defaultConnectedSwitch     = ""     // 追加で接続するSwitchのID
	defaultPrivateIPOnly       = false
	defaultPrivateIP           = "" // 追加で接続するSwitch用NICのIP
	defaultPrivateIPSubnetMask = "" // 追加で接続するSwitch用NICのIP
	defaultGateway             = ""
	defaultDiskPlan            = "4"      // SSD
	defaultDiskSize            = 20480    // 20GB
	defaultDiskName            = ""       // ディスク名
	defaultDiskConnection      = "virtio" // virtio
	defaultGroup               = ""       // グループタグ
	defaultAutoReboot          = false    // 自動再起動
	defaultIgnoreVirtioNet     = false    // virtioNICの無効化
	defaultPacketFilter        = ""
	defaultPrivatePacketFilter = ""
	defaultUpgradeKernel       = false
	defaultEnablePWAuth        = false
)

type SakuraServerConfig struct {
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
	EnablePWAuth        bool
}

func (c *SakuraServerConfig) GetPlanID() int64 {
	planID, _ := strconv.ParseInt(fmt.Sprintf("%d%03d", c.MemorySize, c.Core), 10, 64)
	return planID
}

func (c *SakuraServerConfig) IsDiskNameDefault() bool {
	return c.DiskName == defaultDiskName
}

var DefaultServerConfig = &SakuraServerConfig{
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
	EnablePWAuth:        defaultEnablePWAuth,
}

// Options OptionList
var Options = &option.List{
	Options: []option.Option{
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_ACCESS_TOKEN",
				Name:   "sakuracloud-access-token",
				Usage:  "sakuracloud access token",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_ACCESS_TOKEN_SECRET",
				Name:   "sakuracloud-access-token-secret",
				Usage:  "sakuracloud access token secret",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_REGION",
				Name:   "sakuracloud-region",
				Usage:  "sakuracloud region name[tk1a/is1a/is1b/tk1v]",
				Value:  defaultRegion,
			},
			CandidateFunc: func(c *api.Client) []string {
				return []string{"is1a", "is1b", "tk1a"}
			},
			UsageStringsFunc: func(c *api.Client) string {
				return "[is1a / isab / tk1a]"
			},
		},
		option.Option{
			RawMcnOption: mcnflag.IntFlag{
				EnvVar: "SAKURACLOUD_CORE",
				Name:   "sakuracloud-core",
				Usage:  "sakuracloud number of CPU core",
				Value:  defaultCore,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.IntFlag{
				EnvVar: "SAKURACLOUD_MEMORY_SIZE",
				Name:   "sakuracloud-memory-size",
				Usage:  "sakuracloud memory size(GB)",
				Value:  defaultMemorySize,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_CONNECTED_SWITCH",
				Name:   "sakuracloud-connected-switch",
				Usage:  "sakuracloud connected switch['switch ID']",
				Value:  defaultConnectedSwitch,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.BoolFlag{
				EnvVar: "SAKURACLOUD_PRIVATE_IP_ONLY",
				Name:   "sakuracloud-private-ip-only",
				Usage:  "sakuracloud private ip only flag",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_PRIVATE_IP",
				Name:   "sakuracloud-private-ip",
				Usage:  "sakuracloud private ip['xxx.xxx.xxx.xxx']",
				Value:  defaultPrivateIP,
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_PRIVATE_IP_SUBNET_MASK",
				Name:   "sakuracloud-private-ip-subnet-mask",
				Usage:  "sakuracloud private ip subnetmask['255.255.255.0']",
				Value:  defaultPrivateIPSubnetMask,
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_GATEWAY",
				Name:   "sakuracloud-gateway",
				Usage:  "sakuracloud default gateway ip address",
				Value:  defaultGateway,
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_DISK_PLAN",
				Name:   "sakuracloud-disk-plan",
				Usage:  "sakuracloud disk plan[HDD(2)/SSD(4)]",
				Value:  defaultDiskPlan,
			},
			CandidateFunc: func(c *api.Client) []string {
				return []string{"2", "4"}
			},
			UsageStringsFunc: func(c *api.Client) string {
				return "[2(HDD) / 4(SSD)]"
			}},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_DISK_NAME",
				Name:   "sakuracloud-disk-name",
				Usage:  "sakuracloud disk name",
				Value:  defaultDiskName,
			},
		},
		option.Option{
			RawMcnOption: mcnflag.IntFlag{
				EnvVar: "SAKURACLOUD_DISK_SIZE",
				Name:   "sakuracloud-disk-size",
				Usage:  "sakuracloud disk size(MB)[20480,102400,256000,512000]",
				Value:  defaultDiskSize,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_DISK_CONNECTION",
				Name:   "sakuracloud-disk-connection",
				Usage:  "sakuracloud disk connection[virtio/ide]",
				Value:  defaultDiskConnection,
			},
			CandidateFunc: func(c *api.Client) []string {
				return []string{"virtio", "ide"}
			},
			UsageStringsFunc: func(c *api.Client) string {
				return "[virtio / ide]"
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_PASSWORD",
				Name:   "sakuracloud-password",
				Usage:  "sakuracloud user password",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_GROUP",
				Name:   "sakuracloud-group",
				Usage:  "sakuracloud @group tag [a/b/c/d]",
				Value:  defaultGroup,
			},
			CandidateFunc: func(c *api.Client) []string {
				return []string{"a", "b", "c", "d"}
			},
			UsageStringsFunc: func(c *api.Client) string {
				return "[ a / b / c / d ]"
			},
		},
		option.Option{
			RawMcnOption: mcnflag.BoolFlag{
				EnvVar: "SAKURACLOUD_AUTO_REBOOT",
				Name:   "sakuracloud-auto-reboot",
				Usage:  "sakuracloud @auto-reboot tag flag",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.BoolFlag{
				EnvVar: "SAKURACLOUD_IGNORE_VIRTIO_NET",
				Name:   "sakuracloud-ignore-virtio-net",
				Usage:  "sakuracloud ignore @virtio-net-pci tag flag",
			},
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_PACKET_FILTER",
				Name:   "sakuracloud-packet-filter",
				Usage:  "sakuracloud packet-filter for eth0(shared)[filter ID or NAME]",
				Value:  defaultPacketFilter,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.StringFlag{
				EnvVar: "SAKURACLOUD_PRIVATE_PACKET_FILTER",
				Name:   "sakuracloud-private-packet-filter",
				Usage:  "sakuracloud packet-filter for eth1(private)[filter ID or NAME]",
				Value:  defaultPacketFilter,
			},
			// CandidateFunc: func(c *api.Client) []string {
			// 	return []string{"is1a", "is1b", "tk1a"}
			// },
			// UsageStringsFunc: func(c *api.Client) string {
			// 	return "[is1a / isab / tk1a]"
			// },
		},
		option.Option{
			RawMcnOption: mcnflag.BoolFlag{
				EnvVar: "SAKURACLOUD_ENABLE_PASSWORD_AUTH",
				Name:   "sakuracloud-enable-password-auth",
				Usage:  "sakuracloud enable password auth flag",
			},
		},
	},
}
