package integrationtests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	disk "github.com/kubernetes-csi/csi-proxy/v2/pkg/disk"
	diskapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/disk/hostapi"
	volume "github.com/kubernetes-csi/csi-proxy/v2/pkg/volume"
	volumeapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/volume/hostapi"

	"github.com/stretchr/testify/require"
)

func TestVolume(t *testing.T) {
	t.Run("NegativeVolumeTests", func(t *testing.T) {
		negativeVolumeTests(t)
	})

	// these tests should be considered frozen from the API point of view
	volumeClient, err := volume.New(volumeapi.New())
	require.Nil(t, err)

	diskClient, err := disk.New(diskapi.New())
	require.Nil(t, err)

	t.Run("MountVolume", func(t *testing.T) {
		mountVolumeTests(diskClient, volumeClient, t)
	})

	t.Run("GetClosestVolumeFromTargetPath", func(t *testing.T) {
		getClosestVolumeFromTargetPathTests(diskClient, volumeClient, t)
	})
}

func runNegativeListVolumeRequest(t *testing.T, client volume.Interface, diskNum uint32) {
	listRequest := &volume.ListVolumesOnDiskRequest{
		DiskNumber: diskNum,
	}
	_, err := client.ListVolumesOnDisk(context.TODO(), listRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for disknum:%d", diskNum)
	}
}

func runNegativeIsVolumeFormattedRequest(t *testing.T, client volume.Interface, volumeID string) {
	isVolumeFormattedRequest := &volume.IsVolumeFormattedRequest{
		VolumeID: volumeID,
	}
	_, err := client.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for volumeID: %s", volumeID)
	}
}

func runNegativeFormatVolumeRequest(t *testing.T, client volume.Interface, volumeID string) {
	formatVolumeRequest := &volume.FormatVolumeRequest{
		VolumeID: volumeID,
	}
	_, err := client.FormatVolume(context.TODO(), formatVolumeRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for volume id: %s", volumeID)
	}
}

func runNegativeResizeVolumeRequest(t *testing.T, client volume.Interface, volumeID string, size int64) {
	resizeVolumeRequest := &volume.ResizeVolumeRequest{
		VolumeID:  volumeID,
		SizeBytes: size,
	}
	_, err := client.ResizeVolume(context.TODO(), resizeVolumeRequest)
	if err == nil {
		t.Fatalf("Error empty for volume resize. Volume: %s, Size: %d", volumeID, size)
	}
}

func runNegativeMountVolumeRequest(t *testing.T, client volume.Interface, volumeID, targetPath string) {
	// Mount the volume
	mountVolumeRequest := &volume.MountVolumeRequest{
		VolumeID:   volumeID,
		TargetPath: targetPath,
	}

	_, err := client.MountVolume(context.TODO(), mountVolumeRequest)
	if err == nil {
		t.Fatalf("Error empty for volume(%s) mount to path %s.", volumeID, targetPath)
	}
}

func runNegativeUnmountVolumeRequest(t *testing.T, client volume.Interface, volumeID, targetPath string) {
	// Unmount the volume
	unmountVolumeRequest := &volume.UnmountVolumeRequest{
		VolumeID:   volumeID,
		TargetPath: targetPath,
	}
	_, err := client.UnmountVolume(context.TODO(), unmountVolumeRequest)
	if err == nil {
		t.Fatalf("Empty error. Volume id %s dismount from path %s ", volumeID, targetPath)
	}
}

func runNegativeVolumeStatsRequest(t *testing.T, client volume.Interface, volumeID string) {
	// Get VolumeStats
	volumeStatsRequest := &volume.GetVolumeStatsRequest{
		VolumeID: volumeID,
	}
	_, err := client.GetVolumeStats(context.TODO(), volumeStatsRequest)
	if err == nil {
		t.Errorf("Empty error. VolumeStats for id %s", volumeID)
	}
}

