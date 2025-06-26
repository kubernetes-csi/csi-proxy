//go:build windows
// +build windows

package cim

import (
	"fmt"
	"strconv"

	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/server2019/root/microsoft/windows/storage"
)

const (
	// PartitionStyleUnknown indicates an unknown partition table format
	PartitionStyleUnknown = 0
	// PartitionStyleGPT indicates the disk uses GUID Partition Table (GPT) format
	PartitionStyleGPT = 2

	// GPTPartitionTypeBasicData is the GUID for basic data partitions in GPT
	// Used for general purpose storage partitions
	GPTPartitionTypeBasicData = "{ebd0a0a2-b9e5-4433-87c0-68b6b72699c7}"
	// GPTPartitionTypeMicrosoftReserved is the GUID for Microsoft Reserved Partition (MSR)
	// Reserved by Windows for system use
	GPTPartitionTypeMicrosoftReserved = "{e3c9e316-0b5c-4db8-817d-f92df00215ae}"
)

// QueryDiskByNumber retrieves disk information for a specific disk identified by its number.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM MSFT_Disk
//	  WHERE DiskNumber = '<diskNumber>'
//
// Refer to https://learn.microsoft.com/en-us/windows-hardware/drivers/storage/msft-disk
// for the WMI class definition.
func QueryDiskByNumber(diskNumber uint32, selectorList []string) (*storage.MSFT_Disk, error) {
	diskQuery := query.NewWmiQueryWithSelectList("MSFT_Disk", selectorList, "Number", strconv.Itoa(int(diskNumber)))
	instances, err := QueryInstances(WMINamespaceStorage, diskQuery)
	if err != nil {
		return nil, err
	}

	disk, err := storage.NewMSFT_DiskEx1(instances[0])
	if err != nil {
		return nil, fmt.Errorf("failed to query disk %d. error: %v", diskNumber, err)
	}

	return disk, nil
}

// ListDisks retrieves information about all available disks.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM MSFT_Disk
//
// Refer to https://learn.microsoft.com/en-us/windows-hardware/drivers/storage/msft-disk
// for the WMI class definition.
func ListDisks(selectorList []string) ([]*storage.MSFT_Disk, error) {
	diskQuery := query.NewWmiQueryWithSelectList("MSFT_Disk", selectorList)
	instances, err := QueryInstances(WMINamespaceStorage, diskQuery)
	if IgnoreNotFound(err) != nil {
		return nil, err
	}

	var disks []*storage.MSFT_Disk
	for _, instance := range instances {
		disk, err := storage.NewMSFT_DiskEx1(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to query disk %v. error: %v", instance, err)
		}

		disks = append(disks, disk)
	}

	return disks, nil
}

// GetDiskNumber returns the number of a disk.
func GetDiskNumber(disk *storage.MSFT_Disk) (uint32, error) {
	number, err := disk.GetProperty("Number")
	if err != nil {
		return 0, err
	}
	return uint32(number.(int32)), err
}
