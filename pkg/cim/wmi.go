//go:build windows
// +build windows

package cim

import (
	"fmt"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/errors"
	cim "github.com/microsoft/wmi/pkg/wmiinstance"
	"k8s.io/klog/v2"
)

const (
	WMINamespaceRoot    = "Root\\CimV2"
	WMINamespaceStorage = "Root\\Microsoft\\Windows\\Storage"
	WMINamespaceSmb     = "Root\\Microsoft\\Windows\\Smb"
)

type InstanceHandler func(instance *cim.WmiInstance) (bool, error)

// An InstanceIndexer provides index key to a WMI Instance in a map
type InstanceIndexer func(instance *cim.WmiInstance) (string, error)

// NewWMISession creates a new local WMI session for the given namespace, defaulting
// to root namespace if none specified.
func NewWMISession(namespace string) (*cim.WmiSession, error) {
	if namespace == "" {
		namespace = WMINamespaceRoot
	}

	sessionManager := cim.NewWmiSessionManager()
	defer sessionManager.Dispose()

	session, err := sessionManager.GetLocalSession(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get local WMI session for namespace %s. error: %w", namespace, err)
	}

	connected, err := session.Connect()
	if !connected || err != nil {
		return nil, fmt.Errorf("failed to connect to WMI. error: %w", err)
	}

	return session, nil
}

// QueryFromWMI executes a WMI query in the specified namespace and processes each result
// through the provided handler function. Stops processing if handler returns false or encounters an error.
func QueryFromWMI(namespace string, query *query.WmiQuery, handler InstanceHandler) error {
	session, err := NewWMISession(namespace)
	if err != nil {
		return err
	}

	defer session.Close()

	instances, err := session.QueryInstances(query.String())
	if err != nil {
		return fmt.Errorf("failed to query WMI class %s. error: %w", query.ClassName, err)
	}

	if len(instances) == 0 {
		return errors.NotFound
	}

	var cont bool
	for _, instance := range instances {
		cont, err = handler(instance)
		if err != nil {
			err = fmt.Errorf("failed to query WMI class %s instance (%s). error: %w", query.ClassName, instance.String(), err)
		}
		if !cont {
			break
		}
	}

	return err
}

// QueryInstances retrieves all WMI instances matching the given query in the specified namespace.
func QueryInstances(namespace string, query *query.WmiQuery) ([]*cim.WmiInstance, error) {
	var instances []*cim.WmiInstance
	err := QueryFromWMI(namespace, query, func(instance *cim.WmiInstance) (bool, error) {
		instances = append(instances, instance)
		return true, nil
	})
	return instances, err
}

