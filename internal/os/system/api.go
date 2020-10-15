package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Implements the System OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/system/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) GetBIOSSerialNumber() (string, error) {
	// Taken from Kubernetes vSphere cloud provider
	// https://github.com/kubernetes/kubernetes/blob/103e926604de6f79161b78af3e792d0ed282bc06/staging/src/k8s.io/legacy-cloud-providers/vsphere/vsphere_util_windows.go#L28
	result, err := exec.Command("wmic", "bios", "get", "serialnumber").Output()
	if err != nil {
		return "", err
	}
	lines := strings.FieldsFunc(string(result), func(r rune) bool {
		switch r {
		case '\n', '\r':
			return true
		default:
			return false
		}
	})
	if len(lines) != 2 {
		return "", fmt.Errorf("received unexpected value retrieving host uuid: %q", string(result))
	}
	return lines[1], nil
}
