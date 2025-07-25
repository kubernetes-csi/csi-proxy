package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
	"k8s.io/klog/v2"
)

// Implements the iSCSI OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// pkg/iscsi/iscsi.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type HostAPI interface {
	AddTargetPortal(portal *TargetPortal) error
	DiscoverTargetPortal(portal *TargetPortal) ([]string, error)
	ListTargetPortals() ([]TargetPortal, error)
	RemoveTargetPortal(portal *TargetPortal) error
	ConnectTarget(portal *TargetPortal, iqn string, authType string,
		chapUser string, chapSecret string) error
	DisconnectTarget(portal *TargetPortal, iqn string) error
	GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error)
	SetMutualChapSecret(mutualChapSecret string) error
}

type iscsiAPI struct{}

// check that iscsiAPI implements HostAPI
var _ HostAPI = &iscsiAPI{}

func New() HostAPI {
	return iscsiAPI{}
}

func (iscsiAPI) AddTargetPortal(portal *TargetPortal) error {
	return cim.WithCOMThread(func() error {
		existing, err := cim.QueryISCSITargetPortal(portal.Address, portal.Port, nil)
		if cim.IgnoreNotFound(err) != nil {
			return fmt.Errorf("error query target portal at (%s:%d). err: %v", portal.Address, portal.Port, err)
		}

		if existing != nil {
			klog.V(2).Infof("target portal at (%s:%d) already exists", portal.Address, portal.Port)
			return nil
		}

		_, err = cim.NewISCSITargetPortal(portal.Address, portal.Port, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error adding target portal at (%s:%d). err: %v", portal.Address, portal.Port, err)
		}

		return nil
	})
}

func (iscsiAPI) DiscoverTargetPortal(portal *TargetPortal) ([]string, error) {
	var iqns []string
	err := cim.WithCOMThread(func() error {
		targets, err := cim.ListISCSITargetsByTargetPortalAddressAndPort(portal.Address, portal.Port, nil)
		if err != nil {
			return fmt.Errorf("error list targets by target portal at (%s:%d). err: %v", portal.Address, portal.Port, err)
		}

		for _, target := range targets {
			iqn, err := cim.GetISCSITargetNodeAddress(target)
			if err != nil {
				return fmt.Errorf("failed parsing node address of target %v to target portal at (%s:%d). err: %w", target, portal.Address, portal.Port, err)
			}

			iqns = append(iqns, iqn)
		}

		return nil
	})
	return iqns, err
}

func (iscsiAPI) ListTargetPortals() ([]TargetPortal, error) {
	var portals []TargetPortal
	err := cim.WithCOMThread(func() error {
		instances, err := cim.ListISCSITargetPortals(cim.ISCSITargetPortalDefaultSelectorList)
		if err != nil {
			return fmt.Errorf("error list target portals. err: %v", err)
		}

		for _, instance := range instances {
			address, port, err := cim.ParseISCSITargetPortal(instance)
			if err != nil {
				return fmt.Errorf("failed parsing target portal %v. err: %w", instance, err)
			}

			portals = append(portals, TargetPortal{
				Address: address,
				Port:    port,
			})
		}

		return nil
	})
	return portals, err
}

func (iscsiAPI) RemoveTargetPortal(portal *TargetPortal) error {
	return cim.WithCOMThread(func() error {
		instance, err := cim.QueryISCSITargetPortal(portal.Address, portal.Port, nil)
		if err != nil {
			return fmt.Errorf("error query target portal at (%s:%d). err: %v", portal.Address, portal.Port, err)
		}

		result, err := cim.RemoveISCSITargetPortal(instance)
		if result != 0 || err != nil {
			return fmt.Errorf("error removing target portal at (%s:%d). result: %d, err: %w", portal.Address, portal.Port, result, err)
		}

		return nil
	})
}

func (iscsiAPI) ConnectTarget(portal *TargetPortal, iqn string, authType string, chapUser string, chapSecret string) error {
	return cim.WithCOMThread(func() error {
		target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
		if err != nil {
			return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		connected, err := cim.IsISCSITargetConnected(target)
		if err != nil {
			return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		if connected {
			klog.V(2).Infof("target %s from target portal at (%s:%d) is connected.", iqn, portal.Address, portal.Port)
			return nil
		}

		targetAuthType := strings.ToUpper(strings.ReplaceAll(authType, "_", ""))

		result, err := cim.ConnectISCSITarget(portal.Address, portal.Port, iqn, targetAuthType, &chapUser, &chapSecret)
		if err != nil {
			return fmt.Errorf("error connecting to target portal. result: %d, err: %w", result, err)
		}

		return nil
	})
}

func (iscsiAPI) DisconnectTarget(portal *TargetPortal, iqn string) error {
	return cim.WithCOMThread(func() error {
		target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
		if err != nil {
			return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		connected, err := cim.IsISCSITargetConnected(target)
		if err != nil {
			return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		if !connected {
			klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
			return nil
		}

		// get session
		session, err := cim.QueryISCSISessionByTarget(target)
		if err != nil {
			return fmt.Errorf("error query session of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		sessionIdentifier, err := cim.GetISCSISessionIdentifier(session)
		if err != nil {
			return fmt.Errorf("error query session identifier of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		persistent, err := cim.IsISCSISessionPersistent(session)
		if err != nil {
			return fmt.Errorf("error query session persistency of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		if persistent {
			result, err := cim.UnregisterISCSISession(session)
			if err != nil {
				return fmt.Errorf("error unregister session on target %s from target portal at (%s:%d). result: %d, err: %w", iqn, portal.Address, portal.Port, result, err)
			}
		}

		result, err := cim.DisconnectISCSITarget(target, sessionIdentifier)
		if err != nil {
			return fmt.Errorf("error disconnecting target %s from target portal at (%s:%d). result: %d, err: %w", iqn, portal.Address, portal.Port, result, err)
		}

		return nil
	})
}

func (iscsiAPI) GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error) {
	var ids []string
	err := cim.WithCOMThread(func() error {
		target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
		if err != nil {
			return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		connected, err := cim.IsISCSITargetConnected(target)
		if err != nil {
			return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		if !connected {
			klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
			return nil
		}

		disks, err := cim.ListDisksByTarget(target)
		if err != nil {
			return fmt.Errorf("error getting target disks on target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
		}

		for _, disk := range disks {
			number, err := cim.GetDiskNumber(disk)
			if err != nil {
				return fmt.Errorf("error getting number of disk %v on target %s from target portal at (%s:%d). err: %w", disk, iqn, portal.Address, portal.Port, err)
			}

			ids = append(ids, strconv.Itoa(int(number)))
		}

		return nil
	})
	return ids, err
}

func (iscsiAPI) SetMutualChapSecret(mutualChapSecret string) error {
	return cim.WithCOMThread(func() error {
		result, err := cim.SetISCSISessionChapSecret(mutualChapSecret)
		if err != nil {
			return fmt.Errorf("error setting mutual chap secret. result: %d, err: %v", result, err)
		}

		return nil
	})
}
