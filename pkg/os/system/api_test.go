package system

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/pkg/errors"
)

type MockService struct {
	Name        string
	DisplayName string
	State       string
	StartMode   string
	Dependents  []cim.ServiceInterface

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

func (m *MockService) GetDependents() ([]cim.ServiceInterface, error) {
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

func (m *MockService) Refresh() error {
	return nil
}

type MockServiceFactory struct {
	Services map[string]cim.ServiceInterface
	Err      error
}

func (f *MockServiceFactory) GetService(name string) (cim.ServiceInterface, error) {
	svc, ok := f.Services[name]
	if !ok {
		return nil, fmt.Errorf("service not found: %s", name)
	}
	return svc, f.Err
}

func TestWaitUntilServiceState_Success(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateChanged := false

	stateCheck := func(s cim.ServiceInterface, state string) (bool, string, error) {
		if stateChanged {
			svc.State = serviceStateRunning
			return true, svc.State, nil
		}
		return false, svc.State, nil
	}

	stateTransition := func(s cim.ServiceInterface) error {
		stateChanged = true
		return nil
	}

	impl := ServiceManagerImpl{}
	state, err := impl.WaitUntilServiceState(svc, stateTransition, stateCheck, 10*time.Millisecond, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != serviceStateRunning {
		t.Fatalf("expected state %q, got %q", serviceStateRunning, state)
	}
}

func TestWaitUntilServiceState_Timeout(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateCheck := func(s cim.ServiceInterface, state string) (bool, string, error) {
		return false, svc.State, nil
	}

	stateTransition := func(s cim.ServiceInterface) error {
		return nil
	}

	impl := ServiceManagerImpl{}
	state, err := impl.WaitUntilServiceState(svc, stateTransition, stateCheck, 10*time.Millisecond, 50*time.Millisecond)
	if !errors.Is(err, errTimedOut) {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if state != svc.State {
		t.Fatalf("expected state %q, got %q", svc.State, state)
	}
}

func TestWaitUntilServiceState_TransitionFails(t *testing.T) {
	svc := &MockService{Name: "svc", State: "Stopped"}

	stateCheck := func(s cim.ServiceInterface, state string) (bool, string, error) {
		return false, svc.State, nil
	}

	stateTransition := func(s cim.ServiceInterface) error {
		return fmt.Errorf("transition failed")
	}

	impl := ServiceManagerImpl{}
	_, err := impl.WaitUntilServiceState(svc, stateTransition, stateCheck, 10*time.Millisecond, 50*time.Millisecond)
	if err == nil || err.Error() != "transition failed" {
		t.Fatalf("expected transition error, got %v", err)
	}
}

func TestGetDependentsForService(t *testing.T) {
	// Construct the dependency tree
	svcC := &MockService{Name: "C", State: serviceStateRunning}
	svcB := &MockService{Name: "B", State: serviceStateRunning, Dependents: []cim.ServiceInterface{svcC}}
	svcA := &MockService{Name: "A", State: serviceStateRunning, Dependents: []cim.ServiceInterface{svcB}}

	factory := &MockServiceFactory{
		Services: map[string]cim.ServiceInterface{
			"A": svcA,
			"B": svcB,
			"C": svcC,
		},
	}

	impl := ServiceManagerImpl{
		serviceFactory: factory,
	}

	names, err := impl.GetDependentsForService("A")
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
}

func TestGetDependentsForService_SkipsNonRunning(t *testing.T) {
	svcB := &MockService{Name: "B", State: "Stopped"}
	svcA := &MockService{Name: "A", State: serviceStateRunning, Dependents: []cim.ServiceInterface{svcB}}

	factory := &MockServiceFactory{
		Services: map[string]cim.ServiceInterface{
			"A": svcA,
			"B": svcB,
		},
	}

	impl := ServiceManagerImpl{
		serviceFactory: factory,
	}

	names, err := impl.GetDependentsForService("A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"A"} // B is skipped due to stopped state
	if len(names) != len(expected) {
		t.Fatalf("expected %d services, got %d", len(expected), len(names))
	}
}

func TestGetDependenciesForService_Winmgmt(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_MULTI_SERVICE_DEPENDENTS")) != "true" {
		t.Skipf("Test skipped")
	}

	impl := ServiceManagerImpl{
		serviceFactory: cim.Win32ServiceFactory{},
	}

	serviceName := "Winmgmt"
	names, err := impl.GetDependentsForService(serviceName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 4
	if len(names) != expected || names[len(names)-1] != serviceName {
		t.Fatalf("expected %d services, got %d", expected, len(names))
	}
}
