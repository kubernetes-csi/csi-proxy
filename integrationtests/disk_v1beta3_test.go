package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta3"
	diskv1beta3client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1beta3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func v1beta3DiskTests(t *testing.T) {
	t.Run("ListDiskIDs", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())

		client, err := diskv1beta3client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		// initialize disk
		_, vhdCleanup := diskInit(t)
		defer vhdCleanup()

		diskNumber := 0
		id := "page83"
		listRequest := &v1beta3.ListDiskIDsRequest{}
		diskIDsResponse, err := client.ListDiskIDs(context.TODO(), listRequest)
		require.Nil(t, err)
		t.Logf("diskIDsResponse=%v", diskIDsResponse)

		cmd := "hostname"
		hostname, err := runPowershellCmd(t, cmd)
		if err != nil {
			t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, hostname)
		}

		hostname = strings.TrimSpace(hostname)
		diskIDsMap := diskIDsResponse.DiskIds
		if len(diskIDsMap) == 0 {
			t.Errorf("Expected to get diskIDs, instead got diskIDsResponse.DiskIds=%+v", diskIDsMap)
		}

		if diskIDsMap != nil {
			diskIDs, found := diskIDsMap[uint32(diskNumber)]
			if !found {
				t.Errorf("Cannot find Disk %d", diskNumber)
			}

			page83 := diskIDs.Page83
			if page83 == "" {
				t.Errorf("page83 field of diskNumber=%d should be defined, instead got diskIDs=%v", diskNumber, diskIDs)
			}

			serialNumber := diskIDs.SerialNumber
			if serialNumber == "" {
				t.Errorf("serialNumber field of diskNumber=%d should be defined, instead got diskIDs=%v", diskNumber, diskIDs)
			}

			// In GCE, the page83 ID of Disk 0 contains the hostname
			if !strings.Contains(page83, hostname) {
				t.Errorf("%s ID of Disk %d is incorrect. Expected to contain: %s. Received: %s", id, diskNumber, hostname, page83)
			}
		}
	})

	t.Run("Get/SetAttachState", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())

		client, err := diskv1beta3client.NewClient()
		require.NoError(t, err)

		defer client.Close()

		// initialize disk
		vhd, vhdCleanup := diskInit(t)
		defer vhdCleanup()

		out, err := runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Set-Disk -IsOffline $true", vhd.DiskNumber))
		require.NoError(t, err, "failed setting disk offline, out=%v", out)

		getReq := &v1beta3.GetDiskStateRequest{DiskNumber: vhd.DiskNumber}
		getResp, err := client.GetDiskState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.False(t, getResp.IsOnline, "Expected disk to be offline")
		}

		setReq := &v1beta3.SetDiskStateRequest{DiskNumber: vhd.DiskNumber, IsOnline: true}
		_, err = client.SetDiskState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", vhd.DiskNumber))
		assert.NoError(t, err)

		result, err := strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.False(t, result, "Expected disk to be online")

		getReq = &v1beta3.GetDiskStateRequest{DiskNumber: vhd.DiskNumber}
		getResp, err = client.GetDiskState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.True(t, getResp.IsOnline, "Expected disk is online")
		}

		setReq = &v1beta3.SetDiskStateRequest{DiskNumber: vhd.DiskNumber, IsOnline: false}
		_, err = client.SetDiskState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", vhd.DiskNumber))
		assert.NoError(t, err)

		result, err = strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.True(t, result, "Expected disk to be offline")
	})
}
