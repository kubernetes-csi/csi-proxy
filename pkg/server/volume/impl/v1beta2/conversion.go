package v1beta2

import (
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta2"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/volume/impl"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_v1beta2_ListVolumesOnDiskRequest_To_impl_ListVolumesOnDiskRequest(in *v1beta2.ListVolumesOnDiskRequest, out *impl.ListVolumesOnDiskRequest) error {
	diskIDUint, err := strconv.ParseUint(in.DiskId, 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse diskId: err=%+v", err)
	}
	out.DiskNumber = uint32(diskIDUint)
	return nil
}

func Convert_impl_ListVolumesOnDiskRequest_To_v1beta2_ListVolumesOnDiskRequest(in *impl.ListVolumesOnDiskRequest, out *v1beta2.ListVolumesOnDiskRequest) error {
	out.DiskId = strconv.FormatUint(uint64(in.DiskNumber), 10)
	return nil
}

func Convert_v1beta2_MountVolumeRequest_To_impl_MountVolumeRequest(in *v1beta2.MountVolumeRequest, out *impl.MountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.TargetPath = in.Path
	return nil
}

func Convert_impl_MountVolumeRequest_To_v1beta2_MountVolumeRequest(in *impl.MountVolumeRequest, out *v1beta2.MountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.TargetPath
	return nil
}

func Convert_v1beta2_ResizeVolumeRequest_To_impl_ResizeVolumeRequest(in *v1beta2.ResizeVolumeRequest, out *impl.ResizeVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.SizeBytes = in.Size
	return nil
}

func Convert_impl_ResizeVolumeRequest_To_v1beta2_ResizeVolumeRequest(in *impl.ResizeVolumeRequest, out *v1beta2.ResizeVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Size = in.SizeBytes
	return nil
}
