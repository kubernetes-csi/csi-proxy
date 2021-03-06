// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	v1beta1 "github.com/kubernetes-csi/csi-proxy/client/api/volume/v1beta1"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/volume/impl"
)

func autoConvert_v1beta1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in *v1beta1.DismountVolumeRequest, out *impl.DismountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.Path
	return nil
}

// Convert_v1beta1_DismountVolumeRequest_To_impl_DismountVolumeRequest is an autogenerated conversion function.
func Convert_v1beta1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in *v1beta1.DismountVolumeRequest, out *impl.DismountVolumeRequest) error {
	return autoConvert_v1beta1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in, out)
}

func autoConvert_impl_DismountVolumeRequest_To_v1beta1_DismountVolumeRequest(in *impl.DismountVolumeRequest, out *v1beta1.DismountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.Path
	return nil
}

// Convert_impl_DismountVolumeRequest_To_v1beta1_DismountVolumeRequest is an autogenerated conversion function.
func Convert_impl_DismountVolumeRequest_To_v1beta1_DismountVolumeRequest(in *impl.DismountVolumeRequest, out *v1beta1.DismountVolumeRequest) error {
	return autoConvert_impl_DismountVolumeRequest_To_v1beta1_DismountVolumeRequest(in, out)
}

func autoConvert_v1beta1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in *v1beta1.DismountVolumeResponse, out *impl.DismountVolumeResponse) error {
	return nil
}

// Convert_v1beta1_DismountVolumeResponse_To_impl_DismountVolumeResponse is an autogenerated conversion function.
func Convert_v1beta1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in *v1beta1.DismountVolumeResponse, out *impl.DismountVolumeResponse) error {
	return autoConvert_v1beta1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in, out)
}

func autoConvert_impl_DismountVolumeResponse_To_v1beta1_DismountVolumeResponse(in *impl.DismountVolumeResponse, out *v1beta1.DismountVolumeResponse) error {
	return nil
}

// Convert_impl_DismountVolumeResponse_To_v1beta1_DismountVolumeResponse is an autogenerated conversion function.
func Convert_impl_DismountVolumeResponse_To_v1beta1_DismountVolumeResponse(in *impl.DismountVolumeResponse, out *v1beta1.DismountVolumeResponse) error {
	return autoConvert_impl_DismountVolumeResponse_To_v1beta1_DismountVolumeResponse(in, out)
}

func autoConvert_v1beta1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in *v1beta1.FormatVolumeRequest, out *impl.FormatVolumeRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1beta1_FormatVolumeRequest_To_impl_FormatVolumeRequest is an autogenerated conversion function.
func Convert_v1beta1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in *v1beta1.FormatVolumeRequest, out *impl.FormatVolumeRequest) error {
	return autoConvert_v1beta1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in, out)
}

func autoConvert_impl_FormatVolumeRequest_To_v1beta1_FormatVolumeRequest(in *impl.FormatVolumeRequest, out *v1beta1.FormatVolumeRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_FormatVolumeRequest_To_v1beta1_FormatVolumeRequest is an autogenerated conversion function.
func Convert_impl_FormatVolumeRequest_To_v1beta1_FormatVolumeRequest(in *impl.FormatVolumeRequest, out *v1beta1.FormatVolumeRequest) error {
	return autoConvert_impl_FormatVolumeRequest_To_v1beta1_FormatVolumeRequest(in, out)
}

func autoConvert_v1beta1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in *v1beta1.FormatVolumeResponse, out *impl.FormatVolumeResponse) error {
	return nil
}

// Convert_v1beta1_FormatVolumeResponse_To_impl_FormatVolumeResponse is an autogenerated conversion function.
func Convert_v1beta1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in *v1beta1.FormatVolumeResponse, out *impl.FormatVolumeResponse) error {
	return autoConvert_v1beta1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in, out)
}

func autoConvert_impl_FormatVolumeResponse_To_v1beta1_FormatVolumeResponse(in *impl.FormatVolumeResponse, out *v1beta1.FormatVolumeResponse) error {
	return nil
}

