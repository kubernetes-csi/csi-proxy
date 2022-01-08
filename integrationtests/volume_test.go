package integrationtests

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta3"
	v1beta3client "github.com/kubernetes-csi/csi-proxy/client/groups/volume/v1beta3"
)

func runNegativeListVolumeRequest(t *testing.T, client *v1beta3client.Client, diskNum uint32) {
	listRequest := &v1beta3.ListVolumesOnDiskRequest{
		DiskNumber: diskNum,
	}
	_, err := client.ListVolumesOnDisk(context.TODO(), listRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for disknum:%d", diskNum)
	}
}

func runNegativeIsVolumeFormattedRequest(t *testing.T, client *v1beta3client.Client, volumeID string) {
	isVolumeFormattedRequest := &v1beta3.IsVolumeFormattedRequest{
		VolumeId: volumeID,
	}
	_, err := client.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for volumeID: %s", volumeID)
	}
}

func runNegativeFormatVolumeRequest(t *testing.T, client *v1beta3client.Client, volumeID string) {
	formatVolumeRequest := &v1beta3.FormatVolumeRequest{
		VolumeId: volumeID,
	}
	_, err := client.FormatVolume(context.TODO(), formatVolumeRequest)
	if err == nil {
		t.Fatalf("Empty error. Expected error for volume id: %s", volumeID)
	}
}

func runNegativeResizeVolumeRequest(t *testing.T, client *v1beta3client.Client, volumeID string, size int64) {
	resizeVolumeRequest := &v1beta3.ResizeVolumeRequest{
		VolumeId:  volumeID,
		SizeBytes: size,
	}
	_, err := client.ResizeVolume(context.TODO(), resizeVolumeRequest)
	if err == nil {
		t.Fatalf("Error empty for volume resize. Volume: %s, Size: %d", volumeID, size)
	}
}

func runNegativeMountVolumeRequest(t *testing.T, client *v1beta3client.Client, volumeID, targetPath string) {
	// Mount the volume
	mountVolumeRequest := &v1beta3.MountVolumeRequest{
		VolumeId:   volumeID,
		TargetPath: targetPath,
	}

	_, err := client.MountVolume(context.TODO(), mountVolumeRequest)
	if err == nil {
		t.Fatalf("Error empty for volume(%s) mount to path %s.", volumeID, targetPath)
	}
}

func runNegativeUnmountVolumeRequest(t *testing.T, client *v1beta3client.Client, volumeID, targetPath string) {
	// Unmount the volume
	unmountVolumeRequest := &v1beta3.UnmountVolumeRequest{
		VolumeId:   volumeID,
		TargetPath: targetPath,
	}
	_, err := client.UnmountVolume(context.TODO(), unmountVolumeRequest)
	if err == nil {
		t.Fatalf("Empty error. Volume id %s dismount from path %s ", volumeID, targetPath)
	}
}

func runNegativeVolumeStatsRequest(t *testing.T, client *v1beta3client.Client, volumeID string) {
	// Get VolumeStats
	volumeStatsRequest := &v1beta3.GetVolumeStatsRequest{
		VolumeId: volumeID,
	}
	_, err := client.GetVolumeStats(context.TODO(), volumeStatsRequest)
	if err == nil {
		t.Errorf("Empty error. VolumeStats for id %s", volumeID)
	}
}

func negativeVolumeTests(t *testing.T) {
	var client *v1beta3client.Client
	var err error

	if client, err = v1beta3client.NewClient(); err != nil {
		t.Fatalf("Client new error: %v", err)
	}
	defer client.Close()

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

// sizeIsAround returns true if the actual size is around the expected size
// (considers the fact that some bytes were lost)
func sizeIsAround(t *testing.T, actualSize, expectedSize int64) bool {
	// An upper bound on the number of bytes that are lost when creating or resizing a partition
	var volumeSizeBytesLoss int64 = (20 * 1024 * 1024)
	var lowerBound = expectedSize - volumeSizeBytesLoss
	var upperBound = expectedSize
	t.Logf("Checking that the size is inside the bounds: %d < (actual) %d < %d", lowerBound, actualSize, upperBound)
	return lowerBound <= actualSize && actualSize <= upperBound
}

func negativeDiskTests(t *testing.T) {
	var client *v1beta3client.Client
	var err error

	if client, err = v1beta3client.NewClient(); err != nil {
		t.Fatalf("Client new error: %v", err)
	}
	defer client.Close()
}

func TestVolumeAPIs(t *testing.T) {
	t.Run("NegativeDiskTests", func(t *testing.T) {
		negativeDiskTests(t)
	})
	t.Run("NegativeVolumeTests", func(t *testing.T) {
		negativeVolumeTests(t)
	})

	// TODO: These tests will fail on Github Actions because Hyper-V is disabled
	// see https://github.com/actions/virtual-environments/pull/2525

	// these tests should be considered frozen from the API point of view
	t.Run("v1alpha1Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1alpha1VolumeTests(t)
	})
	t.Run("v1beta1Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1beta1VolumeTests(t)
	})
	t.Run("v1beta2Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1beta2VolumeTests(t)
	})
	t.Run("v1beta3Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1beta3VolumeTests(t)
	})
	t.Run("v1Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1VolumeTests(t)
	})
	t.Run("v2alpha1Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v2alpha1VolumeTests(t)
	})
}
