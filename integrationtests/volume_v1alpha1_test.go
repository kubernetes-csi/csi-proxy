package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1alpha1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/volume/v1alpha1"
)

// v1alpha1VolumeTests tests that the API is compatible with versions that are before
// the latest, e.g. that a v1alpha1 client can still use the server csi-proxy v1beta3
func v1alpha1VolumeTests(t *testing.T) {
	// it's intended for this client to be v1alpha1
	// i.e. don't change it if there are upgrades
	var err error
	v1alpha1Client, err := v1alpha1client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create new v1alpha1 client, err=%+v", err)
	}

	vhd, vhdCleanup := diskInit(t)
	defer vhdCleanup()

	// get first volume
	listRequest := &v1alpha1.ListVolumesOnDiskRequest{
		DiskId: strconv.FormatUint(uint64(vhd.DiskNumber), 10),
	}
	listResponse, err := v1alpha1Client.ListVolumesOnDisk(context.TODO(), listRequest)
	if err != nil {
		t.Fatalf("List response: %v", err)
	}

	volumeIDsLen := len(listResponse.VolumeIds)
	if volumeIDsLen != 1 {
		t.Fatalf("Number of volumes not equal to 1: %d", volumeIDsLen)
	}
	volumeID := listResponse.VolumeIds[0]

	// format volume (skip IsVolumeFormatted calls)
	formatVolumeRequest := &v1alpha1.FormatVolumeRequest{
		VolumeId: volumeID,
	}
	_, err = v1alpha1Client.FormatVolume(context.TODO(), formatVolumeRequest)
	if err != nil {
		t.Fatalf("Volume format failed. Error: %v", err)
	}

	// check that the volume is formatted now
	isVolumeFormattedRequest := &v1alpha1.IsVolumeFormattedRequest{
		VolumeId: volumeID,
	}
	isVolumeFormattedResponse, err := v1alpha1Client.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err != nil {
		t.Fatalf("Is volume formatted request error: %v", err)
	}
	if !isVolumeFormattedResponse.Formatted {
		t.Fatal("Volume must be formatted at this point")
	}

	// Resize the disk to twice its size (from 1GB to 2GB)
	// To resize a volume we need to resize the virtual hard disk first and then the partition
	cmd := fmt.Sprintf("Resize-VHD -Path %s -SizeBytes %d", vhd.Path, int64(vhd.InitialSize*2))
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s.", err, cmd, out)
	}

	// in v1alpha1 there's no stats call, just make sure that we can resize a volume without problems
	oldVolumeSize := vhd.InitialSize
	newVolumeSize := int64(float32(oldVolumeSize) * 1.5)
	resizeVolumeRequest := &v1alpha1.ResizeVolumeRequest{
		VolumeId: volumeID,
		// resize the partition to 1.5x times instead
		Size: newVolumeSize,
	}
	t.Logf("Attempt to resize volume from sizeBytes=%d to sizeBytes=%d", oldVolumeSize, newVolumeSize)
	_, err = v1alpha1Client.ResizeVolume(context.TODO(), resizeVolumeRequest)
	if err != nil {
		t.Fatalf("Volume resize request failed. Error: %v", err)
	}

	// Mount the volume
	mountVolumeRequest := &v1alpha1.MountVolumeRequest{
		VolumeId: volumeID,
		Path:     vhd.Mount,
	}
	_, err = v1alpha1Client.MountVolume(context.TODO(), mountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}

	// Unmount the volume
	dismountVolumeRequest := &v1alpha1.DismountVolumeRequest{
		VolumeId: volumeID,
		Path:     vhd.Mount,
	}
	_, err = v1alpha1Client.DismountVolume(context.TODO(), dismountVolumeRequest)
	if err != nil {
		t.Fatalf("Volume id %s mount to path %s failed. Error: %v", volumeID, vhd.Mount, err)
	}
}
