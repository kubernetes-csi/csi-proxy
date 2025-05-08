//go:build windows
// +build windows

package cim

import (
	"fmt"
	"strconv"

	"github.com/microsoft/wmi/pkg/base/query"
	cim "github.com/microsoft/wmi/pkg/wmiinstance"
	"github.com/microsoft/wmi/server2019/root/microsoft/windows/storage"
)

// ListISCSITargetPortals retrieves a list of iSCSI target portals.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM MSFT_IscsiTargetPortal
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargetportal
// for the WMI class definition.
func ListISCSITargetPortals(selectorList []string) ([]*storage.MSFT_iSCSITargetPortal, error) {
	q := query.NewWmiQueryWithSelectList("MSFT_IscsiTargetPortal", selectorList)
	instances, err := QueryInstances(WMINamespaceStorage, q)
	if IgnoreNotFound(err) != nil {
		return nil, err
	}

	var targetPortals []*storage.MSFT_iSCSITargetPortal
	for _, instance := range instances {
		portal, err := storage.NewMSFT_iSCSITargetPortalEx1(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to query iSCSI target portal %v. error: %v", instance, err)
		}

		targetPortals = append(targetPortals, portal)
	}

	return targetPortals, nil
}

// QueryISCSITargetPortal retrieves information about a specific iSCSI target portal
// identified by its network address and port number.
//
// The equivalent WMI query is:
//
//	SELECT [selectors] FROM MSFT_IscsiTargetPortal
//	  WHERE TargetPortalAddress = '<address>'
//	    AND TargetPortalPortNumber = '<port>'
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargetportal
// for the WMI class definition.
func QueryISCSITargetPortal(address string, port uint32, selectorList []string) (*storage.MSFT_iSCSITargetPortal, error) {
	portalQuery := query.NewWmiQueryWithSelectList(
		"MSFT_iSCSITargetPortal", selectorList,
		"TargetPortalAddress", address,
		"TargetPortalPortNumber", strconv.Itoa(int(port)))
	instances, err := QueryInstances(WMINamespaceStorage, portalQuery)
	if err != nil {
		return nil, err
	}

	targetPortal, err := storage.NewMSFT_iSCSITargetPortalEx1(instances[0])
	if err != nil {
		return nil, fmt.Errorf("failed to query iSCSI target portal at (%s:%d). error: %v", address, port, err)
	}

	return targetPortal, nil
}

// NewISCSITargetPortal creates a new iSCSI target portal.
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargetportal-new
// for the WMI method definition.
func NewISCSITargetPortal(targetPortalAddress string,
	targetPortalPortNumber uint32,
	initiatorInstanceName *string,
	initiatorPortalAddress *string,
	isHeaderDigest *bool,
	isDataDigest *bool) (*storage.MSFT_iSCSITargetPortal, error) {
	params := map[string]interface{}{
		"TargetPortalAddress":    targetPortalAddress,
		"TargetPortalPortNumber": targetPortalPortNumber,
	}
	if initiatorInstanceName != nil {
		params["InitiatorInstanceName"] = *initiatorInstanceName
	}
	if initiatorPortalAddress != nil {
		params["InitiatorPortalAddress"] = *initiatorPortalAddress
	}
	if isHeaderDigest != nil {
		params["IsHeaderDigest"] = *isHeaderDigest
	}
	if isDataDigest != nil {
		params["IsDataDigest"] = *isDataDigest
	}
	result, _, err := InvokeCimMethod(WMINamespaceStorage, "MSFT_iSCSITargetPortal", "New", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create iSCSI target portal with %v. result: %d, error: %v", params, result, err)
	}

	return QueryISCSITargetPortal(targetPortalAddress, targetPortalPortNumber, nil)
}

