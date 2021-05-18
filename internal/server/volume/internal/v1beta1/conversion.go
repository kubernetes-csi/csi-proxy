package v1beta1

import (
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta1"
	internal "github.com/kubernetes-csi/csi-proxy/internal/server/volume/internal"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_v1beta1_ListVolumesOnDiskRequest_To_internal_ListVolumesOnDiskRequest(in *v1beta1.ListVolumesOnDiskRequest, out *internal.ListVolumesOnDiskRequest) error {
	diskIDUint, err := strconv.ParseUint(in.DiskId, 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse diskId: err=%+v", err)
	}
	out.DiskNumber = uint32(diskIDUint)
	return nil
}

func Convert_internal_ListVolumesOnDiskRequest_To_v1beta1_ListVolumesOnDiskRequest(in *internal.ListVolumesOnDiskRequest, out *v1beta1.ListVolumesOnDiskRequest) error {
	out.DiskId = strconv.FormatUint(uint64(in.DiskNumber), 10)
	return nil
}

func Convert_v1beta1_MountVolumeRequest_To_internal_MountVolumeRequest(in *v1beta1.MountVolumeRequest, out *internal.MountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.TargetPath = in.Path
	return nil
}

func Convert_internal_MountVolumeRequest_To_v1beta1_MountVolumeRequest(in *internal.MountVolumeRequest, out *v1beta1.MountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.TargetPath
	return nil
}

func Convert_v1beta1_ResizeVolumeRequest_To_internal_ResizeVolumeRequest(in *v1beta1.ResizeVolumeRequest, out *internal.ResizeVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.SizeBytes = in.Size
	return nil
}

func Convert_internal_ResizeVolumeRequest_To_v1beta1_ResizeVolumeRequest(in *internal.ResizeVolumeRequest, out *v1beta1.ResizeVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Size = in.SizeBytes
	return nil
}
