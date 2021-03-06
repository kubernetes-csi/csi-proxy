// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package impl

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"google.golang.org/grpc"
)

type VersionedAPI interface {
	Register(grpcServer *grpc.Server)
}

// All the functions this group's server needs to define.
type ServerInterface interface {
	NewSmbGlobalMapping(context.Context, *NewSmbGlobalMappingRequest, apiversion.Version) (*NewSmbGlobalMappingResponse, error)
	RemoveSmbGlobalMapping(context.Context, *RemoveSmbGlobalMappingRequest, apiversion.Version) (*RemoveSmbGlobalMappingResponse, error)
}
