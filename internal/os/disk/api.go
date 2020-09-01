// +build windows

package disk

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	shared "github.com/kubernetes-csi/csi-proxy/internal/shared/disk"
	"k8s.io/klog"
)

var (
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")
)

const (
	IOCTL_STORAGE_GET_DEVICE_NUMBER = 0x2D1080
	IOCTL_STORAGE_QUERY_PROPERTY    = 0x002d1400
)

// Implements the OS API calls related to Disk Devices. All code here should be very simple
// pass-through to the OS APIs or cmdlets. Any logic around the APIs/cmdlet invocation
// should go in internal/server/filesystem/disk.go so that logic can be easily unit-tested
// without requiring specific OS environments.
type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

// ListDiskLocations - constructs a map with the disk number as the key and the DiskLocation structure
// as the value. The DiskLocation struct has various fields like the Adapter, Bus, Target and LUNID.
func (APIImplementor) ListDiskLocations() (map[string]shared.DiskLocation, error) {
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

	m := make(map[string]shared.DiskLocation)
	for _, v := range getDisk {
		str := v["location"].(string)
		num := fmt.Sprintf("%d", int(v["number"].(float64)))

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

func (APIImplementor) Rescan() error {
	cmd := "Update-HostStorageCache"
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error updating host storage cache output: %q, err: %v", string(out), err)
	}
	return nil
}

func (APIImplementor) IsDiskInitialized(diskID string) (bool, error) {
	cmd := fmt.Sprintf("Get-Disk -Number %s | Where partitionstyle -eq 'raw'", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking initialized status of disk %s: %v, %v", diskID, out, err)
	}
	if len(out) == 0 {
		// disks with raw initialization not detected
		return true, nil
	}
	return false, nil
}

func (APIImplementor) InitializeDisk(diskID string) error {
	cmd := fmt.Sprintf("Initialize-Disk -Number %s -PartitionStyle MBR", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error initializing disk %s: %v, %v", diskID, out, err)
	}
	return nil
}

func (APIImplementor) PartitionsExist(diskID string) (bool, error) {
	cmd := fmt.Sprintf("Get-Partition | Where DiskNumber -eq %s", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking presence of partitions on disk %s: %v, %v", diskID, out, err)
	}
	if len(out) > 0 {
		// disk has paritions in it
		return true, nil
	}
	return false, nil
}

func (APIImplementor) CreatePartition(diskID string) error {
	cmd := fmt.Sprintf("New-Partition -DiskNumber %s -UseMaximumSize", diskID)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating parition on disk %s: %v, %v", diskID, out, err)
	}
	return nil
}

func (imp APIImplementor) GetDiskNumberByName(diskName string) (string, error) {
	diskNumber, err := imp.GetDiskNumberWithID(diskName)
	return strconv.FormatUint(uint64(diskNumber), 10), err
}

func (APIImplementor) GetDiskNumber(disk syscall.Handle) (uint32, error) {
	var bytes uint32
	devNum := StorageDeviceNumber{}
	buflen := uint32(unsafe.Sizeof(devNum.DeviceType)) + uint32(unsafe.Sizeof(devNum.DeviceNumber)) + uint32(unsafe.Sizeof(devNum.PartitionNumber))

	err := syscall.DeviceIoControl(disk, IOCTL_STORAGE_GET_DEVICE_NUMBER, nil, 0, (*byte)(unsafe.Pointer(&devNum)), buflen, &bytes, nil)

	return devNum.DeviceNumber, err
}

func (APIImplementor) DiskHasPage83ID(disk syscall.Handle, matchID string) (bool, error) {
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
		if pID.CodeSet == StorageIDCodeSetASCII && pID.Association == StorageIDAssocDevice {
			for m = 0; m < pID.IdentifierSize; m++ {
				page83ID = append(page83ID, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&pID.Identifier[0])) + byteSize*uintptr(m))))
			}

			page83IDString := string(page83ID)
			if strings.Contains(page83IDString, matchID) {
				return true, nil
			}
		}
		pID = (*StorageIdentifier)(unsafe.Pointer(uintptr(unsafe.Pointer(pID)) + byteSize*uintptr(pID.NextOffset)))
	}
	return false, nil
}

func (APIImplementor) GetDiskPage83ID(disk syscall.Handle) (string, error) {
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
		if pID.CodeSet == StorageIDCodeSetASCII && pID.Association == StorageIDAssocDevice {
			for m = 0; m < pID.IdentifierSize; m++ {
				page83ID = append(page83ID, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&pID.Identifier[0])) + byteSize*uintptr(m))))
			}

			return string(page83ID), nil
		}
		pID = (*StorageIdentifier)(unsafe.Pointer(uintptr(unsafe.Pointer(pID)) + byteSize*uintptr(pID.NextOffset)))
	}
	return "", nil
}

func (imp APIImplementor) GetDiskNumberWithID(page83ID string) (uint32, error) {
	out, err := exec.Command("powershell.exe", "(get-disk | select Path) | ConvertTo-Json").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("Could not query disk paths")
	}

	outString := string(out)
	diskPaths := []DiskPath{}
	json.Unmarshal([]byte(outString), &diskPaths)

	for i := range diskPaths {
		h, err := syscall.Open(diskPaths[i].Path, syscall.O_RDONLY, 0)
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
func (imp APIImplementor) ListDiskIDs() (map[string]shared.DiskIDs, error) {
	out, err := exec.Command("powershell.exe", "(get-disk | select Path) | ConvertTo-Json").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Could not query disk paths")
	}

	outString := string(out)
	diskPaths := []DiskPath{}
	json.Unmarshal([]byte(outString), &diskPaths)

	m := make(map[string]shared.DiskIDs)

	for i := range diskPaths {
		h, err := syscall.Open(diskPaths[i].Path, syscall.O_RDONLY, 0)
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

		diskNumString := strconv.FormatUint(uint64(diskNumber), 10)

		diskIDs := make(map[string]string)
		diskIDs["page83"] = page83
		m[diskNumString] = shared.DiskIDs{Identifiers: diskIDs}
	}

	return m, nil
}

func (imp APIImplementor) DiskStats(diskID string) (int64, error) {
	cmd := fmt.Sprintf("(Get-Disk -Number %s).Size", diskID)
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
