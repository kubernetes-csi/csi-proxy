package system

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/os/system"
	internal "github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	"k8s.io/klog/v2"
)

type Server struct {
	hostAPI API
}

type API interface {
	GetBIOSSerialNumber() (string, error)
	GetService(name string) (*system.ServiceInfo, error)
	StartService(name string) error
	StopService(name string, force bool) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) GetBIOSSerialNumber(context context.Context, request *internal.GetBIOSSerialNumberRequest, version apiversion.Version) (*internal.GetBIOSSerialNumberResponse, error) {
	klog.V(4).Infof("calling GetBIOSSerialNumber")
	response := &internal.GetBIOSSerialNumberResponse{}
	serialNumber, err := s.hostAPI.GetBIOSSerialNumber()
	if err != nil {
		klog.Errorf("failed GetBIOSSerialNumber: %v", err)
		return response, err
	}

	response.SerialNumber = serialNumber
	return response, nil
}

func (s *Server) GetService(context context.Context, request *internal.GetServiceRequest, version apiversion.Version) (*internal.GetServiceResponse, error) {
	klog.V(4).Infof("calling GetService name=%s", request.Name)
	response := &internal.GetServiceResponse{}
	info, err := s.hostAPI.GetService(request.Name)
	if err != nil {
		klog.Errorf("failed GetService: %v", err)
		return response, err
	}

	response.DisplayName = info.DisplayName
	response.StartType = internal.Startype(info.StartType)
	response.Status = internal.ServiceStatus(info.Status)
	return response, nil
}

func (s *Server) StartService(context context.Context, request *internal.StartServiceRequest, version apiversion.Version) (*internal.StartServiceResponse, error) {
	klog.V(4).Infof("calling StartService name=%s", request.Name)
	response := &internal.StartServiceResponse{}
	err := s.hostAPI.StartService(request.Name)
	if err != nil {
		klog.Errorf("failed StartService: %v", err)
		return response, err
	}

	return response, nil
}

func (s *Server) StopService(context context.Context, request *internal.StopServiceRequest, version apiversion.Version) (*internal.StopServiceResponse, error) {
	klog.V(4).Infof("calling StopService name=%s", request.Name)
	response := &internal.StopServiceResponse{}
	err := s.hostAPI.StopService(request.Name, request.Force)
	if err != nil {
		klog.Errorf("failed StopService: %v", err)
		return response, err
	}

	return response, nil
}
