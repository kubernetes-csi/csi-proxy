package disk

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	shared "github.com/kubernetes-csi/csi-proxy/pkg/shared/disk"
	"github.com/kubernetes-csi/csi-proxy/pkg/wmi"
	"k8s.io/klog/v2"
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
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			// "location":  "PCI Slot 3 : Adapter 0 : Port 0 : Target 1 : LUN 0"
			disks, err := wmi.ListDisks(scope, wmi.DiskSelectorListForDiskNumberAndLocation)
			if err != nil {
				return fmt.Errorf("failed to list disks: %w", err)
			}

			err = wmi.ForEach(disks, func(disk *wmi.COMDispatchObject) error {
				num, err := wmi.GetDiskNumber(disk)
				if err != nil {
					return fmt.Errorf("failed to query disk number: %v, %w", disk, err)
				}

				location, err := wmi.GetDiskLocation(disk)
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
							case "Bus":
								d.Bus = strings.TrimSpace(itemSplit[1])
							default:
								klog.Warningf("Got unknown field : %s=%s", itemSplit[0], itemSplit[1])
							}
						}
					}

					if found {
						m[num] = d
					}
				}
				return nil
			})
			return err
		})
	})
	return m, err
}

func (imp DiskAPI) Rescan() error {
	return wmi.WithCOMThread(func() error {
		err := wmi.RescanDisks()
		if err != nil {
			return fmt.Errorf("error updating host storage cache output. err: %w", err)
		}
		return nil
	})
}

func (imp DiskAPI) IsDiskInitialized(diskNumber uint32) (bool, error) {
	var partitionStyle uint16
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, wmi.DiskSelectorListForPartitionStyle)
			if err != nil {
				return fmt.Errorf("error checking initialized status of disk %d: %w", diskNumber, err)
			}

			partitionStyle, err = wmi.GetDiskPartitionStyle(disk)
			if err != nil {
				return fmt.Errorf("failed to query partition style of disk %d: %w", diskNumber, err)
			}

			return nil
		})
	})
	return partitionStyle != wmi.PartitionStyleUnknown, err
}

func (imp DiskAPI) InitializeDisk(diskNumber uint32) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, nil)
			if err != nil {
				return fmt.Errorf("failed to initializing disk %d. error: %w", diskNumber, err)
			}

			err = wmi.InitializeDisk(disk, wmi.PartitionStyleGPT)
			if err != nil {
				return fmt.Errorf("failed to initializing disk %d: error: %w", diskNumber, err)
			}

			return nil
		})
	})
}

func (imp DiskAPI) BasicPartitionsExist(diskNumber uint32) (bool, error) {
	var exist bool
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			partitions, err := wmi.ListPartitionsWithFilters(scope, nil, wmi.FilterForPartitionOnDisk(diskNumber), wmi.FilterForPartitionsOfTypeNormal())
			if err != nil {
				return fmt.Errorf("error checking presence of partitions on disk %d:, %w", diskNumber, err)
			}

			exist = len(partitions) > 0
			return nil
		})
	})
	return exist, err
}

