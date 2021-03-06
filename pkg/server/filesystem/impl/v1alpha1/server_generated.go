// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/filesystem/impl"
	"google.golang.org/grpc"
)

var version = apiversion.NewVersionOrPanic("v1alpha1")

type versionedAPI struct {
	apiGroupServer impl.ServerInterface
}

func NewVersionedServer(apiGroupServer impl.ServerInterface) impl.VersionedAPI {
	return &versionedAPI{
		apiGroupServer: apiGroupServer,
	}
}

func (s *versionedAPI) Register(grpcServer *grpc.Server) {
	v1alpha1.RegisterFilesystemServer(grpcServer, s)
}

func (s *versionedAPI) IsMountPoint(context context.Context, versionedRequest *v1alpha1.IsMountPointRequest) (*v1alpha1.IsMountPointResponse, error) {
	request := &impl.IsMountPointRequest{}
	if err := Convert_v1alpha1_IsMountPointRequest_To_impl_IsMountPointRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.IsMountPoint(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha1.IsMountPointResponse{}
	if err := Convert_impl_IsMountPointResponse_To_v1alpha1_IsMountPointResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}

func (s *versionedAPI) LinkPath(context context.Context, versionedRequest *v1alpha1.LinkPathRequest) (*v1alpha1.LinkPathResponse, error) {
	request := &impl.LinkPathRequest{}
	if err := Convert_v1alpha1_LinkPathRequest_To_impl_LinkPathRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.LinkPath(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha1.LinkPathResponse{}
	if err := Convert_impl_LinkPathResponse_To_v1alpha1_LinkPathResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}

func (s *versionedAPI) Mkdir(context context.Context, versionedRequest *v1alpha1.MkdirRequest) (*v1alpha1.MkdirResponse, error) {
	request := &impl.MkdirRequest{}
	if err := Convert_v1alpha1_MkdirRequest_To_impl_MkdirRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.Mkdir(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha1.MkdirResponse{}
	if err := Convert_impl_MkdirResponse_To_v1alpha1_MkdirResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}

func (s *versionedAPI) PathExists(context context.Context, versionedRequest *v1alpha1.PathExistsRequest) (*v1alpha1.PathExistsResponse, error) {
	request := &impl.PathExistsRequest{}
	if err := Convert_v1alpha1_PathExistsRequest_To_impl_PathExistsRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.PathExists(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha1.PathExistsResponse{}
	if err := Convert_impl_PathExistsResponse_To_v1alpha1_PathExistsResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}

func (s *versionedAPI) Rmdir(context context.Context, versionedRequest *v1alpha1.RmdirRequest) (*v1alpha1.RmdirResponse, error) {
	request := &impl.RmdirRequest{}
	if err := Convert_v1alpha1_RmdirRequest_To_impl_RmdirRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.Rmdir(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha1.RmdirResponse{}
	if err := Convert_impl_RmdirResponse_To_v1alpha1_RmdirResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}
