package disk

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

type DiskPath struct {
	Path string `json:"Path"`
}
