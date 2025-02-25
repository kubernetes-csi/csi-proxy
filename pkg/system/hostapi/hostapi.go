package api

import (
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/v2/pkg/utils"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/server2019/root/cimv2"
)

// Implements the System OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// pkg/system/system.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type HostAPI interface {
	GetBIOSSerialNumber() (string, error)
	GetService(name string) (*ServiceInfo, error)
	StartService(name string) error
	StopService(name string, force bool) error
}

type systemAPI struct{}

func New() HostAPI {
	return systemAPI{}
}

func (systemAPI) GetBIOSSerialNumber() (string, error) {
	biosQuery := query.NewWmiQueryWithSelectList("CIM_BIOSElement", []string{"SerialNumber"})
	instances, err := cim.QueryInstances("", biosQuery)
	if err != nil {
		return "", err
	}

	bios, err := cimv2.NewCIM_BIOSElementEx1(instances[0])
	if err != nil {
		return "", fmt.Errorf("failed to get BIOS element: %w", err)
	}

	sn, err := bios.GetPropertySerialNumber()
	if err != nil {
		return "", fmt.Errorf("failed to get BIOS serial number property: %w", err)
	}

	return sn, nil
}

func (systemAPI) GetService(name string) (*ServiceInfo, error) {
	serviceQuery := query.NewWmiQueryWithSelectList("Win32_Service", []string{"DisplayName", "State", "StartMode"}, "Name", name)
	instances, err := cim.QueryInstances("", serviceQuery)
	if err != nil {
		return nil, err
	}

	service, err := cimv2.NewWin32_ServiceEx1(instances[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", name, err)
	}

	displayName, err := service.GetPropertyDisplayName()
	if err != nil {
		return nil, fmt.Errorf("failed to get displayName property of service %s: %w", name, err)
	}

	state, err := service.GetPropertyState()
	if err != nil {
		return nil, fmt.Errorf("failed to get state property of service %s: %w", name, err)
	}

	startMode, err := service.GetPropertyStartMode()
	if err != nil {
		return nil, fmt.Errorf("failed to get startMode property of service %s: %w", name, err)
	}

	return &ServiceInfo{
		DisplayName: displayName,
		StartType:   startMode,
		Status:      state,
	}, nil
}

func (systemAPI) StartService(name string) error {
	// Note: both StartService and StopService are not implemented by WMI
	script := `Start-Service -Name $env:ServiceName`
	cmdEnv := fmt.Sprintf("ServiceName=%s", name)
	out, err := utils.RunPowershellCmd(script, cmdEnv)
	if err != nil {
		return fmt.Errorf("error starting service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}

func (systemAPI) StopService(name string, force bool) error {
	script := `Stop-Service -Name $env:ServiceName -Force:$([System.Convert]::ToBoolean($env:Force))`
	out, err := utils.RunPowershellCmd(script, fmt.Sprintf("ServiceName=%s", name), fmt.Sprintf("Force=%t", force))
	if err != nil {
		return fmt.Errorf("error stopping service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}
