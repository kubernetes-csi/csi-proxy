# CSI Proxy Volume v1 API
<a name="top"></a>

## Table of Contents

- [Volume RPCs](#v1.VolumeRPCs)

- [Volume Messages](#v1.VolumeMessages)

<a name="v1.VolumeRPCs"></a>

## v1 Volume RPCs

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListVolumesOnDisk | [ListVolumesOnDiskRequest](#v1.ListVolumesOnDiskRequest) | [ListVolumesOnDiskResponse](#v1.ListVolumesOnDiskResponse) | ListVolumesOnDisk returns the volume IDs (in \\.\Volume{GUID} format) for all volumes from a given disk number and partition number (optional) |
| MountVolume | [MountVolumeRequest](#v1.MountVolumeRequest) | [MountVolumeResponse](#v1.MountVolumeResponse) | MountVolume mounts the volume at the requested global staging path. |
| UnmountVolume | [UnmountVolumeRequest](#v1.UnmountVolumeRequest) | [UnmountVolumeResponse](#v1.UnmountVolumeResponse) | UnmountVolume flushes data cache to disk and removes the global staging path. |
| IsVolumeFormatted | [IsVolumeFormattedRequest](#v1.IsVolumeFormattedRequest) | [IsVolumeFormattedResponse](#v1.IsVolumeFormattedResponse) | IsVolumeFormatted checks if a volume is formatted. |
| FormatVolume | [FormatVolumeRequest](#v1.FormatVolumeRequest) | [FormatVolumeResponse](#v1.FormatVolumeResponse) | FormatVolume formats a volume with NTFS. |
| ResizeVolume | [ResizeVolumeRequest](#v1.ResizeVolumeRequest) | [ResizeVolumeResponse](#v1.ResizeVolumeResponse) | ResizeVolume performs resizing of the partition and file system for a block based volume. |
| GetVolumeStats | [GetVolumeStatsRequest](#v1.GetVolumeStatsRequest) | [GetVolumeStatsResponse](#v1.GetVolumeStatsResponse) | GetVolumeStats gathers total bytes and used bytes for a volume. |
| GetDiskNumberFromVolumeID | [GetDiskNumberFromVolumeIDRequest](#v1.GetDiskNumberFromVolumeIDRequest) | [GetDiskNumberFromVolumeIDResponse](#v1.GetDiskNumberFromVolumeIDResponse) | GetDiskNumberFromVolumeID gets the disk number of the disk where the volume is located. |
| GetVolumeIDFromTargetPath | [GetVolumeIDFromTargetPathRequest](#v1.GetVolumeIDFromTargetPathRequest) | [GetVolumeIDFromTargetPathResponse](#v1.GetVolumeIDFromTargetPathResponse) | GetVolumeIDFromTargetPath gets the volume id for a given target path. |
| WriteVolumeCache | [WriteVolumeCacheRequest](#v1.WriteVolumeCacheRequest) | [WriteVolumeCacheResponse](#v1.WriteVolumeCacheResponse) | WriteVolumeCache write volume cache to disk. |


<a name="v1.VolumeMessages"></a>
<p align="right"><a href="#top">Top</a></p>

## v1 Volume Messages

<a name="v1.FormatVolumeRequest"></a>
### FormatVolumeRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to format. |

<a name="v1.FormatVolumeResponse"></a>
### FormatVolumeResponse
Intentionally empty.

<a name="v1.GetDiskNumberFromVolumeIDRequest"></a>
### GetDiskNumberFromVolumeIDRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to get the disk number for. |

<a name="v1.GetDiskNumberFromVolumeIDResponse"></a>
### GetDiskNumberFromVolumeIDResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Corresponding disk number. |

<a name="v1.GetVolumeIDFromTargetPathRequest"></a>
### GetVolumeIDFromTargetPathRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| target_path | string |  | The target path. |

<a name="v1.GetVolumeIDFromTargetPathResponse"></a>
### GetVolumeIDFromTargetPathResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | The volume device ID. |

<a name="v1.GetVolumeStatsRequest"></a>
### GetVolumeStatsRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device Id of the volume to get the stats for. |

<a name="v1.GetVolumeStatsResponse"></a>
### GetVolumeStatsResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total_bytes | int64 |  | Total bytes |
| used_bytes | int64 |  | Used bytes |

<a name="v1.IsVolumeFormattedRequest"></a>
### IsVolumeFormattedRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to check. |

<a name="v1.IsVolumeFormattedResponse"></a>
### IsVolumeFormattedResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| formatted | bool |  | Is the volume formatted with NTFS. |

<a name="v1.ListVolumesOnDiskRequest"></a>
### ListVolumesOnDiskRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Disk device number of the disk to query for volumes. |
| partition_number | uint32 |  | The partition number (optional), by default it uses the first partition of the disk. |

<a name="v1.ListVolumesOnDiskResponse"></a>
### ListVolumesOnDiskResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_ids | string | repeated | Volume device IDs of volumes on the specified disk. |


<a name="v1.MountVolumeRequest"></a>
### MountVolumeRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to mount. |
| target_path | string |  | Path in the host&#39;s file system where the volume needs to be mounted. |

<a name="v1.MountVolumeResponse"></a>
### MountVolumeResponse
Intentionally empty.

<a name="v1.ResizeVolumeRequest"></a>
### ResizeVolumeRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to resize. |
| size_bytes | int64 |  | New size in bytes of the volume. |

<a name="v1.ResizeVolumeResponse"></a>
### ResizeVolumeResponse
Intentionally empty.

<a name="v1.UnmountVolumeRequest"></a>
### UnmountVolumeRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to dismount. |
| target_path | string |  | Path where the volume has been mounted. |

<a name="v1.UnmountVolumeResponse"></a>
### UnmountVolumeResponse
Intentionally empty.

<a name="v1.WriteVolumeCacheRequest"></a>
### WriteVolumeCacheRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| volume_id | string |  | Volume device ID of the volume to flush the cache. |


<a name="v1.WriteVolumeCacheResponse"></a>
### WriteVolumeCacheResponse
Intentionally empty.
