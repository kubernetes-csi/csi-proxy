package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
	"k8s.io/klog/v2"
)

type stateCheckFunc func(cim.ServiceInterface, string) (bool, string, error)
type stateTransitionFunc func(cim.ServiceInterface) error

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
	serviceStateCheckInternal = 200 * time.Millisecond
	serviceStateCheckTimeout  = 30 * time.Second
	errTimedOut               = errors.New("timed out")
)

type ServiceManager interface {
	WaitUntilServiceState(cim.ServiceInterface, stateTransitionFunc, stateCheckFunc, time.Duration, time.Duration) (string, error)
	GetDependentsForService(string) ([]string, error)
}

type ServiceFactory interface {
	GetService(name string) (cim.ServiceInterface, error)
}

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

type systemAPI struct {
	serviceFactory ServiceFactory
	serviceManager ServiceManager
}

func New() HostAPI {
	serviceFactory := cim.Win32ServiceFactory{}
	return systemAPI{
		serviceFactory: serviceFactory,
		serviceManager: ServiceManagerImpl{
			serviceFactory: serviceFactory,
		},
	}
}

func (systemAPI) GetBIOSSerialNumber() (string, error) {
	var sn string
	err := cim.WithCOMThread(func() error {
		bios, err := cim.QueryBIOSElement(cim.BIOSSelectorList)
		if err != nil {
			return fmt.Errorf("failed to get BIOS element: %w", err)
		}

		sn, err = cim.GetBIOSSerialNumber(bios)
		if err != nil {
			return fmt.Errorf("failed to get BIOS serial number property: %w", err)
		}

		return nil
	})
	return sn, err
}

func (impl systemAPI) GetService(name string) (*ServiceInfo, error) {
	var serviceInfo *ServiceInfo
	err := cim.WithCOMThread(func() error {
		service, err := impl.serviceFactory.GetService(name)
		if err != nil {
			return fmt.Errorf("failed to get service %s. error: %w", name, err)
		}

		displayName, err := cim.GetServiceDisplayName(service)
		if err != nil {
			return fmt.Errorf("failed to get displayName property of service %s: %w", name, err)
		}

		state, err := cim.GetServiceState(service)
		if err != nil {
			return fmt.Errorf("failed to get state property of service %s: %w", name, err)
		}

		startMode, err := cim.GetServiceStartMode(service)
		if err != nil {
			return fmt.Errorf("failed to get startMode property of service %s: %w", name, err)
		}

		serviceInfo = &ServiceInfo{
			DisplayName: displayName,
			StartType:   startMode,
			Status:      state,
		}
		return nil
	})
	return serviceInfo, err
}

