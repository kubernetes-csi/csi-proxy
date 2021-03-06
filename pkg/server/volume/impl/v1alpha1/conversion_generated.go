// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	v1alpha1 "github.com/kubernetes-csi/csi-proxy/client/api/volume/v1alpha1"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/volume/impl"
)

func autoConvert_v1alpha1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in *v1alpha1.DismountVolumeRequest, out *impl.DismountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.Path
	return nil
}

// Convert_v1alpha1_DismountVolumeRequest_To_impl_DismountVolumeRequest is an autogenerated conversion function.
func Convert_v1alpha1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in *v1alpha1.DismountVolumeRequest, out *impl.DismountVolumeRequest) error {
	return autoConvert_v1alpha1_DismountVolumeRequest_To_impl_DismountVolumeRequest(in, out)
}

func autoConvert_impl_DismountVolumeRequest_To_v1alpha1_DismountVolumeRequest(in *impl.DismountVolumeRequest, out *v1alpha1.DismountVolumeRequest) error {
	out.VolumeId = in.VolumeId
	out.Path = in.Path
	return nil
}

// Convert_impl_DismountVolumeRequest_To_v1alpha1_DismountVolumeRequest is an autogenerated conversion function.
func Convert_impl_DismountVolumeRequest_To_v1alpha1_DismountVolumeRequest(in *impl.DismountVolumeRequest, out *v1alpha1.DismountVolumeRequest) error {
	return autoConvert_impl_DismountVolumeRequest_To_v1alpha1_DismountVolumeRequest(in, out)
}

func autoConvert_v1alpha1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in *v1alpha1.DismountVolumeResponse, out *impl.DismountVolumeResponse) error {
	return nil
}

// Convert_v1alpha1_DismountVolumeResponse_To_impl_DismountVolumeResponse is an autogenerated conversion function.
func Convert_v1alpha1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in *v1alpha1.DismountVolumeResponse, out *impl.DismountVolumeResponse) error {
	return autoConvert_v1alpha1_DismountVolumeResponse_To_impl_DismountVolumeResponse(in, out)
}

func autoConvert_impl_DismountVolumeResponse_To_v1alpha1_DismountVolumeResponse(in *impl.DismountVolumeResponse, out *v1alpha1.DismountVolumeResponse) error {
	return nil
}

// Convert_impl_DismountVolumeResponse_To_v1alpha1_DismountVolumeResponse is an autogenerated conversion function.
func Convert_impl_DismountVolumeResponse_To_v1alpha1_DismountVolumeResponse(in *impl.DismountVolumeResponse, out *v1alpha1.DismountVolumeResponse) error {
	return autoConvert_impl_DismountVolumeResponse_To_v1alpha1_DismountVolumeResponse(in, out)
}

func autoConvert_v1alpha1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in *v1alpha1.FormatVolumeRequest, out *impl.FormatVolumeRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1alpha1_FormatVolumeRequest_To_impl_FormatVolumeRequest is an autogenerated conversion function.
func Convert_v1alpha1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in *v1alpha1.FormatVolumeRequest, out *impl.FormatVolumeRequest) error {
	return autoConvert_v1alpha1_FormatVolumeRequest_To_impl_FormatVolumeRequest(in, out)
}

func autoConvert_impl_FormatVolumeRequest_To_v1alpha1_FormatVolumeRequest(in *impl.FormatVolumeRequest, out *v1alpha1.FormatVolumeRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_FormatVolumeRequest_To_v1alpha1_FormatVolumeRequest is an autogenerated conversion function.
func Convert_impl_FormatVolumeRequest_To_v1alpha1_FormatVolumeRequest(in *impl.FormatVolumeRequest, out *v1alpha1.FormatVolumeRequest) error {
	return autoConvert_impl_FormatVolumeRequest_To_v1alpha1_FormatVolumeRequest(in, out)
}

func autoConvert_v1alpha1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in *v1alpha1.FormatVolumeResponse, out *impl.FormatVolumeResponse) error {
	return nil
}

// Convert_v1alpha1_FormatVolumeResponse_To_impl_FormatVolumeResponse is an autogenerated conversion function.
func Convert_v1alpha1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in *v1alpha1.FormatVolumeResponse, out *impl.FormatVolumeResponse) error {
	return autoConvert_v1alpha1_FormatVolumeResponse_To_impl_FormatVolumeResponse(in, out)
}

