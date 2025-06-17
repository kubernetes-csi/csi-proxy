package disk

import (
	"encoding/hex"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	shared "github.com/kubernetes-csi/csi-proxy/pkg/shared/disk"
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
	m := make(map[uint32]shared.DiskLocation)
	err := cim.WithCOMThread(func() error {
		// "location":  "PCI Slot 3 : Adapter 0 : Port 0 : Target 1 : LUN 0"
		disks, err := cim.ListDisks(cim.DiskSelectorListForDiskNumberAndLocation)
		if err != nil {
			return fmt.Errorf("could not query disk locations")
		}

		for _, disk := range disks {
			num, err := cim.GetDiskNumber(disk)
			if err != nil {
				return fmt.Errorf("failed to query disk number: %v, %w", disk, err)
			}

			location, err := cim.GetDiskLocation(disk)
			if err != nil {
				return fmt.Errorf("failed to query disk location: %v, %w", disk, err)
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
					m[num] = d
				}
			}
		}

		return nil
	})
	return m, err
}

func (imp DiskAPI) Rescan() error {
	return cim.WithCOMThread(func() error {
		result, err := cim.RescanDisks()
		if err != nil {
			return fmt.Errorf("error updating host storage cache output. result: %d, err: %v", result, err)
		}
		return nil
	})
}

func (imp DiskAPI) IsDiskInitialized(diskNumber uint32) (bool, error) {
	var partitionStyle int32
	err := cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, cim.DiskSelectorListForPartitionStyle)
		if err != nil {
			return fmt.Errorf("error checking initialized status of disk %d: %v", diskNumber, err)
		}

		partitionStyle, err = cim.GetDiskPartitionStyle(disk)
		if err != nil {
			return fmt.Errorf("failed to query partition style of disk %d: %v", diskNumber, err)
		}

		return nil
	})
	return partitionStyle != cim.PartitionStyleUnknown, err
}

func (imp DiskAPI) InitializeDisk(diskNumber uint32) error {
	return cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, nil)
		if err != nil {
			return fmt.Errorf("failed to initializing disk %d. error: %w", diskNumber, err)
		}

		result, err := cim.InitializeDisk(disk, cim.PartitionStyleGPT)
		if result != 0 || err != nil {
			return fmt.Errorf("failed to initializing disk %d: result %d, error: %w", diskNumber, result, err)
		}

		return nil
	})
}

func (imp DiskAPI) BasicPartitionsExist(diskNumber uint32) (bool, error) {
	var exist bool
	err := cim.WithCOMThread(func() error {
		partitions, err := cim.ListPartitionsWithFilters(nil, cim.FilterForPartitionOnDisk(diskNumber), cim.FilterForPartitionsOfTypeNormal())
		if cim.IgnoreNotFound(err) != nil {
			return fmt.Errorf("error checking presence of partitions on disk %d:, %v", diskNumber, err)
		}

		exist = len(partitions) > 0
		return nil
	})
	return exist, err
}

func (imp DiskAPI) CreateBasicPartition(diskNumber uint32) error {
	return cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, nil)
		if err != nil {
			return err
		}

		result, err := cim.CreatePartition(
			disk,
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
		if (result != 0 && result != cim.ErrorCodeCreatePartitionAccessPathAlreadyInUse) || err != nil {
			return fmt.Errorf("error creating partition on disk %d. result: %d, err: %v", diskNumber, result, err)
		}

		result, _, err = cim.RefreshDisk(disk)
		if result != 0 || err != nil {
			return fmt.Errorf("error rescan disk (%d). result %d, error: %v", diskNumber, result, err)
		}

		partitions, err := cim.ListPartitionsWithFilters(nil, cim.FilterForPartitionOnDisk(diskNumber), cim.FilterForPartitionsOfTypeNormal())
		if err != nil {
			return fmt.Errorf("error query basic partition on disk %d:, %v", diskNumber, err)
		}

		if len(partitions) == 0 {
			return fmt.Errorf("failed to create basic partition on disk %d:, %v", diskNumber, err)
		}

		partition := partitions[0]
		result, status, err := cim.SetPartitionState(partition, true)
		if result != 0 || err != nil {
			return fmt.Errorf("error bring partition %v on disk %d online. result: %d, status %s, err: %v", partition, diskNumber, result, status, err)
		}

		return nil
	})
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
	var diskNumberResult uint32
	err := cim.WithCOMThread(func() error {
		disks, err := cim.ListDisks(cim.DiskSelectorListForPathAndSerialNumber)
		if err != nil {
			return err
		}

		for _, disk := range disks {
			path, err := cim.GetDiskPath(disk)
			if err != nil {
				return fmt.Errorf("failed to query disk path: %v, %w", disk, err)
			}

			diskNumber, diskPage83ID, err := imp.GetDiskNumberAndPage83ID(path)
			if err != nil {
				return err
			}

			if diskPage83ID == page83ID {
				diskNumberResult = diskNumber
				return nil
			}
		}

		return fmt.Errorf("could not find disk with Page83 ID %s", page83ID)
	})
	return diskNumberResult, err
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
	m := make(map[uint32]shared.DiskIDs)
	err := cim.WithCOMThread(func() error {
		disks, err := cim.ListDisks(cim.DiskSelectorListForPathAndSerialNumber)
		if err != nil {
			return err
		}

		for _, disk := range disks {
			path, err := cim.GetDiskPath(disk)
			if err != nil {
				return fmt.Errorf("failed to query disk path: %v, %w", disk, err)
			}

			sn, err := cim.GetDiskSerialNumber(disk)
			if err != nil {
				return fmt.Errorf("failed to query disk serial number: %v, %w", disk, err)
			}

			diskNumber, page83, err := imp.GetDiskNumberAndPage83ID(path)
			if err != nil {
				return err
			}

			m[diskNumber] = shared.DiskIDs{
				Page83:       page83,
				SerialNumber: sn,
			}
		}

		return nil
	})
	return m, err
}

func (imp DiskAPI) GetDiskStats(diskNumber uint32) (int64, error) {
	// TODO: change to uint64 as it does not make sense to use int64 for size
	size := int64(-1)
	err := cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, cim.DiskSelectorListForSize)
		if err != nil {
			return err
		}

		size, err = cim.GetDiskSize(disk)
		if err != nil {
			return fmt.Errorf("failed to query size of disk %d. %v", diskNumber, err)
		}

		return nil
	})
	return size, err
}

func (imp DiskAPI) SetDiskState(diskNumber uint32, isOnline bool) error {
	return cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, cim.DiskSelectorListForIsOffline)
		if err != nil {
			return err
		}

		isOffline, err := cim.IsDiskOffline(disk)
		if err != nil {
			return fmt.Errorf("error setting disk %d attach state. error: %v", diskNumber, err)
		}

		if isOnline == !isOffline {
			return nil
		}

		result, _, err := cim.SetDiskState(disk, isOnline)
		if result != 0 || err != nil {
			return fmt.Errorf("setting disk %d attach state (isOnline: %v): result %d, error: %w", diskNumber, isOnline, result, err)
		}

		return nil
	})
}

func (imp DiskAPI) GetDiskState(diskNumber uint32) (bool, error) {
	var isOffline bool
	err := cim.WithCOMThread(func() error {
		disk, err := cim.QueryDiskByNumber(diskNumber, cim.DiskSelectorListForIsOffline)
		if err != nil {
			return err
		}

		isOffline, err = cim.IsDiskOffline(disk)
		if err != nil {
			return fmt.Errorf("error parsing disk %d state. error: %v", diskNumber, err)
		}

		return nil
	})
	return !isOffline, err
}
