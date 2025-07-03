package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
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

func IsPathValid(path string) (bool, error) {
	pathString, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return false, fmt.Errorf("invalid path: %w", err)
	}

	attrs, err := windows.GetFileAttributes(pathString)
	if err != nil {
		if errors.Is(err, windows.ERROR_PATH_NOT_FOUND) || errors.Is(err, windows.ERROR_FILE_NOT_FOUND) || errors.Is(err, windows.ERROR_INVALID_NAME) {
			return false, nil
		}

		// GetFileAttribute returns user or password incorrect for a disconnected SMB connection after the password is changed
		return false, fmt.Errorf("failed to get path %s attribute: %w", path, err)
	}

	klog.V(6).Infof("Path %s attribute: %d", path, attrs)
	return attrs != windows.INVALID_FILE_ATTRIBUTES, nil
}