func autoConvert_impl_FormatVolumeResponse_To_v1alpha1_FormatVolumeResponse(in *impl.FormatVolumeResponse, out *v1alpha1.FormatVolumeResponse) error {
	return nil
}

// Convert_impl_FormatVolumeResponse_To_v1alpha1_FormatVolumeResponse is an autogenerated conversion function.
func Convert_impl_FormatVolumeResponse_To_v1alpha1_FormatVolumeResponse(in *impl.FormatVolumeResponse, out *v1alpha1.FormatVolumeResponse) error {
	return autoConvert_impl_FormatVolumeResponse_To_v1alpha1_FormatVolumeResponse(in, out)
}

func autoConvert_v1alpha1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in *v1alpha1.IsVolumeFormattedRequest, out *impl.IsVolumeFormattedRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_v1alpha1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest is an autogenerated conversion function.
func Convert_v1alpha1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in *v1alpha1.IsVolumeFormattedRequest, out *impl.IsVolumeFormattedRequest) error {
	return autoConvert_v1alpha1_IsVolumeFormattedRequest_To_impl_IsVolumeFormattedRequest(in, out)
}

func autoConvert_impl_IsVolumeFormattedRequest_To_v1alpha1_IsVolumeFormattedRequest(in *impl.IsVolumeFormattedRequest, out *v1alpha1.IsVolumeFormattedRequest) error {
	out.VolumeId = in.VolumeId
	return nil
}

// Convert_impl_IsVolumeFormattedRequest_To_v1alpha1_IsVolumeFormattedRequest is an autogenerated conversion function.
func Convert_impl_IsVolumeFormattedRequest_To_v1alpha1_IsVolumeFormattedRequest(in *impl.IsVolumeFormattedRequest, out *v1alpha1.IsVolumeFormattedRequest) error {
	return autoConvert_impl_IsVolumeFormattedRequest_To_v1alpha1_IsVolumeFormattedRequest(in, out)
}

func autoConvert_v1alpha1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in *v1alpha1.IsVolumeFormattedResponse, out *impl.IsVolumeFormattedResponse) error {
	out.Formatted = in.Formatted
	return nil
}

// Convert_v1alpha1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse is an autogenerated conversion function.
func Convert_v1alpha1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in *v1alpha1.IsVolumeFormattedResponse, out *impl.IsVolumeFormattedResponse) error {
	return autoConvert_v1alpha1_IsVolumeFormattedResponse_To_impl_IsVolumeFormattedResponse(in, out)
}

func autoConvert_impl_IsVolumeFormattedResponse_To_v1alpha1_IsVolumeFormattedResponse(in *impl.IsVolumeFormattedResponse, out *v1alpha1.IsVolumeFormattedResponse) error {
	out.Formatted = in.Formatted
	return nil
}

// Convert_impl_IsVolumeFormattedResponse_To_v1alpha1_IsVolumeFormattedResponse is an autogenerated conversion function.
func Convert_impl_IsVolumeFormattedResponse_To_v1alpha1_IsVolumeFormattedResponse(in *impl.IsVolumeFormattedResponse, out *v1alpha1.IsVolumeFormattedResponse) error {
	return autoConvert_impl_IsVolumeFormattedResponse_To_v1alpha1_IsVolumeFormattedResponse(in, out)
}

// detected external conversion function
// Convert_v1alpha1_ListVolumesOnDiskRequest_To_impl_ListVolumesOnDiskRequest(in *v1alpha1.ListVolumesOnDiskRequest, out *impl.ListVolumesOnDiskRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_ListVolumesOnDiskRequest_To_v1alpha1_ListVolumesOnDiskRequest(in *impl.ListVolumesOnDiskRequest, out *v1alpha1.ListVolumesOnDiskRequest) error
// skipping generation of the auto function

func autoConvert_v1alpha1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in *v1alpha1.ListVolumesOnDiskResponse, out *impl.ListVolumesOnDiskResponse) error {
	out.VolumeIds = *(*[]string)(unsafe.Pointer(&in.VolumeIds))
	return nil
}

// Convert_v1alpha1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse is an autogenerated conversion function.
func Convert_v1alpha1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in *v1alpha1.ListVolumesOnDiskResponse, out *impl.ListVolumesOnDiskResponse) error {
	return autoConvert_v1alpha1_ListVolumesOnDiskResponse_To_impl_ListVolumesOnDiskResponse(in, out)
}

