package system

import (
	"errors"
	"fmt"
	"time"

	wmi "github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	"k8s.io/klog/v2"
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

type stateCheckFunc func(*wmi.Scope, wmi.ServiceInterface, string) (bool, string, error)
type stateTransitionFunc func(*wmi.Scope, wmi.ServiceInterface) error

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

	serviceStateCheckInternal = 200 * time.Millisecond
	serviceStateCheckTimeout  = 30 * time.Second
	errTimedOut               = errors.New("timed out")
)

func serviceStartModeToStartType(startMode string) uint32 {
	return startModeMappings[startMode]
}

func serviceState(status string) uint32 {
	return stateMappings[status]
}

type ServiceManager interface {
	WaitUntilServiceState(scope *wmi.Scope, service wmi.ServiceInterface, stateTransition stateTransitionFunc, stateCheck stateCheckFunc, interval time.Duration, timeout time.Duration) (string, error)
	GetDependentsForService(scope *wmi.Scope, name string) ([]string, error)
}

type ServiceFactory interface {
	GetService(scope *wmi.Scope, name string) (wmi.ServiceInterface, error)
}

type APIImplementor struct {
	serviceFactory ServiceFactory
	serviceManager ServiceManager
}

func New() APIImplementor {
	serviceFactory := wmi.Win32ServiceFactory{}
	return APIImplementor{
		serviceFactory: serviceFactory,
		serviceManager: ServiceManagerImpl{
			serviceFactory: serviceFactory,
		},
	}
}

func (APIImplementor) GetBIOSSerialNumber() (string, error) {
	var sn string
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			bios, err := wmi.QueryBIOSElement(scope, wmi.BIOSSelectorList)
			if err != nil {
				return fmt.Errorf("failed to get BIOS element: %w", err)
			}

			sn, err = wmi.GetBIOSSerialNumber(bios)
			if err != nil {
				return fmt.Errorf("failed to get BIOS serial number property: %w", err)
			}

			return nil
		})
	})
	return sn, err
}

func (impl APIImplementor) GetService(name string) (*ServiceInfo, error) {
	var serviceInfo *ServiceInfo
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			service, err := impl.serviceFactory.GetService(scope, name)
			if err != nil {
				return fmt.Errorf("failed to get service %s. error: %w", name, err)
			}

			displayName, err := wmi.GetServiceDisplayName(service)
			if err != nil {
				return fmt.Errorf("failed to get displayName property of service %s: %w", name, err)
			}

			state, err := wmi.GetServiceState(service)
			if err != nil {
				return fmt.Errorf("failed to get state property of service %s: %w", name, err)
			}

			startMode, err := wmi.GetServiceStartMode(service)
			if err != nil {
				return fmt.Errorf("failed to get startMode property of service %s: %w", name, err)
			}

			serviceInfo = &ServiceInfo{
				DisplayName: displayName,
				StartType:   serviceStartModeToStartType(startMode),
				Status:      serviceState(state),
			}

			return nil
		})
	})
	return serviceInfo, err
}

func (impl APIImplementor) StartService(name string) error {
	startService := func(_ *wmi.Scope, service wmi.ServiceInterface) error {
		retVal, err := service.StartService()
		if err != nil || (retVal != startServiceErrorCodeAccepted && retVal != startServiceErrorCodeAlreadyRunning) {
			return fmt.Errorf("error starting service name %s. return value: %d, error: %w", name, retVal, err)
		}
		return nil
	}
	serviceRunningCheck := func(scope *wmi.Scope, service wmi.ServiceInterface, state string) (bool, string, error) {
		err := service.Refresh(scope)
		if err != nil {
			return false, "", err
		}

		newState, err := wmi.GetServiceState(service)
		if err != nil {
			return false, state, err
		}

		klog.V(6).Infof("service (%v) state check: %s => %s", service, state, newState)
		return state == serviceStateRunning, newState, err
	}

	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			service, err := impl.serviceFactory.GetService(scope, name)
			if err != nil {
				return fmt.Errorf("failed to get service %s. error: %w", name, err)
			}

			state, err := impl.serviceManager.WaitUntilServiceState(scope, service, startService, serviceRunningCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
			if err != nil && !errors.Is(err, errTimedOut) {
				return fmt.Errorf("failed to wait for service %s state change. error: %w", name, err)
			}

			if state != serviceStateRunning {
				return fmt.Errorf("timed out waiting for service %s to become running", name)
			}

			return nil
		})
	})
}

