package disk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	shared "github.com/kubernetes-csi/csi-proxy/internal/shared/disk"
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
	// PartitionsExist checks if the disk `diskNumber` has any partitions.
	PartitionsExist(diskNumber uint32) (bool, error)
	// CreatePartitoin creates a partition in disk `diskNumber`
	CreatePartition(diskNumber uint32) error
	// Rescan updates the host storage cache (re-enumerates disk, partition and volume objects)
	Rescan() error
	// GetDiskNumberByName gets a disk number by `diskName`.
	GetDiskNumberByName(diskName string) (uint32, error)
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
func (DiskAPI) ListDiskLocations() (map[uint32]shared.DiskLocation, error) {
	// sample response
	// [{
	//    "number":  0,
	//    "location":  "PCI Slot 3 : Adapter 0 : Port 0 : Target 1 : LUN 0"
	// }, ...]
	cmd := fmt.Sprintf("Get-Disk | select number, location | ConvertTo-Json")
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list disk location. cmd: %q, output: %q, err %v", cmd, string(out), err)
	}

	var getDisk []map[string]interface{}
	err = json.Unmarshal(out, &getDisk)
	if err != nil {
		return nil, err
	}

	m := make(map[uint32]shared.DiskLocation)
	for _, v := range getDisk {
		str := v["location"].(string)
		num := v["number"].(uint32)

		found := false
		s := strings.Split(str, ":")
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
	return m, nil
}

func (DiskAPI) Rescan() error {
	cmd := "Update-HostStorageCache"
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error updating host storage cache output: %q, err: %v", string(out), err)
	}
	return nil
}

func (DiskAPI) IsDiskInitialized(diskNumber uint32) (bool, error) {
	cmd := fmt.Sprintf("Get-Disk -Number %d | Where partitionstyle -eq 'raw'", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking initialized status of disk %d: %v, %v", diskNumber, out, err)
	}
	if len(out) == 0 {
		// disks with raw initialization not detected
		return true, nil
	}
	return false, nil
}

func (DiskAPI) InitializeDisk(diskNumber uint32) error {
	cmd := fmt.Sprintf("Initialize-Disk -Number %d -PartitionStyle GPT", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error initializing disk %d: %v, %v", diskNumber, out, err)
	}
	return nil
}

func (DiskAPI) PartitionsExist(diskNumber uint32) (bool, error) {
	cmd := fmt.Sprintf("Get-Partition | Where DiskNumber -eq %d", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking presence of partitions on disk %d: %v, %v", diskNumber, out, err)
	}
	if len(out) > 0 {
		// disk has partitions in it
		return true, nil
	}
	return false, nil
}

func (DiskAPI) CreatePartition(diskNumber uint32) error {
	cmd := fmt.Sprintf("New-Partition -DiskNumber %d -UseMaximumSize", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating parition on disk %d: %v, %v", diskNumber, out, err)
	}
	return nil
}

func (imp DiskAPI) GetDiskNumberByName(diskName string) (uint32, error) {
	diskNumber, err := imp.GetDiskNumberWithID(diskName)
	return diskNumber, err
}

func (DiskAPI) GetDiskNumber(disk syscall.Handle) (uint32, error) {
	var bytes uint32
	devNum := StorageDeviceNumber{}
	buflen := uint32(unsafe.Sizeof(devNum.DeviceType)) + uint32(unsafe.Sizeof(devNum.DeviceNumber)) + uint32(unsafe.Sizeof(devNum.PartitionNumber))

	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_GET_DEVICE_NUMBER, nil, 0, (*byte)(unsafe.Pointer(&devNum)), buflen, &bytes, nil)

	return devNum.DeviceNumber, err
}

func (DiskAPI) DiskHasPage83ID(disk syscall.Handle, matchID string) (bool, error) {
	query := StoragePropertyQuery{}

	bufferSize := uint32(4 * 1024)
	buffer := make([]byte, 4*1024)
	var size uint32
	var n uint32
	var m uint16

	query.QueryType = PropertyStandardQuery
	query.PropertyID = StorageDeviceIDProperty

	querySize := uint32(unsafe.Sizeof(query.PropertyID)) + uint32(unsafe.Sizeof(query.QueryType)) + uint32(unsafe.Sizeof(query.Byte))
	querySize = uint32(unsafe.Sizeof(query))
	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_QUERY_PROPERTY, (*byte)(unsafe.Pointer(&query)), querySize, (*byte)(unsafe.Pointer(&buffer[0])), bufferSize, &size, nil)
	if err != nil {
		return false, fmt.Errorf("IOCTL_STORAGE_QUERY_PROPERTY failed: %v", err)
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

			var page83IDString string
			if pID.CodeSet == StorageIDCodeSetASCII {
				page83IDString = string(page83ID)
			} else if pID.CodeSet == StorageIDCodeSetBinary {
				page83IDString = hex.EncodeToString(page83ID)
			}
			if strings.Contains(page83IDString, matchID) {
				return true, nil
			}
		}
		pID = (*StorageIdentifier)(unsafe.Pointer(uintptr(unsafe.Pointer(pID)) + byteSize*uintptr(pID.NextOffset)))
	}
	return false, nil
}

