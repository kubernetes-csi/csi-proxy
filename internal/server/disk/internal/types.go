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
	DiskLocations map[uint32]*DiskLocation
}

type PartitionDiskRequest struct {
	// Disk device ID of the disk to partition
	DiskNumber uint32
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
	DiskNumber uint32
}

type ListDiskIDsRequest struct {
}

type DiskIDs struct {
	// Map of Disk ID types and Disk ID values
	Page83       string
	SerialNumber string
}

type ListDiskIDsResponse struct {
	// Map of disk device numbers and IDs associated with each disk device
	DiskIDs map[uint32]*DiskIDs
}

type GetDiskStatsRequest struct {
	DiskNumber uint32
}

type GetDiskStatsResponse struct {
	TotalBytes int64
}

type SetDiskStateRequest struct {
	// Disk device ID of the disk which state will change
	DiskNumber uint32

	// Online state to set for the disk. true for online, false for offline
	IsOnline bool
}

type SetDiskStateResponse struct {
}

type GetDiskStateRequest struct {
	// Disk device ID of the disk
	DiskNumber uint32
}

type GetDiskStateResponse struct {
	// Online state of the disk. true for online, false for offline
	IsOnline bool
}

// These structs are used in pre v1beta3 API versions

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
