# CSI Proxy Disk v1 API
<a name="top"></a>

## Table of Contents

- [Disk RPCs](#v1.DiskRPCs)

- [Disk Messages](#v1.DiskMessages)


<a name="v1.DiskRPCs"></a>
## v1 Disk RPCs

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListDiskLocations | [ListDiskLocationsRequest](#v1.ListDiskLocationsRequest) | [ListDiskLocationsResponse](#v1.ListDiskLocationsResponse) | ListDiskLocations returns locations &lt;Adapter, Bus, Target, LUN ID&gt; of all disk devices enumerated by the host. |
| PartitionDisk | [PartitionDiskRequest](#v1.PartitionDiskRequest) | [PartitionDiskResponse](#v1.PartitionDiskResponse) | PartitionDisk initializes and partitions a disk device with the GPT partition style (if the disk has not been partitioned already) and returns the resulting volume device ID. |
| Rescan | [RescanRequest](#v1.RescanRequest) | [RescanResponse](#v1.RescanResponse) | Rescan refreshes the host&#39;s storage cache. |
| ListDiskIDs | [ListDiskIDsRequest](#v1.ListDiskIDsRequest) | [ListDiskIDsResponse](#v1.ListDiskIDsResponse) | ListDiskIDs returns a map of DiskID objects where the key is the disk number. |
| GetDiskStats | [GetDiskStatsRequest](#v1.GetDiskStatsRequest) | [GetDiskStatsResponse](#v1.GetDiskStatsResponse) | GetDiskStats returns the stats of a disk (currently it returns the disk size). |
| SetDiskState | [SetDiskStateRequest](#v1.SetDiskStateRequest) | [SetDiskStateResponse](#v1.SetDiskStateResponse) | SetDiskState sets the offline/online state of a disk. |
| GetDiskState | [GetDiskStateRequest](#v1.GetDiskStateRequest) | [GetDiskStateResponse](#v1.GetDiskStateResponse) | GetDiskState gets the offline/online state of a disk. |


<a name="v1.DiskMessages"></a>
<p align="right"><a href="#top">Top</a></p>

## v1 Disk Messages

<a name="v1.DiskIDs"></a>
### DiskIDs

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page83 | string |  | The disk page83 id. |
| serial_number | string |  | The disk serial number. |

<a name="v1.DiskLocation"></a>
### DiskLocation

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Adapter | string |  |  |
| Bus | string |  |  |
| Target | string |  |  |
| LUNID | string |  |  |

<a name="v1.GetDiskStateRequest"></a>
### GetDiskStateRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Disk device number of the disk. |

<a name="v1.GetDiskStateResponse"></a>
### GetDiskStateResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_online | bool |  | Online state of the disk. true for online, false for offline. |

<a name="v1.GetDiskStatsRequest"></a>
### GetDiskStatsRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Disk device number of the disk to get the stats from. |

<a name="v1.GetDiskStatsResponse"></a>
### GetDiskStatsResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total_bytes | int64 |  | Total size of the volume. |

<a name="v1.ListDiskIDsRequest"></a>

### ListDiskIDsRequest
Intentionally empty.

<a name="v1.ListDiskIDsResponse"></a>
### ListDiskIDsResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| diskIDs | [ListDiskIDsResponse.DiskIDsEntry](#v1.ListDiskIDsResponse.DiskIDsEntry) | repeated | Map of disk numbers and disk identifiers associated with each disk device.

the case is intentional for protoc to generate the field as DiskIDs |

<a name="v1.ListDiskIDsResponse.DiskIDsEntry"></a>
### ListDiskIDsResponse.DiskIDsEntry

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | uint32 |  |  |
| value | [DiskIDs](#v1.DiskIDs) |  |  |

<a name="v1.ListDiskLocationsRequest"></a>

### ListDiskLocationsRequest
Intentionally empty.

<a name="v1.ListDiskLocationsResponse"></a>
### ListDiskLocationsResponse

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_locations | [ListDiskLocationsResponse.DiskLocationsEntry](#v1.ListDiskLocationsResponse.DiskLocationsEntry) | repeated | Map of disk number and &lt;adapter, bus, target, lun ID&gt; associated with each disk device. |

<a name="v1.ListDiskLocationsResponse.DiskLocationsEntry"></a>
### ListDiskLocationsResponse.DiskLocationsEntry

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | uint32 |  |  |
| value | [DiskLocation](#v1.DiskLocation) |  |  |

<a name="v1.PartitionDiskRequest"></a>
### PartitionDiskRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Disk device number of the disk to partition. |

<a name="v1.PartitionDiskResponse"></a>
### PartitionDiskResponse
Intentionally empty.

<a name="v1.RescanRequest"></a>
### RescanRequest
Intentionally empty.

<a name="v1.RescanResponse"></a>

### RescanResponse
Intentionally empty.

<a name="v1.SetDiskStateRequest"></a>
### SetDiskStateRequest

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| disk_number | uint32 |  | Disk device number of the disk. |
| is_online | bool |  | Online state to set for the disk. true for online, false for offline. |

<a name="v1.SetDiskStateResponse"></a>
### SetDiskStateResponse
Intentionally empty.