// Convert_impl_FormatVolumeResponse_To_v1beta1_FormatVolumeResponse is an autogenerated conversion function.
func Convert_impl_FormatVolumeResponse_To_v1beta1_FormatVolumeResponse(in *impl.FormatVolumeResponse, out *v1beta1.FormatVolumeResponse) error {
	return autoConvert_impl_FormatVolumeResponse_To_v1beta1_FormatVolumeResponse(in, out)
}

func autoConvert_v1beta1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in *v1beta1.IsVolumeFormattedRequest, out *impl.IsVolumeFormattedRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1beta1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest is an autogenerated conversion function.
func Convert_v1beta1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in *v1beta1.IsVolumeFormattedRequest, out *impl.IsVolumeFormattedRequest) error {
	return autoConvert_v1beta1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in, out)
}

func autoConvert_impl_IsVolumeFormattedRequest_To_v1beta1_IsVolumeFormattedRequest(in *impl.IsVolumeFormattedRequest, out *v1beta1.IsVolumeFormattedRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_IsVolumeFormattedRequest_To_v1beta1_IsVolumeFormattedRequest is an autogenerated conversion function.
func Convert_impl_IsVolumeFormattedRequest_To_v1beta1_IsVolumeFormattedRequest(in *impl.IsVolumeFormattedRequest, out *v1beta1.IsVolumeFormattedRequest) error {
	return autoConvert_impl_IsVolumeFormattedRequest_To_v1beta1_IsVolumeFormattedRequest(in, out)
}

func autoConvert_v1beta1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in *v1beta1.IsVolumeFormattedResponse, out *impl.IsVolumeFormattedResponse) error {
	out.Formatted = in.Formatted
	return nil
}

// Convert_v1beta1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse is an autogenerated conversion function.
func Convert_v1beta1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in *v1beta1.IsVolumeFormattedResponse, out *impl.IsVolumeFormattedResponse) error {
	return autoConvert_v1beta1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in, out)
}

func autoConvert_impl_IsVolumeFormattedResponse_To_v1beta1_IsVolumeFormattedResponse(in *impl.IsVolumeFormattedResponse, out *v1beta1.IsVolumeFormattedResponse) error {
	out.Formatted = in.Formatted
	return nil
}

// Convert_impl_IsVolumeFormattedResponse_To_v1beta1_IsVolumeFormattedResponse is an autogenerated conversion function.
func Convert_impl_IsVolumeFormattedResponse_To_v1beta1_IsVolumeFormattedResponse(in *impl.IsVolumeFormattedResponse, out *v1beta1.IsVolumeFormattedResponse) error {
	return autoConvert_impl_IsVolumeFormattedResponse_To_v1beta1_IsVolumeFormattedResponse(in, out)
}

// detected external conversion function
// Convert_v1beta1_ListVolumesOnDiskRequest_To_impl_ListVolumesOnDiskRequest(in *v1beta1.ListVolumesOnDiskRequest, out *impl.ListVolumesOnDiskRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_ListVolumesOnDiskRequest_To_v1beta1_ListVolumesOnDiskRequest(in *impl.ListVolumesOnDiskRequest, out *v1beta1.ListVolumesOnDiskRequest) error
// skipping generation of the auto function

func autoConvert_v1beta1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in *v1beta1.ListVolumesOnDiskResponse, out *impl.ListVolumesOnDiskResponse) error {
	out.VolumeIds = *(*[]string)(unsafe.Pointer(&in.VolumeIds))
	return nil
}

// Convert_v1beta1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse is an autogenerated conversion function.
func Convert_v1beta1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in *v1beta1.ListVolumesOnDiskResponse, out *impl.ListVolumesOnDiskResponse) error {
	return autoConvert_v1beta1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in, out)
}

func autoConvert_impl_ListVolumesOnDiskResponse_To_v1beta1_ListVolumesOnDiskResponse(in *impl.ListVolumesOnDiskResponse, out *v1beta1.ListVolumesOnDiskResponse) error {
	out.VolumeIds = *(*[]string)(unsafe.Pointer(&in.VolumeIds))
	return nil
}

