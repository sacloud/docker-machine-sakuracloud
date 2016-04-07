package resources

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

// DiskEditValue type of disk edit request value
type DiskEditValue struct {
	Password      string
	SSHKey        SSHKey
	DisablePWAuth bool
	Notes         []Resource
}
