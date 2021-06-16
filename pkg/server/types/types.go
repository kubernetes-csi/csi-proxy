package types

import (
	"google.golang.org/grpc"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
)

// VersionedAPI represents an API group & version.
type VersionedAPI struct {
	Group      string
	Version    apiversion.Version
	Registrant func(*grpc.Server)
}

// APIGroup represents an API group.
type APIGroup interface {
	VersionedAPIs() []*VersionedAPI
}
