package smb

import (
	"fmt"
	"strings"

	wmi "github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/utils"
)

type API interface {
	IsSmbMapped(remotePath string) (bool, error)
	NewSmbLink(remotePath, localPath string) error
	NewSmbGlobalMapping(remotePath, username, password string) error
	RemoveSmbGlobalMapping(remotePath string) error
}

type SmbAPI struct {
	RequirePrivacy bool
}

var _ API = &SmbAPI{}

func New(requirePrivacy bool) *SmbAPI {
	return &SmbAPI{
		RequirePrivacy: requirePrivacy,
	}
}

func (*SmbAPI) IsSmbMapped(remotePath string) (bool, error) {
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

// NewSmbLink - creates a directory symbolic link to the remote share.
// The os.Symlink was having issue for cases where the destination was an SMB share - the container
// runtime would complain stating "Access Denied".
func (*SmbAPI) NewSmbLink(remotePath, localPath string) error {
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

func (api *SmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	return wmi.WithCOMThread(func() error {
		err := wmi.NewSmbGlobalMapping(remotePath, username, password, api.RequirePrivacy)
		if err != nil {
			return fmt.Errorf("NewSmbGlobalMapping failed. err: %w", err)
		}
		return nil
	})
}

func (*SmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
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
