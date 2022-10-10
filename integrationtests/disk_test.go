package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/pkg/disk"
	diskapi "github.com/kubernetes-csi/csi-proxy/pkg/disk/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
// Skip on Github Actions as it is expected to fail
func TestDisk(t *testing.T) {
	t.Run("ListDiskIDs,ListDiskLocations", func(t *testing.T) {
		// even though this test doesn't need the VHD API it failed in Github Actions
		//     disk_v1beta3_test.go:30:
		// Error Trace:	disk_v1beta3_test.go:30
		// Error:      	Expected nil, but got: &status.statusError{state:impl.MessageState{NoUnkeyedLiterals:pragma.NoUnkeyedLiterals{}, DoNotCompare:pragma.DoNotCompare{}, DoNotCopy:pragma.DoNotCopy{}, atomicMessageInfo:(*impl.MessageInfo)(nil)}, sizeCache:0, unknownFields:[]uint8(nil), Code:2, Message:"Could not get page83 ID: IOCTL_STORAGE_QUERY_PROPERTY failed: Incorrect function.", Details:[]*anypb.Any(nil)}
		// Test:       	TestDiskAPIGroup/v1beta3Tests/ListDiskIDs,ListDiskLocations
		skipTestOnCondition(t, isRunningOnGhActions())

		client, err := disk.New(diskapi.New())
		require.Nil(t, err)

		listRequest := &disk.ListDiskIDsRequest{}
		diskIDsResponse, err := client.ListDiskIDs(context.TODO(), listRequest)
		require.Nil(t, err)

		// example output for GCE (0 is ok, others are virtual disks)
		// diskIDs:{key:0  value:{page83:"Google  persistent-disk-0"  serial_number:"                    "}}
		// diskIDs:{key:1  value:{page83:"4d53465420202020328d59b360875845ac645473be8267bf"}}
		// diskIDs:{key:2  value:{page83:"4d534654202020208956a91dadfe3d48865f9b9bcbdb8d3e"}}
		// diskIDs:{key:3  value:{page83:"4d534654202020207a3d18d72787ee47bdc127cb4f06403a"}}
		t.Logf("diskIDsResponse=%v", diskIDsResponse)

		cmd := "hostname"
		hostname, err := runPowershellCmd(t, cmd)
		if err != nil {
			t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, hostname)
		}

		diskIDsMap := diskIDsResponse.DiskIDs
		if len(diskIDsMap) == 0 {
			t.Errorf("Expected to get at least one diskIDs, instead got diskIDsResponse.DiskIDs=%+v", diskIDsMap)
		}

		// some disks may have the field Page83, if it's a GCE Persistent disk
		// it'll have a nonempty SerialNumber
		// first disk is the VM disk (other disks might be VHD)
		for diskNumber, diskIDs := range diskIDsMap {
			if len(diskIDs.SerialNumber) > 0 {
				// the nvme disks don't have a Page83 number
				if strings.HasPrefix(diskIDs.SerialNumber, "nvme") {
					continue
				}
				page83 := diskIDs.Page83
				if page83 == "" {
					t.Errorf("page83 field of diskNumber=%d should be defined, instead got diskIDs=%v", diskNumber, diskIDs)
				}
			}
		}

		listDiskLocationsRequest := &disk.ListDiskLocationsRequest{}
		listDiskLocationsResponse, err := client.ListDiskLocations(context.TODO(), listDiskLocationsRequest)
		require.Nil(t, err)
		t.Logf("listDiskLocationsResponse=%v", listDiskLocationsResponse)
		if len(listDiskLocationsResponse.DiskLocations) == 0 {
			t.Errorf("Expected to get at least one diskLocation, instead got DiskLocations=%+v", listDiskLocationsResponse.DiskLocations)
		}
	})

	t.Run("Get/SetDiskState", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())

		client, err := disk.New(diskapi.New())
		require.Nil(t, err)

		// initialize disk
		vhd, vhdCleanup := diskInit(t)
		defer vhdCleanup()

		// disk stats
		diskStatsRequest := &disk.GetDiskStatsRequest{
			DiskNumber: vhd.DiskNumber,
		}
		diskStatsResponse, err := client.GetDiskStats(context.TODO(), diskStatsRequest)
		require.NoError(t, err)
		if !sizeIsAround(t, diskStatsResponse.TotalBytes, vhd.InitialSize) {
			t.Fatalf("DiskStats doesn't have the expected size, wanted (close to)=%d got=%d", vhd.InitialSize, diskStatsResponse.TotalBytes)
		}

		// Rescan
		_, err = client.Rescan(context.TODO(), &disk.RescanRequest{})
		require.NoError(t, err)

		// change disk state
		out, err := runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Set-Disk -IsOffline $true", vhd.DiskNumber))
		require.NoError(t, err, "failed setting disk offline, out=%v", out)

		getReq := &disk.GetDiskStateRequest{DiskNumber: vhd.DiskNumber}
		getResp, err := client.GetDiskState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.False(t, getResp.IsOnline, "Expected disk to be offline")
		}

		setReq := &disk.SetDiskStateRequest{DiskNumber: vhd.DiskNumber, IsOnline: true}
		_, err = client.SetDiskState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", vhd.DiskNumber))
		assert.NoError(t, err)

		result, err := strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.False(t, result, "Expected disk to be online")

		getReq = &disk.GetDiskStateRequest{DiskNumber: vhd.DiskNumber}
		getResp, err = client.GetDiskState(context.TODO(), getReq)

		if assert.NoError(t, err) {
			assert.True(t, getResp.IsOnline, "Expected disk is online")
		}

		setReq = &disk.SetDiskStateRequest{DiskNumber: vhd.DiskNumber, IsOnline: false}
		_, err = client.SetDiskState(context.TODO(), setReq)
		assert.NoError(t, err)

		out, err = runPowershellCmd(t, fmt.Sprintf("Get-Disk -Number %d | Select-Object -ExpandProperty IsOffline", vhd.DiskNumber))
		assert.NoError(t, err)

		result, err = strconv.ParseBool(strings.TrimSpace(out))
		assert.NoError(t, err)
		assert.True(t, result, "Expected disk to be offline")
	})

	t.Run("PartitionDisk", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())

		var err error
		client, err := disk.New(diskapi.New())
		require.Nil(t, err)

		// initialize disk but don't partition it using `diskInit`
		s1 := rand.NewSource(time.Now().UTC().UnixNano())
		r1 := rand.New(s1)

		testPluginPath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\", r1.Intn(100))
		mountPath := fmt.Sprintf("%smount-%d", testPluginPath, r1.Intn(100))
		vhdxPath := fmt.Sprintf("%sdisk-%d.vhdx", testPluginPath, r1.Intn(100))

		var cmd, out string
		const initialSize = 1 * 1024 * 1024 * 1024

		cmd = fmt.Sprintf("mkdir %s", mountPath)
		if out, err = runPowershellCmd(t, cmd); err != nil {
			t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
		}
		cmd = fmt.Sprintf("New-VHD -Path %s -SizeBytes %d", vhdxPath, initialSize)
		if out, err = runPowershellCmd(t, cmd); err != nil {
			t.Fatalf("Error: %v. Command: %q. Out: %s.", err, cmd, out)
		}
		cmd = fmt.Sprintf("Mount-VHD -Path %s", vhdxPath)
		if out, err = runPowershellCmd(t, cmd); err != nil {
			t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
		}

		var diskNum uint64
		var diskNumUnparsed string
		cmd = fmt.Sprintf("(Get-VHD -Path %s).DiskNumber", vhdxPath)
		if diskNumUnparsed, err = runPowershellCmd(t, cmd); err != nil {
			t.Fatalf("Error: %v. Command: %s", err, cmd)
		}
		if diskNum, err = strconv.ParseUint(strings.TrimRight(diskNumUnparsed, "\r\n"), 10, 32); err != nil {
			t.Fatalf("Error: %v", err)
		}

		// make disk partition request
		diskPartitionRequest := &disk.PartitionDiskRequest{
			DiskNumber: uint32(diskNum),
		}
		_, err = client.PartitionDisk(context.TODO(), diskPartitionRequest)
		require.NoError(t, err)
	})
}

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
// Skip on Github Actions as it is expected to fail
func TestDiskAPIGroup(t *testing.T) {
	t.Run("v1Tests", func(t *testing.T) {
		v1DiskTests(t)
	})
	t.Run("v1beta3Tests", func(t *testing.T) {
		v1beta3DiskTests(t)
	})
	t.Run("v1beta2Tests", func(t *testing.T) {
		v1beta2DiskTests(t)
	})
	t.Run("v1beta1Tests", func(t *testing.T) {
		v1beta1DiskTests(t)
	})
	t.Run("v1alpha1Tests", func(t *testing.T) {
		v1alpha1DiskTests(t)
	})
}
