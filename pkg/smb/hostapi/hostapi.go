package api

import (
	"fmt"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/v2/pkg/utils"
)

type HostAPI interface {
	IsSMBMapped(remotePath string) (bool, error)
	NewSMBLink(remotePath, localPath string) error
	NewSMBGlobalMapping(remotePath, username, password string) error
	RemoveSMBGlobalMapping(remotePath string) error
}

type smbAPI struct{}

var _ HostAPI = &smbAPI{}

func New() HostAPI {
	return smbAPI{}
}

func remotePathForQuery(remotePath string) string {
	return strings.ReplaceAll(remotePath, "\\", "\\\\")
}

func (smbAPI) IsSMBMapped(remotePath string) (bool, error) {
	var isMapped bool
	err := cim.WithCOMThread(func() error {
		inst, err := cim.QuerySmbGlobalMappingByRemotePath(remotePath)
		if err != nil {
			return err
		}

		status, err := cim.GetSmbGlobalMappingStatus(inst)
		if err != nil {
			return err
		}

		isMapped = status == cim.SmbMappingStatusOK
		return nil
	})
	return isMapped, cim.IgnoreNotFound(err)
}

// NewSMBLink - creates a directory symbolic link to the remote share.
// The os.Symlink was having issue for cases where the destination was an SMB share - the container
// runtime would complain stating "Access Denied".
func (smbAPI) NewSMBLink(remotePath, localPath string) error {
	if !strings.HasSuffix(remotePath, "\\") {
		// Golang has issues resolving paths mapped to file shares if they do not end in a trailing \
		// so add one if needed.
		remotePath = remotePath + "\\"
	}
	longRemotePath := utils.EnsureLongPath(remotePath)
	longLocalPath := utils.EnsureLongPath(localPath)

	err := utils.CreateSymlink(longLocalPath, longRemotePath, true)
	if err != nil {
		return fmt.Errorf("error linking %s to %s. err: %v", remotePath, localPath, err)
	}

	return nil
}

func (smbAPI) NewSMBGlobalMapping(remotePath, username, password string) error {
	return cim.WithCOMThread(func() error {
		result, err := cim.NewSmbGlobalMapping(remotePath, username, password, true)
		if err != nil {
			return fmt.Errorf("NewSmbGlobalMapping failed. result: %d, err: %v", result, err)
		}
		return nil
	})
}

func (smbAPI) RemoveSMBGlobalMapping(remotePath string) error {
	return cim.WithCOMThread(func() error {
		err := cim.RemoveSmbGlobalMappingByRemotePath(remotePath)
		if err != nil {
			return fmt.Errorf("error remove smb mapping '%s'. err: %v", remotePath, err)
		}
		return nil
	})
}
