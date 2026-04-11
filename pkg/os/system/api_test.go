package system

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	wmi "github.com/kubernetes-csi/csi-proxy/pkg/cim"
)

type MockService struct {
	Name        string
	DisplayName string
	State       string
	StartMode   string
	Dependents  []wmi.ServiceInterface

	StartResult uint32
	StopResult  uint32

	Err error
}

func (m *MockService) GetPropertyName() (string, error) {
	return m.Name, m.Err
}

func (m *MockService) GetPropertyDisplayName() (string, error) {
	return m.DisplayName, m.Err
}

func (m *MockService) GetPropertyState() (string, error) {
	return m.State, m.Err
}

func (m *MockService) GetPropertyStartMode() (string, error) {
	return m.StartMode, m.Err
}

func (m *MockService) GetDependents(_ *wmi.Scope) ([]wmi.ServiceInterface, error) {
	return m.Dependents, m.Err
}

func (m *MockService) StartService() (uint32, error) {
	m.State = "Running"
	return m.StartResult, m.Err
}

func (m *MockService) StopService() (uint32, error) {
	m.State = "Stopped"
	return m.StopResult, m.Err
}

func (m *MockService) Refresh(_ *wmi.Scope) error {
	return nil
}

var _ wmi.ServiceInterface = &MockService{}

type MockServiceFactory struct {
	Services map[string]wmi.ServiceInterface
	Err      error
}

func (f *MockServiceFactory) GetService(_ *wmi.Scope, name string) (wmi.ServiceInterface, error) {
	svc, ok := f.Services[name]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	return svc, f.Err
}

var _ ServiceFactory = &MockServiceFactory{}

func TestWaitUntilServiceState_Success(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateChanged := false

	stateCheck := func(_ *wmi.Scope, _ wmi.ServiceInterface, _ string) (bool, string, error) {
		if stateChanged {
			svc.State = serviceStateRunning
			return true, svc.State, nil
		}
		return false, svc.State, nil
	}

	stateTransition := func(_ *wmi.Scope, _ wmi.ServiceInterface) error {
		stateChanged = true
		return nil
	}

	impl := ServiceManagerImpl{}
	err := wmi.WithScope(func(scope *wmi.Scope) error {
		state, err := impl.WaitUntilServiceState(scope, svc, stateTransition, stateCheck, 10*time.Millisecond, 500*time.Millisecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if state != serviceStateRunning {
			t.Fatalf("expected state %q, got %q", serviceStateRunning, state)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitUntilServiceState_Timeout(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateCheck := func(_ *wmi.Scope, _ wmi.ServiceInterface, _ string) (bool, string, error) {
		return false, svc.State, nil
	}

	stateTransition := func(_ *wmi.Scope, _ wmi.ServiceInterface) error {
		return nil
	}

	impl := ServiceManagerImpl{}
	err := wmi.WithScope(func(scope *wmi.Scope) error {
		state, err := impl.WaitUntilServiceState(scope, svc, stateTransition, stateCheck, 10*time.Millisecond, 50*time.Millisecond)
		if !errors.Is(err, errTimedOut) {
			t.Fatalf("expected timeout error, got %v", err)
		}
		if state != svc.State {
			t.Fatalf("expected state %q, got %q", svc.State, state)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitUntilServiceState_TransitionFails(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateCheck := func(_ *wmi.Scope, _ wmi.ServiceInterface, _ string) (bool, string, error) {
		return false, svc.State, nil
	}

	stateTransition := func(_ *wmi.Scope, _ wmi.ServiceInterface) error {
		return fmt.Errorf("transition failed")
	}

	impl := ServiceManagerImpl{}
	err := wmi.WithScope(func(scope *wmi.Scope) error {
		_, err := impl.WaitUntilServiceState(scope, svc, stateTransition, stateCheck, 10*time.Millisecond, 50*time.Millisecond)
		if err == nil || !strings.Contains(err.Error(), "transition failed") {
			t.Fatalf("expected transition error, got %v", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDependentsForService(t *testing.T) {
	// Construct the dependency tree
	svcC := &MockService{Name: "C", State: serviceStateRunning}
	svcB := &MockService{Name: "B", State: serviceStateRunning, Dependents: []wmi.ServiceInterface{svcC}}
	svcA := &MockService{Name: "A", State: serviceStateRunning, Dependents: []wmi.ServiceInterface{svcB}}

	factory := &MockServiceFactory{
		Services: map[string]wmi.ServiceInterface{
			"A": svcA,
			"B": svcB,
			"C": svcC,
		},
	}

	impl := ServiceManagerImpl{
		serviceFactory: factory,
	}

	err := wmi.WithScope(func(scope *wmi.Scope) error {
		names, err := impl.GetDependentsForService(scope, "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []string{"C", "B", "A"}
		if len(names) != len(expected) {
			t.Fatalf("expected %d services, got %d", len(expected), len(names))
		}
		for i, name := range expected {
			if names[i] != name {
				t.Errorf("expected %s at position %d, got %s", name, i, names[i])
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDependentsForService_SkipsNonRunning(t *testing.T) {
	svcB := &MockService{Name: "B", State: "Stopped"}
	svcA := &MockService{Name: "A", State: serviceStateRunning, Dependents: []wmi.ServiceInterface{svcB}}

	factory := &MockServiceFactory{
		Services: map[string]wmi.ServiceInterface{
			"A": svcA,
			"B": svcB,
		},
	}

	impl := ServiceManagerImpl{
		serviceFactory: factory,
	}

	err := wmi.WithScope(func(scope *wmi.Scope) error {
		names, err := impl.GetDependentsForService(scope, "A")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []string{"A"} // B is skipped due to stopped state
		if len(names) != len(expected) {
			t.Fatalf("expected %d services, got %d", len(expected), len(names))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDependenciesForService_Winmgmt(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("Test skipped on non-Windows platform")
	}
	if strings.ToLower(os.Getenv("TEST_MULTI_SERVICE_DEPENDENTS")) != "true" {
		t.Skipf("Test skipped")
	}

	impl := ServiceManagerImpl{
		serviceFactory: wmi.Win32ServiceFactory{},
	}

	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			serviceName := "Winmgmt"
			names, err := impl.GetDependentsForService(scope, serviceName)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := 4
			if len(names) != expected || names[len(names)-1] != serviceName {
				t.Fatalf("expected %d services, got %d", expected, len(names))
			}
			return nil
		})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
