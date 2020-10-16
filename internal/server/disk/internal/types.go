package internal

type DiskLocation struct {
	Adapter string
	Bus     string
	Target  string
	LUNID   string
}

type ListDiskLocationsRequest struct {
}

type ListDiskLocationsResponse struct {
	// Map of disk device IDs and <adapter, bus, target, lun ID> associated with each disk device
	DiskLocations map[string]*DiskLocation
}

type PartitionDiskRequest struct {
	// Disk device ID of the disk to partition
	DiskID string
}

type PartitionDiskResponse struct {
}

type RescanRequest struct {
}

type RescanResponse struct {
}

type GetDiskNumberByNameRequest struct {
	// Disk device ID of the disk to partition
	DiskName string
}

type GetDiskNumberByNameResponse struct {
	DiskNumber string
}

type DiskIDs struct {
	// Map of Disk ID types and Disk ID values
	Identifiers map[string]string
}

type ListDiskIDsRequest struct {
}

type ListDiskIDsResponse struct {
	// Map of disk device numbers and IDs <page83> associated with each disk device
	DiskIDs map[string]*DiskIDs
}

type DiskStatsRequest struct {
	DiskID string
}

type DiskStatsResponse struct {
	DiskSize int64
}

type SetAttachStateRequest struct {
	// Disk device ID of the disk which state will change
	DiskID string

	// Online state to set for the disk. true for online, false for offline
	IsOnline bool
}

type SetAttachStateResponse struct {
}

type GetAttachStateRequest struct {
	// Disk device ID of the disk
	DiskID string
}

type GetAttachStateResponse struct {
	// Online state of the disk. true for online, false for offline
	IsOnline bool
}
