package resources

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