func (DiskAPI) GetDiskPage83ID(disk syscall.Handle) (string, error) {
	query := StoragePropertyQuery{}

	bufferSize := uint32(4 * 1024)
	buffer := make([]byte, 4*1024)
	var size uint32
	var n uint32
	var m uint16

	query.QueryType = PropertyStandardQuery
	query.PropertyID = StorageDeviceIDProperty

	querySize := uint32(unsafe.Sizeof(query.PropertyID)) + uint32(unsafe.Sizeof(query.QueryType)) + uint32(unsafe.Sizeof(query.Byte))
	querySize = uint32(unsafe.Sizeof(query))
	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_QUERY_PROPERTY, (*byte)(unsafe.Pointer(&query)), querySize, (*byte)(unsafe.Pointer(&buffer[0])), bufferSize, &size, nil)
	if err != nil {
		return "", fmt.Errorf("IOCTL_STORAGE_QUERY_PROPERTY failed: %v", err)
	}

	devIDDesc := (*StorageDeviceIDDescriptor)(unsafe.Pointer(&buffer[0]))

	pID := (*StorageIdentifier)(unsafe.Pointer(&devIDDesc.Identifiers[0]))

	page83ID := []byte{}
	byteSize := unsafe.Sizeof(byte(0))
	for n = 0; n < devIDDesc.NumberOfIdentifiers; n++ {
		if pID.Association == StorageIDAssocDevice {
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
	out, err := exec.Command("powershell.exe", "(get-disk | select Path) | ConvertTo-Json").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("Could not query disk paths")
	}

	outString := string(out)
	disks := []Disk{}
	json.Unmarshal([]byte(outString), &disks)

	for i := range disks {
		h, err := syscall.Open(disks[i].Path, syscall.O_RDONLY, 0)
		if err != nil {
			return 0, err
		}

		found, err := imp.DiskHasPage83ID(h, page83ID)
		if found {
			return imp.GetDiskNumber(h)
		}
	}

	return 0, fmt.Errorf("Could not find disk with Page83 ID %s", page83ID)
}

// ListDiskIDs - constructs a map with the disk number as the key and the DiskID structure
// as the value. The DiskID struct has a field for the page83 ID.
func (imp DiskAPI) ListDiskIDs() (map[uint32]shared.DiskIDs, error) {
	out, err := exec.Command("powershell.exe", "(get-disk | select Path, SerialNumber) | ConvertTo-Json").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Could not query disk paths")
	}

	outString := string(out)
	disks := []Disk{}
	json.Unmarshal([]byte(outString), &disks)

	m := make(map[uint32]shared.DiskIDs)

	for i := range disks {
		h, err := syscall.Open(disks[i].Path, syscall.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		page83, err := imp.GetDiskPage83ID(h)
		if err != nil {
			return m, fmt.Errorf("Could not get page83 ID: %v", err)
		}

		diskNumber, err := imp.GetDiskNumber(h)
		if err != nil {
			return m, fmt.Errorf("Could not get disk number: %v", err)
		}

		m[diskNumber] = shared.DiskIDs{
			Page83:       page83,
			SerialNumber: disks[i].SerialNumber,
		}
	}

	return m, nil
}

func (imp DiskAPI) GetDiskStats(diskNumber uint32) (int64, error) {
	cmd := fmt.Sprintf("(Get-Disk -Number %d).Size", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil || len(out) == 0 {
		return -1, fmt.Errorf("error getting size of disk. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return -1, fmt.Errorf("error compiling regex. err: %v", err)
	}
	diskSizeOutput := reg.ReplaceAllString(string(out), "")

	diskSize, err := strconv.ParseInt(diskSizeOutput, 10, 64)

	if err != nil {
		return -1, fmt.Errorf("error parsing size of disk. cmd: %s, output: %s, error: %v", cmd, diskSizeOutput, err)
	}

	return diskSize, nil
}

func (imp DiskAPI) SetDiskState(diskNumber uint32, isOnline bool) error {
	cmd := fmt.Sprintf("(Get-Disk -Number %d) | Set-Disk -IsOffline $%t", diskNumber, !isOnline)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error setting disk attach state. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	return nil
}

func (imp DiskAPI) GetDiskState(diskNumber uint32) (bool, error) {
	cmd := fmt.Sprintf("(Get-Disk -Number %d) | Select-Object -ExpandProperty IsOffline", diskNumber)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error getting disk state. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	sout := strings.TrimSpace(string(out))
	isOffline, err := strconv.ParseBool(sout)
	if err != nil {
		return false, fmt.Errorf("error parsing disk state. output: %s, error: %v", sout, err)
	}

	return !isOffline, nil
}