// TODO: fix the panic in microsoft/wmi library and remove this workaround
// Refer to https://github.com/microsoft/wmi/issues/167
func executeClassMethodParam(classInst *cim.WmiInstance, method *cim.WmiMethod, inParam, outParam cim.WmiMethodParamCollection) (result *cim.WmiMethodResult, err error) {
	klog.V(6).Infof("[WMI] - Executing Method [%s]\n", method.Name)

	iDispatchInstance := classInst.GetIDispatch()
	if iDispatchInstance == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "InvalidInstance")
	}
	rawResult, err := iDispatchInstance.GetProperty("Methods_")
	if err != nil {
		return nil, err
	}
	defer rawResult.Clear()
	// Retrieve the method
	rawMethod, err := rawResult.ToIDispatch().CallMethod("Item", method.Name)
	if err != nil {
		return nil, err
	}
	defer rawMethod.Clear()

	addInParam := func(inparamVariant *ole.VARIANT, paramName string, paramValue interface{}) error {
		rawProperties, err := inparamVariant.ToIDispatch().GetProperty("Properties_")
		if err != nil {
			return err
		}
		defer rawProperties.Clear()
		rawProperty, err := rawProperties.ToIDispatch().CallMethod("Item", paramName)
		if err != nil {
			return err
		}
		defer rawProperty.Clear()

		p, err := rawProperty.ToIDispatch().PutProperty("Value", paramValue)
		if err != nil {
			return err
		}
		defer p.Clear()
		return nil
	}

	params := []interface{}{method.Name}
	if len(inParam) > 0 {
		inparamsRaw, err := rawMethod.ToIDispatch().GetProperty("InParameters")
		if err != nil {
			return nil, err
		}
		defer inparamsRaw.Clear()

		inparams, err := oleutil.CallMethod(inparamsRaw.ToIDispatch(), "SpawnInstance_")
		if err != nil {
			return nil, err
		}
		defer inparams.Clear()

		for _, inp := range inParam {
			addInParam(inparams, inp.Name, inp.Value)
		}

		params = append(params, inparams)
	}

	result = &cim.WmiMethodResult{
		OutMethodParams: map[string]*cim.WmiMethodParam{},
	}
	outparams, err := classInst.GetIDispatch().CallMethod("ExecMethod_", params...)
	if err != nil {
		return
	}
	defer outparams.Clear()
	returnRaw, err := outparams.ToIDispatch().GetProperty("ReturnValue")
	if err != nil {
		return
	}
	defer returnRaw.Clear()
	if returnRaw.Value() != nil {
		result.ReturnValue = returnRaw.Value().(int32)
		klog.V(6).Infof("[WMI] - Return [%d] ", result.ReturnValue)
	}

	for _, outp := range outParam {
		returnRawIn, err1 := outparams.ToIDispatch().GetProperty(outp.Name)
		if err1 != nil {
			err = err1
			return
		}
		defer returnRawIn.Clear()

		value, err1 := cim.GetVariantValue(returnRawIn)
		if err1 != nil {
			err = err1
			return
		}

		result.OutMethodParams[outp.Name] = cim.NewWmiMethodParam(outp.Name, value)
	}
	return
}

// InvokeCimMethod calls a static method on a specific WMI class with given input parameters,
// returning the method's return value, output parameters, and any error encountered.
func InvokeCimMethod(namespace, class, methodName string, inputParameters map[string]interface{}) (int, map[string]interface{}, error) {
	session, err := NewWMISession(namespace)
	if err != nil {
		return -1, nil, err
	}

	defer session.Close()

	rawResult, err := session.Session.CallMethod("Get", class)
	if err != nil {
		return -1, nil, err
	}

	classInst, err := cim.CreateWmiInstance(rawResult, session)
	if err != nil {
		return -1, nil, err
	}

	method, err := cim.NewWmiMethod(methodName, classInst)
	if err != nil {
		return -1, nil, err
	}

	var inParam cim.WmiMethodParamCollection
	for k, v := range inputParameters {
		inParam = append(inParam, &cim.WmiMethodParam{
			Name:  k,
			Value: v,
		})
	}

	var outParam cim.WmiMethodParamCollection
	var result *cim.WmiMethodResult
	result, err = executeClassMethodParam(classInst, method, inParam, outParam)
	if err != nil {
		return -1, nil, err
	}

	outputParameters := make(map[string]interface{})
	for _, v := range result.OutMethodParams {
		outputParameters[v.Name] = v.Value
	}

	return int(result.ReturnValue), outputParameters, nil
}

// IgnoreNotFound returns nil if the error is nil or a "not found" error,
// otherwise returns the original error.
func IgnoreNotFound(err error) error {
	if err == nil || errors.IsNotFound(err) {
		return nil
	}
	return err
}

// parseObjectRef extracts the object ID from a WMI object reference string.
// The result string is in this format
// {1}\\WIN-8E2EVAQ9QSB\ROOT/Microsoft/Windows/Storage/Providers_v2\WSP_Partition.ObjectId="{b65bb3cd-da86-11ee-854b-806e6f6e6963}:PR:{00000000-0000-0000-0000-100000000000}\\?\scsi#disk&ven_vmware&prod_virtual_disk#4&2c28f6c4&0&000000#{53f56307-b6bf-11d0-94f2-00a0c91efb8b}"
// from an escape string
func parseObjectRef(input, objectClass, refName string) (string, error) {
	tokens := strings.Split(input, fmt.Sprintf("%s.%s=", objectClass, refName))
	if len(tokens) < 2 {
		return "", fmt.Errorf("invalid object ID value: %s", input)
	}

	objectID := tokens[1]
	objectID = strings.ReplaceAll(objectID, "\\\"", "\"")
	objectID = strings.ReplaceAll(objectID, "\\\\", "\\")
	objectID = objectID[1 : len(objectID)-1]
	return objectID, nil
}

