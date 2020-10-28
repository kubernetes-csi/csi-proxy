package iscsi

import (
	"context"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/os/iscsi"
	"github.com/kubernetes-csi/csi-proxy/internal/server/iscsi/internal"
	"k8s.io/klog"
)

const defaultIscsiPort = 3260

type Server struct {
	hostAPI API
}

type API interface {
	AddTargetPortal(portal *iscsi.TargetPortal) error
	DiscoverTargetPortal(portal *iscsi.TargetPortal) ([]string, error)
	ListTargetPortals() ([]iscsi.TargetPortal, error)
	RemoveTargetPortal(portal *iscsi.TargetPortal) error
	ConnectTarget(portal *iscsi.TargetPortal, iqn string, authType string,
		chapUser string, chapSecret string) error
	DisconnectTarget(portal *iscsi.TargetPortal, iqn string) error
	GetTargetDisks(portal *iscsi.TargetPortal, iqn string) ([]string, error)
	SetMutualChapSecret(mutualChapSecret string) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) requestTPtoAPITP(portal *internal.TargetPortal) *iscsi.TargetPortal {
	port := portal.TargetPort
	if port == 0 {
		port = defaultIscsiPort
	}
	return &iscsi.TargetPortal{Address: portal.TargetAddress, Port: port}
}

func (s *Server) AddTargetPortal(context context.Context, request *internal.AddTargetPortalRequest, version apiversion.Version) (*internal.AddTargetPortalResponse, error) {
	klog.V(4).Infof("calling AddTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &internal.AddTargetPortalResponse{}
	err := s.hostAPI.AddTargetPortal(s.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed AddTargetPortal %v", err)
		return response, err
	}

	return response, nil
}

func AuthTypeToString(authType internal.AuthenticationType) (string, error) {
	switch authType {
	case internal.NONE:
		return "NONE", nil
	case internal.ONE_WAY_CHAP:
		return "ONEWAYCHAP", nil
	case internal.MUTUAL_CHAP:
		return "MUTUALCHAP", nil
	default:
		return "", fmt.Errorf("invalid authentication type authType=%v", authType)
	}
}

func (s *Server) ConnectTarget(context context.Context, req *internal.ConnectTargetRequest, version apiversion.Version) (*internal.ConnectTargetResponse, error) {
	klog.V(4).Infof("calling ConnectTarget with portal %s:%d and iqn %s"+
		" auth=%v chapuser=%v", req.TargetPortal.TargetAddress,
		req.TargetPortal.TargetPort, req.Iqn, req.AuthType, req.ChapUsername)

	response := &internal.ConnectTargetResponse{}
	authType, err := AuthTypeToString(req.AuthType)
	if err != nil {
		klog.Errorf("Error parsing parameters: %v", err)
		return response, err
	}

	err = s.hostAPI.ConnectTarget(s.requestTPtoAPITP(req.TargetPortal), req.Iqn,
		authType, req.ChapUsername, req.ChapSecret)
	if err != nil {
		klog.Errorf("failed ConnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (s *Server) DisconnectTarget(context context.Context, request *internal.DisconnectTargetRequest, version apiversion.Version) (*internal.DisconnectTargetResponse, error) {
	klog.V(4).Infof("calling DisconnectTarget with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.Iqn)

	response := &internal.DisconnectTargetResponse{}
	err := s.hostAPI.DisconnectTarget(s.requestTPtoAPITP(request.TargetPortal), request.Iqn)
	if err != nil {
		klog.Errorf("failed DisconnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (s *Server) DiscoverTargetPortal(context context.Context, request *internal.DiscoverTargetPortalRequest, version apiversion.Version) (*internal.DiscoverTargetPortalResponse, error) {
	klog.V(4).Infof("calling DiscoverTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &internal.DiscoverTargetPortalResponse{}
	iqns, err := s.hostAPI.DiscoverTargetPortal(s.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed DiscoverTargetPortal %v", err)
		return response, err
	}

	response.Iqns = iqns
	return response, nil
}

func (s *Server) GetTargetDisks(context context.Context, request *internal.GetTargetDisksRequest, version apiversion.Version) (*internal.GetTargetDisksResponse, error) {
	klog.V(4).Infof("calling GetTargetDisks with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.Iqn)
	response := &internal.GetTargetDisksResponse{}
	disks, err := s.hostAPI.GetTargetDisks(s.requestTPtoAPITP(request.TargetPortal), request.Iqn)
	if err != nil {
		klog.Errorf("failed GetTargetDisks %v", err)
		return response, err
	}

	result := make([]string, 0, len(disks))
	for _, d := range disks {
		result = append(result, d)
	}

	response.DiskIDs = result

	return response, nil
}

func (s *Server) ListTargetPortals(context context.Context, request *internal.ListTargetPortalsRequest, version apiversion.Version) (*internal.ListTargetPortalsResponse, error) {
	klog.V(4).Infof("calling ListTargetPortals")
	response := &internal.ListTargetPortalsResponse{}
	portals, err := s.hostAPI.ListTargetPortals()
	if err != nil {
		klog.Errorf("failed ListTargetPortals %v", err)
		return response, err
	}

	result := make([]*internal.TargetPortal, 0, len(portals))
	for _, p := range portals {
		result = append(result, &internal.TargetPortal{
			TargetAddress: p.Address,
			TargetPort:    p.Port,
		})
	}

	response.TargetPortals = result

	return response, nil
}

func (s *Server) RemoveTargetPortal(context context.Context, request *internal.RemoveTargetPortalRequest, version apiversion.Version) (*internal.RemoveTargetPortalResponse, error) {
	klog.V(4).Infof("calling RemoveTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &internal.RemoveTargetPortalResponse{}
	err := s.hostAPI.RemoveTargetPortal(s.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed RemoveTargetPortal %v", err)
		return response, err
	}

	return response, nil
}

func (s *Server) SetMutualChapSecret(context context.Context, request *internal.SetMutualChapSecretRequest, version apiversion.Version) (*internal.SetMutualChapSecretResponse, error) {
	klog.V(4).Info("calling SetMutualChapSecret")

	minimumVersion := apiversion.NewVersionOrPanic("v1alpha2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("SetMutualChapSecret requires CSI-Proxy API version v1alpha2 or greater")
	}

	response := &internal.SetMutualChapSecretResponse{}
	err := s.hostAPI.SetMutualChapSecret(request.MutualChapSecret)
	if err != nil {
		klog.Errorf("failed SetMutualChapSecret %v", err)
		return response, err
	}

	return response, nil
}
