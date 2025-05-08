package system

import (
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	"github.com/kubernetes-csi/csi-proxy/pkg/utils"
)

// Implements the System OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/system/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type ServiceInfo struct {
	// Service display name
	DisplayName string `json:"DisplayName"`

	// Service start type
	StartType uint32 `json:"StartType"`

	// Service status
	Status uint32 `json:"Status"`
}

var (
	startModeMappings = map[string]uint32{
		"Boot":     impl.START_TYPE_BOOT,
		"System":   impl.START_TYPE_SYSTEM,
		"Auto":     impl.START_TYPE_AUTOMATIC,
		"Manual":   impl.START_TYPE_MANUAL,
		"Disabled": impl.START_TYPE_DISABLED,
	}

	statusMappings = map[string]uint32{
		"Unknown":          impl.SERVICE_STATUS_UNKNOWN,
		"Stopped":          impl.SERVICE_STATUS_STOPPED,
		"Start Pending":    impl.SERVICE_STATUS_START_PENDING,
		"Stop Pending":     impl.SERVICE_STATUS_STOP_PENDING,
		"Running":          impl.SERVICE_STATUS_RUNNING,
		"Continue Pending": impl.SERVICE_STATUS_CONTINUE_PENDING,
		"Pause Pending":    impl.SERVICE_STATUS_PAUSE_PENDING,
		"Paused":           impl.SERVICE_STATUS_PAUSED,
	}
)

func serviceStartModeToStartType(startMode string) uint32 {
	return startModeMappings[startMode]
}

func serviceState(status string) uint32 {
	return statusMappings[status]
}

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) GetBIOSSerialNumber() (string, error) {
	bios, err := cim.QueryBIOSElement([]string{"SerialNumber"})
	if err != nil {
		return "", fmt.Errorf("failed to get BIOS element: %w", err)
	}

	sn, err := bios.GetPropertySerialNumber()
	if err != nil {
		return "", fmt.Errorf("failed to get BIOS serial number property: %w", err)
	}

	return sn, nil
}

func (APIImplementor) GetService(name string) (*ServiceInfo, error) {
	service, err := cim.QueryServiceByName(name, []string{"DisplayName", "State", "StartMode"})
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
		StartType:   serviceStartModeToStartType(startMode),
		Status:      serviceState(state),
	}, nil
}

func (APIImplementor) StartService(name string) error {
	// Note: both StartService and StopService are not implemented by WMI
	script := `Start-Service -Name $env:ServiceName`
	cmdEnv := fmt.Sprintf("ServiceName=%s", name)
	out, err := utils.RunPowershellCmd(script, cmdEnv)
	if err != nil {
		return fmt.Errorf("error starting service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}

func (APIImplementor) StopService(name string, force bool) error {
	script := `Stop-Service -Name $env:ServiceName -Force:$([System.Convert]::ToBoolean($env:Force))`
	out, err := utils.RunPowershellCmd(script, fmt.Sprintf("ServiceName=%s", name), fmt.Sprintf("Force=%t", force))
	if err != nil {
		return fmt.Errorf("error stopping service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}
