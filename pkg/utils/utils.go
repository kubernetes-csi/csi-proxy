package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
	"k8s.io/klog/v2"
)

const (
	MaxPathLengthWindows = 260

	// LongPathPrefix is the prefix of Windows long path
	LongPathPrefix = `\\?\`
)

func RunPowershellCmd(command string, envs ...string) ([]byte, error) {
	command = fmt.Sprintf("$global:ProgressPreference = 'SilentlyContinue'; %s", command)
	cmd := exec.Command("powershell", "-Mta", "-NoProfile", "-Command", command)
	cmd.Env = append(os.Environ(), envs...)
	klog.V(8).Infof("Executing command: %q", cmd.String())
	out, err := cmd.CombinedOutput()
	return out, err
}

func EnsureLongPath(path string) string {
	if !strings.HasPrefix(path, LongPathPrefix) {
		path = LongPathPrefix + path
	}
	return path
}

// IsMountedFolder checks whether the `path` is a mounted folder.
func IsMountedFolder(path string) (bool, error) {
	// https://learn.microsoft.com/en-us/windows/win32/fileio/determining-whether-a-directory-is-a-volume-mount-point
	utf16Path, _ := windows.UTF16PtrFromString(path)
	attrs, err := windows.GetFileAttributes(utf16Path)
	if err != nil {
		return false, err
	}

	if (attrs & windows.FILE_ATTRIBUTE_REPARSE_POINT) == 0 {
		return false, nil
	}

	var findData windows.Win32finddata
	findHandle, err := windows.FindFirstFile(utf16Path, &findData)
	if err != nil && !errors.Is(err, windows.ERROR_NO_MORE_FILES) {
		return false, err
	}

	for err == nil {
		if findData.Reserved0&windows.IO_REPARSE_TAG_MOUNT_POINT != 0 {
			return true, nil
		}

		err = windows.FindNextFile(findHandle, &findData)
		if err != nil && !errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			return false, err
		}
	}

	return false, nil
}

func IsPathSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	// for windows NTFS, check if the path is symlink instead of directory.
	isSymlink := fi.Mode()&os.ModeSymlink != 0 || fi.Mode()&os.ModeIrregular != 0
	return isSymlink, nil
}
