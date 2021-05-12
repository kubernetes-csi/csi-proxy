// Defines all the structs that the server is aware of, because all
// the apis are included as the target for code generation it also
// has definition for older APIs e.g. volume/v1alpha1, volume/v1beta1, etc
// Because of this some structs are needed but are no longer used

package internal

type ListVolumesOnDiskRequest struct {
	PartitionNumber int64
	DiskNumber      int64
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
	DiskNumber int64
}

type GetVolumeIDFromTargetPathRequest struct {
	TargetPath string
}

type GetVolumeIDFromTargetPathResponse struct {
	VolumeId string
}

// These fields are deprecated but are needed because the generator needs them
type DismountVolumeRequest struct{}
type DismountVolumeResponse struct{}
type VolumeDiskNumberRequest struct{}
type VolumeDiskNumberResponse struct{}
type VolumeStatsRequest struct{}
type VolumeStatsResponse struct{}
type VolumeIDFromMountRequest struct{}
type VolumeIDFromMountResponse struct{}
