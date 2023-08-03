package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1alpha1"
	diskv1alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1alpha1"
	"github.com/stretchr/testify/require"
)

func v1alpha1DiskTests(t *testing.T) {
	t.Run("ListDiskLocations", func(t *testing.T) {
		// fails in Github Actions with
		// Error:     disk_v1alpha1_test.go:25: listDiskLocationsResponse=
		// Error:     disk_v1alpha1_test.go:27: Expected to get at least one diskLocation, instead got DiskLocations=map[]
		skipTestOnCondition(t, isRunningOnGhActions())

		client, err := diskv1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		listDiskLocationsRequest := &v1alpha1.ListDiskLocationsRequest{}
		listDiskLocationsResponse, err := client.ListDiskLocations(context.TODO(), listDiskLocationsRequest)
		require.Nil(t, err)
		t.Logf("listDiskLocationsResponse=%v", listDiskLocationsResponse)
		if len(listDiskLocationsResponse.DiskLocations) == 0 {
			t.Errorf("Expected to get at least one diskLocation, instead got DiskLocations=%+v", listDiskLocationsResponse.DiskLocations)
		}
	})

	t.Run("Rescan", func(t *testing.T) {
		client, err := diskv1alpha1client.NewClient()
		require.NoError(t, err)

		defer client.Close()

		// Rescan
		_, err = client.Rescan(context.TODO(), &v1alpha1.RescanRequest{})
		require.NoError(t, err)
	})

	t.Run("PartitionDisk", func(t *testing.T) {
		var err error
		client, err := diskv1alpha1client.NewClient()
		require.NoError(t, err)
		defer client.Close()

		// initialize disk but don't partition it using `diskInit`
		s1 := rand.NewSource(time.Now().UTC().UnixNano())
		r1 := rand.New(s1)

		testPluginPath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\", r1.Intn(100))
		mountPath := fmt.Sprintf("%smount-%d", testPluginPath, r1.Intn(100))
		vhdxPath := fmt.Sprintf("%sdisk-%d.vhdx", testPluginPath, r1.Intn(100))

		var cmd, out string
		const initialSize = 1 * 1024 * 1024 * 1024
		const partitionStyle = "GPT"

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

		var diskNumUnparsed string
		cmd = fmt.Sprintf("(Get-VHD -Path %s).DiskNumber", vhdxPath)
		if diskNumUnparsed, err = runPowershellCmd(t, cmd); err != nil {
			t.Fatalf("Error: %v. Command: %s", err, cmd)
		}

		// make disk partition request
		diskPartitionRequest := &v1alpha1.PartitionDiskRequest{
			DiskID: strings.TrimSpace(diskNumUnparsed),
		}
		_, err = client.PartitionDisk(context.TODO(), diskPartitionRequest)
		require.NoError(t, err)
	})
}
