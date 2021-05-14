package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta3"
	v1beta3client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1beta3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
// Skip on Github Actions as it is expected to fail
func TestDiskAPIGroup(t *testing.T) {
	t.Run("ListDiskIDs", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		client, err := v1beta3client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		diskNumber := 0
		id := "page83"
		listRequest := &v1beta3.ListDiskIDsRequest{}
		diskIDsResponse, err := client.ListDiskIDs(context.TODO(), listRequest)
		require.Nil(t, err)

		cmd := "hostname"
		hostname, err := runPowershellCmd(cmd)
		if err != nil {
			t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, hostname)
		}

		hostname = strings.TrimSpace(hostname)

		diskIDsMap := diskIDsResponse.GetDiskIDs()
		if diskIDsMap != nil {
			diskIDs, found := diskIDsMap[strconv.Itoa(diskNumber)]
			if !found {
				t.Errorf("Cannot find Disk %d", diskNumber)
			}

			idValue, found := diskIDs.Identifiers[id]
			if !found {
				t.Errorf("Cannot find %s ID of Disk %d", id, diskNumber)
			}

			// In GCE, the page83 ID of Disk 0 contains the hostname
			if !strings.Contains(idValue, hostname) {
				t.Errorf("%s ID of Disk %d is incorrect. Expected to contain: %s. Received: %s", id, diskNumber, hostname, idValue)
			}
		}
	})

	t.Run("Get/SetAttachState", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		client, err := v1beta3client.NewClient()
		require.NoError(t, err)

		defer client.Close()

		s1 := rand.NewSource(time.Now().UTC().UnixNano())
		r1 := rand.New(s1)

		testPluginPath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\", r1.Intn(100))
		mountPath := fmt.Sprintf("%smount-%d", testPluginPath, r1.Intn(100))
		vhdxPath := fmt.Sprintf("%sdisk-%d.vhdx", testPluginPath, r1.Intn(100))

		defer diskCleanup(t, vhdxPath, mountPath, testPluginPath)
		diskNum := diskInit(t, vhdxPath, mountPath, testPluginPath)
		diskNumAsString := strconv.FormatUint(uint64(diskNum), 10)

		out, err := runPowershellCmd(fmt.Sprintf("Get-Disk -Number %d | Set-Disk -IsOffline $true", diskNum))
		require.NoError(t, err, "failed setting disk offline, out=%v", out)

		getReq := &v1beta3.GetAttachStateRequest{DiskID: diskNumAsString}
		getResp, err := client.GetAttachState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.False(t, getResp.IsOnline, "Expected disk to be offline")
		}

		setReq := &v1beta3.SetAttachStateRequest{DiskID: diskNumAsString, IsOnline: true}
		_, err = client.SetAttachState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", diskNum))
		assert.NoError(t, err)

		result, err := strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.False(t, result, "Expected disk to be online")

		getReq = &v1beta3.GetAttachStateRequest{DiskID: diskNumAsString}
		getResp, err = client.GetAttachState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.True(t, getResp.IsOnline, "Expected disk is online")
		}

		setReq = &v1beta3.SetAttachStateRequest{DiskID: diskNumAsString, IsOnline: false}
		_, err = client.SetAttachState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", diskNum))
		assert.NoError(t, err)

		result, err = strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.True(t, result, "Expected disk to be offline")
	})
}
