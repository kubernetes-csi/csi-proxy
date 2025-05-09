package disk

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	shared "github.com/kubernetes-csi/csi-proxy/pkg/shared/disk"
	"github.com/microsoft/wmi/pkg/base/query"
	"k8s.io/klog/v2"
)

var (
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
)

const (
	IOCTL_STORAGE_GET_DEVICE_NUMBER = 0x2D1080
	IOCTL_STORAGE_QUERY_PROPERTY    = 0x002d1400
)

// API declares the interface exposed by the internal API
type API interface {
	// ListDiskLocations - constructs a map with the disk number as the key and the DiskLocation structure
	// as the value. The DiskLocation struct has various fields like the Adapter, Bus, Target and LUNID.
	ListDiskLocations() (map[uint32]shared.DiskLocation, error)
	// IsDiskInitialized returns true if the disk identified by `diskNumber` is initialized.
	IsDiskInitialized(diskNumber uint32) (bool, error)
	// InitializeDisk initializes the disk `diskNumber`
	InitializeDisk(diskNumber uint32) error
	// BasicPartitionsExist checks if the disk `diskNumber` has any basic partitions.
	BasicPartitionsExist(diskNumber uint32) (bool, error)
	// CreateBasicPartition creates a partition in disk `diskNumber`
	CreateBasicPartition(diskNumber uint32) error
	// Rescan updates the host storage cache (re-enumerates disk, partition and volume objects)
	Rescan() error
	// GetDiskNumberByName gets a disk number by page83 ID (disk name)
	GetDiskNumberByName(page83ID string) (uint32, error)
	// ListDiskIDs list all disks by disk number.
	ListDiskIDs() (map[uint32]shared.DiskIDs, error)
	// GetDiskStats gets the disk stats of the disk `diskNumber`.
	GetDiskStats(diskNumber uint32) (int64, error)
	// SetDiskState sets the offline/online state of the disk `diskNumber`.
	SetDiskState(diskNumber uint32, isOnline bool) error
	// GetDiskState gets the offline/online state of the disk `diskNumber`.
	GetDiskState(diskNumber uint32) (bool, error)
}

// DiskAPI implements the OS API calls related to Disk Devices. All code here should be very simple
// pass-through to the OS APIs or cmdlets. Any logic around the APIs/cmdlet invocation
// should go in internal/server/filesystem/disk.go so that logic can be easily unit-tested
// without requiring specific OS environments.
type DiskAPI struct{}

// ensure that DiskAPI implements the exposed API
var _ API = &DiskAPI{}

func New() DiskAPI {
	return DiskAPI{}
}

// ListDiskLocations - constructs a map with the disk number as the key and the DiskLocation structure
// as the value. The DiskLocation struct has various fields like the Adapter, Bus, Target and LUNID.
func (imp DiskAPI) ListDiskLocations() (map[uint32]shared.DiskLocation, error) {
	// "location":  "PCI Slot 3 : Adapter 0 : Port 0 : Target 1 : LUN 0"
	disks, err := cim.ListDisks([]string{"Number", "Location"})
	if err != nil {
		return nil, fmt.Errorf("could not query disk locations")
	}

	m := make(map[uint32]shared.DiskLocation)
	for _, disk := range disks {
		num, err := disk.GetProperty("Number")
		if err != nil {
			return m, fmt.Errorf("failed to query disk number: %v, %w", disk, err)
		}

		location, err := disk.GetPropertyLocation()
		if err != nil {
			return m, fmt.Errorf("failed to query disk location: %v, %w", disk, err)
		}

		found := false
		s := strings.Split(location, ":")
		if len(s) >= 5 {
			var d shared.DiskLocation
			for _, item := range s {
				item = strings.TrimSpace(item)
				itemSplit := strings.Split(item, " ")
				if len(itemSplit) == 2 {
					found = true
					switch strings.TrimSpace(itemSplit[0]) {
					case "Adapter":
						d.Adapter = strings.TrimSpace(itemSplit[1])
					case "Target":
						d.Target = strings.TrimSpace(itemSplit[1])
					case "LUN":
						d.LUNID = strings.TrimSpace(itemSplit[1])
					default:
						klog.Warningf("Got unknown field : %s=%s", itemSplit[0], itemSplit[1])
					}
				}
			}

			if found {
				m[uint32(num.(int32))] = d
			}
		}
	}

	return m, nil
}

func (imp DiskAPI) Rescan() error {
	result, _, err := cim.InvokeCimMethod(cim.WMINamespaceStorage, "MSFT_StorageSetting", "UpdateHostStorageCache", nil)
	if err != nil {
		return fmt.Errorf("error updating host storage cache output. result: %d, err: %v", result, err)
	}
	return nil
}

func (imp DiskAPI) IsDiskInitialized(diskNumber uint32) (bool, error) {
	var partitionStyle int32
	disk, err := cim.QueryDiskByNumber(diskNumber, []string{"PartitionStyle"})
	if err != nil {
		return false, fmt.Errorf("error checking initialized status of disk %d. %v", diskNumber, err)
	}

	retValue, err := disk.GetProperty("PartitionStyle")
	if err != nil {
		return false, fmt.Errorf("failed to query partition style of disk %d: %w", diskNumber, err)
	}

	partitionStyle = retValue.(int32)
	return partitionStyle != cim.PartitionStyleUnknown, nil
}

