package iscsi

import (
	"context"
	"fmt"

	iscsiapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/iscsi/hostapi"
	"k8s.io/klog/v2"
)

const defaultISCSIPort = 3260

type ISCSI struct {
	hostAPI iscsiapi.HostAPI
}

type Interface interface {
	// AddTargetPortal registers an iSCSI target network address for later
	// discovery.
	// AddTargetPortal currently does not support selecting different NICs or
	// a different iSCSI initiator (e.g a hardware initiator). This means that
	// Windows will select the initiator NIC and instance on its own.
	AddTargetPortal(context.Context, *AddTargetPortalRequest) (*AddTargetPortalResponse, error)

	// ConnectTarget connects to an iSCSI Target
	ConnectTarget(context.Context, *ConnectTargetRequest) (*ConnectTargetResponse, error)

	// DisconnectTarget disconnects from an iSCSI Target
	DisconnectTarget(context.Context, *DisconnectTargetRequest) (*DisconnectTargetResponse, error)

	// DiscoverTargetPortal initiates discovery on an iSCSI target network address
	// and returns discovered IQNs.
	DiscoverTargetPortal(context.Context, *DiscoverTargetPortalRequest) (*DiscoverTargetPortalResponse, error)

	// GetTargetDisks returns the disk addresses that correspond to an iSCSI
	// target
	GetTargetDisks(context.Context, *GetTargetDisksRequest) (*GetTargetDisksResponse, error)

	// ListTargetPortal lists all currently registered iSCSI target network
	// addresses.
	ListTargetPortals(context.Context, *ListTargetPortalsRequest) (*ListTargetPortalsResponse, error)

	// RemoveTargetPortal removes an iSCSI target network address registration.
	RemoveTargetPortal(context.Context, *RemoveTargetPortalRequest) (*RemoveTargetPortalResponse, error)

	// SetMutualChapSecret sets the default CHAP secret that all initiators on
	// this machine (node) use to authenticate the target on mutual CHAP
	// authentication.
	// NOTE: This method affects global node state and should only be used
	//       with consideration to other CSI drivers that run concurrently.
	SetMutualChapSecret(context.Context, *SetMutualChapSecretRequest) (*SetMutualChapSecretResponse, error)
}

var _ Interface = &ISCSI{}

func New(hostAPI iscsiapi.HostAPI) (*ISCSI, error) {
	return &ISCSI{
		hostAPI: hostAPI,
	}, nil
}

func (ic *ISCSI) requestTPtoAPITP(portal *TargetPortal) *iscsiapi.TargetPortal {
	port := portal.TargetPort
	if port == 0 {
		port = defaultISCSIPort
	}
	return &iscsiapi.TargetPortal{Address: portal.TargetAddress, Port: port}
}

func (ic *ISCSI) AddTargetPortal(context context.Context, request *AddTargetPortalRequest) (*AddTargetPortalResponse, error) {
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

func (ic *ISCSI) ConnectTarget(context context.Context, req *ConnectTargetRequest) (*ConnectTargetResponse, error) {
	klog.V(4).Infof("calling ConnectTarget with portal %s:%d and iqn %s"+
		" auth=%v chapuser=%v", req.TargetPortal.TargetAddress,
		req.TargetPortal.TargetPort, req.IQN, req.AuthType, req.ChapUsername)

	response := &ConnectTargetResponse{}
	authType, err := AuthTypeToString(req.AuthType)
	if err != nil {
		klog.Errorf("Error parsing parameters: %v", err)
		return response, err
	}

	err = ic.hostAPI.ConnectTarget(ic.requestTPtoAPITP(req.TargetPortal), req.IQN,
		authType, req.ChapUsername, req.ChapSecret)
	if err != nil {
		klog.Errorf("failed ConnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (ic *ISCSI) DisconnectTarget(context context.Context, request *DisconnectTargetRequest) (*DisconnectTargetResponse, error) {
	klog.V(4).Infof("calling DisconnectTarget with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.IQN)

	response := &DisconnectTargetResponse{}
	err := ic.hostAPI.DisconnectTarget(ic.requestTPtoAPITP(request.TargetPortal), request.IQN)
	if err != nil {
		klog.Errorf("failed DisconnectTarget %v", err)
		return response, err
	}

	return response, nil
}

func (ic *ISCSI) DiscoverTargetPortal(context context.Context, request *DiscoverTargetPortalRequest) (*DiscoverTargetPortalResponse, error) {
	klog.V(4).Infof("calling DiscoverTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &DiscoverTargetPortalResponse{}
	iqns, err := ic.hostAPI.DiscoverTargetPortal(ic.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed DiscoverTargetPortal %v", err)
		return response, err
	}

	response.IQNs = iqns
	return response, nil
}

func (ic *ISCSI) GetTargetDisks(context context.Context, request *GetTargetDisksRequest) (*GetTargetDisksResponse, error) {
	klog.V(4).Infof("calling GetTargetDisks with portal %s:%d and iqn %s",
		request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort, request.IQN)
	response := &GetTargetDisksResponse{}
	disks, err := ic.hostAPI.GetTargetDisks(ic.requestTPtoAPITP(request.TargetPortal), request.IQN)
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

func (ic *ISCSI) ListTargetPortals(context context.Context, request *ListTargetPortalsRequest) (*ListTargetPortalsResponse, error) {
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

func (ic *ISCSI) RemoveTargetPortal(context context.Context, request *RemoveTargetPortalRequest) (*RemoveTargetPortalResponse, error) {
	klog.V(4).Infof("calling RemoveTargetPortal with portal %s:%d", request.TargetPortal.TargetAddress, request.TargetPortal.TargetPort)
	response := &RemoveTargetPortalResponse{}
	err := ic.hostAPI.RemoveTargetPortal(ic.requestTPtoAPITP(request.TargetPortal))
	if err != nil {
		klog.Errorf("failed RemoveTargetPortal %v", err)
		return response, err
	}

	return response, nil
}

func (ic *ISCSI) SetMutualChapSecret(context context.Context, request *SetMutualChapSecretRequest) (*SetMutualChapSecretResponse, error) {
	klog.V(4).Info("calling SetMutualChapSecret")

	response := &SetMutualChapSecretResponse{}
	err := ic.hostAPI.SetMutualChapSecret(request.MutualChapSecret)
	if err != nil {
		klog.Errorf("failed SetMutualChapSecret %v", err)
		return response, err
	}

	return response, nil
}
