package client

import (
	"errors"
	"os"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
)

const (
	// pipePrefix is the prefix for Windows named pipes' names
	pipePrefix = `\\.\\pipe\\`

	// CsiProxyNamedPipePrefix is the prefix for the named pipes the proxy creates.
	// The suffix will be the API group and version,
	// e.g. "\\.\\pipe\\csi-proxy-iscsi-v1", "\\.\\pipe\\csi-proxy-filesystem-v2alpha1", etc.
	csiProxyNamedPipePrefix = "csi-proxy-"
)

func PipePath(apiGroupName string, apiVersion apiversion.Version) string {
	return pipePrefix + csiProxyNamedPipePrefix + apiGroupName + "-" + apiVersion.String()
}

// FindFirstNamedPipe finds the first named pipe available at the Windows named pipe locations.
// Only the first version from `apiVersions` that exists will be returned, this mechanism
// allows picking the first version available in the case where there's a beta and v1 API
func FindFirstNamedPipe(apiGroupName string, apiVersions []apiversion.Version) (string, error) {
	for _, version := range apiVersions {
		namedPipe := PipePath(apiGroupName, version)
		if _, err := os.Lstat(namedPipe); err != nil {
			return namedPipe, nil
		}
	}
	return "", errors.New("Couldn't find a valid named pipe")
}