func negativeVolumeTests(t *testing.T) {
	client, err := volume.New(volumeapi.New())
	require.Nil(t, err)

	// Empty volume test
	runNegativeIsVolumeFormattedRequest(t, client, "")
	// -ve volume id
	runNegativeIsVolumeFormattedRequest(t, client, "-1")

	// Format volume negative tests
	runNegativeFormatVolumeRequest(t, client, "")
	runNegativeFormatVolumeRequest(t, client, "-1")

	// Resize volume negative tests
	runNegativeResizeVolumeRequest(t, client, "", 2*1024*1024)
	runNegativeResizeVolumeRequest(t, client, "-1", 2*1024*1024)

	// Mount volume negative tests
	runNegativeMountVolumeRequest(t, client, "", "")
	runNegativeMountVolumeRequest(t, client, "-1", "")

	// Unmount volume negative tests
	runNegativeUnmountVolumeRequest(t, client, "", "")
	runNegativeUnmountVolumeRequest(t, client, "-1", "")

	runNegativeVolumeStatsRequest(t, client, "")
	runNegativeVolumeStatsRequest(t, client, "-1")
}

func getClosestVolumeFromTargetPathTests(diskClient disk.Interface, volumeClient volume.Interface, t *testing.T) {
	t.Run("DriveLetterVolume", func(t *testing.T) {
		vhd, _, vhdCleanup := volumeInit(volumeClient, t)
		defer vhdCleanup()

		// vhd.Mount dir exists, because there are no volumes above it should return the C:\ volume
		request := &volume.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: vhd.Mount,
		}
		response, err := volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		// the C drive volume
		targetb, err := runPowershellCmd(t, `(Get-Partition -DriveLetter C | Get-Volume).UniqueId`)
		if err != nil {
			t.Fatalf("Failed to get the C: drive volume")
		}
		cDriveVolume := strings.TrimSpace(string(targetb))

		if response.VolumeID != cDriveVolume {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the C: drive volume")
		}
	})
	t.Run("AncestorVolumeFromNestedDirectory", func(t *testing.T) {
		var err error
		vhd, volumeID, vhdCleanup := volumeInit(volumeClient, t)
		defer vhdCleanup()

		// Mount the volume
		mountVolumeRequest := &volume.MountVolumeRequest{
			VolumeID:   volumeID,
			TargetPath: vhd.Mount,
		}
		_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
		if err != nil {
			t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
		}

		// Unmount the volume
		defer func() {
			unmountVolumeRequest := &volume.UnmountVolumeRequest{
				VolumeID:   volumeID,
				TargetPath: vhd.Mount,
			}
			_, err = volumeClient.UnmountVolume(context.TODO(), unmountVolumeRequest)
			if err != nil {
				t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
			}
		}()

		nestedDirectory := filepath.Join(vhd.Mount, "foo/bar")
		err = os.MkdirAll(nestedDirectory, os.ModeDir)
		if err != nil {
			t.Fatalf("Failed to create directory=%s", nestedDirectory)
		}

		// the volume returned should be the VHD volume
		request := &volume.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: nestedDirectory,
		}
		response, err := volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		if response.VolumeID != volumeID {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the VHD volume=%s", volumeID)
		}
	})

	t.Run("SymlinkToVolume", func(t *testing.T) {
		var err error
		vhd, volumeID, vhdCleanup := volumeInit(volumeClient, t)
		defer vhdCleanup()

		// Mount the volume
		mountVolumeRequest := &volume.MountVolumeRequest{
			VolumeID:   volumeID,
			TargetPath: vhd.Mount,
		}
		_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
		if err != nil {
			t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
		}

		// Unmount the volume
		defer func() {
			unmountVolumeRequest := &volume.UnmountVolumeRequest{
				VolumeID:   volumeID,
				TargetPath: vhd.Mount,
			}
			_, err = volumeClient.UnmountVolume(context.TODO(), unmountVolumeRequest)
			if err != nil {
				t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
			}
		}()

		testPluginPath, _ := getTestPluginPath()
		err = os.MkdirAll(testPluginPath, os.ModeDir)
		if err != nil {
			t.Fatalf("Failed to create directory=%s", testPluginPath)
		}

		sourceSymlink := filepath.Join(testPluginPath, "source")
		err = os.Symlink(vhd.Mount, sourceSymlink)
		if err != nil {
			t.Fatalf("Failed to create the symlink=%s", sourceSymlink)
		}

		// the volume returned should be the VHD volume
		var request *volume.GetClosestVolumeIDFromTargetPathRequest
		var response *volume.GetClosestVolumeIDFromTargetPathResponse
		request = &volume.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: sourceSymlink,
		}
		response, err = volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		if response.VolumeID != volumeID {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the VHD volume=%s", volumeID)
		}
	})
}