var (
	// Indexes iSCSI targets by their Object ID specified in node address
	mappingISCSITargetIndexer = mappingObjectRefIndexer("iSCSITarget", "MSFT_iSCSITarget", "NodeAddress")
	// Indexes iSCSI target portals by their Object ID specified in portal address
	mappingISCSITargetPortalIndexer = mappingObjectRefIndexer("iSCSITargetPortal", "MSFT_iSCSITargetPortal", "TargetPortalAddress")
	// Indexes iSCSI connections by their Object ID specified in connection identifier
	mappingISCSIConnectionIndexer = mappingObjectRefIndexer("iSCSIConnection", "MSFT_iSCSIConnection", "ConnectionIdentifier")
	// Indexes iSCSI sessions by their Object ID specified in session identifier
	mappingISCSISessionIndexer = mappingObjectRefIndexer("iSCSISession", "MSFT_iSCSISession", "SessionIdentifier")

	// Indexes iSCSI targets by their node address
	iscsiTargetIndexer = stringPropertyIndexer("NodeAddress")
	// Indexes iSCSI targets by their target portal address
	iscsiTargetPortalIndexer = stringPropertyIndexer("TargetPortalAddress")
	// Indexes iSCSI connections by their connection identifier
	iscsiConnectionIndexer = stringPropertyIndexer("ConnectionIdentifier")
	// Indexes iSCSI sessions by their session identifier
	iscsiSessionIndexer = stringPropertyIndexer("SessionIdentifier")
)

// ListISCSITargetToISCSITargetPortalMapping builds a mapping between iSCSI target and iSCSI target portal with iSCSI target as the key.
//
// The equivalent WMI query is:
//
//		SELECT [selectors] FROM MSFT_iSCSITargetToiSCSITargetPortal
//
//	  iSCSITarget                                                                  | iSCSITargetPortal
//	  -----------                                                                  | -----------------
//	  MSFT_iSCSITarget (NodeAddress = "iqn.1991-05.com.microsoft:win-8e2evaq9q...) | MSFT_iSCSITargetPortal (TargetPortalAdd...
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargettoiscsitargetportal
// for the WMI class definition.
func ListISCSITargetToISCSITargetPortalMapping() (map[string]string, error) {
	return ListWMIInstanceMappings(WMINamespaceStorage, "MSFT_iSCSITargetToiSCSITargetPortal", nil, mappingISCSITargetIndexer, mappingISCSITargetPortalIndexer)
}

// ListISCSIConnectionToISCSITargetMapping builds a mapping between iSCSI connection and iSCSI target with iSCSI connection as the key.
//
// The equivalent WMI query is:
//
//			SELECT [selectors] FROM MSFT_iSCSITargetToiSCSIConnection
//
//	   iSCSIConnection                                                     | iSCSITarget
//	   ---------------                                                     | -----------
//	   MSFT_iSCSIConnection (ConnectionIdentifier = "ffffac0cacbff010-15") | MSFT_iSCSITarget (NodeAddress = "iqn.1991-05.com...
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargettoiscsitargetportal
// for the WMI class definition.
func ListISCSIConnectionToISCSITargetMapping() (map[string]string, error) {
	return ListWMIInstanceMappings(WMINamespaceStorage, "MSFT_iSCSITargetToiSCSIConnection", nil, mappingISCSIConnectionIndexer, mappingISCSITargetIndexer)
}

// ListISCSISessionToISCSITargetMapping builds a mapping between iSCSI session and iSCSI target with iSCSI session as the key.
//
// The equivalent WMI query is:
//
//			SELECT [selectors] FROM MSFT_iSCSITargetToiSCSISession
//
//	   iSCSISession                                                                | iSCSITarget
//	   ------------                                                                | -----------
//	   MSFT_iSCSISession (SessionIdentifier = "ffffac0cacbff010-4000013700000016") | MSFT_iSCSITarget (NodeAddress = "iqn.199...
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitargettoiscsisession
// for the WMI class definition.
func ListISCSISessionToISCSITargetMapping() (map[string]string, error) {
	return ListWMIInstanceMappings(WMINamespaceStorage, "MSFT_iSCSITargetToiSCSISession", nil, mappingISCSISessionIndexer, mappingISCSITargetIndexer)
}

