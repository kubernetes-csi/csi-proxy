package smb

import (
	"context"
	//"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
)

type Server struct{
	hostAPI API
}

type API interface {
	MountSmbShare(remotePath, localPath, username, password string, readOnly bool) error
	UnmountSmbShare(remotePath, localPath string) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) MountSmbShare(context context.Context, request *internal.MountSmbShareRequest, version apiversion.Version) (*internal.MountSmbShareResponse, error) {
	response := &internal.MountSmbShareResponse{}
	remotePath := request.RemotePath

	if remotePath == "" {
		response.Error = "remote path is empty"
	}

	localPath := request.LocalPath
	if localPath == "" {
		response.Error = "local path is empty"
	}                                                                                                                                             

	err := s.hostAPI.MountSmbShare(remotePath, localPath, request.Username, request.Password, request.ReadOnly)
	if err != nil {
		response.Error = err.Error()
	}
	return response, nil
}


func (s *Server) UnmountSmbShare(context context.Context, request *internal.UnmountSmbShareRequest, version apiversion.Version) (*internal.UnmountSmbShareResponse, error) {
	response := &internal.UnmountSmbShareResponse{}
	remotePath := request.RemotePath

	if remotePath == "" {
		response.Error = "remote path is empty"
	}

	localPath := request.LocalPath
	if localPath == "" {
		response.Error = "local path is empty"
	}                                                                                                                                             

	err := s.hostAPI.UnmountSmbShare(remotePath, localPath)
	if err != nil {
		response.Error = err.Error()
	}
	return response, nil
}