func autoConvert_impl_ListVolumesOnDiskResponse_To_v1alpha1_ListVolumesOnDiskResponse(in *impl.ListVolumesOnDiskResponse, out *v1alpha1.ListVolumesOnDiskResponse) error {
	out.VolumeIds = *(*[]string)(unsafe.Pointer(&in.VolumeIds))
	return nil
}

// Convert_impl_ListVolumesOnDiskResponse_To_v1alpha1_ListVolumesOnDiskResponse is an autogenerated conversion function.
func Convert_impl_ListVolumesOnDiskResponse_To_v1alpha1_ListVolumesOnDiskResponse(in *impl.ListVolumesOnDiskResponse, out *v1alpha1.ListVolumesOnDiskResponse) error {
	return autoConvert_impl_ListVolumesOnDiskResponse_To_v1alpha1_ListVolumesOnDiskResponse(in, out)
}

// detected external conversion function
// Convert_v1alpha1_MountVolumeRequest_To_impl_MountVolumeRequest(in *v1alpha1.MountVolumeRequest, out *impl.MountVolumeRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_MountVolumeRequest_To_v1alpha1_MountVolumeRequest(in *impl.MountVolumeRequest, out *v1alpha1.MountVolumeRequest) error
// skipping generation of the auto function

func autoConvert_v1alpha1_MountVolumeResponse_To_impl_MountVolumeResponse(in *v1alpha1.MountVolumeResponse, out *impl.MountVolumeResponse) error {
	return nil
}

// Convert_v1alpha1_MountVolumeResponse_To_impl_MountVolumeResponse is an autogenerated conversion function.
func Convert_v1alpha1_MountVolumeResponse_To_impl_MountVolumeResponse(in *v1alpha1.MountVolumeResponse, out *impl.MountVolumeResponse) error {
	return autoConvert_v1alpha1_MountVolumeResponse_To_impl_MountVolumeResponse(in, out)
}

func autoConvert_impl_MountVolumeResponse_To_v1alpha1_MountVolumeResponse(in *impl.MountVolumeResponse, out *v1alpha1.MountVolumeResponse) error {
	return nil
}

// Convert_impl_MountVolumeResponse_To_v1alpha1_MountVolumeResponse is an autogenerated conversion function.
func Convert_impl_MountVolumeResponse_To_v1alpha1_MountVolumeResponse(in *impl.MountVolumeResponse, out *v1alpha1.MountVolumeResponse) error {
	return autoConvert_impl_MountVolumeResponse_To_v1alpha1_MountVolumeResponse(in, out)
}

// detected external conversion function
// Convert_v1alpha1_ResizeVolumeRequest_To_impl_ResizeVolumeRequest(in *v1alpha1.ResizeVolumeRequest, out *impl.ResizeVolumeRequest) error
// skipping generation of the auto function

// detected external conversion function
// Convert_impl_ResizeVolumeRequest_To_v1alpha1_ResizeVolumeRequest(in *impl.ResizeVolumeRequest, out *v1alpha1.ResizeVolumeRequest) error
// skipping generation of the auto function

func autoConvert_v1alpha1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in *v1alpha1.ResizeVolumeResponse, out *impl.ResizeVolumeResponse) error {
	return nil
}

// Convert_v1alpha1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse is an autogenerated conversion function.
func Convert_v1alpha1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in *v1alpha1.ResizeVolumeResponse, out *impl.ResizeVolumeResponse) error {
	return autoConvert_v1alpha1_ResizeVolumeResponse_To_impl_ResizeVolumeResponse(in, out)
}

func autoConvert_impl_ResizeVolumeResponse_To_v1alpha1_ResizeVolumeResponse(in *impl.ResizeVolumeResponse, out *v1alpha1.ResizeVolumeResponse) error {
	return nil
}

// Convert_impl_ResizeVolumeResponse_To_v1alpha1_ResizeVolumeResponse is an autogenerated conversion function.
func Convert_impl_ResizeVolumeResponse_To_v1alpha1_ResizeVolumeResponse(in *impl.ResizeVolumeResponse, out *v1alpha1.ResizeVolumeResponse) error {
	return autoConvert_impl_ResizeVolumeResponse_To_v1alpha1_ResizeVolumeResponse(in, out)
}
