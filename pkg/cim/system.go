//go:build windows
// +build windows

package cim

import (
	"fmt"

	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/server2019/root/cimv2"
)

// QueryBIOSElement retrieves the BIOS element.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM CIM_BIOSElement
//
// Refer to https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/cim-bioselement
// for the WMI class definition.
func QueryBIOSElement(selectorList []string) (*cimv2.CIM_BIOSElement, error) {
	biosQuery := query.NewWmiQueryWithSelectList("CIM_BIOSElement", selectorList)
	instances, err := QueryInstances("", biosQuery)
	if err != nil {
		return nil, err
	}

	bios, err := cimv2.NewCIM_BIOSElementEx1(instances[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get BIOS element: %w", err)
	}

	return bios, err
}

// QueryServiceByName retrieves a specific service by its name.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM Win32_Service
//
// Refer to https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-service
// for the WMI class definition.
func QueryServiceByName(name string, selectorList []string) (*cimv2.Win32_Service, error) {
	serviceQuery := query.NewWmiQueryWithSelectList("Win32_Service", selectorList, "Name", name)
	instances, err := QueryInstances("", serviceQuery)
	if err != nil {
		return nil, err
	}

	service, err := cimv2.NewWin32_ServiceEx1(instances[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", name, err)
	}

	return service, err
}
