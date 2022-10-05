package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	system "github.com/kubernetes-csi/csi-proxy/v2/pkg/system"
	systemapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/system/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {
	t.Run("GetBIOSSerialNumber", func(t *testing.T) {
		client, err := system.New(systemapi.New())
		require.Nil(t, err)

		request := &system.GetBIOSSerialNumberRequest{}
		response, err := client.GetBIOSSerialNumber(context.TODO(), request)
		require.Nil(t, err)
		require.NotNil(t, response)

		result, err := exec.Command("wmic", "bios", "get", "serialnumber").Output()
		require.Nil(t, err)

		t.Logf("The serial number is %s", response.SerialNumber)

		resultString := string(result)
		require.True(t, strings.Contains(resultString, response.SerialNumber))
	})

	t.Run("GetService", func(t *testing.T) {
		const ServiceName = "MSiSCSI"
		client, err := system.New(systemapi.New())
		require.Nil(t, err)

		// Make sure service is stopped
		_, err = runPowershellCmd(t, fmt.Sprintf(`Stop-Service -Name "%s"`, ServiceName))
		require.NoError(t, err)
		assertServiceStopped(t, ServiceName)

		request := &system.GetServiceRequest{Name: ServiceName}
		response, err := client.GetService(context.TODO(), request)
		require.NoError(t, err)
		require.NotNil(t, response)

		out, err := runPowershellCmd(t, fmt.Sprintf(`Get-Service -Name "%s" `+
			`| Select-Object DisplayName, Status, StartType | ConvertTo-Json`,
			ServiceName))
		require.NoError(t, err)

		var serviceInfo = struct {
			DisplayName string `json:"DisplayName"`
			Status      uint32 `json:"Status"`
			StartType   uint32 `json:"StartType"`
		}{}

		err = json.Unmarshal([]byte(out), &serviceInfo)
		require.NoError(t, err, "failed unmarshalling json out=%v", out)

		assert.Equal(t, serviceInfo.Status, uint32(response.Status))
		assert.Equal(t, system.SERVICE_STATUS_STOPPED, response.Status)
		assert.Equal(t, serviceInfo.StartType, uint32(response.StartType))
		assert.Equal(t, serviceInfo.DisplayName, response.DisplayName)
	})

	t.Run("Stop/Start Service", func(t *testing.T) {
		const ServiceName = "MSiSCSI"
		client, err := system.New(systemapi.New())
		require.Nil(t, err)

		_, err = runPowershellCmd(t, fmt.Sprintf(`Stop-Service -Name "%s"`, ServiceName))
		require.NoError(t, err)
		assertServiceStopped(t, ServiceName)

		startReq := &system.StartServiceRequest{Name: ServiceName}
		startResp, err := client.StartService(context.TODO(), startReq)

		assert.NoError(t, err)
		assert.NotNil(t, startResp)
		assertServiceStarted(t, ServiceName)

		stopReq := &system.StopServiceRequest{Name: ServiceName}
		stopResp, err := client.StopService(context.TODO(), stopReq)

		assert.NoError(t, err)
		assert.NotNil(t, stopResp)
		assertServiceStopped(t, ServiceName)
	})
}

func assertServiceStarted(t *testing.T, serviceName string) {
	assertServiceStatus(t, serviceName, "Running")
}

func assertServiceStopped(t *testing.T, serviceName string) {
	assertServiceStatus(t, serviceName, "Stopped")
}

func assertServiceStatus(t *testing.T, serviceName string, status string) {
	out, err := runPowershellCmd(t, fmt.Sprintf(`Get-Service -Name "%s" | `+
		`Select-Object -ExpandProperty Status`, serviceName))
	if !assert.NoError(t, err, "Failed getting service out=%s", out) {
		return
	}

	assert.Equal(t, strings.TrimSpace(out), status)
}
