package smb

import (
	"context"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	fsserver "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
	"k8s.io/klog"
)

type Server struct {
	hostAPI  API
	fsServer *fsserver.Server
}

type API interface {
	IsSmbMapped(remotePath string) (bool, error)
	SMBLink(remotePath, localPath string) error
	NewSmbGlobalMapping(remotePath, localPath, username, password string) error
	RemoveSmbGlobalMapping(remotePath string) error
}

func NewServer(hostAPI API, fsServer *fsserver.Server) (*Server, error) {
	return &Server{
		hostAPI:  hostAPI,
		fsServer: fsServer,
	}, nil
}

func (s *Server) NewSmbGlobalMapping(context context.Context, request *internal.NewSmbGlobalMappingRequest, version apiversion.Version) (*internal.NewSmbGlobalMappingResponse, error) {
	response := &internal.NewSmbGlobalMappingResponse{}
	remotePath := request.RemotePath
	localPath := request.LocalPath

	if remotePath == "" {
		return response, fmt.Errorf("remote path is empty")
	}

	isMapped, err := s.hostAPI.IsSmbMapped(remotePath)
	if err != nil {
		isMapped = false
	}

	if !isMapped {
		klog.V(4).Infof("Remote %s not mapped. Mapping now!", remotePath)
		err := s.hostAPI.NewSmbGlobalMapping(remotePath, localPath, request.Username, request.Password)
		if err != nil {
			return response, err
		}
	}

	if len(localPath) != 0 {
		err = s.fsServer.ValidateSMBLinkPath(localPath)
		if err != nil {
			return response, err
		}
		err = s.hostAPI.SMBLink(remotePath, localPath)
		if err != nil {
			return response, fmt.Errorf("creating link %s to %s failed with error: %v", localPath, remotePath, err)
		}
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
