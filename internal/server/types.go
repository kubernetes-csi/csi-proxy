package server

import (
	"google.golang.org/grpc"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
)

type VersionedAPI struct {
	Group      string
	Version    apiversion.Version
	Registrant func(*grpc.Server)
}

type APIGroup interface {
	VersionedAPIs() []*VersionedAPI
}
