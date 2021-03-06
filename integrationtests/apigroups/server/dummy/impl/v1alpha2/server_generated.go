// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package v1alpha2

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/api/dummy/v1alpha2"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/server/dummy/impl"
	"google.golang.org/grpc"
)

var version = apiversion.NewVersionOrPanic("v1alpha2")

type versionedAPI struct {
	apiGroupServer impl.ServerInterface
}

func NewVersionedServer(apiGroupServer impl.ServerInterface) impl.VersionedAPI {
	return &versionedAPI{
		apiGroupServer: apiGroupServer,
	}
}

func (s *versionedAPI) Register(grpcServer *grpc.Server) {
	v1alpha2.RegisterDummyServer(grpcServer, s)
}

func (s *versionedAPI) ComputeDouble(context context.Context, versionedRequest *v1alpha2.ComputeDoubleRequest) (*v1alpha2.ComputeDoubleResponse, error) {
	request := &impl.ComputeDoubleRequest{}
	if err := Convert_v1alpha2_ComputeDoubleRequest_To_impl_ComputeDoubleRequest(versionedRequest, request); err != nil {
		return nil, err
	}

	response, err := s.apiGroupServer.ComputeDouble(context, request, version)
	if err != nil {
		return nil, err
	}

	versionedResponse := &v1alpha2.ComputeDoubleResponse{}
	if err := Convert_impl_ComputeDoubleResponse_To_v1alpha2_ComputeDoubleResponse(response, versionedResponse); err != nil {
		return nil, err
	}

	return versionedResponse, err
}
