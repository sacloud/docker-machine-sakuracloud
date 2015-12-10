package resources

// Request type of SakuraCloud API Request
type Request struct {
	// *SakuraCloudResources
	Server    *Server                `json:",omitempty"`
	Disk      *Disk                  `json:",omitempty"`
	Note      *Note                  `json:",omitempty"`
	Interface *Interface             `json:",omitempty"`
	From      int                    `json:",omitempty"`
	Count     int                    `json:",omitempty"`
	Sort      []string               `json:",omitempty"`
	Filter    map[string]interface{} `json:",omitempty"`
	Exclude   []string               `json:",omitempty"`
	Include   []string               `json:",omitempty"`
}
