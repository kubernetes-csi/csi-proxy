//go:build windows
// +build windows

package cim

import (
	"fmt"
	"strconv"

	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/server2019/root/microsoft/windows/storage"
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
