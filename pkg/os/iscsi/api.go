package iscsi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/microsoft/wmi/server2019/root/microsoft/windows/storage"
	"k8s.io/klog/v2"
)

// Implements the iSCSI OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/iscsi/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func parseTargetPortal(instance *storage.MSFT_iSCSITargetPortal) (string, uint32, error) {
	portalAddress, err := instance.GetPropertyTargetPortalAddress()
	if err != nil {
		return "", 0, fmt.Errorf("failed parsing target portal address %v. err: %w", instance, err)
	}

	portalPort, err := instance.GetProperty("TargetPortalPortNumber")
	if err != nil {
		return "", 0, fmt.Errorf("failed parsing target portal port number %v. err: %w", instance, err)
	}

	return portalAddress, uint32(portalPort.(int32)), nil
}

func (APIImplementor) AddTargetPortal(portal *TargetPortal) error {
	existing, err := cim.QueryISCSITargetPortal(portal.Address, portal.Port, nil)
	if cim.IgnoreNotFound(err) != nil {
		return err
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
}

func (APIImplementor) DiscoverTargetPortal(portal *TargetPortal) ([]string, error) {
	instance, err := cim.QueryISCSITargetPortal(portal.Address, portal.Port, nil)
	if err != nil {
		return nil, err
	}

	targets, err := cim.ListISCSITargetsByTargetPortalWithFilters(nil, []*storage.MSFT_iSCSITargetPortal{instance})
	if err != nil {
		return nil, err
	}

	var iqns []string
	for _, target := range targets {
		iqn, err := target.GetProperty("NodeAddress")
		if err != nil {
			return nil, fmt.Errorf("failed parsing node address of target %v to target portal at (%s:%d). err: %w", target, portal.Address, portal.Port, err)
		}

		iqns = append(iqns, iqn.(string))
	}

	return iqns, nil
}

func (APIImplementor) ListTargetPortals() ([]TargetPortal, error) {
	instances, err := cim.ListISCSITargetPortals([]string{"TargetPortalAddress", "TargetPortalPortNumber"})
	if err != nil {
		return nil, err
	}

	var portals []TargetPortal
	for _, instance := range instances {
		address, port, err := parseTargetPortal(instance)
		if err != nil {
			return nil, fmt.Errorf("failed parsing target portal %v. err: %w", instance, err)
		}

		portals = append(portals, TargetPortal{
			Address: address,
			Port:    port,
		})
	}

	return portals, nil
}

func (APIImplementor) RemoveTargetPortal(portal *TargetPortal) error {
	instance, err := cim.QueryISCSITargetPortal(portal.Address, portal.Port, nil)
	if err != nil {
		return err
	}

	address, port, err := parseTargetPortal(instance)
	if err != nil {
		return fmt.Errorf("failed to parse target portal %v. error: %v", instance, err)
	}

	result, err := instance.InvokeMethodWithReturn("Remove",
		nil,
		nil,
		int(port),
		address,
	)
	if result != 0 || err != nil {
		return fmt.Errorf("error removing target portal at (%s:%d). result: %d, err: %w", address, port, result, err)
	}

	return nil
}

func (APIImplementor) ConnectTarget(portal *TargetPortal, iqn string, authType string, chapUser string, chapSecret string) error {
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn, nil)
	if err != nil {
		return err
	}

	connected, err := target.GetPropertyIsConnected()
	if err != nil {
		return err
	}

	if connected {
		klog.V(2).Infof("target %s from target portal at (%s:%d) is connected.", iqn, portal.Address, portal.Port)
		return nil
	}

	targetAuthType := strings.ToUpper(strings.ReplaceAll(authType, "_", ""))

	result, _, err := cim.ConnectISCSITarget(portal.Address, portal.Port, iqn, targetAuthType, &chapUser, &chapSecret)
	if err != nil {
		return fmt.Errorf("error connecting to target portal. result: %d, err: %w", result, err)
	}

	return nil
}

func (APIImplementor) DisconnectTarget(portal *TargetPortal, iqn string) error {
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn, nil)
	if err != nil {
		return err
	}

	connected, err := target.GetPropertyIsConnected()
	if err != nil {
		return fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	if !connected {
		klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
		return nil
	}

	// get session
	session, err := cim.QueryISCSISessionByTarget(target, nil)
	if err != nil {
		return fmt.Errorf("error query session of  target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	sessionIdentifier, err := session.GetPropertySessionIdentifier()
	if err != nil {
		return fmt.Errorf("error query session identifier of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	persistent, err := session.GetPropertyIsPersistent()
	if err != nil {
		return fmt.Errorf("error query session persistency of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	if persistent {
		result, err := session.InvokeMethodWithReturn("Unregister")
		if err != nil {
			return fmt.Errorf("error unregister session on target %s from target portal at (%s:%d). result: %d, err: %w", iqn, portal.Address, portal.Port, result, err)
		}
	}

	result, err := target.InvokeMethodWithReturn("Disconnect", sessionIdentifier)
	if err != nil {
		return fmt.Errorf("error disconnecting target %s from target portal at (%s:%d). result: %d, err: %w", iqn, portal.Address, portal.Port, result, err)
	}

	return nil
}

func (APIImplementor) GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error) {
	// Converting DiskNumber to string for compatibility with disk api group
	// Not using pipeline in order to validate that items are non-empty
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn, nil)
	if err != nil {
		return nil, err
	}

	connected, err := target.GetPropertyIsConnected()
	if err != nil {
		return nil, fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	if !connected {
		klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
		return nil, nil
	}

	disks, err := cim.ListDisksByTarget(target, []string{})

	if err != nil {
		return nil, fmt.Errorf("error getting target disks on target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	var ids []string
	for _, disk := range disks {
		number, err := disk.GetProperty("Number")
		if err != nil {
			return nil, fmt.Errorf("error getting number of disk %v on target %s from target portal at (%s:%d). err: %w", disk, iqn, portal.Address, portal.Port, err)
		}

		ids = append(ids, strconv.Itoa(int(number.(int32))))
	}
	return ids, nil
}

func (APIImplementor) SetMutualChapSecret(mutualChapSecret string) error {
	result, _, err := cim.InvokeCimMethod(cim.WMINamespaceStorage, "MSFT_iSCSISession", "SetCHAPSecret", map[string]interface{}{"ChapSecret": mutualChapSecret})
	if err != nil {
		return fmt.Errorf("error setting mutual chap secret. result: %d, err: %v", result, err)
	}

	return nil
}