// Convert_impl_ListVolumesOnDiskResponse_To_v1beta1_ListVolumesOnDiskResponse is an autogenerated conversion function.
func Convert_impl_ListVolumesOnDiskResponse_To_v1beta1_ListVolumesOnDiskResponse(in *impl.ListVolumesOnDiskResponse, out *v1beta1.ListVolumesOnDiskResponse) error {
	return autoConvert_impl_ListVolumesOnDiskResponse_To_v1beta1_ListVolumesOnDiskResponse(in, out)
}

// detected external conversion function
// Convert_v1beta1_MountVolumeRequest_To_impl_MountVolumeRequest(in *v1beta1.MountVolumeRequest, out *impl.MountVolumeRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_MountVolumeRequest_To_v1beta1_MountVolumeRequest(in *impl.MountVolumeRequest, out *v1beta1.MountVolumeRequest) error
// skipping generation of the auto function

func autoConvert_v1beta1_MountVolumeResponse_To_impl_MountVolumeResponse(in *v1beta1.MountVolumeResponse, out *impl.MountVolumeResponse) error {
	return nil
}

// Convert_v1beta1_MountVolumeResponse_To_impl_MountVolumeResponse is an autogenerated conversion function.
func Convert_v1beta1_MountVolumeResponse_To_impl_MountVolumeResponse(in *v1beta1.MountVolumeResponse, out *impl.MountVolumeResponse) error {
	return autoConvert_v1beta1_MountVolumeResponse_To_impl_MountVolumeResponse(in, out)
}

func autoConvert_impl_MountVolumeResponse_To_v1beta1_MountVolumeResponse(in *impl.MountVolumeResponse, out *v1beta1.MountVolumeResponse) error {
	return nil
}

// Convert_impl_MountVolumeResponse_To_v1beta1_MountVolumeResponse is an autogenerated conversion function.
func Convert_impl_MountVolumeResponse_To_v1beta1_MountVolumeResponse(in *impl.MountVolumeResponse, out *v1beta1.MountVolumeResponse) error {
	return autoConvert_impl_MountVolumeResponse_To_v1beta1_MountVolumeResponse(in, out)
}

// detected external conversion function
// Convert_v1beta1_ResizeVolumeRequest_To_impl_ResizeVolumeRequest(in *v1beta1.ResizeVolumeRequest, out *impl.ResizeVolumeRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_ResizeVolumeRequest_To_v1beta1_ResizeVolumeRequest(in *impl.ResizeVolumeRequest, out *v1beta1.ResizeVolumeRequest) error
// skipping generation of the auto function

func autoConvert_v1beta1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in *v1beta1.ResizeVolumeResponse, out *impl.ResizeVolumeResponse) error {
	return nil
}

// Convert_v1beta1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse is an autogenerated conversion function.
func Convert_v1beta1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in *v1beta1.ResizeVolumeResponse, out *impl.ResizeVolumeResponse) error {
	return autoConvert_v1beta1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in, out)
}

func autoConvert_impl_ResizeVolumeResponse_To_v1beta1_ResizeVolumeResponse(in *impl.ResizeVolumeResponse, out *v1beta1.ResizeVolumeResponse) error {
	return nil
}

// Convert_impl_ResizeVolumeResponse_To_v1beta1_ResizeVolumeResponse is an autogenerated conversion function.
func Convert_impl_ResizeVolumeResponse_To_v1beta1_ResizeVolumeResponse(in *impl.ResizeVolumeResponse, out *v1beta1.ResizeVolumeResponse) error {
	return autoConvert_impl_ResizeVolumeResponse_To_v1beta1_ResizeVolumeResponse(in, out)
}

