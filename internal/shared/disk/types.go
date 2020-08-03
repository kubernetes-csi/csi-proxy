package shared

// Shared location to include between the os and the server code to comply with
// go's inclusion rules around 'internal'.

// DiskLocation definition
type DiskLocation struct {
	Adapter string
	Bus     string
	Target  string
	LUNID   string
}

// DiskID definition
type DiskIDs struct {
	// Map of Disk ID types and Disk ID values
	Identifiers map[string]string
}