func (imp DiskAPI) CreateBasicPartition(diskNumber uint32) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, nil)
			if err != nil {
				return err
			}

			err = wmi.CreatePartition(
				disk,
				nil,                           // Size
				true,                          // UseMaximumSize
				nil,                           // Offset
				nil,                           // Alignment
				nil,                           // DriveLetter
				false,                         // AssignDriveLetter
				nil,                           // MbrType,
				wmi.GPTPartitionTypeBasicData, // GPT Type
				false,                         // IsHidden
				false,                         // IsActive,
			)
			if err != nil {
				var werr *wmi.WMIError
				if !errors.As(err, &werr) || werr.Code != wmi.ErrorCodeCreatePartitionAccessPathAlreadyInUse {
					return fmt.Errorf("error creating partition on disk %d. err: %w", diskNumber, err)
				}
			}

			_, err = wmi.RefreshDisk(disk)
			if err != nil {
				return fmt.Errorf("error rescan disk (%d). error: %w", diskNumber, err)
			}

			partitions, err := wmi.ListPartitionsWithFilters(scope, nil, wmi.FilterForPartitionOnDisk(diskNumber), wmi.FilterForPartitionsOfTypeNormal())
			if err != nil {
				return fmt.Errorf("error query basic partition on disk %d:, %w", diskNumber, err)
			}

			if len(partitions) == 0 {
				return fmt.Errorf("failed to create basic partition on disk %d:, %w", diskNumber, err)
			}

			partition := partitions[0]
			status, err := wmi.SetPartitionState(partition, true)
			if err != nil {
				return fmt.Errorf("error bring partition %v on disk %d online. status %s, err: %w", partition, diskNumber, status, err)
			}

			return nil
		})
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
		return "", fmt.Errorf("IOCTL_STORAGE_QUERY_PROPERTY failed: %w", err)
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
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disks, err := wmi.ListDisks(scope, wmi.DiskSelectorListForPathAndSerialNumber)
			if err != nil {
				return fmt.Errorf("failed to list disks: %w", err)
			}

			found := false
			err = wmi.ForEach(disks, func(disk *wmi.COMDispatchObject) error {
				path, err := wmi.GetDiskPath(disk)
				if err != nil {
					return fmt.Errorf("failed to query disk path: %v, %w", disk, err)
				}

				diskNumber, diskPage83ID, err := imp.GetDiskNumberAndPage83ID(path)
				if err != nil {
					return err
				}

				if diskPage83ID == page83ID {
					diskNumberResult = diskNumber
					found = true
					return wmi.ErrStopIteration
				}
				return nil
			})
			if err != nil {
				return err
			}

			if !found {
				return fmt.Errorf("could not find disk with Page83 ID %s: %w", page83ID, wmi.ErrNotFound)
			}
			return nil
		})
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
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disks, err := wmi.ListDisks(scope, wmi.DiskSelectorListForPathAndSerialNumber)
			if err != nil {
				return fmt.Errorf("failed to list disks: %w", err)
			}

			err = wmi.ForEach(disks, func(disk *wmi.COMDispatchObject) error {
				path, err := wmi.GetDiskPath(disk)
				if err != nil {
					return fmt.Errorf("failed to query disk path: %v, %w", disk, err)
				}

				sn, err := wmi.GetDiskSerialNumber(disk)
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
				return nil
			})
			return err
		})
	})
	return m, err
}

func (imp DiskAPI) GetDiskStats(diskNumber uint32) (size int64, err error) {
	// TODO: change to uint64 as it does not make sense to use int64 for size
	size = -1
	err = wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, wmi.DiskSelectorListForSize)
			if err != nil {
				return err
			}

			sz, err := wmi.GetDiskSize(disk)
			if err != nil {
				return fmt.Errorf("failed to query size of disk %d. %w", diskNumber, err)
			}

			size = int64(sz)
			return nil
		})
	})
	return size, err
}

func (imp DiskAPI) SetDiskState(diskNumber uint32, isOnline bool) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, wmi.DiskSelectorListForIsOffline)
			if err != nil {
				return err
			}

			isOffline, err := wmi.IsDiskOffline(disk)
			if err != nil {
				return fmt.Errorf("error setting disk %d attach state. error: %w", diskNumber, err)
			}

			if isOnline == !isOffline {
				klog.V(2).Infof("Disk %d is already in the desired state", diskNumber)
				return nil
			}

			_, err = wmi.SetDiskState(disk, isOnline)
			if err != nil {
				return fmt.Errorf("setting disk %d attach state (isOnline: %v): error: %w", diskNumber, isOnline, err)
			}

			return nil
		})
	})
}

func (imp DiskAPI) GetDiskState(diskNumber uint32) (bool, error) {
	var isOffline bool
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			disk, err := wmi.QueryDiskByNumber(scope, diskNumber, wmi.DiskSelectorListForIsOffline)
			if err != nil {
				return err
			}

			isOffline, err = wmi.IsDiskOffline(disk)
			if err != nil {
				return fmt.Errorf("error parsing disk %d state. error: %w", diskNumber, err)
			}

			return nil
		})
	})
	return !isOffline, err
}
