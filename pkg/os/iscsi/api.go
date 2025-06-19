package iscsi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
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
	targets, err := cim.ListISCSITargetsByTargetPortalAddressAndPort(portal.Address, portal.Port, nil)
	if err != nil {
		return nil, err
	}

	var iqns []string
	for _, target := range targets {
		iqn, err := cim.GetISCSITargetNodeAddress(target)
		if err != nil {
			return nil, fmt.Errorf("failed parsing node address of target %v to target portal at (%s:%d). err: %w", target, portal.Address, portal.Port, err)
		}

		iqns = append(iqns, iqn)
	}

	return iqns, nil
}

func (APIImplementor) ListTargetPortals() ([]TargetPortal, error) {
	instances, err := cim.ListISCSITargetPortals(cim.ISCSITargetPortalDefaultSelectorList)
	if err != nil {
		return nil, err
	}

	var portals []TargetPortal
	for _, instance := range instances {
		address, port, err := cim.ParseISCSITargetPortal(instance)
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

	result, err := cim.RemoveISCSITargetPortal(instance)
	if result != 0 || err != nil {
		return fmt.Errorf("error removing target portal at (%s:%d). result: %d, err: %w", portal.Address, portal.Port, result, err)
	}

	return nil
}

func (APIImplementor) ConnectTarget(portal *TargetPortal, iqn string, authType string, chapUser string, chapSecret string) error {
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
	if err != nil {
		return err
	}

	connected, err := cim.IsISCSITargetConnected(target)
	if err != nil {
		return err
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
}

func (APIImplementor) DisconnectTarget(portal *TargetPortal, iqn string) error {
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
	if err != nil {
		return err
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
		return fmt.Errorf("error query session of  target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
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
}

func (APIImplementor) GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error) {
	target, err := cim.QueryISCSITarget(portal.Address, portal.Port, iqn)
	if err != nil {
		return nil, err
	}

	connected, err := cim.IsISCSITargetConnected(target)
	if err != nil {
		return nil, fmt.Errorf("error query connected of target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	if !connected {
		klog.V(2).Infof("target %s from target portal at (%s:%d) is not connected.", iqn, portal.Address, portal.Port)
		return nil, nil
	}

	disks, err := cim.ListDisksByTarget(target)
	if err != nil {
		return nil, fmt.Errorf("error getting target disks on target %s from target portal at (%s:%d). err: %w", iqn, portal.Address, portal.Port, err)
	}

	var ids []string
	for _, disk := range disks {
		number, err := cim.GetDiskNumber(disk)
		if err != nil {
			return nil, fmt.Errorf("error getting number of disk %v on target %s from target portal at (%s:%d). err: %w", disk, iqn, portal.Address, portal.Port, err)
		}

		ids = append(ids, strconv.Itoa(int(number)))
	}
	return ids, nil
}

func (APIImplementor) SetMutualChapSecret(mutualChapSecret string) error {
	result, err := cim.SetISCSISessionChapSecret(mutualChapSecret)
	if err != nil {
		return fmt.Errorf("error setting mutual chap secret. result: %d, err: %v", result, err)
	}

	return nil
}
