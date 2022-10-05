package api

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/utils"
)

// Implements the System OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// pkg/system/system.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type API interface {
	GetBIOSSerialNumber() (string, error)
	GetService(name string) (*ServiceInfo, error)
	StartService(name string) error
	StopService(name string, force bool) error
}

type systemAPI struct{}

func New() API {
	return systemAPI{}
}

func (systemAPI) GetBIOSSerialNumber() (string, error) {
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

func (systemAPI) GetService(name string) (*ServiceInfo, error) {
	script := `Get-Service -Name $env:ServiceName | Select-Object DisplayName, Status, StartType | ` +
		`ConvertTo-JSON`
	cmdEnv := fmt.Sprintf("ServiceName=%s", name)
	out, err := utils.RunPowershellCmd(script, cmdEnv)
	if err != nil {
		return nil, fmt.Errorf("error querying service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	var serviceInfo ServiceInfo
	err = json.Unmarshal(out, &serviceInfo)
	if err != nil {
		return nil, err
	}

	return &serviceInfo, nil
}

func (systemAPI) StartService(name string) error {
	script := `Start-Service -Name $env:ServiceName`
	cmdEnv := fmt.Sprintf("ServiceName=%s", name)
	out, err := utils.RunPowershellCmd(script, cmdEnv)
	if err != nil {
		return fmt.Errorf("error starting service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}

func (systemAPI) StopService(name string, force bool) error {
	script := `Stop-Service -Name $env:ServiceName -Force:$([System.Convert]::ToBoolean($env:Force))`
	out, err := utils.RunPowershellCmd(script, fmt.Sprintf("ServiceName=%s", name), fmt.Sprintf("Force=%t", force))
	if err != nil {
		return fmt.Errorf("error stopping service name=%s. cmd: %s, output: %s, error: %v", name, script, string(out), err)
	}

	return nil
}
