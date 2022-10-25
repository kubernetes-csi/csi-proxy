package system

import (
	"context"

	systemapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/system/hostapi"
	"k8s.io/klog/v2"
)

type System struct {
	hostAPI systemapi.HostAPI
}

type Interface interface {
	// GetBIOSSerialNumber returns the device's serial number
	GetBIOSSerialNumber(context.Context, *GetBIOSSerialNumberRequest) (*GetBIOSSerialNumberResponse, error)

	// GetService queries a Windows service state
	GetService(context.Context, *GetServiceRequest) (*GetServiceResponse, error)

	// StartService starts a Windows service
	// NOTE: This method affects global node state and should only be used
	//       with consideration to other CSI drivers that run concurrently.
	StartService(context.Context, *StartServiceRequest) (*StartServiceResponse, error)

	// StopService stops a Windows service
	// NOTE: This method affects global node state and should only be used
	//       with consideration to other CSI drivers that run concurrently.
	StopService(context.Context, *StopServiceRequest) (*StopServiceResponse, error)
}

// check that System implements Interface
var _ Interface = &System{}

func New(hostAPI systemapi.HostAPI) (*System, error) {
	return &System{
		hostAPI: hostAPI,
	}, nil
}

func (s *System) GetBIOSSerialNumber(context context.Context, request *GetBIOSSerialNumberRequest) (*GetBIOSSerialNumberResponse, error) {
	klog.V(4).Infof("calling GetBIOSSerialNumber")
	response := &GetBIOSSerialNumberResponse{}
	serialNumber, err := s.hostAPI.GetBIOSSerialNumber()
	if err != nil {
		klog.Errorf("failed GetBIOSSerialNumber: %v", err)
		return response, err
	}

	response.SerialNumber = serialNumber
	return response, nil
}

func (s *System) GetService(context context.Context, request *GetServiceRequest) (*GetServiceResponse, error) {
	klog.V(4).Infof("calling GetService name=%s", request.Name)
	response := &GetServiceResponse{}
	info, err := s.hostAPI.GetService(request.Name)
	if err != nil {
		klog.Errorf("failed GetService: %v", err)
		return response, err
	}

	response.DisplayName = info.DisplayName
	response.StartType = Startype(info.StartType)
	response.Status = ServiceStatus(info.Status)
	return response, nil
}

func (s *System) StartService(context context.Context, request *StartServiceRequest) (*StartServiceResponse, error) {
	klog.V(4).Infof("calling StartService name=%s", request.Name)
	response := &StartServiceResponse{}
	err := s.hostAPI.StartService(request.Name)
	if err != nil {
		klog.Errorf("failed StartService: %v", err)
		return response, err
	}

	return response, nil
}

func (s *System) StopService(context context.Context, request *StopServiceRequest) (*StopServiceResponse, error) {
	klog.V(4).Infof("calling StopService name=%s", request.Name)
	response := &StopServiceResponse{}
	err := s.hostAPI.StopService(request.Name, request.Force)
	if err != nil {
		klog.Errorf("failed StopService: %v", err)
		return response, err
	}

	return response, nil
}
