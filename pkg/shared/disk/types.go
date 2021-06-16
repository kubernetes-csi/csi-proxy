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

// DiskIDs definition
type DiskIDs struct {
	Page83       string
	SerialNumber string
}
