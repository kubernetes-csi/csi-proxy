package v1beta1

// Add manual conversion functions here to override automatic conversion functions

import (
	v1beta1 "github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta1"
	internal "github.com/kubernetes-csi/csi-proxy/internal/server/volume/internal"
)

func Convert_internal_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse(in *internal.VolumeStatsResponse, out *v1beta1.VolumeStatsResponse) error {
	out.DiskSize = in.DiskSize
	out.VolumeSize = in.VolumeSize
	out.VolumeUsedSize = in.VolumeUsedSize
	return nil
}

func Convert_v1beta1_VolumeStatsRequest_To_internal_VolumeStatsRequest(in *v1beta1.VolumeStatsRequest, out *internal.VolumeStatsRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}
