package main

import (
	"github.com/kubernetes-csi/csi-proxy/internal/server"
	filesystem "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
)

func main() {
	s := server.NewServer(apiGroups()...)
	if err := s.Start(nil); err != nil {
		panic(err)
	}
}

// apiGroups returns the list of enabled API groups.
func apiGroups() []server.APIGroup {
	return []server.APIGroup{
		&filesystem.Server{},
	}
}
