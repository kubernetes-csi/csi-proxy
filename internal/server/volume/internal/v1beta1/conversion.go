package v1beta1

// Add manual conversion functions here to override automatic conversion functions

import (
	"github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta1"
	"github.com/kubernetes-csi/csi-proxy/internal/server/volume/internal"
)

func Convert_internal_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse(in *internal.VolumeStatsResponse, out *v1beta1.VolumeStatsResponse) error {
	out.VolumeSize = in.VolumeSize
	out.VolumeUsedSize = in.VolumeUsedSize
	return nil
}

func Convert_v1beta1_VolumeStatsRequest_To_internal_VolumeStatsRequest(in *v1beta1.VolumeStatsRequest, out *internal.VolumeStatsRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

func Convert_internal_VolumeDiskNumberResponse_To_v1beta1_VolumeDiskNumberResponse(in *internal.VolumeDiskNumberResponse, out *v1beta1.VolumeDiskNumberResponse) error {
	out.DiskNumber = in.DiskNumber
	return nil
}

func Convert_v1beta1_VolumeDiskNumberRequest_To_internal_VolumeDiskNumberRequest(in *v1beta1.VolumeDiskNumberRequest, out *internal.VolumeDiskNumberRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

func Convert_internal_VolumeIDFromMountResponse_To_v1beta1_VolumeIDFromMountResponse(in *internal.VolumeIDFromMountResponse, out *v1beta1.VolumeIDFromMountResponse) error {
	out.VolumeId = in.VolumeId
	return nil
}

func Convert_v1beta1_VolumeIDFromMountRequest_To_internal_VolumeIDFromMountRequest(in *v1beta1.VolumeIDFromMountRequest, out *internal.VolumeIDFromMountRequest) error {
	out.Mount = in.Mount
	return nil
}
