package resources

// Note type of startup script
type Note struct {
	*Resource
	Name    string
	Content string
	*EAvailability
}