// ListDiskToISCSIConnectionMapping builds a mapping between disk and iSCSI connection with disk Object ID as the key.
//
// The equivalent WMI query is:
//
//			SELECT [selectors] FROM MSFT_iSCSIConnectionToDisk
//
//	   Disk                                                               | iSCSIConnection
//	   ----                                                               | ---------------
//	   MSFT_Disk (ObjectId = "{1}\\WIN-8E2EVAQ9QSB\root/Microsoft/Win...) | MSFT_iSCSIConnection (ConnectionIdentifier = "fff...
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsiconnectiontodisk
// for the WMI class definition.
func ListDiskToISCSIConnectionMapping() (map[string]string, error) {
	return ListWMIInstanceMappings(WMINamespaceStorage, "MSFT_iSCSIConnectionToDisk", nil, mappingObjectRefIndexer("Disk", "MSFT_Disk", "ObjectId"), mappingISCSIConnectionIndexer)
}

// ListISCSITargetsByTargetPortalWithFilters retrieves all iSCSI targets from the specified iSCSI target portal and conditions by query filters.
//
// It lists all the iSCSI targets via the following WMI query
//
//	SELECT [selectors] FROM MSFT_iSCSITarget
//
// Then find all iSCSITarget objects from MSFT_iSCSITargetToiSCSITargetPortal mapping.
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitarget
// for the WMI class definition.
func ListISCSITargetsByTargetPortalWithFilters(targetSelectorList []string, portals []*storage.MSFT_iSCSITargetPortal, filters ...*query.WmiQueryFilter) ([]*storage.MSFT_iSCSITarget, error) {
	targetQuery := query.NewWmiQueryWithSelectList("MSFT_iSCSITarget", targetSelectorList)
	targetQuery.Filters = append(targetQuery.Filters, filters...)
	instances, err := QueryInstances(WMINamespaceStorage, targetQuery)
	if err != nil {
		return nil, err
	}

	var portalInstances []*cim.WmiInstance
	for _, portal := range portals {
		portalInstances = append(portalInstances, portal.WmiInstance)
	}

	targetToTargetPortalMapping, err := ListISCSITargetToISCSITargetPortalMapping()
	if err != nil {
		return nil, err
	}

	targetInstances, err := FindInstancesByMapping(instances, iscsiTargetIndexer, portalInstances, iscsiTargetPortalIndexer, targetToTargetPortalMapping)
	if err != nil {
		return nil, err
	}

	var targets []*storage.MSFT_iSCSITarget
	for _, instance := range targetInstances {
		target, err := storage.NewMSFT_iSCSITargetEx1(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to query iSCSI target %v. %v", instance, err)
		}

		targets = append(targets, target)
	}

	return targets, nil
}

// QueryISCSITarget retrieves the iSCSI target from the specified portal address, portal and node address.
func QueryISCSITarget(address string, port uint32, nodeAddress string, selectorList []string) (*storage.MSFT_iSCSITarget, error) {
	portal, err := QueryISCSITargetPortal(address, port, nil)
	if err != nil {
		return nil, err
	}

	targets, err := ListISCSITargetsByTargetPortalWithFilters(selectorList, []*storage.MSFT_iSCSITargetPortal{portal},
		query.NewWmiQueryFilter("NodeAddress", nodeAddress, query.Equals))
	if err != nil {
		return nil, err
	}

	return targets[0], nil
}

// QueryISCSISessionByTarget retrieves the iSCSI session from the specified iSCSI target.
//
// It lists all the iSCSI sessions via the following WMI query
//
//	SELECT [selectors] FROM MSFT_iSCSISession
//
// Then find all MSFT_iSCSISession objects from MSFT_iSCSITargetToiSCSISession mapping.
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsisession
// for the WMI class definition.
func QueryISCSISessionByTarget(target *storage.MSFT_iSCSITarget, selectorList []string) (*storage.MSFT_iSCSISession, error) {
	sessionQuery := query.NewWmiQueryWithSelectList("MSFT_iSCSISession", selectorList)
	sessionInstances, err := QueryInstances(WMINamespaceStorage, sessionQuery)
	if err != nil {
		return nil, err
	}

	targetToTargetSessionMapping, err := ListISCSISessionToISCSITargetMapping()
	if err != nil {
		return nil, err
	}

	filtered, err := FindInstancesByMapping(sessionInstances, iscsiSessionIndexer, []*cim.WmiInstance{target.WmiInstance}, iscsiTargetIndexer, targetToTargetSessionMapping)
	if err != nil {
		return nil, err
	}

	session, err := storage.NewMSFT_iSCSISessionEx1(filtered[0])
	return session, err
}

