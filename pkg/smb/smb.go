package smb

import (
	"context"
	"fmt"
	"strings"

	fs "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem"
	smbapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/smb/api"
	"k8s.io/klog/v2"
)

type Smb struct {
	hostAPI smbapi.API
	fs      fs.Interface
}

type Interface interface {
	NewSmbGlobalMapping(context.Context, *NewSmbGlobalMappingRequest) (*NewSmbGlobalMappingResponse, error)
	RemoveSmbGlobalMapping(context.Context, *RemoveSmbGlobalMappingRequest) (*RemoveSmbGlobalMappingResponse, error)
}

// check that Smb implements the Interface
var _ Interface = &Smb{}

func normalizeWindowsPath(path string) string {
	normalizedPath := strings.Replace(path, "/", "\\", -1)
	return normalizedPath
}

func getRootMappingPath(path string) (string, error) {
	items := strings.Split(path, "\\")
	parts := []string{}
	for _, s := range items {
		if len(s) > 0 {
			parts = append(parts, s)
			if len(parts) == 2 {
				break
			}
		}
	}
	if len(parts) != 2 {
		klog.Errorf("remote path (%s) is invalid", path)
		return "", fmt.Errorf("remote path (%s) is invalid", path)
	}
	// parts[0] is a smb host name
	// parts[1] is a smb share name
	return strings.ToLower("\\\\" + parts[0] + "\\" + parts[1]), nil
}

func New(hostAPI smbapi.API, fsClient fs.Interface) (*Smb, error) {
	return &Smb{
		hostAPI: hostAPI,
		fs:      fsClient,
	}, nil
}

func (s *Smb) NewSmbGlobalMapping(context context.Context, request *NewSmbGlobalMappingRequest) (*NewSmbGlobalMappingResponse, error) {
	klog.V(2).Infof("calling NewSmbGlobalMapping with remote path %q", request.RemotePath)
	response := &NewSmbGlobalMappingResponse{}
	remotePath := normalizeWindowsPath(request.RemotePath)
	localPath := request.LocalPath

	if remotePath == "" {
		klog.Errorf("remote path is empty")
		return response, fmt.Errorf("remote path is empty")
	}

	mappingPath, err := getRootMappingPath(remotePath)
	if err != nil {
		return response, err
	}

	isMapped, err := s.hostAPI.IsSmbMapped(mappingPath)
	if err != nil {
		isMapped = false
	}

	if isMapped {
		klog.V(4).Infof("Remote %s already mapped. Validating...", mappingPath)

		validResp, err := s.fs.PathValid(context, &fs.PathValidRequest{Path: mappingPath})
		if err != nil {
			klog.Warningf("PathValid(%s) failed with %v, ignore error", mappingPath, err)
		}

		if !validResp.Valid {
			klog.V(4).Infof("RemotePath %s is not valid, removing now", mappingPath)
			err := s.hostAPI.RemoveSmbGlobalMapping(mappingPath)
			if err != nil {
				klog.Errorf("RemoveSmbGlobalMapping(%s) failed with %v", mappingPath, err)
				return response, err
			}
			isMapped = false
		} else {
			klog.V(4).Infof("RemotePath %s is valid", mappingPath)
		}
	}

	if !isMapped {
		klog.V(4).Infof("Remote %s not mapped. Mapping now!", mappingPath)
		err = s.hostAPI.NewSmbGlobalMapping(mappingPath, request.Username, request.Password)
		if err != nil {
			klog.Errorf("failed NewSmbGlobalMapping %v", err)
			return response, err
		}
	}

	if len(localPath) != 0 {
		klog.V(4).Infof("ValidatePathWindows: '%s'", localPath)
		err = fs.ValidatePathWindows(localPath)
		if err != nil {
			klog.Errorf("failed validate plugin path %v", err)
			return response, err
		}
		err = s.hostAPI.NewSmbLink(remotePath, localPath)
		if err != nil {
			klog.Errorf("failed NewSmbLink %v", err)
			return response, fmt.Errorf("creating link %s to %s failed with error: %v", localPath, remotePath, err)
		}
	}

	klog.V(2).Infof("NewSmbGlobalMapping on remote path %q is completed", request.RemotePath)
	return response, nil
}

func (s *Smb) RemoveSmbGlobalMapping(context context.Context, request *RemoveSmbGlobalMappingRequest) (*RemoveSmbGlobalMappingResponse, error) {
	klog.V(2).Infof("calling RemoveSmbGlobalMapping with remote path %q", request.RemotePath)
	response := &RemoveSmbGlobalMappingResponse{}
	remotePath := normalizeWindowsPath(request.RemotePath)

	if remotePath == "" {
		klog.Errorf("remote path is empty")
		return response, fmt.Errorf("remote path is empty")
	}

	mappingPath, err := getRootMappingPath(remotePath)
	if err != nil {
		return response, err
	}

	err = s.hostAPI.RemoveSmbGlobalMapping(mappingPath)
	if err != nil {
		klog.Errorf("failed RemoveSmbGlobalMapping %v", err)
		return response, err
	}

	klog.V(2).Infof("RemoveSmbGlobalMapping on remote path %q is completed", request.RemotePath)
	return response, nil
}
