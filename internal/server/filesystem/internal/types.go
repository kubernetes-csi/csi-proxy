package internal

import (
	"github.com/kubernetes-csi/csi-proxy/client/api"
)

// PathExistsRequest is the internal representation of requests to the PathExists endpoint.
type PathExistsRequest struct {
	// The path to check in the host filesystem.
	Path string
}

// PathExistsResponse is the internal representation of responses from the PathExists endpoint.
type PathExistsResponse struct {
	Success     bool
	CmdletError *api.CmdletError
	Exists      bool
}
