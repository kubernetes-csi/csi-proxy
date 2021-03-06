// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1beta2

import (
	v1beta2 "github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1beta2"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/filesystem/impl"
)

func autoConvert_v1beta2_CreateSymlinkRequest_To_impl_CreateSymlinkRequest(in *v1beta2.CreateSymlinkRequest, out *impl.CreateSymlinkRequest) error {
	out.SourcePath = in.SourcePath
	out.TargetPath = in.TargetPath
	return nil
}

// Convert_v1beta2_CreateSymlinkRequest_To_impl_CreateSymlinkRequest is an autogenerated conversion function.
func Convert_v1beta2_CreateSymlinkRequest_To_impl_CreateSymlinkRequest(in *v1beta2.CreateSymlinkRequest, out *impl.CreateSymlinkRequest) error {
	return autoConvert_v1beta2_CreateSymlinkRequest_To_impl_CreateSymlinkRequest(in, out)
}

func autoConvert_impl_CreateSymlinkRequest_To_v1beta2_CreateSymlinkRequest(in *impl.CreateSymlinkRequest, out *v1beta2.CreateSymlinkRequest) error {
	out.SourcePath = in.SourcePath
	out.TargetPath = in.TargetPath
	return nil
}

// Convert_impl_CreateSymlinkRequest_To_v1beta2_CreateSymlinkRequest is an autogenerated conversion function.
func Convert_impl_CreateSymlinkRequest_To_v1beta2_CreateSymlinkRequest(in *impl.CreateSymlinkRequest, out *v1beta2.CreateSymlinkRequest) error {
	return autoConvert_impl_CreateSymlinkRequest_To_v1beta2_CreateSymlinkRequest(in, out)
}

func autoConvert_v1beta2_CreateSymlinkResponse_To_impl_CreateSymlinkResponse(in *v1beta2.CreateSymlinkResponse, out *impl.CreateSymlinkResponse) error {
	return nil
}

// Convert_v1beta2_CreateSymlinkResponse_To_impl_CreateSymlinkResponse is an autogenerated conversion function.
func Convert_v1beta2_CreateSymlinkResponse_To_impl_CreateSymlinkResponse(in *v1beta2.CreateSymlinkResponse, out *impl.CreateSymlinkResponse) error {
	return autoConvert_v1beta2_CreateSymlinkResponse_To_impl_CreateSymlinkResponse(in, out)
}

func autoConvert_impl_CreateSymlinkResponse_To_v1beta2_CreateSymlinkResponse(in *impl.CreateSymlinkResponse, out *v1beta2.CreateSymlinkResponse) error {
	return nil
}

// Convert_impl_CreateSymlinkResponse_To_v1beta2_CreateSymlinkResponse is an autogenerated conversion function.
func Convert_impl_CreateSymlinkResponse_To_v1beta2_CreateSymlinkResponse(in *impl.CreateSymlinkResponse, out *v1beta2.CreateSymlinkResponse) error {
	return autoConvert_impl_CreateSymlinkResponse_To_v1beta2_CreateSymlinkResponse(in, out)
}

