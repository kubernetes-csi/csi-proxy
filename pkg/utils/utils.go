package utils

import (
	"os"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"
)

const (
	MaxPathLengthWindows = 260

	// LongPathPrefix is the prefix of Windows long path
	LongPathPrefix = `\\?\`
)

func EnsureLongPath(path string) string {
	if !strings.HasPrefix(path, LongPathPrefix) {
		path = LongPathPrefix + path
	}
	return path
}

func RunPowershellCmd(command string, envs ...string) ([]byte, error) {
	cmd := exec.Command("powershell", "-Mta", "-NoProfile", "-Command", command)
	cmd.Env = append(os.Environ(), envs...)
	klog.V(8).Infof("Executing command: %q", cmd.String())
	out, err := cmd.CombinedOutput()
	return out, err
}
