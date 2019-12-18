package smb

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"k8s.io/klog"
	"k8s.io/utils/keymutex"
)

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

// acquire lock for smb mount
var getSMBMountMutex = keymutex.NewHashed(0)

func (APIImplementor) MountSmbShare(remotePath, localPath, username, password string, readOnly bool) error {
	localPath = normalizeWindowsPath(localPath)
	klog.V(4).Infof("MountSmbShare: options(%q) remotePath:%q, localPath:%q, readOnly:%t", remotePath, localPath, readOnly)

	// lock smb mount for the same source
	getSMBMountMutex.LockKey(remotePath)
	defer getSMBMountMutex.UnlockKey(remotePath)

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
		return fmt.Errorf("SmbGlobaNewSmbGlobalMappinglMapping failed: %v, output: %q", err, string(output))
	}

	if output, err := exec.Command("cmd", "/c", "mklink", "/D", localPath, remotePath).CombinedOutput(); err != nil {
		klog.Errorf("mklink failed: %v, source(%q) target(%q) output: %q", err, remotePath, localPath, string(output))
		return fmt.Errorf("MountSmbShare failed: %v, output: %q", err, string(output))
	}
	return nil
}

func (APIImplementor) UnmountSmbShare(remotePath, localPath string) error {
	klog.V(4).Infof("Unmount target (%q)", localPath)
	target := normalizeWindowsPath(localPath)
	if output, err := exec.Command("cmd", "/c", "rmdir", target).CombinedOutput(); err != nil {
		klog.Errorf("rmdir %q failed: %v, output: %q", target, err, string(output))
		return fmt.Errorf("UnmountSmbShare failed: %v, output: %q", err, output)
	}
	if output, err := removeSmbMapping(remotePath); err != nil {
		klog.Errorf("Remove-SmbGlobalMapping failed: %v, output: %q", err, output)
		return fmt.Errorf("UnmountSmbShare failed: %v, output: %q", err, output)
	}
	return nil
}

func normalizeWindowsPath(path string) string {
	normalizedPath := strings.Replace(path, "/", "\\", -1)
	if strings.HasPrefix(normalizedPath, "\\") {
		normalizedPath = "c:" + normalizedPath
	}
	return normalizedPath
}

// remove SMB mapping
func removeSmbMapping(remotePath string) (string, error) {
	cmd := exec.Command("powershell", "/c", `Remove-SmbGlobalMapping -RemotePath $Env:smbremotepath -Force`)
	cmd.Env = append(os.Environ(), fmt.Sprintf("smbremotepath=%s", remotePath))
	output, err := cmd.CombinedOutput()
	return string(output), err
}
