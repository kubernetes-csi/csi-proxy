package integrationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/system/v1alpha1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/system/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBIOSSerialNumber(t *testing.T) {
	t.Run("GetBIOSSerialNumber", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		request := &v1alpha1.GetBIOSSerialNumberRequest{}
		response, err := client.GetBIOSSerialNumber(context.TODO(), request)
		require.Nil(t, err)
		require.NotNil(t, response)

		serialNumber, err := runPowershellCmd(t, fmt.Sprintf(`(Get-CimInstance -ClassName Win32_BIOS).SerialNumber`))
		if err != nil {
			t.Fatalf("command to get serial number failed: %v", err)
		}
		t.Logf("The serial number is %s", serialNumber)

		require.True(t, strings.Contains(serialNumber, response.SerialNumber))
	})
}

func TestServiceCommands(t *testing.T) {
	t.Run("GetService", func(t *testing.T) {
		const ServiceName = "MSiSCSI"
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		// Make sure service is stopped
		_, err = runPowershellCmd(t, fmt.Sprintf(`Stop-Service -Name "%s"`, ServiceName))
		require.NoError(t, err)
		assertServiceStopped(t, ServiceName)

		request := &v1alpha1.GetServiceRequest{Name: ServiceName}
		response, err := client.GetService(context.TODO(), request)
		require.NoError(t, err)
		require.NotNil(t, response)

		out, err := runPowershellCmd(t, fmt.Sprintf(`Get-Service -Name "%s" `+
			`| Select-Object DisplayName, Status, StartType | ConvertTo-Json`,
			ServiceName))
		require.NoError(t, err)

		serviceInfo := struct {
			DisplayName string `json:"DisplayName"`
			Status      uint32 `json:"Status"`
			StartType   uint32 `json:"StartType"`
		}{}

		err = json.Unmarshal([]byte(out), &serviceInfo)
		require.NoError(t, err, "failed unmarshalling json out=%v", out)

		assert.Equal(t, serviceInfo.Status, uint32(response.Status))
		assert.Equal(t, v1alpha1.ServiceStatus_STOPPED, response.Status)
		assert.Equal(t, serviceInfo.StartType, uint32(response.StartType))
		assert.Equal(t, serviceInfo.DisplayName, response.DisplayName)
	})

	t.Run("Stop/Start Service", func(t *testing.T) {
		const ServiceName = "MSiSCSI"
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		_, err = runPowershellCmd(t, fmt.Sprintf(`Stop-Service -Name "%s"`, ServiceName))
		require.NoError(t, err)
		assertServiceStopped(t, ServiceName)

		startReq := &v1alpha1.StartServiceRequest{Name: ServiceName}
		startResp, err := client.StartService(context.TODO(), startReq)

		assert.NoError(t, err)
		assert.NotNil(t, startResp)
		assertServiceStarted(t, ServiceName)

		stopReq := &v1alpha1.StopServiceRequest{Name: ServiceName}
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
