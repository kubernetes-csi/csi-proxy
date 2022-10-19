package smb

import (
	"context"
	"fmt"
	"strings"

	fs "github.com/kubernetes-csi/csi-proxy/pkg/filesystem"
	smbapi "github.com/kubernetes-csi/csi-proxy/pkg/smb/hostapi"
	"k8s.io/klog/v2"
)

type SMB struct {
	hostAPI smbapi.HostAPI
	fs      fs.Interface
}

type Interface interface {
	// NewSMBGlobalMapping creates an SMB mapping on the SMB client to an SMB share.
	NewSMBGlobalMapping(context.Context, *NewSMBGlobalMappingRequest) (*NewSMBGlobalMappingResponse, error)

	// RemoveSMBGlobalMapping removes the SMB mapping to an SMB share.
	RemoveSMBGlobalMapping(context.Context, *RemoveSMBGlobalMappingRequest) (*RemoveSMBGlobalMappingResponse, error)
}

// check that SMB implements the Interface
var _ Interface = &SMB{}

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
	// parts[0] is a SMB host name
	// parts[1] is a SMB share name
	return strings.ToLower("\\\\" + parts[0] + "\\" + parts[1]), nil
}

func New(hostAPI smbapi.HostAPI, fsClient fs.Interface) (*SMB, error) {
	return &SMB{
		hostAPI: hostAPI,
		fs:      fsClient,
	}, nil
}

func (s *SMB) NewSMBGlobalMapping(context context.Context, request *NewSMBGlobalMappingRequest) (*NewSMBGlobalMappingResponse, error) {
	klog.V(2).Infof("calling NewSMBGlobalMapping with remote path %q", request.RemotePath)
	response := &NewSMBGlobalMappingResponse{}
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

	isMapped, err := s.hostAPI.IsSMBMapped(mappingPath)
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
			err := s.hostAPI.RemoveSMBGlobalMapping(mappingPath)
			if err != nil {
				klog.Errorf("RemoveSMBGlobalMapping(%s) failed with %v", mappingPath, err)
				return response, err
			}
			isMapped = false
		} else {
			klog.V(4).Infof("RemotePath %s is valid", mappingPath)
		}
	}

	if !isMapped {
		klog.V(4).Infof("Remote %s not mapped. Mapping now!", mappingPath)
		err = s.hostAPI.NewSMBGlobalMapping(mappingPath, request.Username, request.Password)
		if err != nil {
			klog.Errorf("failed NewSMBGlobalMapping %v", err)
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
		err = s.hostAPI.NewSMBLink(remotePath, localPath)
		if err != nil {
			klog.Errorf("failed NewSMBLink %v", err)
			return response, fmt.Errorf("creating link %s to %s failed with error: %v", localPath, remotePath, err)
		}
	}

	klog.V(2).Infof("NewSMBGlobalMapping on remote path %q is completed", request.RemotePath)
	return response, nil
}

func (s *SMB) RemoveSMBGlobalMapping(context context.Context, request *RemoveSMBGlobalMappingRequest) (*RemoveSMBGlobalMappingResponse, error) {
	klog.V(2).Infof("calling RemoveSMBGlobalMapping with remote path %q", request.RemotePath)
	response := &RemoveSMBGlobalMappingResponse{}
	remotePath := normalizeWindowsPath(request.RemotePath)

	if remotePath == "" {
		klog.Errorf("remote path is empty")
		return response, fmt.Errorf("remote path is empty")
	}

	mappingPath, err := getRootMappingPath(remotePath)
	if err != nil {
		return response, err
	}

	err = s.hostAPI.RemoveSMBGlobalMapping(mappingPath)
	if err != nil {
		klog.Errorf("failed RemoveSMBGlobalMapping %v", err)
		return response, err
	}

	klog.V(2).Infof("RemoveSMBGlobalMapping on remote path %q is completed", request.RemotePath)
	return response, nil
}