func (imp DiskAPI) InitializeDisk(diskNumber uint32) error {
	disk, err := cim.QueryDiskByNumber(diskNumber, nil)
	if err != nil {
		return fmt.Errorf("failed to initializing disk %d. error: %w", diskNumber, err)
	}

	result, err := disk.InvokeMethodWithReturn("Initialize", int32(cim.PartitionStyleGPT))
	if result != 0 || err != nil {
		return fmt.Errorf("failed to initializing disk %d: result %d, error: %w", diskNumber, result, err)
	}

	return nil
}

func (imp DiskAPI) BasicPartitionsExist(diskNumber uint32) (bool, error) {
	partitions, err := cim.ListPartitionsWithFilters(nil,
		query.NewWmiQueryFilter("DiskNumber", strconv.Itoa(int(diskNumber)), query.Equals),
		query.NewWmiQueryFilter("GptType", cim.GPTPartitionTypeMicrosoftReserved, query.NotEquals))
	if cim.IgnoreNotFound(err) != nil {
		return false, fmt.Errorf("error checking presence of partitions on disk %d:, %v", diskNumber, err)
	}

	return len(partitions) > 0, nil
}

func (imp DiskAPI) CreateBasicPartition(diskNumber uint32) error {
	disk, err := cim.QueryDiskByNumber(diskNumber, nil)
	if err != nil {
		return err
	}

	result, err := disk.InvokeMethodWithReturn(
		"CreatePartition",
		nil,                           // Size
		true,                          // UseMaximumSize
		nil,                           // Offset
		nil,                           // Alignment
		nil,                           // DriveLetter
		false,                         // AssignDriveLetter
		nil,                           // MbrType,
		cim.GPTPartitionTypeBasicData, // GPT Type
		false,                         // IsHidden
		false,                         // IsActive,
	)
	// 42002 is returned by driver letter failed to assign after partition
	if (result != 0 && result != 42002) || err != nil {
		return fmt.Errorf("error creating partition on disk %d. result: %d, err: %v", diskNumber, result, err)
	}

	var status string
	result, err = disk.InvokeMethodWithReturn("Refresh", &status)
	if result != 0 || err != nil {
		return fmt.Errorf("error rescan disk (%d). result %d, error: %v", diskNumber, result, err)
	}

	partitions, err := cim.ListPartitionsWithFilters(nil,
		query.NewWmiQueryFilter("DiskNumber", strconv.Itoa(int(diskNumber)), query.Equals),
		query.NewWmiQueryFilter("GptType", cim.GPTPartitionTypeMicrosoftReserved, query.NotEquals))
	if err != nil {
		return fmt.Errorf("error query basic partition on disk %d:, %v", diskNumber, err)
	}

	if len(partitions) == 0 {
		return fmt.Errorf("failed to create basic partition on disk %d:, %v", diskNumber, err)
	}

	partition := partitions[0]
	result, err = partition.InvokeMethodWithReturn("Online", status)
	if result != 0 || err != nil {
		return fmt.Errorf("error bring partition %v on disk %d online. result: %d, status %s, err: %v", partition, diskNumber, result, status, err)
	}

	err = partition.Refresh()
	return err
}

func (imp DiskAPI) GetDiskNumberByName(page83ID string) (uint32, error) {
	diskNumber, err := imp.GetDiskNumberWithID(page83ID)
	return diskNumber, err
}

func (imp DiskAPI) GetDiskNumber(disk syscall.Handle) (uint32, error) {
	var bytes uint32
	devNum := StorageDeviceNumber{}
	buflen := uint32(unsafe.Sizeof(devNum.DeviceType)) + uint32(unsafe.Sizeof(devNum.DeviceNumber)) + uint32(unsafe.Sizeof(devNum.PartitionNumber))

	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_GET_DEVICE_NUMBER, nil, 0, (*byte)(unsafe.Pointer(&devNum)), buflen, &bytes, nil)

	return devNum.DeviceNumber, err
}