// ListWMIInstanceMappings queries WMI instances and creates a map using custom indexing functions
// to extract keys and values from each instance.
func ListWMIInstanceMappings(namespace, mappingClassName string, selectorList []string, keyIndexer InstanceIndexer, valueIndexer InstanceIndexer) (map[string]string, error) {
	q := query.NewWmiQueryWithSelectList(mappingClassName, selectorList)
	mappingInstances, err := QueryInstances(namespace, q)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, mapping := range mappingInstances {
		key, err := keyIndexer(mapping)
		if err != nil {
			return nil, err
		}

		value, err := valueIndexer(mapping)
		if err != nil {
			return nil, err
		}

		result[key] = value
	}

	return result, nil
}

// FindInstancesByMapping filters instances based on a mapping relationship,
// matching instances through custom indexing and mapping functions.
func FindInstancesByMapping(instanceToFind []*cim.WmiInstance, instanceToFindIndex InstanceIndexer, associatedInstances []*cim.WmiInstance, associatedInstanceIndexer InstanceIndexer, instanceMappings map[string]string) ([]*cim.WmiInstance, error) {
	associatedInstanceObjectIDMapping := map[string]*cim.WmiInstance{}
	for _, inst := range associatedInstances {
		key, err := associatedInstanceIndexer(inst)
		if err != nil {
			return nil, err
		}

		associatedInstanceObjectIDMapping[key] = inst
	}

	var filtered []*cim.WmiInstance
	for _, inst := range instanceToFind {
		key, err := instanceToFindIndex(inst)
		if err != nil {
			return nil, err
		}

		valueObjectID, ok := instanceMappings[key]
		if !ok {
			continue
		}

		_, ok = associatedInstanceObjectIDMapping[strings.ToUpper(valueObjectID)]
		if !ok {
			continue
		}
		filtered = append(filtered, inst)
	}

	if len(filtered) == 0 {
		return nil, errors.NotFound
	}

	return filtered, nil
}

// mappingObjectRefIndexer indexes an WMI object by the Object ID reference from a specified property.
func mappingObjectRefIndexer(propertyName, className, refName string) InstanceIndexer {
	return func(instance *cim.WmiInstance) (string, error) {
		valueVal, err := instance.GetProperty(propertyName)
		if err != nil {
			return "", err
		}

		refValue, err := parseObjectRef(valueVal.(string), className, refName)
		return strings.ToUpper(refValue), err
	}
}

// stringPropertyIndexer indexes a WMI object from a string property.
func stringPropertyIndexer(propertyName string) InstanceIndexer {
	return func(instance *cim.WmiInstance) (string, error) {
		valueVal, err := instance.GetProperty(propertyName)
		if err != nil {
			return "", err
		}

		return strings.ToUpper(valueVal.(string)), err
	}
}

var (
	// objectIDPropertyIndexer indexes a WMI object from its ObjectId property.
	objectIDPropertyIndexer = stringPropertyIndexer("ObjectId")
)

// FindInstancesByObjectIDMapping filters instances based on ObjectId mapping
// between two sets of WMI instances.
func FindInstancesByObjectIDMapping(instanceToFind []*cim.WmiInstance, associatedInstances []*cim.WmiInstance, instanceMappings map[string]string) ([]*cim.WmiInstance, error) {
	return FindInstancesByMapping(instanceToFind, objectIDPropertyIndexer, associatedInstances, objectIDPropertyIndexer, instanceMappings)
}
