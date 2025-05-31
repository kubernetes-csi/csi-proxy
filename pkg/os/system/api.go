package system

import (
	"fmt"
	"time"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	"github.com/microsoft/wmi/pkg/errors"
	wmiinst "github.com/microsoft/wmi/pkg/wmiinstance"
	"github.com/microsoft/wmi/server2019/root/cimv2"
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

type periodicalCheckFunc func() (bool, error)

const (
	// startServiceErrorCodeAccepted indicates the request is accepted
	startServiceErrorCodeAccepted = 0

	// startServiceErrorCodeAlreadyRunning indicates a service is already running
	startServiceErrorCodeAlreadyRunning = 10

	// stopServiceErrorCodeAccepted indicates the request is accepted
	stopServiceErrorCodeAccepted = 0

	// stopServiceErrorCodeStopPending indicates the request cannot be sent to the service because the state of the service is 0,1,2 (pending)
	stopServiceErrorCodeStopPending = 5

	// stopServiceErrorCodeDependentRunning indicates a service cannot be stopped as its dependents may still be running
	stopServiceErrorCodeDependentRunning = 3

	serviceStateRunning = "Running"
	serviceStateStopped = "Stopped"
)

var (
	startModeMappings = map[string]uint32{
		"Boot":     impl.START_TYPE_BOOT,
		"System":   impl.START_TYPE_SYSTEM,
		"Auto":     impl.START_TYPE_AUTOMATIC,
		"Manual":   impl.START_TYPE_MANUAL,
		"Disabled": impl.START_TYPE_DISABLED,
	}

	stateMappings = map[string]uint32{
		"Unknown":           impl.SERVICE_STATUS_UNKNOWN,
		serviceStateStopped: impl.SERVICE_STATUS_STOPPED,
		"Start Pending":     impl.SERVICE_STATUS_START_PENDING,
		"Stop Pending":      impl.SERVICE_STATUS_STOP_PENDING,
		serviceStateRunning: impl.SERVICE_STATUS_RUNNING,
		"Continue Pending":  impl.SERVICE_STATUS_CONTINUE_PENDING,
		"Pause Pending":     impl.SERVICE_STATUS_PAUSE_PENDING,
		"Paused":            impl.SERVICE_STATUS_PAUSED,
	}

	serviceStateCheckInternal = 500 * time.Millisecond
	serviceStateCheckTimeout  = 5 * time.Second
)

func serviceStartModeToStartType(startMode string) uint32 {
	return startModeMappings[startMode]
}

func serviceState(status string) uint32 {
	return stateMappings[status]
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

func waitForServiceState(serviceCheck periodicalCheckFunc, interval time.Duration, timeout time.Duration) error {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return errors.Timedout
		case <-ticker.C:
			done, err := serviceCheck()
			if err != nil {
				return err
			}

			if done {
				return nil
			}
		}
	}
}

func getServiceState(name string) (string, *cimv2.Win32_Service, error) {
	service, err := cim.QueryServiceByName(name, nil)
	if err != nil {
		return "", nil, err
	}

	state, err := service.GetPropertyState()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get state property of service %s: %w", name, err)
	}

	return state, service, nil
}

func (APIImplementor) StartService(name string) error {
	state, service, err := getServiceState(name)
	if err != nil {
		return err
	}

	if state != serviceStateRunning {
		var retVal uint32
		retVal, err = service.StartService()
		if err != nil || (retVal != startServiceErrorCodeAccepted && retVal != startServiceErrorCodeAlreadyRunning) {
			return fmt.Errorf("error starting service name %s. return value: %d, error: %v", name, retVal, err)
		}

		err = waitForServiceState(func() (bool, error) {
			state, service, err = getServiceState(name)
			if err != nil {
				return false, err
			}

			return state == serviceStateRunning, nil

		}, serviceStateCheckInternal, serviceStateCheckTimeout)
		if err != nil {
			return fmt.Errorf("error waiting service %s become running. error: %v", name, err)
		}
	}

	if state != serviceStateRunning {
		return fmt.Errorf("error starting service name %s. current state: %s", name, state)
	}

	return nil
}

func (APIImplementor) StopService(name string, force bool) error {
	state, service, err := getServiceState(name)
	if err != nil {
		return err
	}

	if state == serviceStateStopped {
		return nil
	}

	stopSingleService := func(name string, service *wmiinst.WmiInstance) (bool, error) {
		retVal, err := service.InvokeMethodWithReturn("StopService")
		if err != nil || (retVal != stopServiceErrorCodeAccepted && retVal != stopServiceErrorCodeStopPending) {
			if retVal == stopServiceErrorCodeDependentRunning {
				return true, fmt.Errorf("error stopping service %s as dependent services are not stopped", name)
			}
			return false, fmt.Errorf("error stopping service %s. return value: %d, error: %v", name, retVal, err)
		}

		var serviceState string
		err = waitForServiceState(func() (bool, error) {
			serviceState, _, err = getServiceState(name)
			if err != nil {
				return false, err
			}

			return serviceState == serviceStateStopped, nil

		}, serviceStateCheckInternal, serviceStateCheckTimeout)
		if err != nil {
			return false, fmt.Errorf("error waiting service %s become stopped. error: %v", name, err)
		}

		if serviceState != serviceStateStopped {
			return false, fmt.Errorf("error stopping service name %s. current state: %s", name, serviceState)
		}

		return false, nil
	}

	dependentRunning, err := stopSingleService(name, service.WmiInstance)
	if !force || err == nil || !dependentRunning {
		return err
	}

	var serviceNames []string
	var servicesToCheck wmiinst.WmiInstanceCollection
	servicesByName := map[string]*wmiinst.WmiInstance{}

	servicesToCheck = append(servicesToCheck, service.WmiInstance)
	i := 0
	for i < len(servicesToCheck) {
		current := servicesToCheck[i]
		i += 1

		currentNameVal, err := current.GetProperty("Name")
		if err != nil {
			return err
		}

		currentName := currentNameVal.(string)
		if _, ok := servicesByName[currentName]; ok {
			continue
		}

		currentStateVal, err := current.GetProperty("State")
		if err != nil {
			return err
		}

		currentState := currentStateVal
		if currentState != serviceStateRunning {
			continue
		}

		servicesByName[currentName] = current
		serviceNames = append(serviceNames, currentName)

		dependents, err := current.GetAssociated("Win32_DependentService", "Win32_Service", "Dependent", "Antecedent")
		if err != nil {
			return err
		}

		servicesToCheck = append(servicesToCheck, dependents...)
	}

	i = len(serviceNames) - 1
	for i >= 0 {
		serviceName := serviceNames[i]
		i -= 1

		state, service, err := getServiceState(serviceName)
		if err != nil {
			return err
		}

		if state == serviceStateStopped {
			continue
		}

		_, err = stopSingleService(serviceName, service.WmiInstance)
		if err != nil {
			return err
		}
	}

	return nil
}