func (imp DiskAPI) GetDiskPage83ID(disk syscall.Handle) (string, error) {
	query := StoragePropertyQuery{}

	bufferSize := uint32(4 * 1024)
	buffer := make([]byte, 4*1024)
	var size uint32
	var n uint32
	var m uint16

	query.QueryType = PropertyStandardQuery
	query.PropertyID = StorageDeviceIDProperty

	querySize := uint32(unsafe.Sizeof(query))
	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_QUERY_PROPERTY, (*byte)(unsafe.Pointer(&query)), querySize, (*byte)(unsafe.Pointer(&buffer[0])), bufferSize, &size, nil)
	if err != nil {
		return "", fmt.Errorf("IOCTL_STORAGE_QUERY_PROPERTY failed: %v", err)
	}

	devIDDesc := (*StorageDeviceIDDescriptor)(unsafe.Pointer(&buffer[0]))

	pID := (*StorageIdentifier)(unsafe.Pointer(&devIDDesc.Identifiers[0]))

	page83ID := []byte{}
	byteSize := unsafe.Sizeof(byte(0))
	for n = 0; n < devIDDesc.NumberOfIdentifiers; n++ {
		if pID.Association == StorageIDAssocDevice && (pID.CodeSet == StorageIDCodeSetBinary || pID.CodeSet == StorageIDCodeSetASCII) {
			for m = 0; m < pID.IdentifierSize; m++ {
				page83ID = append(page83ID, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&pID.Identifier[0])) + byteSize*uintptr(m))))
			}

			if pID.CodeSet == StorageIDCodeSetASCII {
				return string(page83ID), nil
			} else if pID.CodeSet == StorageIDCodeSetBinary {
				return hex.EncodeToString(page83ID), nil
			}
		}
		pID = (*StorageIdentifier)(unsafe.Pointer(uintptr(unsafe.Pointer(pID)) + byteSize*uintptr(pID.NextOffset)))
	}
	return "", nil
}

func (imp DiskAPI) GetDiskNumberWithID(page83ID string) (uint32, error) {
	disks, err := cim.ListDisks([]string{"Path", "SerialNumber"})
	if err != nil {
		return 0, err
	}

	for _, disk := range disks {
		path, err := disk.GetPropertyPath()
		if err != nil {
			return 0, fmt.Errorf("failed to query disk path: %v, %w", disk, err)
		}

		diskNumber, diskPage83ID, err := imp.GetDiskNumberAndPage83ID(path)
		if err != nil {
			return 0, err
		}

		if diskPage83ID == page83ID {
			return diskNumber, nil
		}
	}

	return 0, fmt.Errorf("could not find disk with Page83 ID %s", page83ID)
}

func (imp DiskAPI) GetDiskNumberAndPage83ID(path string) (uint32, string, error) {
	h, err := syscall.Open(path, syscall.O_RDONLY, 0)
	defer syscall.Close(h)
	if err != nil {
		return 0, "", err
	}

	diskNumber, err := imp.GetDiskNumber(h)
	if err != nil {
		return 0, "", err
	}

	page83ID, err := imp.GetDiskPage83ID(h)
	if err != nil {
		return 0, "", err
	}

	return diskNumber, page83ID, nil
}

// ListDiskIDs - constructs a map with the disk number as the key and the DiskID structure
// as the value. The DiskID struct has a field for the page83 ID.
func (imp DiskAPI) ListDiskIDs() (map[uint32]shared.DiskIDs, error) {
	disks, err := cim.ListDisks([]string{"Path", "SerialNumber"})
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]shared.DiskIDs)
	for _, disk := range disks {
		path, err := disk.GetPropertyPath()
		if err != nil {
			return m, fmt.Errorf("failed to query disk path: %v, %w", disk, err)
		}

		sn, err := disk.GetPropertySerialNumber()
		if err != nil {
			return m, fmt.Errorf("failed to query disk serial number: %v, %w", disk, err)
		}

		diskNumber, page83, err := imp.GetDiskNumberAndPage83ID(path)
		if err != nil {
			return m, err
		}

		m[diskNumber] = shared.DiskIDs{
			Page83:       page83,
			SerialNumber: sn,
		}
	}
	return m, nil
}

func (imp DiskAPI) GetDiskStats(diskNumber uint32) (int64, error) {
	// TODO: change to uint64 as it does not make sense to use int64 for size
	var size int64
	disk, err := cim.QueryDiskByNumber(diskNumber, []string{"Size"})
	if err != nil {
		return -1, err
	}

	sz, err := disk.GetProperty("Size")
	if err != nil {
		return -1, fmt.Errorf("failed to query size of disk %d. %v", diskNumber, err)
	}

	size, err = strconv.ParseInt(sz.(string), 10, 64)
	return size, err
}

func (imp DiskAPI) SetDiskState(diskNumber uint32, isOnline bool) error {
	disk, err := cim.QueryDiskByNumber(diskNumber, []string{"IsOffline"})
	if err != nil {
		return err
	}

	offline, err := disk.GetPropertyIsOffline()
	if err != nil {
		return fmt.Errorf("error setting disk %d attach state. error: %v", diskNumber, err)
	}

	if isOnline == !offline {
		return nil
	}

	method := "Offline"
	if isOnline {
		method = "Online"
	}

	result, err := disk.InvokeMethodWithReturn(method)
	if result != 0 || err != nil {
		return fmt.Errorf("setting disk %d attach state %s: result %d, error: %w", diskNumber, method, result, err)
	}

	return nil
}

func (imp DiskAPI) GetDiskState(diskNumber uint32) (bool, error) {
	disk, err := cim.QueryDiskByNumber(diskNumber, []string{"IsOffline"})
	if err != nil {
		return false, err
	}

	isOffline, err := disk.GetPropertyIsOffline()
	if err != nil {
		return false, fmt.Errorf("error parsing disk %d state. error: %v", diskNumber, err)
	}

	return !isOffline, nil
}
