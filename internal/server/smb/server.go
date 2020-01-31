package smb

import (
	"context"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
)

type Server struct {
	hostAPI API
}

type API interface {
	NewSmbGlobalMapping(remotePath, username, password string) error
	RemoveSmbGlobalMapping(remotePath string) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) NewSmbGlobalMapping(context context.Context, request *internal.NewSmbGlobalMappingRequest, version apiversion.Version) (*internal.NewSmbGlobalMappingResponse, error) {
	response := &internal.NewSmbGlobalMappingResponse{}
	remotePath := request.RemotePath

	if remotePath == "" {
		return response, fmt.Errorf("remote path is empty")
	}

	err := s.hostAPI.NewSmbGlobalMapping(remotePath, request.Username, request.Password)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Server) RemoveSmbGlobalMapping(context context.Context, request *internal.RemoveSmbGlobalMappingRequest, version apiversion.Version) (*internal.RemoveSmbGlobalMappingResponse, error) {
	response := &internal.RemoveSmbGlobalMappingResponse{}
	remotePath := request.RemotePath

	if remotePath == "" {
		return response, fmt.Errorf("remote path is empty")
	}

	err := s.hostAPI.RemoveSmbGlobalMapping(remotePath)
	if err != nil {
		return response, err
	}
	return response, nil
}
