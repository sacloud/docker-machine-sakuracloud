package resources

// Interface type of server nic
type Interface struct {
	*Resource
	IPAddress     string `json:",omitempty"`
	UserIPAddress string `json:",omitempty"`
	MACAddress    string `json:",omitempty"`
}