func (impl APIImplementor) stopSingleService(scope *wmi.Scope, name string) (bool, error) {
	var dependentRunning bool
	stopService := func(_ *wmi.Scope, service wmi.ServiceInterface) error {
		retVal, err := service.StopService()
		if err != nil || (retVal != stopServiceErrorCodeAccepted && retVal != stopServiceErrorCodeStopPending) {
			if retVal == stopServiceErrorCodeDependentRunning {
				dependentRunning = true
				return fmt.Errorf("error stopping service %s as dependent services are not stopped", name)
			}
			return fmt.Errorf("error stopping service %s. return value: %d, error: %w", name, retVal, err)
		}
		return nil
	}
	serviceStoppedCheck := func(scope *wmi.Scope, service wmi.ServiceInterface, state string) (bool, string, error) {
		err := service.Refresh(scope)
		if err != nil {
			return false, "", fmt.Errorf("error refresh service %s instance. error: %w", name, err)
		}

		newState, err := wmi.GetServiceState(service)
		if err != nil {
			return false, state, fmt.Errorf("error getting service %s state. error: %w", name, err)
		}

		klog.V(6).Infof("service (%v) state check: %s => %s", service, state, newState)
		return newState == serviceStateStopped, newState, nil
	}

	service, err := impl.serviceFactory.GetService(scope, name)
	if err != nil {
		return dependentRunning, fmt.Errorf("failed to get service %s. error: %w", name, err)
	}

	state, err := impl.serviceManager.WaitUntilServiceState(scope, service, stopService, serviceStoppedCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
	if err != nil && !errors.Is(err, errTimedOut) {
		return dependentRunning, fmt.Errorf("error stopping service name %s. current state: %s", name, state)
	}

	if state != serviceStateStopped {
		return dependentRunning, fmt.Errorf("timed out waiting for service %s to stop", name)
	}

	return dependentRunning, nil
}

func (impl APIImplementor) StopService(name string, force bool) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			dependentRunning, err := impl.stopSingleService(scope, name)
			if err == nil {
				return nil
			}
			if !dependentRunning || !force {
				return fmt.Errorf("failed to stop service %s. error: %w", name, err)
			}

			serviceNames, err := impl.serviceManager.GetDependentsForService(scope, name)
			if err != nil {
				return fmt.Errorf("error getting dependent services for service name %s: %w", name, err)
			}

			for _, serviceName := range serviceNames {
				_, err = impl.stopSingleService(scope, serviceName)
				if err != nil {
					return fmt.Errorf("failed to stop service %s. error: %w", name, err)
				}
			}

			return nil
		})
	})
}

type ServiceManagerImpl struct {
	serviceFactory ServiceFactory
}

func (impl ServiceManagerImpl) WaitUntilServiceState(scope *wmi.Scope, service wmi.ServiceInterface, stateTransition stateTransitionFunc, stateCheck stateCheckFunc, interval time.Duration, timeout time.Duration) (string, error) {
	done, state, err := stateCheck(scope, service, "")
	if err != nil {
		return state, fmt.Errorf("service %v state check failed: %w", service, err)
	}
	if done {
		return state, nil
	}

	// Perform transition if not already in desired state
	if err := stateTransition(scope, service); err != nil {
		return state, fmt.Errorf("service %v state transition failed: %w", service, err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			klog.V(6).Infof("Checking service (%v) state...", service)
			done, state, err = stateCheck(scope, service, state)
			if err != nil {
				return state, fmt.Errorf("service %v state check failed: %w", service, err)
			}
			if done {
				klog.V(6).Infof("service (%v) state is %s and transition done.", service, state)
				return state, nil
			}
		case <-timeoutChan:
			done, state, err = stateCheck(scope, service, state)
			return state, errTimedOut
		}
	}
}

func (impl ServiceManagerImpl) GetDependentsForService(scope *wmi.Scope, name string) ([]string, error) {
	var serviceNames []string
	var servicesToCheck []wmi.ServiceInterface
	servicesByName := map[string]string{}

	service, err := impl.serviceFactory.GetService(scope, name)
	if err != nil {
		return serviceNames, fmt.Errorf("failed to get service %s. error: %w", name, err)
	}

	servicesToCheck = append(servicesToCheck, service)

	i := 0
	for i < len(servicesToCheck) {
		service = servicesToCheck[i]
		i += 1

		serviceName, err := wmi.GetServiceName(service)
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service name %v. error: %w", service, err)
		}

		currentState, err := wmi.GetServiceState(service)
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service %s state. error: %w", serviceName, err)
		}

		if currentState != serviceStateRunning {
			continue
		}

		servicesByName[serviceName] = serviceName
		// prepend the current service to the front
		serviceNames = append([]string{serviceName}, serviceNames...)

		dependents, err := service.GetDependents(scope)
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service %s dependents. error: %w", serviceName, err)
		}

		servicesToCheck = append(servicesToCheck, dependents...)
	}

	return serviceNames, nil
}
