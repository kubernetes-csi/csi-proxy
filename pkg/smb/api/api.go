package api

import (
	"fmt"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/v2/pkg/utils"
)

type API interface {
	IsSmbMapped(remotePath string) (bool, error)
	NewSmbLink(remotePath, localPath string) error
	NewSmbGlobalMapping(remotePath, username, password string) error
	RemoveSmbGlobalMapping(remotePath string) error
}

type smbAPI struct{}

var _ API = &smbAPI{}

func New() API {
	return smbAPI{}
}

func (smbAPI) IsSmbMapped(remotePath string) (bool, error) {
	cmdLine := `$(Get-SmbGlobalMapping -RemotePath $Env:smbremotepath -ErrorAction Stop).Status `
	cmdEnv := fmt.Sprintf("smbremotepath=%s", remotePath)
	out, err := utils.RunPowershellCmd(cmdLine, cmdEnv)
	if err != nil {
		return false, fmt.Errorf("error checking smb mapping. cmd %s, output: %s, err: %v", remotePath, string(out), err)
	}

	if len(out) == 0 || !strings.EqualFold(strings.TrimSpace(string(out)), "OK") {
		return false, nil
	}
	return true, nil
}

// NewSmbLink - creates a directory symbolic link to the remote share.
// The os.Symlink was having issue for cases where the destination was an SMB share - the container
// runtime would complain stating "Access Denied". Because of this, we had to perform
// this operation with powershell commandlet creating an directory softlink.
// Since os.Symlink is currently being used in working code paths, no attempt is made in
// alpha to merge the paths.
// TODO (for beta release): Merge the link paths - os.Symlink and Powershell link path.
func (smbAPI) NewSmbLink(remotePath, localPath string) error {
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

func (smbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	// use PowerShell Environment Variables to store user input string to prevent command line injection
	// https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-5.1
	cmdLine := fmt.Sprintf(`$PWord = ConvertTo-SecureString -String $Env:smbpassword -AsPlainText -Force` +
		`;$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $Env:smbuser, $PWord` +
		`;New-SmbGlobalMapping -RemotePath $Env:smbremotepath -Credential $Credential -RequirePrivacy $true`)

	if output, err := utils.RunPowershellCmd(cmdLine, fmt.Sprintf("smbuser=%s", username),
		fmt.Sprintf("smbpassword=%s", password),
		fmt.Sprintf("smbremotepath=%s", remotePath)); err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed. output: %q, err: %v", string(output), err)
	}
	return nil
}

func (smbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	cmd := `Remove-SmbGlobalMapping -RemotePath $Env:smbremotepath -Force`
	if output, err := utils.RunPowershellCmd(cmd, fmt.Sprintf("smbremotepath=%s", remotePath)); err != nil {
		return fmt.Errorf("UnmountSmbShare failed. output: %q, err: %v", string(output), err)
	}
	return nil
}
