package generators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// a groupDefinition represents an API group definition.
type groupDefinition struct {
	name          string
	apiBasePkg    string
	serverBasePkg string
	clientBasePkg string
	versions      []*apiVersion
	// serverCallbacks are all the callbacks the internal server needs, across all versions.
	serverCallbacks orderedCallbacks
}

// an apiVersion represents a single version of a given API group.
type apiVersion struct {
	*types.Package
	// serverCallbacks are the version's specific callbacks.
	serverCallbacks orderedCallbacks
}

type namedCallback struct {
	name     string
	callback *types.Type
}

// orderedCallbacks is an alphabetically sorted list of named callbacks.
// Sorting alphabetically allows the generation to be deterministic.
// This implementation's performance is not perfect, but good enough given
// that no group/version will ever have more than a few dozens callbacks at worst.
type orderedCallbacks []namedCallback

// getOrInsert either:
// * if no callback with the same name already exists, inserts it and returns nil
// * otherwise, simply returns the existing callback with the same name
func (oc *orderedCallbacks) getOrInsert(callback namedCallback) *namedCallback {
	existing, pos := oc.getWithPosition(callback.name)

	if existing != nil {
		// found, return the existing callback
		return existing
	}

	// insert
	*oc = append((*oc)[:pos], append([]namedCallback{callback}, (*oc)[pos:]...)...)
	return nil
}

// get returns the named callback with the given name, if present.
func (oc orderedCallbacks) get(name string) *namedCallback {
	callback, _ := oc.getWithPosition(name)
	return callback
}

// getWithPosition returns the named callback with the given name, if present, with
// its position; otherwise returns the position at which it should be inserted.
func (oc orderedCallbacks) getWithPosition(name string) (callback *namedCallback, pos int) {
	pos = sort.Search(len(oc), func(i int) bool {
		return oc[i].name >= name
	})

	if pos < len(oc) && oc[pos].name == name {
		callback = &oc[pos]
	}

	return
}

func newGroupDefinition(name, apiBasePkg string) *groupDefinition {
	return &groupDefinition{
		name:          name,
		apiBasePkg:    apiBasePkg,
		serverBasePkg: defaultServerBasePkg,
		clientBasePkg: defaultClientBasePkg,
	}
}

func (d *groupDefinition) addVersion(versionPkg *types.Package) {
	serverInterface, present := versionPkg.Types[d.serverInterfaceName()]
	if !present {
		klog.Fatalf("did not find interface %s in package %s", d.serverInterfaceName(), versionPkg.Path)
	}
	if serverInterface.Kind != types.Interface {
		klog.Fatalf("type %s in package %s should be an interface, it actually is a %s",
			d.serverInterfaceName(), versionPkg.Path, serverInterface.Kind)
	}

	version := &apiVersion{
		Package: versionPkg,
	}
	d.versions = append(d.versions, version)

	for callbackName, versionedCallback := range serverInterface.Methods {
		d.validateServerCallback(callbackName, versionedCallback, version)

		version.serverCallbacks.getOrInsert(namedCallback{
			name:     callbackName,
			callback: versionedCallback,
		})

		namedServerCallback := namedCallback{
			name:     callbackName,
			callback: replaceTypesPackage(versionedCallback, versionPkg.Path, pkgPlaceholder),
		}

		if previousCallback := d.serverCallbacks.getOrInsert(namedServerCallback); previousCallback != nil {
			if namedServerCallback.callback.String() != previousCallback.callback.String() {
				errorMsg := fmt.Sprintf("Endpoint %s in API group %s inconsistent across versions:", callbackName, d.name)
				for _, vsn := range d.versions {
					if vsnCallback := vsn.serverCallbacks.get(callbackName); vsnCallback != nil {
						errorMsg += fmt.Sprintf("\n  - in version %s: %s", vsn.Name, vsnCallback.callback)
					}
				}
				errorMsg += fmt.Sprintf("\nYields 2 different signatures for the internal server callback:\n%s\nand\n%s",
					previousCallback.callback, namedServerCallback.callback)
				klog.Fatalf(errorMsg)
			}
		}
	}
}

