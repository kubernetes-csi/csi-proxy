package api

import (
	"fmt"
	"strconv"
	"strings"

	wmi "github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
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
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			existing, err := wmi.QueryISCSITargetPortal(scope, portal.Address, uint16(portal.Port), nil)
			if wmi.IgnoreNotFound(err) != nil {
				return fmt.Errorf("error query target portal at (%s:%d). err: %w", portal.Address, portal.Port, err)
			}

			if existing != nil {
				klog.V(2).Infof("target portal at (%s:%d) already exists", portal.Address, portal.Port)
				return nil
			}

			err = wmi.NewISCSITargetPortal(portal.Address, uint16(portal.Port), nil, nil, nil, nil)
			if err != nil {
				return fmt.Errorf("error adding target portal at (%s:%d). err: %w", portal.Address, portal.Port, err)
			}

			return nil
		})
	})
}

func (iscsiAPI) DiscoverTargetPortal(portal *TargetPortal) ([]string, error) {
	var iqns []string
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			targets, err := wmi.ListISCSITargetsByTargetPortalAddressAndPort(scope, portal.Address, uint16(portal.Port), nil)
			if err != nil {
				return fmt.Errorf("error list targets by target portal at (%s:%d). err: %w", portal.Address, portal.Port, err)
			}

			err = wmi.ForEach(targets, func(target *wmi.COMDispatchObject) error {
				iqn, err := wmi.GetISCSITargetNodeAddress(target)
				if err != nil {
					return fmt.Errorf("failed parsing node address of target %v to target portal at (%s:%d). err: %w", target, portal.Address, portal.Port, err)
				}

				iqns = append(iqns, iqn)
				return nil
			})
			return err
		})
	})
	return iqns, err
}

func (iscsiAPI) ListTargetPortals() ([]TargetPortal, error) {
	var portals []TargetPortal
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			instances, err := wmi.ListISCSITargetPortals(scope, wmi.ISCSITargetPortalDefaultSelectorList)
			if err != nil {
				return fmt.Errorf("error list target portals. err: %w", err)
			}

			err = wmi.ForEach(instances, func(instance *wmi.COMDispatchObject) error {
				address, port, err := wmi.ParseISCSITargetPortal(instance)
				if err != nil {
					return fmt.Errorf("failed parsing target portal %v. err: %w", instance, err)
				}

				portals = append(portals, TargetPortal{
					Address: address,
					Port:    uint32(port),
				})
				return nil
			})
			return err
		})
	})
	return portals, err
}

func (iscsiAPI) RemoveTargetPortal(portal *TargetPortal) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			instance, err := wmi.QueryISCSITargetPortal(scope, portal.Address, uint16(portal.Port), nil)
			if err != nil {
				return fmt.Errorf("error query target portal at (%s:%d). err: %w", portal.Address, portal.Port, err)
			}

			err = wmi.RemoveISCSITargetPortal(instance)
			if err != nil {
				return fmt.Errorf("error removing target portal at (%s:%d). err: %w", portal.Address, portal.Port, err)
			}

			return nil
		})
	})
}

func (iscsiAPI) ConnectTarget(portal *TargetPortal, iqn string, authType string, chapUser string, chapSecret string) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			target, err := wmi.QueryISCSITarget(scope, portal.Address, uint16(portal.Port), iqn)
			if err != nil {
				return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			connected, err := wmi.IsISCSITargetConnected(target)
			if err != nil {
				return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			if connected {
				klog.V(2).Infof("target %s from target portal at (%s:%d) is connected.", iqn, portal.Address, portal.Port)
				return nil
			}

			targetAuthType := strings.ToUpper(strings.ReplaceAll(authType, "_", ""))

			err = wmi.ConnectISCSITarget(portal.Address, uint16(portal.Port), iqn, targetAuthType, &chapUser, &chapSecret)
			if err != nil {
				return fmt.Errorf("error connecting to target portal. err: %w", err)
			}

			return nil
		})
	})
}

func (iscsiAPI) DisconnectTarget(portal *TargetPortal, iqn string) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			target, err := wmi.QueryISCSITarget(scope, portal.Address, uint16(portal.Port), iqn)
			if err != nil {
				return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			connected, err := wmi.IsISCSITargetConnected(target)
			if err != nil {
				return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}
			if !connected {
				klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
				return nil
			}

			session, err := wmi.QueryISCSISessionByTarget(scope, target)
			if err != nil {
				return fmt.Errorf("error query session of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			sessionIdentifier, err := wmi.GetISCSISessionIdentifier(session)
			if err != nil {
				return fmt.Errorf("error query session identifier of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			persistent, err := wmi.IsISCSISessionPersistent(session)
			if err != nil {
				return fmt.Errorf("error query session persistency of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			if persistent {
				if err = wmi.UnregisterISCSISession(session); err != nil {
					return fmt.Errorf("error unregister session on target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
				}
			}

			err = wmi.DisconnectISCSITarget(target, sessionIdentifier)
			if err != nil {
				return fmt.Errorf("error disconnecting target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			return nil
		})
	})
}

func (iscsiAPI) GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error) {
	var ids []string
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			target, err := wmi.QueryISCSITarget(scope, portal.Address, uint16(portal.Port), iqn)
			if err != nil {
				return fmt.Errorf("error query target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			connected, err := wmi.IsISCSITargetConnected(target)
			if err != nil {
				return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			if !connected {
				klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
				return nil
			}

			disks, err := wmi.ListDisksByTarget(scope, target)
			if err != nil {
				return fmt.Errorf("error getting target disks on target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
			}

			err = wmi.ForEach(disks, func(disk *wmi.COMDispatchObject) error {
				number, err := wmi.GetDiskNumber(disk)
				if err != nil {
					return fmt.Errorf("error getting number of disk %v on target %s from target portal at (%s:%d). err: %w", disk, iqn, portal.Address, portal.Port, err)
				}

				ids = append(ids, strconv.FormatUint(uint64(number), 10))
				return nil
			})
			return err
		})
	})
	return ids, err
}

func (iscsiAPI) SetMutualChapSecret(mutualChapSecret string) error {
	return wmi.WithCOMThread(func() error {
		err := wmi.SetISCSISessionChapSecret(mutualChapSecret)
		if err != nil {
			return fmt.Errorf("error setting mutual chap secret. err: %w", err)
		}

		return nil
	})
}
