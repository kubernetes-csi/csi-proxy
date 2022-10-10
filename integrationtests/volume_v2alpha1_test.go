package integrationtests

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	diskv1 "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1"
	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v2alpha1"
	diskv1client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1"
	v2alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/volume/v2alpha1"
)

// volumeInitV2Alpha1 initializes a volume, it creates a VHD, initializes it,
// creates a partition with the max size and formats the volume corresponding to that partition
func volumeInitV2Alpha1(volumeClient *v2alpha1client.Client, t *testing.T) (*VirtualHardDisk, string, func()) {
	vhd, vhdCleanup := diskInit(t)

	listRequest := &v2alpha1.ListVolumesOnDiskRequest{
		DiskNumber: vhd.DiskNumber,
	}
	listResponse, err := volumeClient.ListVolumesOnDisk(context.TODO(), listRequest)
	if err != nil {
		t.Fatalf("List response: %v", err)
	}

	volumeIDsLen := len(listResponse.VolumeIds)
	if volumeIDsLen != 1 {
		t.Fatalf("Number of volumes not equal to 1: %d", volumeIDsLen)
	}
	volumeID := listResponse.VolumeIds[0]
	t.Logf("VolumeId %v", volumeID)

	isVolumeFormattedRequest := &v2alpha1.IsVolumeFormattedRequest{
		VolumeId: volumeID,
	}
	isVolumeFormattedResponse, err := volumeClient.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err != nil {
		t.Fatalf("Is volume formatted request error: %v", err)
	}
	if isVolumeFormattedResponse.Formatted {
		t.Fatal("Volume formatted. Unexpected !!")
	}

	formatVolumeRequest := &v2alpha1.FormatVolumeRequest{
		VolumeId: volumeID,
	}
	_, err = volumeClient.FormatVolume(context.TODO(), formatVolumeRequest)
	if err != nil {
		t.Fatalf("Volume format failed. Error: %v", err)
	}

	isVolumeFormattedResponse, err = volumeClient.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err != nil {
		t.Fatalf("Is volume formatted request error: %v", err)
	}
	if !isVolumeFormattedResponse.Formatted {
		t.Fatal("Volume should be formatted. Unexpected !!")
	}
	return vhd, volumeID, vhdCleanup
}

