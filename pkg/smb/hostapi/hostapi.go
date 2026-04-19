package api

import (
	"fmt"
	"strings"

	wmi "github.com/kubernetes-csi/csi-proxy/v2/pkg/cim"
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

func (smbAPI) IsSMBMapped(remotePath string) (bool, error) {
	var isMapped bool
	err := wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			inst, err := wmi.QuerySmbGlobalMappingByRemotePath(scope, remotePath)
			if err != nil {
				return err
			}

			status, err := wmi.GetSmbGlobalMappingStatus(inst)
			if err != nil {
				return err
			}

			isMapped = status == wmi.SmbMappingStatusOK
			return nil
		})
	})
	return isMapped, wmi.IgnoreNotFound(err)
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
		return fmt.Errorf("error linking %s to %s. err: %w", remotePath, localPath, err)
	}

	return nil
}

func (smbAPI) NewSMBGlobalMapping(remotePath, username, password string) error {
	requirePrivacy := true
	return wmi.WithCOMThread(func() error {
		err := wmi.NewSmbGlobalMapping(remotePath, username, password, requirePrivacy)
		if err != nil {
			return fmt.Errorf("NewSmbGlobalMapping failed. err: %w", err)
		}
		return nil
	})
}

func (smbAPI) RemoveSMBGlobalMapping(remotePath string) error {
	return wmi.WithCOMThread(func() error {
		return wmi.WithScope(func(scope *wmi.Scope) error {
			err := wmi.RemoveSmbGlobalMappingByRemotePath(scope, remotePath)
			if err != nil {
				return fmt.Errorf("error remove smb mapping '%s'. err: %w", remotePath, err)
			}
			return nil
		})
	})
}
