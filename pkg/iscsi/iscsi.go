package iscsi

import (
	"context"
	"fmt"

	iscsiapi "github.com/kubernetes-csi/csi-proxy/pkg/iscsi/api"
	"k8s.io/klog/v2"
)

const defaultIscsiPort = 3260

type IsCSI struct {
	hostAPI iscsiapi.API
}

type Interface interface {
	AddTargetPortal(context.Context, *AddTargetPortalRequest) (*AddTargetPortalResponse, error)
	ConnectTarget(context.Context, *ConnectTargetRequest) (*ConnectTargetResponse, error)
	DisconnectTarget(context.Context, *DisconnectTargetRequest) (*DisconnectTargetResponse, error)
	DiscoverTargetPortal(context.Context, *DiscoverTargetPortalRequest) (*DiscoverTargetPortalResponse, error)
	GetTargetDisks(context.Context, *GetTargetDisksRequest) (*GetTargetDisksResponse, error)
	ListTargetPortals(context.Context, *ListTargetPortalsRequest) (*ListTargetPortalsResponse, error)
	RemoveTargetPortal(context.Context, *RemoveTargetPortalRequest) (*RemoveTargetPortalResponse, error)
	SetMutualChapSecret(context.Context, *SetMutualChapSecretRequest) (*SetMutualChapSecretResponse, error)
}

var _ Interface = &IsCSI{}

func New(hostAPI iscsiapi.API) (*IsCSI, error) {
	return &IsCSI{
		hostAPI: hostAPI,
	}, nil
}

func (ic *IsCSI) requestTPtoAPITP(portal *TargetPortal) *iscsiapi.TargetPortal {
	port := portal.TargetPort
	if port == 0 {
		port = defaultIscsiPort
	}
	return &iscsiapi.TargetPortal{Address: portal.TargetAddress, Port: port}
}

func (ic *IsCSI) AddTargetPortal(context context.Context, request *AddTargetPortalRequest) (*AddTargetPortalResponse, error) {
	klog.V(4).Infof("calling AddTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &AddTargetPortalResponse{}
	err := ic.hostAPI.AddTargetPortal(ic.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed AddTargetPortal %v", err)
		return response, err
	}

	return response, nil
}

func AuthTypeToString(authType AuthenticationType) (string, error) {
	switch authType {
	case NONE:
		return "NONE", nil
	case ONE_WAY_CHAP:
		return "ONEWAYCHAP", nil
	case MUTUAL_CHAP:
		return "MUTUALCHAP", nil
	default:
		return "", fmt.Errorf("invalid authentication type authType=%v", authType)
	}
}

func (ic *IsCSI) ConnectTarget(context context.Context, req *ConnectTargetRequest) (*ConnectTargetResponse, error) {
	klog.V(4).Infof("calling ConnectTarget with portal %s:%d and iqn %s"+
		" auth=%v chapuser=%v", req.TargetPortal.TargetAddress,
		req.TargetPortal.TargetPort, req.Iqn, req.AuthType, req.ChapUsername)

	response := &ConnectTargetResponse{}
	authType, err := AuthTypeToString(req.AuthType)
	if err != nil {
		klog.Errorf("Error parsing parameters: %v", err)
		return response, err
	}

	err = ic.hostAPI.ConnectTarget(ic.requestTPtoAPITP(req.TargetPortal), req.Iqn,
		authType, req.ChapUsername, req.ChapSecret)
	if err != nil {
		klog.Errorf("failed ConnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (ic *IsCSI) DisconnectTarget(context context.Context, request *DisconnectTargetRequest) (*DisconnectTargetResponse, error) {
	klog.V(4).Infof("calling DisconnectTarget with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.Iqn)

	response := &DisconnectTargetResponse{}
	err := ic.hostAPI.DisconnectTarget(ic.requestTPtoAPITP(request.TargetPortal), request.Iqn)
	if err != nil {
		klog.Errorf("failed DisconnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (ic *IsCSI) DiscoverTargetPortal(context context.Context, request *DiscoverTargetPortalRequest) (*DiscoverTargetPortalResponse, error) {
	klog.V(4).Infof("calling DiscoverTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &DiscoverTargetPortalResponse{}
	iqns, err := ic.hostAPI.DiscoverTargetPortal(ic.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed DiscoverTargetPortal %v", err)
		return response, err
	}

	response.Iqns = iqns
	return response, nil
}

func (ic *IsCSI) GetTargetDisks(context context.Context, request *GetTargetDisksRequest) (*GetTargetDisksResponse, error) {
	klog.V(4).Infof("calling GetTargetDisks with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.Iqn)
	response := &GetTargetDisksResponse{}
	disks, err := ic.hostAPI.GetTargetDisks(ic.requestTPtoAPITP(request.TargetPortal), request.Iqn)
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

func (ic *IsCSI) ListTargetPortals(context context.Context, request *ListTargetPortalsRequest) (*ListTargetPortalsResponse, error) {
	klog.V(4).Infof("calling ListTargetPortals")
	response := &ListTargetPortalsResponse{}
	portals, err := ic.hostAPI.ListTargetPortals()
	if err != nil {
		klog.Errorf("failed ListTargetPortals %v", err)
		return response, err
	}

	result := make([]*TargetPortal, 0, len(portals))
	for _, p := range portals {
		result = append(result, &TargetPortal{
			TargetAddress: p.Address,
			TargetPort:    p.Port,
		})
	}

	response.TargetPortals = result

	return response, nil
}

func (ic *IsCSI) RemoveTargetPortal(context context.Context, request *RemoveTargetPortalRequest) (*RemoveTargetPortalResponse, error) {
	klog.V(4).Infof("calling RemoveTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &RemoveTargetPortalResponse{}
	err := ic.hostAPI.RemoveTargetPortal(ic.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed RemoveTargetPortal %v", err)
		return response, err
	}

	return response, nil
}

func (ic *IsCSI) SetMutualChapSecret(context context.Context, request *SetMutualChapSecretRequest) (*SetMutualChapSecretResponse, error) {
	klog.V(4).Info("calling SetMutualChapSecret")

	response := &SetMutualChapSecretResponse{}
	err := ic.hostAPI.SetMutualChapSecret(request.MutualChapSecret)
	if err != nil {
		klog.Errorf("failed SetMutualChapSecret %v", err)
		return response, err
	}

	return response, nil
}
