package iscsi

import (
	"encoding/json"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/pkg/utils"
)

// Implements the iSCSI OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/iscsi/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) AddTargetPortal(portal *TargetPortal) error {
	cmdLine := fmt.Sprintf(
		`New-IscsiTargetPortal -TargetPortalAddress ${Env:iscsi_tp_address} ` +
			`-TargetPortalPortNumber ${Env:iscsi_tp_port}`)
	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port)})
	if err != nil {
		return fmt.Errorf("error adding target portal. cmd %s, output: %s, err: %v", cmdLine, string(out), err)
	}

	return nil
}

func (APIImplementor) DiscoverTargetPortal(portal *TargetPortal) ([]string, error) {
	// ConvertTo-Json is not part of the pipeline because powershell converts an
	// array with one element to a single element
	cmdLine := fmt.Sprintf(
		`ConvertTo-Json -InputObject @(Get-IscsiTargetPortal -TargetPortalAddress ` +
			`${Env:iscsi_tp_address} -TargetPortalPortNumber ${Env:iscsi_tp_port} | ` +
			`Get-IscsiTarget | Select-Object -ExpandProperty NodeAddress)`)
	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port)})
	if err != nil {
		return nil, fmt.Errorf("error discovering target portal. cmd: %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	var iqns []string
	err = json.Unmarshal(out, &iqns)
	if err != nil {
		return nil, fmt.Errorf("failed parsing iqn list. cmd: %s output: %s, err: %w", cmdLine, string(out), err)
	}

	return iqns, nil
}

func (APIImplementor) ListTargetPortals() ([]TargetPortal, error) {
	cmdLine := fmt.Sprintf(
		`ConvertTo-Json -InputObject @(Get-IscsiTargetPortal | ` +
			`Select-Object TargetPortalAddress, TargetPortalPortNumber)`)

	out, err := utils.RunPowershellCmd(cmdLine)
	if err != nil {
		return nil, fmt.Errorf("error listing target portals. cmd %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	var portals []TargetPortal
	err = json.Unmarshal(out, &portals)
	if err != nil {
		return nil, fmt.Errorf("failed parsing target portal list. cmd: %s output: %s, err: %w", cmdLine, string(out), err)
	}

	return portals, nil
}

func (APIImplementor) RemoveTargetPortal(portal *TargetPortal) error {
	cmdLine := fmt.Sprintf(
		`Get-IscsiTargetPortal -TargetPortalAddress ${Env:iscsi_tp_address} ` +
			`-TargetPortalPortNumber ${Env:iscsi_tp_port} | Remove-IscsiTargetPortal ` +
			`-Confirm:$false`)

	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port)})
	if err != nil {
		return fmt.Errorf("error removing target portal. cmd %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	return nil
}

func (APIImplementor) ConnectTarget(portal *TargetPortal, iqn string,
	authType string, chapUser string, chapSecret string) error {
	// Not using InputObject as Connect-IscsiTarget's InputObject does not work.
	// This is due to being a static WMI method together with a bug in the
	// powershell version of the API.
	cmdLine := fmt.Sprintf(
		`Connect-IscsiTarget -TargetPortalAddress ${Env:iscsi_tp_address}` +
			` -TargetPortalPortNumber ${Env:iscsi_tp_port} -NodeAddress ${Env:iscsi_target_iqn}` +
			` -AuthenticationType ${Env:iscsi_auth_type}`)

	if chapUser != "" {
		cmdLine += fmt.Sprintf(` -ChapUsername ${Env:iscsi_chap_user}`)
	}

	if chapSecret != "" {
		cmdLine += fmt.Sprintf(` -ChapSecret ${Env:iscsi_chap_secret}`)
	}

	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port),
		fmt.Sprintf("iscsi_target_iqn=%s", iqn),
		fmt.Sprintf("iscsi_auth_type=%s", authType),
		fmt.Sprintf("iscsi_chap_user=%s", chapUser),
		fmt.Sprintf("iscsi_chap_secret=%s", chapSecret)})
	if err != nil {
		return fmt.Errorf("error connecting to target portal. cmd %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	return nil
}

func (APIImplementor) DisconnectTarget(portal *TargetPortal, iqn string) error {
	// Using InputObject instead of pipe to verify input is not empty
	cmdLine := fmt.Sprintf(
		`Disconnect-IscsiTarget -InputObject (Get-IscsiTargetPortal ` +
			`-TargetPortalAddress ${Env:iscsi_tp_address} -TargetPortalPortNumber ${Env:iscsi_tp_port} ` +
			` | Get-IscsiTarget | Where-Object { $_.NodeAddress -eq ${Env:iscsi_target_iqn} }) ` +
			`-Confirm:$false`)

	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port),
		fmt.Sprintf("iscsi_target_iqn=%s", iqn)})
	if err != nil {
		return fmt.Errorf("error disconnecting from target portal. cmd %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	return nil
}

func (APIImplementor) GetTargetDisks(portal *TargetPortal, iqn string) ([]string, error) {
	// Converting DiskNumber to string for compatibility with disk api group
	// Not using pipeline in order to validate that items are non-empty
	cmdLine := fmt.Sprintf(
		`$ErrorActionPreference = "Stop"; ` +
			`$tp = Get-IscsiTargetPortal -TargetPortalAddress ${Env:iscsi_tp_address} -TargetPortalPortNumber ${Env:iscsi_tp_port}; ` +
			`$t = $tp | Get-IscsiTarget | Where-Object { $_.NodeAddress -eq ${Env:iscsi_target_iqn} }; ` +
			`$c = Get-IscsiConnection -IscsiTarget $t; ` +
			`$ids = $c | Get-Disk | Select -ExpandProperty Number | Out-String -Stream; ` +
			`ConvertTo-Json -InputObject @($ids)`)

	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_tp_address=%s", portal.Address),
		fmt.Sprintf("iscsi_tp_port=%d", portal.Port),
		fmt.Sprintf("iscsi_target_iqn=%s", iqn)})
	if err != nil {
		return nil, fmt.Errorf("error getting target disks. cmd %s, output: %s, err: %w", cmdLine, string(out), err)
	}

	var ids []string
	err = json.Unmarshal(out, &ids)
	if err != nil {
		return nil, fmt.Errorf("error parsing iqn target disks. cmd: %s output: %s, err: %w", cmdLine, string(out), err)
	}

	return ids, nil
}

func (APIImplementor) SetMutualChapSecret(mutualChapSecret string) error {
	cmdLine := fmt.Sprintf(
		`Set-IscsiChapSecret -ChapSecret ${Env:iscsi_mutual_chap_secret}`)
	out, err := utils.RunPowershellCmdWithEnvs(cmdLine, []string{fmt.Sprintf("iscsi_mutual_chap_secret=%s", mutualChapSecret)})
	if err != nil {
		return fmt.Errorf("error setting mutual chap secret. cmd %s,"+
			" output: %s, err: %v", cmdLine, string(out), err)
	}

	return nil
}
