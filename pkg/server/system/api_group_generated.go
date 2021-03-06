// Code generated by csi-proxy-api-gen. DO NOT EDIT.

package system

import (
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/system/impl/v1alpha1"
	srvtypes "github.com/kubernetes-csi/csi-proxy/pkg/server/types"
)

const name = "system"

// ensure the server defines all the required methods
var _ impl.ServerInterface = &Server{}

func (s *Server) VersionedAPIs() []*srvtypes.VersionedAPI {
	v1alpha1Server := v1alpha1.NewVersionedServer(s)

	return []*srvtypes.VersionedAPI{
		{
			Group:      name,
			Version:    apiversion.NewVersionOrPanic("v1alpha1"),
			Registrant: v1alpha1Server.Register,
		},
	}
}
