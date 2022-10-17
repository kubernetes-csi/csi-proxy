package integrationtests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func installISCSITarget() error {
	_, err := runPowershellScript(iSCSITargetInstallScript)
	if err != nil {
		return fmt.Errorf("failed installing iSCSI target. err=%v", err)
	}

	return nil
}

const iSCSITargetInstallScript = `
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Install iSCSI Target
Install-WindowsFeature FS-iSCSITarget-Server

# Setup for loopback usage
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\iSCSI Target" -Name AllowLoopBack -Value 1
Restart-Service WinTarget
`

type iSCSISetupConfig struct {
	Iqn string `json:"iqn"`
	Ip  string `json:"ip"`
}

const iSCSIEnvironmentSetupScript = `
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$targetName = "%s"

# Get local IPv4 (e.g. 10.30.1.15, not 127.0.0.1)
$address = $(Get-NetIPAddress | Where-Object { $_.InterfaceAlias -eq "Ethernet" -and $_.AddressFamily -eq "IPv4" }).IPAddress

# Create virtual disk in RAM
New-IscsiVirtualDisk -Path "ramdisk:scratch-${targetName}.vhdx" -Size 100MB | Out-Null

# Create a target that allows all initiator IQNs and map a disk to the new target
$target = New-IscsiServerTarget -TargetName $targetName -InitiatorIds @("Iqn:*")
Add-IscsiVirtualDiskTargetMapping -TargetName $targetName -DevicePath "ramdisk:scratch-${targetName}.vhdx" | Out-Null

$output = @{
  "iqn" = "$($target.TargetIqn)"
  "ip"  = $address
}

$output | ConvertTo-Json | Write-Output
`

const iSCSISetChapScript = `
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$targetName = "%s"
$username = "%s"
$password = "%s"
$securestring = ConvertTo-SecureString -String $password -AsPlainText -Force
$chap = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList ($username, $securestring)
Set-IscsiServerTarget -TargetName $targetName -EnableChap $true -Chap $chap
`

func setChap(targetName string, username string, password string) error {
	script := fmt.Sprintf(iSCSISetChapScript, targetName, username, password)
	_, err := runPowershellScript(script)
	if err != nil {
		return fmt.Errorf("failed setting CHAP on iSCSI target=%v. err=%v", targetName, err)
	}

	return nil
}

const iSCSISetReverseChapScript = `
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$targetName = "%s"
$password = "%s"
$username = "doesnt-matter"
$securestring = ConvertTo-SecureString -String $password -AsPlainText -Force

# Windows initiator does not uses the username for mutual authentication
$chap = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList ($username, $securestring)
Set-IscsiServerTarget -TargetName $targetName -EnableReverseChap $true -ReverseChap $chap
`

func setReverseChap(targetName string, password string) error {
	script := fmt.Sprintf(iSCSISetReverseChapScript, targetName, password)
	_, err := runPowershellScript(script)
	if err != nil {
		return fmt.Errorf("failed setting reverse CHAP on iSCSI target=%v. err=%v", targetName, err)
	}

	return nil
}

func cleanup() error {
	_, err := runPowershellScript(iSCSICleanupScript)
	if err != nil {
		return fmt.Errorf("failed cleaning up environment. err=%v", err)
	}

	return nil
}

func requireCleanup(t *testing.T) {
	err := cleanup()
	if err != nil {
		t.Fatal(err)
	}
}

const iSCSICleanupScript = `
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Clean initiator
Get-Disk | Where-Object {$_.Bustype -eq "iSCSI"} | Set-Disk -IsOffline:$true
Get-IscsiTarget | Disconnect-IscsiTarget -Confirm:$false
Get-IscsiTargetPortal | Remove-IscsiTargetPortal -confirm:$false

# Clean target
Get-IscsiServerTarget | Remove-IscsiServerTarget
Get-IscsiVirtualDisk | Remove-IscsiVirtualDisk

# Stop iSCSI initiator
Get-Service "MsiSCSI" | Stop-Service
`

func writeTempFile(text string, extension string) (string, error) {
	pattern := fmt.Sprintf("*.%s", extension)
	tempfile, err := ioutil.TempFile(os.TempDir(), pattern)
	if err != nil {
		return "", fmt.Errorf("failed creating temp file pattern=%v: %w", pattern, err)
	}

	defer tempfile.Close()

	_, err = tempfile.WriteString(text)
	if err != nil {
		return "", fmt.Errorf("failed writing to temp file name=%v: %w", tempfile.Name(), err)
	}

	return tempfile.Name(), nil
}

func runPowershellScript(script string) (string, error) {
	path, err := writeTempFile(script, "ps1")
	if err != nil {
		return "", err
	}

	defer os.Remove(path)

	cmd := exec.Command("powershell", "-File", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running powershell script. path %s, output: %s, err: %w", path, string(out), err)
	}

	return string(out), nil
}

func setupEnv(targetName string) (*iSCSISetupConfig, error) {
	script := fmt.Sprintf(iSCSIEnvironmentSetupScript, targetName)
	out, err := runPowershellScript(script)
	if err != nil {
		return nil, fmt.Errorf("failed setting up environment. err=%v", err)
	}

	config := iSCSISetupConfig{}
	err = json.Unmarshal([]byte(out), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