func autoConvert_v1beta2_IsSymlinkRequest_To_impl_IsSymlinkRequest(in *v1beta2.IsSymlinkRequest, out *impl.IsSymlinkRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_v1beta2_IsSymlinkRequest_To_impl_IsSymlinkRequest is an autogenerated conversion function.
func Convert_v1beta2_IsSymlinkRequest_To_impl_IsSymlinkRequest(in *v1beta2.IsSymlinkRequest, out *impl.IsSymlinkRequest) error {
	return autoConvert_v1beta2_IsSymlinkRequest_To_impl_IsSymlinkRequest(in, out)
}

func autoConvert_impl_IsSymlinkRequest_To_v1beta2_IsSymlinkRequest(in *impl.IsSymlinkRequest, out *v1beta2.IsSymlinkRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_impl_IsSymlinkRequest_To_v1beta2_IsSymlinkRequest is an autogenerated conversion function.
func Convert_impl_IsSymlinkRequest_To_v1beta2_IsSymlinkRequest(in *impl.IsSymlinkRequest, out *v1beta2.IsSymlinkRequest) error {
	return autoConvert_impl_IsSymlinkRequest_To_v1beta2_IsSymlinkRequest(in, out)
}

func autoConvert_v1beta2_IsSymlinkResponse_To_impl_IsSymlinkResponse(in *v1beta2.IsSymlinkResponse, out *impl.IsSymlinkResponse) error {
	out.IsSymlink = in.IsSymlink
	return nil
}

// Convert_v1beta2_IsSymlinkResponse_To_impl_IsSymlinkResponse is an autogenerated conversion function.
func Convert_v1beta2_IsSymlinkResponse_To_impl_IsSymlinkResponse(in *v1beta2.IsSymlinkResponse, out *impl.IsSymlinkResponse) error {
	return autoConvert_v1beta2_IsSymlinkResponse_To_impl_IsSymlinkResponse(in, out)
}

func autoConvert_impl_IsSymlinkResponse_To_v1beta2_IsSymlinkResponse(in *impl.IsSymlinkResponse, out *v1beta2.IsSymlinkResponse) error {
	out.IsSymlink = in.IsSymlink
	return nil
}

// Convert_impl_IsSymlinkResponse_To_v1beta2_IsSymlinkResponse is an autogenerated conversion function.
func Convert_impl_IsSymlinkResponse_To_v1beta2_IsSymlinkResponse(in *impl.IsSymlinkResponse, out *v1beta2.IsSymlinkResponse) error {
	return autoConvert_impl_IsSymlinkResponse_To_v1beta2_IsSymlinkResponse(in, out)
}

func autoConvert_v1beta2_MkdirRequest_To_impl_MkdirRequest(in *v1beta2.MkdirRequest, out *impl.MkdirRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_v1beta2_MkdirRequest_To_impl_MkdirRequest is an autogenerated conversion function.
func Convert_v1beta2_MkdirRequest_To_impl_MkdirRequest(in *v1beta2.MkdirRequest, out *impl.MkdirRequest) error {
	return autoConvert_v1beta2_MkdirRequest_To_impl_MkdirRequest(in, out)
}

func autoConvert_impl_MkdirRequest_To_v1beta2_MkdirRequest(in *impl.MkdirRequest, out *v1beta2.MkdirRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_impl_MkdirRequest_To_v1beta2_MkdirRequest is an autogenerated conversion function.
func Convert_impl_MkdirRequest_To_v1beta2_MkdirRequest(in *impl.MkdirRequest, out *v1beta2.MkdirRequest) error {
	return autoConvert_impl_MkdirRequest_To_v1beta2_MkdirRequest(in, out)
}

func autoConvert_v1beta2_MkdirResponse_To_impl_MkdirResponse(in *v1beta2.MkdirResponse, out *impl.MkdirResponse) error {
	return nil
}

// Convert_v1beta2_MkdirResponse_To_impl_MkdirResponse is an autogenerated conversion function.
func Convert_v1beta2_MkdirResponse_To_impl_MkdirResponse(in *v1beta2.MkdirResponse, out *impl.MkdirResponse) error {
	return autoConvert_v1beta2_MkdirResponse_To_impl_MkdirResponse(in, out)
}

func autoConvert_impl_MkdirResponse_To_v1beta2_MkdirResponse(in *impl.MkdirResponse, out *v1beta2.MkdirResponse) error {
	return nil
}

// Convert_impl_MkdirResponse_To_v1beta2_MkdirResponse is an autogenerated conversion function.
func Convert_impl_MkdirResponse_To_v1beta2_MkdirResponse(in *impl.MkdirResponse, out *v1beta2.MkdirResponse) error {
	return autoConvert_impl_MkdirResponse_To_v1beta2_MkdirResponse(in, out)
}

func autoConvert_v1beta2_PathExistsRequest_To_impl_PathExistsRequest(in *v1beta2.PathExistsRequest, out *impl.PathExistsRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_v1beta2_PathExistsRequest_To_impl_PathExistsRequest is an autogenerated conversion function.
func Convert_v1beta2_PathExistsRequest_To_impl_PathExistsRequest(in *v1beta2.PathExistsRequest, out *impl.PathExistsRequest) error {
	return autoConvert_v1beta2_PathExistsRequest_To_impl_PathExistsRequest(in, out)
}

func autoConvert_impl_PathExistsRequest_To_v1beta2_PathExistsRequest(in *impl.PathExistsRequest, out *v1beta2.PathExistsRequest) error {
	out.Path = in.Path
	return nil
}

// Convert_impl_PathExistsRequest_To_v1beta2_PathExistsRequest is an autogenerated conversion function.
func Convert_impl_PathExistsRequest_To_v1beta2_PathExistsRequest(in *impl.PathExistsRequest, out *v1beta2.PathExistsRequest) error {
	return autoConvert_impl_PathExistsRequest_To_v1beta2_PathExistsRequest(in, out)
}

func autoConvert_v1beta2_PathExistsResponse_To_impl_PathExistsResponse(in *v1beta2.PathExistsResponse, out *impl.PathExistsResponse) error {
	out.Exists = in.Exists
	return nil
}

// Convert_v1beta2_PathExistsResponse_To_impl_PathExistsResponse is an autogenerated conversion function.
func Convert_v1beta2_PathExistsResponse_To_impl_PathExistsResponse(in *v1beta2.PathExistsResponse, out *impl.PathExistsResponse) error {
	return autoConvert_v1beta2_PathExistsResponse_To_impl_PathExistsResponse(in, out)
}

func autoConvert_impl_PathExistsResponse_To_v1beta2_PathExistsResponse(in *impl.PathExistsResponse, out *v1beta2.PathExistsResponse) error {
	out.Exists = in.Exists
	return nil
}

// Convert_impl_PathExistsResponse_To_v1beta2_PathExistsResponse is an autogenerated conversion function.
func Convert_impl_PathExistsResponse_To_v1beta2_PathExistsResponse(in *impl.PathExistsResponse, out *v1beta2.PathExistsResponse) error {
	return autoConvert_impl_PathExistsResponse_To_v1beta2_PathExistsResponse(in, out)
}

func autoConvert_v1beta2_RmdirRequest_To_impl_RmdirRequest(in *v1beta2.RmdirRequest, out *impl.RmdirRequest) error {
	out.Path = in.Path
	out.Force = in.Force
	return nil
}

// Convert_v1beta2_RmdirRequest_To_impl_RmdirRequest is an autogenerated conversion function.
func Convert_v1beta2_RmdirRequest_To_impl_RmdirRequest(in *v1beta2.RmdirRequest, out *impl.RmdirRequest) error {
	return autoConvert_v1beta2_RmdirRequest_To_impl_RmdirRequest(in, out)
}

func autoConvert_impl_RmdirRequest_To_v1beta2_RmdirRequest(in *impl.RmdirRequest, out *v1beta2.RmdirRequest) error {
	out.Path = in.Path
	out.Force = in.Force
	return nil
}

// Convert_impl_RmdirRequest_To_v1beta2_RmdirRequest is an autogenerated conversion function.
func Convert_impl_RmdirRequest_To_v1beta2_RmdirRequest(in *impl.RmdirRequest, out *v1beta2.RmdirRequest) error {
	return autoConvert_impl_RmdirRequest_To_v1beta2_RmdirRequest(in, out)
}

func autoConvert_v1beta2_RmdirResponse_To_impl_RmdirResponse(in *v1beta2.RmdirResponse, out *impl.RmdirResponse) error {
	return nil
}

// Convert_v1beta2_RmdirResponse_To_impl_RmdirResponse is an autogenerated conversion function.
func Convert_v1beta2_RmdirResponse_To_impl_RmdirResponse(in *v1beta2.RmdirResponse, out *impl.RmdirResponse) error {
	return autoConvert_v1beta2_RmdirResponse_To_impl_RmdirResponse(in, out)
}

func autoConvert_impl_RmdirResponse_To_v1beta2_RmdirResponse(in *impl.RmdirResponse, out *v1beta2.RmdirResponse) error {
	return nil
}

// Convert_impl_RmdirResponse_To_v1beta2_RmdirResponse is an autogenerated conversion function.
func Convert_impl_RmdirResponse_To_v1beta2_RmdirResponse(in *impl.RmdirResponse, out *v1beta2.RmdirResponse) error {
	return autoConvert_impl_RmdirResponse_To_v1beta2_RmdirResponse(in, out)
}