func (impl systemAPI) StartService(name string) error {
	startService := func(service cim.ServiceInterface) error {
		retVal, err := service.StartService()
		if err != nil || (retVal != startServiceErrorCodeAccepted && retVal != startServiceErrorCodeAlreadyRunning) {
			return fmt.Errorf("error starting service name %s. return value: %d, error: %w", name, retVal, err)
		}
		return nil
	}
	serviceRunningCheck := func(service cim.ServiceInterface, state string) (bool, string, error) {
		err := service.Refresh()
		if err != nil {
			return false, "", err
		}

		newState, err := cim.GetServiceState(service)
		if err != nil {
			return false, state, err
		}

		klog.V(6).Infof("service (%v) state check: %s => %s", service, state, newState)
		return state == serviceStateRunning, newState, err
	}

	return cim.WithCOMThread(func() error {
		service, err := impl.serviceFactory.GetService(name)
		if err != nil {
			return fmt.Errorf("failed to get service %s. error: %w", name, err)
		}

		state, err := impl.serviceManager.WaitUntilServiceState(service, startService, serviceRunningCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
		if err != nil && !errors.Is(err, errTimedOut) {
			return fmt.Errorf("failed to wait for service %s state change. error: %w", name, err)
		}

		if state != serviceStateRunning {
			return fmt.Errorf("timed out waiting for service %s to become running", name)
		}

		return nil
	})
}

func (impl systemAPI) stopSingleService(name string) (bool, error) {
	var dependentRunning bool
	stopService := func(service cim.ServiceInterface) error {
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
	serviceStoppedCheck := func(service cim.ServiceInterface, state string) (bool, string, error) {
		err := service.Refresh()
		if err != nil {
			return false, "", fmt.Errorf("error refresh service %s instance. error: %w", name, err)
		}

		newState, err := cim.GetServiceState(service)
		if err != nil {
			return false, state, fmt.Errorf("error getting service %s state. error: %w", name, err)
		}

		klog.V(6).Infof("service (%v) state check: %s => %s", service, state, newState)
		return newState == serviceStateStopped, newState, nil
	}

	service, err := impl.serviceFactory.GetService(name)
	if err != nil {
		return dependentRunning, fmt.Errorf("failed to get service %s. error: %w", name, err)
	}

	state, err := impl.serviceManager.WaitUntilServiceState(service, stopService, serviceStoppedCheck, serviceStateCheckInternal, serviceStateCheckTimeout)
	if err != nil && !errors.Is(err, errTimedOut) {
		return dependentRunning, fmt.Errorf("error stopping service name %s. current state: %s", name, state)
	}

	if state != serviceStateStopped {
		return dependentRunning, fmt.Errorf("timed out waiting for service %s to stop", name)
	}

	return dependentRunning, nil
}

func (impl systemAPI) StopService(name string, force bool) error {
	return cim.WithCOMThread(func() error {
		dependentRunning, err := impl.stopSingleService(name)
		if err == nil {
			return nil
		}
		if !dependentRunning || !force {
			return fmt.Errorf("failed to stop service %s. error: %w", name, err)
		}

		serviceNames, err := impl.serviceManager.GetDependentsForService(name)
		if err != nil {
			return fmt.Errorf("error getting dependent services for service name %s", name)
		}

		for _, serviceName := range serviceNames {
			_, err = impl.stopSingleService(serviceName)
			if err != nil {
				return fmt.Errorf("failed to stop service %s. error: %w", name, err)
			}
		}

		return nil
	})
}

type ServiceManagerImpl struct {
	serviceFactory ServiceFactory
}

func (impl ServiceManagerImpl) WaitUntilServiceState(service cim.ServiceInterface, stateTransition stateTransitionFunc, stateCheck stateCheckFunc, interval time.Duration, timeout time.Duration) (string, error) {
	done, state, err := stateCheck(service, "")
	if err != nil {
		return state, fmt.Errorf("service %v state check failed: %w", service, err)
	}
	if done {
		return state, nil
	}

	// Perform transition if not already in desired state
	if err := stateTransition(service); err != nil {
		return state, fmt.Errorf("service %v state transition failed: %w", service, err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			klog.V(6).Infof("Checking service (%v) state...", service)
			done, state, err = stateCheck(service, state)
			if err != nil {
				return state, fmt.Errorf("service %v state check failed: %w", service, err)
			}
			if done {
				klog.V(6).Infof("service (%v) state is %s and transition done.", service, state)
				return state, nil
			}
		case <-timeoutChan:
			done, state, err = stateCheck(service, state)
			return state, errTimedOut
		}
	}
}

func (impl ServiceManagerImpl) GetDependentsForService(name string) ([]string, error) {
	var serviceNames []string
	var servicesToCheck []cim.ServiceInterface
	servicesByName := map[string]string{}

	service, err := impl.serviceFactory.GetService(name)
	if err != nil {
		return serviceNames, fmt.Errorf("failed to get service %s. error: %w", name, err)
	}

	servicesToCheck = append(servicesToCheck, service)
	i := 0
	for i < len(servicesToCheck) {
		service = servicesToCheck[i]
		i += 1

		serviceName, err := cim.GetServiceName(service)
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service name %v. error: %w", service, err)
		}

		currentState, err := cim.GetServiceState(service)
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service %s state. error: %w", serviceName, err)
		}

		if currentState != serviceStateRunning {
			continue
		}

		servicesByName[serviceName] = serviceName
		// prepend the current service to the front
		serviceNames = append([]string{serviceName}, serviceNames...)

		dependents, err := service.GetDependents()
		if err != nil {
			return serviceNames, fmt.Errorf("error getting service %s dependents. error: %w", serviceName, err)
		}

		servicesToCheck = append(servicesToCheck, dependents...)
	}

	return serviceNames, nil
}
