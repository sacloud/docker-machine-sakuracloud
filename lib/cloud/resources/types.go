package resources

// Resource type of sakuracloud resource(have ID:string)
type Resource struct {
	ID string `json:",omitempty"`
}

// NumberResource type of sakuracloud resource(int64)
type NumberResource struct {
	ID int64 `json:",omitempty"`
}

// EAvailability Enum of sakuracloud
type EAvailability struct {
	Availability string `json:",omitempty"`
}

// IsAvailable Return availability
func IsAvailable(a *EAvailability) bool {
	return a.Availability == "available"
}

// SakuraCloudResources type of resources
type SakuraCloudResources struct {
	Server            *Server            `json:",omitempty"`
	Disk              *Disk              `json:",omitempty"`
	Note              *Note              `json:",omitempty"`
	PacketFilter      *PacketFilter      `json:",omitempty"`
	ServerPlan        *ServerPlan        `json:",omitempty"`
	CommonServiceItem *CommonServiceItem `json:",omitempty"`
}

// SakuraCloudResourceList type of resources
type SakuraCloudResourceList struct {
	Servers            []Server            `json:",omitempty"`
	Notes              []Note              `json:",omitempty"`
	Archives           []Archive           `json:",omitempty"`
	PacketFilters      []PacketFilter      `json:",omitempty"`
	ServerPlans        []ServerPlan        `json:",omitempty"`
	CommonServiceItems []CommonServiceItem `json:",omitempty"`
}

// Server type of create server request values
type Server struct {
	*Resource
	Name              string
	HostName          string              `json:",omitempty"`
	Icon              NumberResource      `json:",omitempty"`
	Description       string              `json:",omitempty"`
	ServerPlan        NumberResource      `json:",omitempty"`
	Tags              []string            `json:",omitempty"`
	ConnectedSwitches []map[string]string `json:",omitempty"`
	Disks             []Disk              `json:",omitempty"`
	Interfaces        []Interface         `json:",omitempty"`
	Instance          *Instance           `json:",omitempty"`
}

// Instance type of instance
type Instance struct {
	Status string `json:",omitempty"`
}

// Disk type of disk
type Disk struct {
	*Resource
	Name           string
	Plan           NumberResource `json:",omitempty"`
	SizeMB         int
	Connection     string   `json:",omitempty"`
	SourceArchive  Resource `json:",omitempty"`
	ReinstallCount int      `json:",omitempty"`
	*EAvailability
}

// SSHKey type of sshkey
type SSHKey struct {
	PublicKey string
}

// DiskEditValue type of disk edit request value
type DiskEditValue struct {
	Password      string
	SSHKey        SSHKey
	DisablePWAuth bool
	Notes         []Resource
}

// Interface type of server nic
type Interface struct {
	*Resource
	IPAddress     string `json:",omitempty"`
	UserIPAddress string `json:",omitempty"`
	MACAddress    string `json:",omitempty"`
}

// Note type of startup script
type Note struct {
	*Resource
	Name    string
	Content string
	*EAvailability
}

// Archive type of Public Archive
type Archive struct {
	*Resource
}

// PacketFilter type of PacketFilter
type PacketFilter struct {
	*Resource
	Name                string
	Description         string `json:",omitempty"`
	RequiredHostVersion string `json:",omitempty"`
	//	Expression          string `json:",omitempty"`
	Notice string `json:",omitempty"`
}

// ServerPlan type of ServerPlan
type ServerPlan struct {
	*NumberResource
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	CPU          int    `json:",omitempty"`
	MemoryMB     int    `json:",omitempty"`
	ServiceClass string `json:",omitempty"`
	*EAvailability
}

// CommonServiceItem type of CommonServiceItem
type CommonServiceItem struct {
	*Resource
	Name        string
	Description string      `json:",omitempty"`
	Settings    DnsSetting  `json:",omitempty"`
	Status      DnsStatus   `json:",omitempty"`
	Provider    DnsProvider `json:",omitempty"`
}

// CreateNewCommonServiceItem Create new CommonServiceItem
func CreateNewCommonServiceItem(zoneName string) *CommonServiceItem {
	return &CommonServiceItem{
		Resource: &Resource{ID: ""},
		Name:     zoneName,
		Status: DnsStatus{
			Zone: zoneName,
		},
		Settings: DnsSetting{
			DNS: DnsRecordSets{},
		},
		Provider: DnsProvider{
			Class: "dns",
		},
	}

}

func (d *CommonServiceItem) HasRecord() bool {
	return len(d.Settings.DNS.ResourceRecordSets) > 0
}

type DnsSetting struct {
	DNS DnsRecordSets `json:",omitempty"`
}

type DnsRecordSets struct {
	ResourceRecordSets []DnsRecordSet
}

func (d *DnsRecordSets) AddDnsRecordSet(name string, ip string) {
	var record DnsRecordSet
	var isExist = false
	for i := range d.ResourceRecordSets {
		if d.ResourceRecordSets[i].Name == name && d.ResourceRecordSets[i].Type == "A" {
			d.ResourceRecordSets[i].RData = ip
			isExist = true
		}
	}

	if !isExist {
		record = DnsRecordSet{
			Name:  name,
			Type:  "A",
			RData: ip,
		}
		d.ResourceRecordSets = append(d.ResourceRecordSets, record)
	}
}

func (d *DnsRecordSets) DeleteDnsRecordSet(name string, ip string) {
	res := []DnsRecordSet{}
	for i := range d.ResourceRecordSets {
		if d.ResourceRecordSets[i].Name != name || d.ResourceRecordSets[i].Type != "A" || d.ResourceRecordSets[i].RData != ip {
			res = append(res, d.ResourceRecordSets[i])
		}
	}

	d.ResourceRecordSets = res
}

type DnsRecordSet struct {
	Name  string `json:",omitempty"`
	Type  string `json:",omitempty"`
	RData string `json:",omitempty"`
	TTL   string `json:",omitempty"`
}

type DnsStatus struct {
	Zone string   `json:",omitempty"`
	NS   []string `json:",omitempty"`
}

type DnsProvider struct {
	Class string `json:",omitempty"`
}