func mountVolumeTests(diskClient disk.Interface, volumeClient volume.Interface, t *testing.T) {
	vhd, volumeID, vhdCleanup := volumeInit(volumeClient, t)
	defer vhdCleanup()

	volumeStatsRequest := &volume.GetVolumeStatsRequest{
		VolumeID: volumeID,
	}

	volumeStatsResponse, err := volumeClient.GetVolumeStats(context.TODO(), volumeStatsRequest)
	if err != nil {
		t.Fatalf("VolumeStats request error: %v", err)
	}
	// For a volume formatted with 1GB it should be around 1GB, in practice it was 1056947712 bytes or 0.9844GB
	// let's compare with a range of +- 20MB
	if !sizeIsAround(t, volumeStatsResponse.TotalBytes, vhd.InitialSize) {
		t.Fatalf("volumeStatsResponse.TotalBytes reported is not valid, it is %v", volumeStatsResponse.TotalBytes)
	}

	// Resize the disk to twice its size (from 1GB to 2GB)
	// To resize a volume we need to resize the virtual hard disk first and then the partition
	cmd := fmt.Sprintf("Resize-VHD -Path %s -SizeBytes %d", vhd.Path, int64(vhd.InitialSize*2))
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s.", err, cmd, out)
	}

	// Resize the volume to 1.5GB
	oldVolumeSize := volumeStatsResponse.TotalBytes
	newVolumeSize := int64(float32(oldVolumeSize) * 1.5)

	// This is the max partition size when doing a resize to 2GB
	//
	//    Get-PartitionSupportedSize -DiskNumber 7 -PartitionNumber 2 | ConvertTo-Json
	//    {
	//    	"SizeMin":  404725760,
	//    	"SizeMax":  2130689536
	//    }
	resizeVolumeRequest := &volume.ResizeVolumeRequest{
		VolumeID: volumeID,
		// resize the partition to 1.5x times instead
		SizeBytes: newVolumeSize,
	}

	t.Logf("Attempt to resize volume from sizeBytes=%d to sizeBytes=%d", oldVolumeSize, newVolumeSize)

	_, err = volumeClient.ResizeVolume(context.TODO(), resizeVolumeRequest)
	if err != nil {
		t.Fatalf("Volume resize request failed. Error: %v", err)
	}

	volumeStatsResponse, err = volumeClient.GetVolumeStats(context.TODO(), volumeStatsRequest)
	if err != nil {
		t.Fatalf("VolumeStats request after resize error: %v", err)
	}
	// resizing from 1GB to approximately 1.5GB
	if !sizeIsAround(t, volumeStatsResponse.TotalBytes, newVolumeSize) {
		t.Fatalf("VolumeSize reported should be greater than the old size, it is %v", volumeStatsResponse.TotalBytes)
	}

	volumeDiskNumberRequest := &volume.GetDiskNumberFromVolumeIDRequest{
		VolumeID: volumeID,
	}

	volumeDiskNumberResponse, err := volumeClient.GetDiskNumberFromVolumeID(context.TODO(), volumeDiskNumberRequest)
	if err != nil {
		t.Fatalf("GetDiskNumberFromVolumeID failed: %v", err)
	}

	diskStatsRequest := &disk.GetDiskStatsRequest{
		DiskNumber: volumeDiskNumberResponse.DiskNumber,
	}

	diskStatsResponse, err := diskClient.GetDiskStats(context.TODO(), diskStatsRequest)
	if err != nil {
		t.Fatalf("DiskStats request error: %v", err)
	}

	if diskStatsResponse.TotalBytes < 0 {
		t.Fatalf("Invalid disk size was returned %v", diskStatsResponse.TotalBytes)
	}

	// Mount the volume
	mountVolumeRequest := &volume.MountVolumeRequest{
		VolumeID:   volumeID,
		TargetPath: vhd.Mount,
	}
	_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}

	// Unmount the volume
	unmountVolumeRequest := &volume.UnmountVolumeRequest{
		VolumeID:   volumeID,
		TargetPath: vhd.Mount,
	}
	_, err = volumeClient.UnmountVolume(context.TODO(), unmountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}
}
