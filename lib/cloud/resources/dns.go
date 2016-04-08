package resources

// CommonServiceItem type of CommonServiceItem
type CommonServiceDnsItem struct {
	*Resource
	Name        string
	Description string                   `json:",omitempty"`
	Status      CommonServiceDnsStatus   `json:",omitempty"`
	Provider    CommonServiceDnsProvider `json:",omitempty"`
	Settings    CommonServiceDnsSettings `json:",omitempty"`
}

type CommonServiceDnsSettings struct {
	DNS DnsRecordSets `json:",omitempty"`
}

type CommonServiceDnsStatus struct {
	Zone string   `json:",omitempty"`
	NS   []string `json:",omitempty"`
}

type CommonServiceDnsProvider struct {
	Class string `json:",omitempty"`
}

// CreateNewDnsCommonServiceItem Create new CommonServiceItem
func CreateNewDnsCommonServiceItem(zoneName string) *CommonServiceDnsItem {
	return &CommonServiceDnsItem{
		Resource: &Resource{ID: ""},
		Name:     zoneName,
		Status: CommonServiceDnsStatus{
			Zone: zoneName,
		},
		Provider: CommonServiceDnsProvider{
			Class: "dns",
		},
		Settings: CommonServiceDnsSettings{
			DNS: DnsRecordSets{},
		},
	}

}

func (d *CommonServiceDnsItem) HasDnsRecord() bool {
	return len(d.Settings.DNS.ResourceRecordSets) > 0
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
	TTL   int    `json:",omitempty"`
}