func autoConvert_v1beta1_VolumeDiskNumberRequest_To_impl_VolumeDiskNumberRequest(in *v1beta1.VolumeDiskNumberRequest, out *impl.VolumeDiskNumberRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1beta1_VolumeDiskNumberRequest_To_impl_VolumeDiskNumberRequest is an autogenerated conversion function.
func Convert_v1beta1_VolumeDiskNumberRequest_To_impl_VolumeDiskNumberRequest(in *v1beta1.VolumeDiskNumberRequest, out *impl.VolumeDiskNumberRequest) error {
	return autoConvert_v1beta1_VolumeDiskNumberRequest_To_impl_VolumeDiskNumberRequest(in, out)
}

func autoConvert_impl_VolumeDiskNumberRequest_To_v1beta1_VolumeDiskNumberRequest(in *impl.VolumeDiskNumberRequest, out *v1beta1.VolumeDiskNumberRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_VolumeDiskNumberRequest_To_v1beta1_VolumeDiskNumberRequest is an autogenerated conversion function.
func Convert_impl_VolumeDiskNumberRequest_To_v1beta1_VolumeDiskNumberRequest(in *impl.VolumeDiskNumberRequest, out *v1beta1.VolumeDiskNumberRequest) error {
	return autoConvert_impl_VolumeDiskNumberRequest_To_v1beta1_VolumeDiskNumberRequest(in, out)
}

func autoConvert_v1beta1_VolumeDiskNumberResponse_To_impl_VolumeDiskNumberResponse(in *v1beta1.VolumeDiskNumberResponse, out *impl.VolumeDiskNumberResponse) error {
	out.DiskNumber = in.DiskNumber
	return nil
}

// Convert_v1beta1_VolumeDiskNumberResponse_To_impl_VolumeDiskNumberResponse is an autogenerated conversion function.
func Convert_v1beta1_VolumeDiskNumberResponse_To_impl_VolumeDiskNumberResponse(in *v1beta1.VolumeDiskNumberResponse, out *impl.VolumeDiskNumberResponse) error {
	return autoConvert_v1beta1_VolumeDiskNumberResponse_To_impl_VolumeDiskNumberResponse(in, out)
}

func autoConvert_impl_VolumeDiskNumberResponse_To_v1beta1_VolumeDiskNumberResponse(in *impl.VolumeDiskNumberResponse, out *v1beta1.VolumeDiskNumberResponse) error {
	out.DiskNumber = in.DiskNumber
	return nil
}

// Convert_impl_VolumeDiskNumberResponse_To_v1beta1_VolumeDiskNumberResponse is an autogenerated conversion function.
func Convert_impl_VolumeDiskNumberResponse_To_v1beta1_VolumeDiskNumberResponse(in *impl.VolumeDiskNumberResponse, out *v1beta1.VolumeDiskNumberResponse) error {
	return autoConvert_impl_VolumeDiskNumberResponse_To_v1beta1_VolumeDiskNumberResponse(in, out)
}

func autoConvert_v1beta1_VolumeIDFromMountRequest_To_impl_VolumeIDFromMountRequest(in *v1beta1.VolumeIDFromMountRequest, out *impl.VolumeIDFromMountRequest) error {
	out.Mount = in.Mount
	return nil
}

// Convert_v1beta1_VolumeIDFromMountRequest_To_impl_VolumeIDFromMountRequest is an autogenerated conversion function.
func Convert_v1beta1_VolumeIDFromMountRequest_To_impl_VolumeIDFromMountRequest(in *v1beta1.VolumeIDFromMountRequest, out *impl.VolumeIDFromMountRequest) error {
	return autoConvert_v1beta1_VolumeIDFromMountRequest_To_impl_VolumeIDFromMountRequest(in, out)
}

func autoConvert_impl_VolumeIDFromMountRequest_To_v1beta1_VolumeIDFromMountRequest(in *impl.VolumeIDFromMountRequest, out *v1beta1.VolumeIDFromMountRequest) error {
	out.Mount = in.Mount
	return nil
}

// Convert_impl_VolumeIDFromMountRequest_To_v1beta1_VolumeIDFromMountRequest is an autogenerated conversion function.
func Convert_impl_VolumeIDFromMountRequest_To_v1beta1_VolumeIDFromMountRequest(in *impl.VolumeIDFromMountRequest, out *v1beta1.VolumeIDFromMountRequest) error {
	return autoConvert_impl_VolumeIDFromMountRequest_To_v1beta1_VolumeIDFromMountRequest(in, out)
}

func autoConvert_v1beta1_VolumeIDFromMountResponse_To_impl_VolumeIDFromMountResponse(in *v1beta1.VolumeIDFromMountResponse, out *impl.VolumeIDFromMountResponse) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1beta1_VolumeIDFromMountResponse_To_impl_VolumeIDFromMountResponse is an autogenerated conversion function.
func Convert_v1beta1_VolumeIDFromMountResponse_To_impl_VolumeIDFromMountResponse(in *v1beta1.VolumeIDFromMountResponse, out *impl.VolumeIDFromMountResponse) error {
	return autoConvert_v1beta1_VolumeIDFromMountResponse_To_impl_VolumeIDFromMountResponse(in, out)
}

func autoConvert_impl_VolumeIDFromMountResponse_To_v1beta1_VolumeIDFromMountResponse(in *impl.VolumeIDFromMountResponse, out *v1beta1.VolumeIDFromMountResponse) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_VolumeIDFromMountResponse_To_v1beta1_VolumeIDFromMountResponse is an autogenerated conversion function.
func Convert_impl_VolumeIDFromMountResponse_To_v1beta1_VolumeIDFromMountResponse(in *impl.VolumeIDFromMountResponse, out *v1beta1.VolumeIDFromMountResponse) error {
	return autoConvert_impl_VolumeIDFromMountResponse_To_v1beta1_VolumeIDFromMountResponse(in, out)
}

func autoConvert_v1beta1_VolumeStatsRequest_To_impl_VolumeStatsRequest(in *v1beta1.VolumeStatsRequest, out *impl.VolumeStatsRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1beta1_VolumeStatsRequest_To_impl_VolumeStatsRequest is an autogenerated conversion function.
func Convert_v1beta1_VolumeStatsRequest_To_impl_VolumeStatsRequest(in *v1beta1.VolumeStatsRequest, out *impl.VolumeStatsRequest) error {
	return autoConvert_v1beta1_VolumeStatsRequest_To_impl_VolumeStatsRequest(in, out)
}

func autoConvert_impl_VolumeStatsRequest_To_v1beta1_VolumeStatsRequest(in *impl.VolumeStatsRequest, out *v1beta1.VolumeStatsRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_VolumeStatsRequest_To_v1beta1_VolumeStatsRequest is an autogenerated conversion function.
func Convert_impl_VolumeStatsRequest_To_v1beta1_VolumeStatsRequest(in *impl.VolumeStatsRequest, out *v1beta1.VolumeStatsRequest) error {
	return autoConvert_impl_VolumeStatsRequest_To_v1beta1_VolumeStatsRequest(in, out)
}

func autoConvert_v1beta1_VolumeStatsResponse_To_impl_VolumeStatsResponse(in *v1beta1.VolumeStatsResponse, out *impl.VolumeStatsResponse) error {
	out.VolumeSize = in.VolumeSize
	out.VolumeUsedSize = in.VolumeUsedSize
	return nil
}

// Convert_v1beta1_VolumeStatsResponse_To_impl_VolumeStatsResponse is an autogenerated conversion function.
func Convert_v1beta1_VolumeStatsResponse_To_impl_VolumeStatsResponse(in *v1beta1.VolumeStatsResponse, out *impl.VolumeStatsResponse) error {
	return autoConvert_v1beta1_VolumeStatsResponse_To_impl_VolumeStatsResponse(in, out)
}

func autoConvert_impl_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse(in *impl.VolumeStatsResponse, out *v1beta1.VolumeStatsResponse) error {
	out.VolumeSize = in.VolumeSize
	out.VolumeUsedSize = in.VolumeUsedSize
	return nil
}

// Convert_impl_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse is an autogenerated conversion function.
func Convert_impl_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse(in *impl.VolumeStatsResponse, out *v1beta1.VolumeStatsResponse) error {
	return autoConvert_impl_VolumeStatsResponse_To_v1beta1_VolumeStatsResponse(in, out)
}
