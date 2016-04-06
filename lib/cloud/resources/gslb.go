package resources

// CommonServiceItem type of CommonServiceItem
type CommonServiceGslbItem struct {
	*Resource
	Name        string
	Description string                    `json:",omitempty"`
	Status      CommonServiceGslbStatus   `json:",omitempty"`
	Provider    CommonServiceGslbProvider `json:",omitempty"`
	Settings    CommonServiceGslbSettings `json:",omitempty"`
}

type CommonServiceGslbSettings struct {
	GSLB GslbRecordSets `json:",omitempty"`
}

type CommonServiceGslbStatus struct {
	FQDN string `json:",omitempty`
}

type CommonServiceGslbProvider struct {
	Class string `json:",omitempty"`
}

// CreateNewGslbCommonServiceItem Create new CommonServiceItem
func CreateNewGslbCommonServiceItem(gslbName string) *CommonServiceGslbItem {
	return &CommonServiceGslbItem{
		Resource: &Resource{ID: ""},
		Name:     gslbName,
		Provider: CommonServiceGslbProvider{
			Class: "gslb",
		},
		Settings: CommonServiceGslbSettings{
			GSLB: GslbRecordSets{
				DelayLoop:   "10",
				HealthCheck: defaultGslbHealthCheck,
				Weighted:    "True",
			},
		},
	}

}

func (d *CommonServiceGslbItem) HasGslbServer() bool {
	return len(d.Settings.GSLB.Servers) > 0
}

type GslbRecordSets struct {
	DelayLoop   string          `json:",omitempty"`
	HealthCheck GslbHealthCheck `json:",omitempty"`
	Weighted    string          `json:",omitempty"`
	Servers     []GslbServer    `json:",omitempty"`
}

func (g *GslbRecordSets) AddServer(ip string) {
	var record GslbServer
	var isExist = false
	for i := range g.Servers {
		if g.Servers[i].IPAddress == ip {
			isExist = true
		}
	}

	if !isExist {
		record = GslbServer{
			IPAddress: ip,
			Enabled:   "True",
			Weight:    "1",
		}
		g.Servers = append(g.Servers, record)
	}
}

func (g *GslbRecordSets) DeleteServer(ip string) {
	res := []GslbServer{}
	for i := range g.Servers {
		if g.Servers[i].IPAddress != ip {
			res = append(res, g.Servers[i])
		}
	}

	g.Servers = res
}

type GslbServer struct {
	IPAddress string `json:",omitempty"`
	Enabled   string `json:",omitempty`
	Weight    string `json:omitempty`
}

type GslbHealthCheck struct {
	Protocol string `json:",omitempty"`
	Path     string `json:",omitempty"`
	Status   string `json:",omitempty"`
}

var defaultGslbHealthCheck = GslbHealthCheck{
	Protocol: "http",
	Path:     "/",
	Status:   "200",
}