func v2alpha1GetClosestVolumeFromTargetPathTests(diskClient *diskv1client.Client, volumeClient *v2alpha1client.Client, t *testing.T) {
	t.Run("DriveLetterVolume", func(t *testing.T) {
		vhd, _, vhdCleanup := volumeInitV2Alpha1(volumeClient, t)
		defer vhdCleanup()

		// vhd.Mount dir exists, because there are no volumes above it should return the C:\ volume
		var request *v2alpha1.GetClosestVolumeIDFromTargetPathRequest
		var response *v2alpha1.GetClosestVolumeIDFromTargetPathResponse
		request = &v2alpha1.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: vhd.Mount,
		}
		response, err := volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		// the C drive volume
		cmd := exec.Command("powershell", "/c", `(Get-Partition -DriveLetter C | Get-Volume).UniqueId`)
		targetb, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to get the C: drive volume")
		}
		cDriveVolume := strings.TrimSpace(string(targetb))

		if response.VolumeId != cDriveVolume {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the C: drive volume")
		}
	})
	t.Run("AncestorVolumeFromNestedDirectory", func(t *testing.T) {
		var err error
		vhd, volumeID, vhdCleanup := volumeInitV2Alpha1(volumeClient, t)
		defer vhdCleanup()

		// Mount the volume
		mountVolumeRequest := &v2alpha1.MountVolumeRequest{
			VolumeId:   volumeID,
			TargetPath: vhd.Mount,
		}
		_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
		if err != nil {
			t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
		}

		// Unmount the volume
		defer func() {
			unmountVolumeRequest := &v2alpha1.UnmountVolumeRequest{
				VolumeId:   volumeID,
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
		var request *v2alpha1.GetClosestVolumeIDFromTargetPathRequest
		var response *v2alpha1.GetClosestVolumeIDFromTargetPathResponse
		request = &v2alpha1.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: nestedDirectory,
		}
		response, err = volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		if response.VolumeId != volumeID {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the VHD volume=%s", volumeID)
		}
	})

	t.Run("SymlinkToVolume", func(t *testing.T) {
		var err error
		vhd, volumeID, vhdCleanup := volumeInitV2Alpha1(volumeClient, t)
		defer vhdCleanup()

		// Mount the volume
		mountVolumeRequest := &v2alpha1.MountVolumeRequest{
			VolumeId:   volumeID,
			TargetPath: vhd.Mount,
		}
		_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
		if err != nil {
			t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
		}

		// Unmount the volume
		defer func() {
			unmountVolumeRequest := &v2alpha1.UnmountVolumeRequest{
				VolumeId:   volumeID,
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
		var request *v2alpha1.GetClosestVolumeIDFromTargetPathRequest
		var response *v2alpha1.GetClosestVolumeIDFromTargetPathResponse
		request = &v2alpha1.GetClosestVolumeIDFromTargetPathRequest{
			TargetPath: sourceSymlink,
		}
		response, err = volumeClient.GetClosestVolumeIDFromTargetPath(context.TODO(), request)
		if err != nil {
			t.Fatalf("GetClosestVolumeIDFromTargetPath request error, err=%v", err)
		}

		if response.VolumeId != volumeID {
			t.Fatalf("The volume from GetClosestVolumeIDFromTargetPath doesn't match the VHD volume=%s", volumeID)
		}
	})
}

func v2alpha1MountVolumeTests(diskClient *diskv1client.Client, volumeClient *v2alpha1client.Client, t *testing.T) {
	vhd, volumeID, vhdCleanup := volumeInitV2Alpha1(volumeClient, t)
	defer vhdCleanup()

	volumeStatsRequest := &v2alpha1.GetVolumeStatsRequest{
		VolumeId: volumeID,
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
	resizeVolumeRequest := &v2alpha1.ResizeVolumeRequest{
		VolumeId: volumeID,
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

	volumeDiskNumberRequest := &v2alpha1.GetDiskNumberFromVolumeIDRequest{
		VolumeId: volumeID,
	}

	volumeDiskNumberResponse, err := volumeClient.GetDiskNumberFromVolumeID(context.TODO(), volumeDiskNumberRequest)
	if err != nil {
		t.Fatalf("GetDiskNumberFromVolumeID failed: %v", err)
	}

	diskStatsRequest := &diskv1.GetDiskStatsRequest{
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
	mountVolumeRequest := &v2alpha1.MountVolumeRequest{
		VolumeId:   volumeID,
		TargetPath: vhd.Mount,
	}
	_, err = volumeClient.MountVolume(context.TODO(), mountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}

	// Unmount the volume
	unmountVolumeRequest := &v2alpha1.UnmountVolumeRequest{
		VolumeId:   volumeID,
		TargetPath: vhd.Mount,
	}
	_, err = volumeClient.UnmountVolume(context.TODO(), unmountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}
}

func v2alpha1VolumeTests(t *testing.T) {
	var volumeClient *v2alpha1client.Client
	var diskClient *diskv1client.Client
	var err error

	if volumeClient, err = v2alpha1client.NewClient(); err != nil {
		t.Fatalf("Client new error: %v", err)
	}
	defer volumeClient.Close()

	if diskClient, err = diskv1client.NewClient(); err != nil {
		t.Fatalf("DiskClient new error: %v", err)
	}
	defer diskClient.Close()

	t.Run("MountVolume", func(t *testing.T) {
		v2alpha1MountVolumeTests(diskClient, volumeClient, t)
	})
	t.Run("GetClosestVolumeFromTargetPath", func(t *testing.T) {
		v2alpha1GetClosestVolumeFromTargetPathTests(diskClient, volumeClient, t)
	})
}
