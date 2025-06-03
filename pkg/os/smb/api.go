package smb

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/utils"
	"golang.org/x/sys/windows"
)

const (
	credentialDelimiter = ":"
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

func remotePathForQuery(remotePath string) string {
	return strings.ReplaceAll(remotePath, "\\", "\\\\")
}

func escapeUserName(userName string) string {
	// refer to https://github.com/PowerShell/PowerShell/blob/9303de597da55963a6e26a8fe164d0b256ca3d4d/src/Microsoft.PowerShell.Commands.Management/cimSupport/cmdletization/cim/cimConverter.cs#L169-L170
	escaped := strings.ReplaceAll(userName, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, credentialDelimiter, "\\"+credentialDelimiter)
	return escaped
}

func createSymlink(link, target string, isDir bool) error {
	linkPtr, err := syscall.UTF16PtrFromString(link)
	if err != nil {
		return err
	}
	targetPtr, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return err
	}

	var flags uint32
	if isDir {
		flags = windows.SYMBOLIC_LINK_FLAG_DIRECTORY
	}

	err = windows.CreateSymbolicLink(
		linkPtr,
		targetPtr,
		flags,
	)
	return err
}

func (*SmbAPI) IsSmbMapped(remotePath string) (bool, error) {
	inst, err := cim.QuerySmbGlobalMappingByRemotePath(remotePathForQuery(remotePath))
	if err != nil {
		return false, cim.IgnoreNotFound(err)
	}

	status, err := inst.GetProperty("Status")
	if err != nil {
		return false, err
	}

	return status.(int32) == cim.SmbMappingStatusOK, nil
}

// NewSmbLink - creates a directory symbolic link to the remote share.
// The os.Symlink was having issue for cases where the destination was an SMB share - the container
// runtime would complain stating "Access Denied". Because of this, we had to perform
// this operation with powershell commandlet creating an directory softlink.
// Since os.Symlink is currently being used in working code paths, no attempt is made in
// alpha to merge the paths.
// TODO (for beta release): Merge the link paths - os.Symlink and Powershell link path.
func (*SmbAPI) NewSmbLink(remotePath, localPath string) error {
	if !strings.HasSuffix(remotePath, "\\") {
		// Golang has issues resolving paths mapped to file shares if they do not end in a trailing \
		// so add one if needed.
		remotePath = remotePath + "\\"
	}
	longRemotePath := utils.EnsureLongPath(remotePath)
	longLocalPath := utils.EnsureLongPath(localPath)

	err := createSymlink(longLocalPath, longRemotePath, true)
	if err != nil {
		return fmt.Errorf("error linking %s to %s. err: %v", remotePath, localPath, err)
	}

	return nil
}

func (api *SmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	params := map[string]interface{}{
		"RemotePath":     remotePath,
		"RequirePrivacy": api.RequirePrivacy,
	}
	if username != "" {
		// refer to https://github.com/PowerShell/PowerShell/blob/9303de597da55963a6e26a8fe164d0b256ca3d4d/src/Microsoft.PowerShell.Commands.Management/cimSupport/cmdletization/cim/cimConverter.cs#L166-L178
		// on how SMB credential is handled in PowerShell
		params["Credential"] = escapeUserName(username) + credentialDelimiter + password
	}

	result, _, err := cim.InvokeCimMethod(cim.WMINamespaceSmb, "MSFT_SmbGlobalMapping", "Create", params)
	if err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed. result: %d, err: %v", result, err)
	}

	return nil
}

func (*SmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	err := cim.RemoveSmbGlobalMappingByRemotePath(remotePathForQuery(remotePath))
	if err != nil {
		return fmt.Errorf("error remove smb mapping '%s'. err: %v", remotePath, err)
	}

	return nil
}