// ListDisksByTarget lists all the disks on the specified iSCSI target.
//
// It lists all the iSCSI connections via the following WMI query
//
//	SELECT [selectors] FROM MSFT_iSCSIConnection
//
// Then find all MSFT_iSCSIConnection objects from MSFT_iSCSITargetToiSCSIConnection mapping,
// locate the MSFT_Disk objects using MSFT_iSCSIConnectionToDisk mapping.
//
// Refer to https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsiconnection
// for the WMI class definition.
func ListDisksByTarget(target *storage.MSFT_iSCSITarget, selectorList []string) ([]*storage.MSFT_Disk, error) {
	// list connections to the given iSCSI target
	connectionQuery := query.NewWmiQueryWithSelectList("MSFT_iSCSIConnection", selectorList)
	connectionInstances, err := QueryInstances(WMINamespaceStorage, connectionQuery)
	if err != nil {
		return nil, err
	}

	connectionToTargetMapping, err := ListISCSIConnectionToISCSITargetMapping()
	if err != nil {
		return nil, err
	}

	connectionsToTarget, err := FindInstancesByMapping(connectionInstances, iscsiConnectionIndexer, []*cim.WmiInstance{target.WmiInstance}, iscsiTargetIndexer, connectionToTargetMapping)
	if err != nil {
		return nil, err
	}

	disks, err := ListDisks(selectorList)
	if err != nil {
		return nil, err
	}

	var diskInstances []*cim.WmiInstance
	for _, disk := range disks {
		diskInstances = append(diskInstances, disk.WmiInstance)
	}

	diskToConnectionMapping, err := ListDiskToISCSIConnectionMapping()
	if err != nil {
		return nil, err
	}

	filtered, err := FindInstancesByMapping(diskInstances, objectIDPropertyIndexer, connectionsToTarget, iscsiConnectionIndexer, diskToConnectionMapping)
	if err != nil {
		return nil, err
	}

	var filteredDisks []*storage.MSFT_Disk
	for _, instance := range filtered {
		disk, err := storage.NewMSFT_DiskEx1(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to query disk %v. error: %v", disk, err)
		}

		filteredDisks = append(filteredDisks, disk)
	}
	return filteredDisks, err
}

// ConnectISCSITarget establishes a connection to an iSCSI target with optional CHAP authentication credential.
//
// Refer https://learn.microsoft.com/en-us/previous-versions/windows/desktop/iscsidisc/msft-iscsitarget-connect
// for the WMI method definition.
func ConnectISCSITarget(portalAddress string, portalPortNumber uint32, nodeAddress string, authType string, chapUsername *string, chapSecret *string) (int, map[string]interface{}, error) {
	inParams := map[string]interface{}{
		"NodeAddress":            nodeAddress,
		"TargetPortalAddress":    portalAddress,
		"TargetPortalPortNumber": int(portalPortNumber),
		"AuthenticationType":     authType,
	}
	// InitiatorPortalAddress
	// IsDataDigest
	// IsHeaderDigest
	// ReportToPnP
	if chapUsername != nil {
		inParams["ChapUsername"] = *chapUsername
	}
	if chapSecret != nil {
		inParams["ChapSecret"] = *chapSecret
	}

	result, outParams, err := InvokeCimMethod(WMINamespaceStorage, "MSFT_iSCSITarget", "Connect", inParams)
	return result, outParams, err
}
