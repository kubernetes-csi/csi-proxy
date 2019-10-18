package filesystem

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/api"
	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/filesystem/internal"
)

// Server is the internal server for the filesytem API group.
type Server struct{}

// PathExists checks if the given path exists on the host.
func (s *Server) PathExists(ctx context.Context, request *internal.PathExistsRequest, version apiversion.Version) (*internal.PathExistsResponse, error) {
	// FIXME: actually implement this!
	return &internal.PathExistsResponse{
		Success: false,
		CmdletError: &api.CmdletError{
			CmdletName: "dummy",
			Code:       12,
			Message:    "hey there " + request.Path,
		},
	}, nil
}