// validateServerCallback checks that server callbacks have the expected shape, i.e.:
// * all versioned (i.e. in the same package) parameter should be pointers
// * return values should all be pointers, except for the last one, which must be an error
// These assumptions are necessary for some of the generators in this package.
func (d *groupDefinition) validateServerCallback(callbackName string, callback *types.Type, version *apiVersion) {
	for _, param := range callback.Signature.Parameters {
		if isVersionedVariable(param, version) && param.Kind != types.Pointer {
			klog.Fatalf("Server callback %s in API %s version %s has a non-pointer versioned parameter: %v",
				callbackName, d.name, version.Name, param)
		}
	}
	for i, returnValue := range callback.Signature.Results {
		if i == len(callback.Signature.Results)-1 {
			if !isBuiltInErrorType(returnValue) {
				klog.Fatalf("The last returned value for server callback %s in API %s version %s should be an error, found %v instead",
					callbackName, d.name, version.Name, returnValue)
			}
		} else if returnValue.Kind != types.Pointer {
			klog.Fatalf("Server callback %s in API %s version %s has a non-pointer return value: %v",
				callbackName, d.name, version.Name, returnValue)
		}
	}
}

// isBuiltInErrorType returns true if type t is the built-in type "error".
func isBuiltInErrorType(t *types.Type) bool {
	return t.Kind == types.Interface && t.Name.Name == "error" && t.Name.Package == ""
}

// serverInterfaceName is the name of the server interface for this API group
// that we expect to find in each version's package.
func (d *groupDefinition) serverInterfaceName() string {
	return fmt.Sprintf("%sServer", strcase.ToCamel(d.name))
}

// serverPkg returns the path of the server package, e.g.
// github.com/kubernetes-csi/csi-proxy/internal/server/<api_group_name>
func (d *groupDefinition) serverPkg() string {
	return fmt.Sprintf("%s/%s", d.serverBasePkg, d.name)
}

// internalServerPkg returns the path of the internal server package, e.g.
// github.com/kubernetes-csi/csi-proxy/internal/server/<api_group_name>/internal
func (d *groupDefinition) internalServerPkg() string {
	return fmt.Sprintf("%s/%s/internal", d.serverBasePkg, d.name)
}

// versionedServerPkg returns the path of the versioned server package, e.g.
// github.com/kubernetes-csi/csi-proxy/internal/server/<api_group_name>/internal/<version>
func (d *groupDefinition) versionedServerPkg(version string) string {
	return fmt.Sprintf("%s/%s/internal/%s", d.serverBasePkg, d.name, version)
}

// versionedClientPkg returns the path of the versioned client package, e.g.
// github.com/kubernetes-csi/csi-proxy/client/groups/<api_group_name>/<version>
func (d *groupDefinition) versionedClientPkg(version string) string {
	return fmt.Sprintf("%s/%s/%s", d.clientBasePkg, d.name, version)
}

// versionedAPIPkg returns the path to the versioned API package, e.g.
// github.com/kubernetes-csi/csi-proxy/client/api/<api_group_name>/<version>
func (d *groupDefinition) versionedAPIPkg(version string) string {
	return fmt.Sprintf("%s/%s", d.apiBasePkg, version)
}

// handy for logging/debugging
func (d *groupDefinition) String() string {
	if d == nil {
		return "<nil>"
	}

	result := fmt.Sprintf("{name: %q", d.name)
	if d.serverBasePkg != "" && d.serverBasePkg != defaultServerBasePkg {
		result += fmt.Sprintf(", serverBasePkg: %q", d.serverBasePkg)
	}
	if d.clientBasePkg != "" && d.clientBasePkg != defaultClientBasePkg {
		result += fmt.Sprintf(", clientBasePkg: %q", d.clientBasePkg)
	}
	if len(d.versions) != 0 {
		result += ", versions: ["
		for _, version := range d.versions {
			if version == nil {
				result += "<nil> "
			} else {
				result += version.Name + " "
			}
		}
		result = result[:len(result)-1] + "]"
	}
	return result + "}"
}

// isVersionedVariable returns true iff t belongs to the version package.
func isVersionedVariable(t *types.Type, version *apiVersion) bool {
	return strings.Contains(t.Name.Name, version.Path) ||
		strings.Contains(t.Name.Package, version.Path)
}
