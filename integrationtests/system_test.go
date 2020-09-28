package integrationtests

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/system/v1alpha1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/system/v1alpha1"
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

		result, err := exec.Command("wmic", "bios", "get", "serialnumber").Output()
		require.Nil(t, err)

		t.Logf("The serial number is %s", response.SerialNumber)

		resultString := string(result)
		require.True(t, strings.Contains(resultString, response.SerialNumber))
	})
}
