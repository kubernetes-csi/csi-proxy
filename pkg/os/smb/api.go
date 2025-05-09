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

func remotePathForQuery(remotePath string) string {
	return strings.ReplaceAll(remotePath, "\\", "\\\\")
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

	cmdLine := `New-Item -ItemType SymbolicLink $Env:smblocalPath -Target $Env:smbremotepath`
	output, err := utils.RunPowershellCmd(cmdLine, fmt.Sprintf("smbremotepath=%s", remotePath), fmt.Sprintf("smblocalpath=%s", localPath))
	if err != nil {
		return fmt.Errorf("error linking %s to %s. output: %s, err: %v", remotePath, localPath, string(output), err)
	}

	return nil
}

func (api *SmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	// use PowerShell Environment Variables to store user input string to prevent command line injection
	// https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-5.1
	cmdLine := fmt.Sprintf(`$PWord = ConvertTo-SecureString -String $Env:smbpassword -AsPlainText -Force`+
		`;$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $Env:smbuser, $PWord`+
		`;New-SmbGlobalMapping -RemotePath $Env:smbremotepath -Credential $Credential -RequirePrivacy $%t`, api.RequirePrivacy)

	if output, err := utils.RunPowershellCmd(cmdLine,
		fmt.Sprintf("smbuser=%s", username),
		fmt.Sprintf("smbpassword=%s", password),
		fmt.Sprintf("smbremotepath=%s", remotePath)); err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed. output: %q, err: %v", string(output), err)
	}
	return nil
	//TODO: move to use WMI when the credentials could be correctly handled
	//params := map[string]interface{}{
	//	"RemotePath":     remotePath,
	//	"RequirePrivacy": api.RequirePrivacy,
	//}
	//if username != "" {
	//	params["Credential"] = fmt.Sprintf("%s:%s", username, password)
	//}
	//result, _, err := cim.InvokeCimMethod(cim.WMINamespaceSmb, "MSFT_SmbGlobalMapping", "Create", params)
	//if err != nil {
	//	return fmt.Errorf("NewSmbGlobalMapping failed. result: %d, err: %v", result, err)
	//}
	//return nil
}

func (*SmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	err := cim.RemoveSmbGlobalMappingByRemotePath(remotePathForQuery(remotePath))
	if err != nil {
		return fmt.Errorf("error remove smb mapping '%s'. err: %v", remotePath, err)
	}

	return nil
}
