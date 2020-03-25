package smb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/klog"
)

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) IsSmbMapped(remotePath string) (bool, error) {
	cmdLine := fmt.Sprintf(`$(Get-SmbGlobalMapping -RemotePath $Env:smbremotepath -ErrorAction Stop).Status `)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("smbremotepath=%s", remotePath))

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error checking smb mapping %s, err: %v", remotePath, err)
	}

	if len(out) == 0 || !strings.EqualFold(strings.TrimSpace(string(out)), "OK") {
		return false, nil
	}
	return true, nil
}

// SMBLink - creates a direcotry symbolic link to the remote share.
// The os.Symlink was having issue for cases where the destination was an SMB share - the container
// runtime would complain stating "Access Denied". Because of this, we had to perform
// this operation with powershell commandlet creating an directory softlink.
// Since os.Symlink is currently being used in working code paths, no attempt is made in
// alpha to merge the paths.
// TODO (for beta release): Merge the link paths - os.Symlink and Powershell link path.
func (APIImplementor) SMBLink(remotePath, localPath string) error {
	cmdLine := fmt.Sprintf(`New-Item -ItemType SymbolicLink $Env:smblocalPath -Target $Env:smbremotepath`)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("smbremotepath=%s", remotePath),
		fmt.Sprintf("smblocalpath=%s", localPath),
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error linking %s to %s, err: %v", remotePath, localPath, err)
	}

	return nil
}

func (APIImplementor) NewSmbGlobalMapping(remotePath, localPath, username, password string) error {
	klog.V(4).Infof("NewSmbGlobalMapping: remotePath:%q, localPath:%q", remotePath, localPath)

	// use PowerShell Environment Variables to store user input string to prevent command line injection
	// https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-5.1
	cmdLine := fmt.Sprintf(`$PWord = ConvertTo-SecureString -String $Env:smbpassword -AsPlainText -Force` +
		`;$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $Env:smbuser, $PWord` +
		`;New-SmbGlobalMapping -RemotePath $Env:smbremotepath -Credential $Credential`)

	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("smbuser=%s", username),
		fmt.Sprintf("smbpassword=%s", password),
		fmt.Sprintf("smbremotepath=%s", remotePath))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("NewSmbGlobalMapping failed: %v, output: %q", err, string(output))
	}
	return nil
}

func (APIImplementor) RemoveSmbGlobalMapping(remotePath string) error {
	klog.V(4).Infof("RemoveSmbGlobalMapping remotePath (%q)", remotePath)
	cmd := exec.Command("powershell", "/c", `Remove-SmbGlobalMapping -RemotePath $Env:smbremotepath -Force`)
	cmd.Env = append(os.Environ(), fmt.Sprintf("smbremotepath=%s", remotePath))
	if output, err := cmd.CombinedOutput(); err != nil {
		klog.Errorf("Remove-SmbGlobalMapping failed: %v, output: %q", err, output)
		return fmt.Errorf("UnmountSmbShare failed: %v, output: %q", err, output)
	}
	return nil
}
