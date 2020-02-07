package disk

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	shared "github.com/kubernetes-csi/csi-proxy/internal/shared/disk"
	"k8s.io/klog"
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
		return nil, err
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
	_, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error updating host storage cache %v", err)
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
		// disks with raw initializtion not detected
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
