package utils

import (
	"os"
	"os/exec"

	"k8s.io/klog/v2"
)

const MaxPathLengthWindows = 260

func RunPowershellCmd(command string) ([]byte, error) {
	return RunPowershellCmdWithEnvs(command, nil)
}

func RunPowershellCmdWithEnv(command string, env string) ([]byte, error) {
	return RunPowershellCmdWithEnvs(command, []string{env})
}

func RunPowershellCmdWithEnvs(command string, envs []string) ([]byte, error) {
	cmd := exec.Command("powershell", "-Mta", "-NoProfile", "-Command", command)
	if envs != nil {
		cmd.Env = append(os.Environ(), envs...)
	}
	klog.V(8).Infof("Executing command: %q", cmd.String())
	out, err := cmd.CombinedOutput()
	return out, err
}
