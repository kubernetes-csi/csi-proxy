package smb

import (
	"fmt"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
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
	inst, err := cim.QuerySmbGlobalMappingByRemotePath(remotePath)
	if err != nil {
		return false, cim.IgnoreNotFound(err)
	}

	status, err := cim.GetSmbGlobalMappingStatus(inst)
	if err != nil {
		return false, err
	}

	return status == cim.SmbMappingStatusOK, nil
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
		return fmt.Errorf("error linking %s to %s. err: %v", remotePath, localPath, err)
	}

	return nil
}

func (api *SmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	result, err := cim.NewSmbGlobalMapping(remotePath, username, password, api.RequirePrivacy)
	if err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed. result: %d, err: %v", result, err)
	}

	return nil
}

func (*SmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	err := cim.RemoveSmbGlobalMappingByRemotePath(remotePath)
	if err != nil {
		return fmt.Errorf("error remove smb mapping '%s'. err: %v", remotePath, err)
	}

	return nil
}
