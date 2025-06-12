package system

import (
	"fmt"
	"time"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	wmierrors "github.com/microsoft/wmi/pkg/errors"
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

type stateCheckFunc func() (bool, string, ServiceInterface, error)
type stateTransitionFunc func(ServiceInterface) error

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
	serviceStateCheckTimeout  = 30 * time.Second
)

func serviceStartModeToStartType(startMode string) uint32 {
	return startModeMappings[startMode]
}

func serviceState(status string) uint32 {
	return stateMappings[status]
}

type ServiceInterface interface {
	GetPropertyName() (string, error)
	GetPropertyDisplayName() (string, error)
	GetPropertyState() (string, error)
	GetPropertyStartMode() (string, error)
	GetDependents() ([]ServiceInterface, error)
	StartService() (result uint32, err error)
	StopService() (result uint32, err error)
}

type ServiceManager interface {
	WaitUntilServiceState(stateTransition stateTransitionFunc, stateCheck stateCheckFunc, interval time.Duration, timeout time.Duration) (string, error)
	GetDependentsForService(name string) ([]string, error)
}

type ServiceFactory interface {
	GetService(name string) (ServiceInterface, error)
}

type APIImplementor struct {
	serviceFactory ServiceFactory
	serviceManager ServiceManager
}

func New() APIImplementor {
	serviceFactory := Win32ServiceFactory{}
	return APIImplementor{
		serviceFactory: serviceFactory,
		serviceManager: ServiceManagerImpl{
			serviceFactory: serviceFactory,
		},
	}
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

func (impl APIImplementor) StartService(name string) error {
	startService := func(service ServiceInterface) error {
		retVal, err := service.StartService()
		if err != nil || (retVal != startServiceErrorCodeAccepted && retVal != startServiceErrorCodeAlreadyRunning) {
			return fmt.Errorf("error starting service name %s. return value: %d, error: %v", name, retVal, err)
		}
		return nil
	}
	serviceRunningCheck := func() (bool, string, ServiceInterface, error) {
		service, err := impl.serviceFactory.GetService(name)
		if err != nil {
			return false, "", nil, err
		}

		state, err := service.GetPropertyState()
		if err != nil {
			return false, state, service, err
		}

		return state == serviceStateRunning, state, service, err
	}

	state, err := impl.serviceManager.WaitUntilServiceState(startService, serviceRunningCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
	if err != nil && !wmierrors.IsTimedout(err) {
		return err
	}

	if state != serviceStateRunning {
		return fmt.Errorf("timed out waiting for service %s to become running", name)
	}

	return nil
}

func (impl APIImplementor) stopSingleService(name string) (bool, error) {
	var dependentRunning bool
	stopService := func(service ServiceInterface) error {
		retVal, err := service.StopService()
		if err != nil || (retVal != stopServiceErrorCodeAccepted && retVal != stopServiceErrorCodeStopPending) {
			if retVal == stopServiceErrorCodeDependentRunning {
				dependentRunning = true
				return fmt.Errorf("error stopping service %s as dependent services are not stopped", name)
			}
			return fmt.Errorf("error stopping service %s. return value: %d, error: %v", name, retVal, err)
		}
		return nil
	}
	serviceStoppedCheck := func() (bool, string, ServiceInterface, error) {
		service, err := impl.serviceFactory.GetService(name)
		if err != nil {
			return false, "", nil, err
		}

		state, err := service.GetPropertyState()
		if err != nil {
			return false, state, service, err
		}

		return state == serviceStateStopped, state, service, err
	}

	state, err := impl.serviceManager.WaitUntilServiceState(stopService, serviceStoppedCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
	if err != nil && !wmierrors.IsTimedout(err) {
		return dependentRunning, fmt.Errorf("error stopping service name %s. current state: %s", name, state)
	}

	if state != serviceStateStopped {
		return dependentRunning, fmt.Errorf("timed out waiting for service %s to stop", name)
	}

	return dependentRunning, nil
}

func (impl APIImplementor) StopService(name string, force bool) error {
	dependentRunning, err := impl.stopSingleService(name)
	if err == nil || !dependentRunning || !force {
		return err
	}

	serviceNames, err := impl.serviceManager.GetDependentsForService(name)
	if err != nil {
		return fmt.Errorf("error getting dependent services for service name %s", name)
	}

	for _, serviceName := range serviceNames {
		_, err = impl.stopSingleService(serviceName)
		if err != nil {
			return err
		}
	}

	return nil
}

type Win32Service struct {
	*cimv2.Win32_Service
}

func (s *Win32Service) GetDependents() ([]ServiceInterface, error) {
	collection, err := s.GetAssociated("Win32_DependentService", "Win32_Service", "Dependent", "Antecedent")
	if err != nil {
		return nil, err
	}

	var result []ServiceInterface
	for _, coll := range collection {
		service, err := cimv2.NewWin32_ServiceEx1(coll)
		if err != nil {
			return nil, err
		}

		result = append(result, &Win32Service{
			service,
		})
	}
	return result, nil
}

type Win32ServiceFactory struct {
}

func (impl Win32ServiceFactory) GetService(name string) (ServiceInterface, error) {
	service, err := cim.QueryServiceByName(name, nil)
	if err != nil {
		return nil, err
	}

	return &Win32Service{Win32_Service: service}, nil
}

type ServiceManagerImpl struct {
	serviceFactory ServiceFactory
}

func (impl ServiceManagerImpl) WaitUntilServiceState(stateTransition stateTransitionFunc, stateCheck stateCheckFunc, interval time.Duration, timeout time.Duration) (string, error) {
	done, state, service, err := stateCheck()
	if err != nil {
		return state, err
	}
	if done {
		return state, err
	}

	// Perform transition if not already in desired state
	if err := stateTransition(service); err != nil {
		return state, err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			done, state, service, err = stateCheck()
			if err != nil {
				return state, fmt.Errorf("check failed: %w", err)
			}
			if done {
				return state, nil
			}
		case <-timeoutChan:
			done, state, _, err = stateCheck()
			return state, wmierrors.Timedout
		}
	}
}

func (impl ServiceManagerImpl) GetDependentsForService(name string) ([]string, error) {
	var serviceNames []string
	var servicesToCheck []ServiceInterface
	servicesByName := map[string]string{}

	service, err := impl.serviceFactory.GetService(name)
	if err != nil {
		return serviceNames, err
	}

	servicesToCheck = append(servicesToCheck, service)
	i := 0
	for i < len(servicesToCheck) {
		service = servicesToCheck[i]
		i += 1

		serviceName, err := service.GetPropertyName()
		if err != nil {
			return serviceNames, err
		}

		currentState, err := service.GetPropertyState()
		if err != nil {
			return serviceNames, err
		}

		if currentState != serviceStateRunning {
			continue
		}

		servicesByName[serviceName] = serviceName
		// prepend the current service to the front
		serviceNames = append([]string{serviceName}, serviceNames...)

		dependents, err := service.GetDependents()
		if err != nil {
			return serviceNames, err
		}

		servicesToCheck = append(servicesToCheck, dependents...)
	}

	return serviceNames, nil
}
