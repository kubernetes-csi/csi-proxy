package integrationtests

import (
	"context"
	"strconv"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta2"
	v1beta2client "github.com/kubernetes-csi/csi-proxy/client/groups/volume/v1beta2"
)

// v1beta2VolumeTests tests that the API is compatible with versions that are before
// the latest, e.g. that a v1beta2 client can still use the server csi-proxy v1beta3
func v1beta2VolumeTests(t *testing.T) {
	// it's intended for this client to be before v1beta3!
	// i.e. don't change it if there are upgrades
	var err error
	v1beta2Client, err := v1beta2client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create new v1beta2 client, err=%+v", err)
	}

	vhd, vhdCleanup := diskInit(t)
	defer vhdCleanup()

	// get first volume
	listRequest := &v1beta2.ListVolumesOnDiskRequest{
		DiskId: strconv.FormatUint(uint64(vhd.DiskNumber), 10),
	}
	listResponse, err := v1beta2Client.ListVolumesOnDisk(context.TODO(), listRequest)
	if err != nil {
		t.Fatalf("List response: %v", err)
	}

	volumeIDsLen := len(listResponse.VolumeIds)
	if volumeIDsLen != 1 {
		t.Fatalf("Number of volumes not equal to 1: %d", volumeIDsLen)
	}
	volumeID := listResponse.VolumeIds[0]

	// format volume (skip IsVolumeFormatted calls)
	formatVolumeRequest := &v1beta2.FormatVolumeRequest{
		VolumeId: volumeID,
	}
	_, err = v1beta2Client.FormatVolume(context.TODO(), formatVolumeRequest)
	if err != nil {
		t.Fatalf("Volume format failed. Error: %v", err)
	}

	// VolumeStats (volume was formatted to 1GB)
	volumeStatsRequest := &v1beta2.VolumeStatsRequest{
		VolumeId: volumeID,
	}
	volumeStatsResponse, err := v1beta2Client.VolumeStats(context.TODO(), volumeStatsRequest)
	if err != nil {
		t.Fatalf("VolumeStats request error: %v", err)
	}
	// For a volume formatted with 1GB it should be around 1GB, in practice it was 1056947712 bytes or 0.9844GB
	// let's compare with a range of - 20MB
	if !sizeIsAround(t, volumeStatsResponse.VolumeSize, vhd.InitialSize) {
		t.Fatalf("volumeStatsResponse.VolumeSize reported is not valid, it is %v", volumeStatsResponse.VolumeSize)
	}

	volumeDiskNumberRequest := &v1beta2.VolumeDiskNumberRequest{
		VolumeId: volumeID,
	}
	_, err = v1beta2Client.GetVolumeDiskNumber(context.TODO(), volumeDiskNumberRequest)
	if err != nil {
		t.Fatalf("GetVolumeDiskNumber failed: %v", err)
	}

	// Resize the volume to 1.5GB
	oldVolumeSize := volumeStatsResponse.VolumeSize
	newVolumeSize := int64(float32(oldVolumeSize) * 1.5)
	resizeVolumeRequest := &v1beta2.ResizeVolumeRequest{
		VolumeId: volumeID,
		Size:     newVolumeSize,
	}
	t.Logf("Attempt to resize volume from sizeBytes=%d to sizeBytes=%d", oldVolumeSize, newVolumeSize)
	_, err = v1beta2Client.ResizeVolume(context.TODO(), resizeVolumeRequest)
	if err != nil {
		t.Fatalf("Volume resize request failed. Error: %v", err)
	}
	volumeStatsResponse, err = v1beta2Client.VolumeStats(context.TODO(), volumeStatsRequest)
	if err != nil {
		t.Fatalf("VolumeStats request after resize error: %v", err)
	}
	if !sizeIsAround(t, volumeStatsResponse.VolumeSize, newVolumeSize) {
		t.Fatalf("VolumeSize reported should be greater than the old size, it is %v", volumeStatsResponse.VolumeSize)
	}

	// Mount the volume
	mountVolumeRequest := &v1beta2.MountVolumeRequest{
		VolumeId: volumeID,
		Path:     vhd.Mount,
	}
	_, err = v1beta2Client.MountVolume(context.TODO(), mountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}

	// Unmount the volume
	dismountVolumeRequest := &v1beta2.DismountVolumeRequest{
		VolumeId: volumeID,
		Path:     vhd.Mount,
	}
	_, err = v1beta2Client.DismountVolume(context.TODO(), dismountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}
}
