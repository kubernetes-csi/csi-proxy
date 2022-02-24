package disk

import "golang.org/x/sys/windows"

type StorageDeviceNumber struct {
	DeviceType      DeviceType
	DeviceNumber    uint32
	PartitionNumber uint32
}
type DeviceType uint32

type StoragePropertyID uint32

const (
	StorageDeviceProperty                  StoragePropertyID = 0
	StorageAdapterProperty                                   = 1
	StorageDeviceIDProperty                                  = 2
	StorageDeviceUniqueIDProperty                            = 3
	StorageDeviceWriteCacheProperty                          = 4
	StorageMiniportProperty                                  = 5
	StorageAccessAlignmentProperty                           = 6
	StorageDeviceSeekPenaltyProperty                         = 7
	StorageDeviceTrimProperty                                = 8
	StorageDeviceWriteAggregationProperty                    = 9
	StorageDeviceDeviceTelemetryProperty                     = 10
	StorageDeviceLBProvisioningProperty                      = 11
	StorageDevicePowerProperty                               = 12
	StorageDeviceCopyOffloadProperty                         = 13
	StorageDeviceResiliencyProperty                          = 14
	StorageDeviceMediumProductType                           = 15
	StorageAdapterRpmbProperty                               = 16
	StorageAdapterCryptoProperty                             = 17
	StorageDeviceIoCapabilityProperty                        = 18
	StorageAdapterProtocolSpecificProperty                   = 19
	StorageDeviceProtocolSpecificProperty                    = 20
	StorageAdapterTemperatureProperty                        = 21
	StorageDeviceTemperatureProperty                         = 22
	StorageAdapterPhysicalTopologyProperty                   = 23
	StorageDevicePhysicalTopologyProperty                    = 24
	StorageDeviceAttributesProperty                          = 25
	StorageDeviceManagementStatus                            = 26
	StorageAdapterSerialNumberProperty                       = 27
	StorageDeviceLocationProperty                            = 28
	StorageDeviceNumaProperty                                = 29
	StorageDeviceZonedDeviceProperty                         = 30
	StorageDeviceUnsafeShutdownCount                         = 31
	StorageDeviceEnduranceProperty                           = 32
)

type StorageQueryType uint32

const (
	PropertyStandardQuery StorageQueryType = iota
	PropertyExistsQuery
	PropertyMaskQuery
	PropertyQueryMaxDefined
)

type StoragePropertyQuery struct {
	PropertyID StoragePropertyID
	QueryType  StorageQueryType
	Byte       []AdditionalParameters
}

type AdditionalParameters byte

type StorageDeviceIDDescriptor struct {
	Version             uint32
	Size                uint32
	NumberOfIdentifiers uint32
	Identifiers         [1]byte
}

type StorageIdentifierCodeSet uint32

const (
	StorageIDCodeSetReserved StorageIdentifierCodeSet = 0
	StorageIDCodeSetBinary                            = 1
	StorageIDCodeSetASCII                             = 2
	StorageIDCodeSetUtf8                              = 3
)

type StorageIdentifierType uint32

const (
	StorageIdTypeVendorSpecific           StorageIdentifierType = 0
	StorageIDTypeVendorID                                       = 1
	StorageIDTypeEUI64                                          = 2
	StorageIDTypeFCPHName                                       = 3
	StorageIDTypePortRelative                                   = 4
	StorageIDTypeTargetPortGroup                                = 5
	StorageIDTypeLogicalUnitGroup                               = 6
	StorageIDTypeMD5LogicalUnitIdentifier                       = 7
	StorageIDTypeScsiNameString                                 = 8
)

type StorageAssociationType uint32

const (
	StorageIDAssocDevice StorageAssociationType = 0
	StorageIDAssocPort                          = 1
	StorageIDAssocTarget                        = 2
)

type StorageIdentifier struct {
	CodeSet        StorageIdentifierCodeSet
	Type           StorageIdentifierType
	IdentifierSize uint16
	NextOffset     uint16
	Association    StorageAssociationType
	Identifier     [1]byte
}

type Disk struct {
	Path         string `json:"Path"`
	SerialNumber string `json:"SerialNumber"`
}

type SetStorageDeviceAttributes struct {
	Version        uint32
	Persist        bool
	Reserved1      uint32
	Attributes     uint64
	AttributesMask uint64
	Reserved2      uint64
}

type GetStorageDeviceAtrributes struct {
	Version    uint32
	Reserved1  uint32
	Attributes uint64
}

type StoragePartitionStyle uint32

const (
	PartitionStyleMbr StoragePartitionStyle = iota
	PartitionStyleGpt
	PartitionStyleRaw
)

type PartitionInfoMbr struct {
	PartitionType       byte
	BootIndicator       bool
	RecognizedPartition bool
	HiddenSectors       uint32
	PartitionId         windows.GUID
}

type PartitionInfoGpt struct {
	PartitionType windows.GUID
	PartitionId   windows.GUID
	Attributes    uint64
	Name          [36]uint16
}

type StoragePartitionInfo struct {
	PartitionStyle     StoragePartitionStyle
	StartingOffset     int64
	PartitionLength    int64
	PartitionNumber    uint32
	RewritePartition   bool
	IsServicePartition bool
	DummmyUnionName    struct {
		PartitionInfoMbr
		PartitionInfoGpt
	}
}

type DriveLayoutInfoMbr struct {
	Signature uint32
	CheckSum  uint32
}

type DriveLayoutInfoGpt struct {
	DiskId               windows.GUID
	StartingUsableOffset int64
	UsableLength         int64
	MaxPartitionCount    uint32
}

type StorageDriveLayoutInfo struct {
	PartitionStyle  StoragePartitionStyle
	PartitionCount  uint32
	DummmyUnionName struct {
		DriveLayoutInfoMbr
		DriveLayoutInfoGpt
	}
	PartitionEntry []StoragePartitionInfo
}
