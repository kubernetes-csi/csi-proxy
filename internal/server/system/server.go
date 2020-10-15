package system

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/system/internal"
	"k8s.io/klog"
)

type Server struct {
	hostAPI API
}

type API interface {
	GetBIOSSerialNumber() (string, error)
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
