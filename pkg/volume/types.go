package volume

type ListVolumesOnDiskRequest struct {
	DiskNumber      uint32
	PartitionNumber uint32
}

type ListVolumesOnDiskResponse struct {
	VolumeIds []string
}

type MountVolumeRequest struct {
	VolumeId   string
	TargetPath string
}

type MountVolumeResponse struct {
}

type IsVolumeFormattedRequest struct {
	VolumeId string
}

type IsVolumeFormattedResponse struct {
	Formatted bool
}

type FormatVolumeRequest struct {
	VolumeId string
}

type FormatVolumeResponse struct {
}

type WriteVolumeCacheRequest struct {
	VolumeId string
}

type WriteVolumeCacheResponse struct {
}

type UnmountVolumeRequest struct {
	VolumeId   string
	TargetPath string
}

type UnmountVolumeResponse struct {
}

type ResizeVolumeRequest struct {
	VolumeId  string
	SizeBytes int64
}

type ResizeVolumeResponse struct {
}

type GetVolumeStatsRequest struct {
	VolumeId string
}

type GetVolumeStatsResponse struct {
	TotalBytes int64
	UsedBytes  int64
}

type GetDiskNumberFromVolumeIDRequest struct {
	VolumeId string
}

type GetDiskNumberFromVolumeIDResponse struct {
	DiskNumber uint32
}

type GetVolumeIDFromTargetPathRequest struct {
	TargetPath string
}

type GetVolumeIDFromTargetPathResponse struct {
	VolumeId string
}

type GetClosestVolumeIDFromTargetPathRequest struct {
	TargetPath string
}

type GetClosestVolumeIDFromTargetPathResponse struct {
	VolumeId string
}
