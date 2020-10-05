package internal

type ListVolumesOnDiskRequest struct {
	DiskId string
}

type ListVolumesOnDiskResponse struct {
	VolumeIds []string
}

type MountVolumeRequest struct {
	VolumeId string
	Path     string
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

type DismountVolumeRequest struct {
	VolumeId string
	Path     string
}

type DismountVolumeResponse struct {
}

type ResizeVolumeRequest struct {
	VolumeId string
	Size     int64
}

type ResizeVolumeResponse struct {
}

type VolumeStatsRequest struct {
	VolumeId string
	Path     string
}

type VolumeStatsResponse struct {
	VolumeSize     int64
	VolumeUsedSize int64
}

type VolumeDiskNumberRequest struct {
	VolumeId string
}

type VolumeDiskNumberResponse struct {
	DiskNumber int64
}

type VolumeIDFromMountRequest struct {
	Mount string
}

type VolumeIDFromMountResponse struct {
	VolumeId string
}
