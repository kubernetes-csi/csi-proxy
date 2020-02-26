package filesystem

import (
	"fmt"
	"os"
)

// Implements the Filesystem OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/filesystem/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (APIImplementor) PathExists(path string) (bool, error) {
	return pathExists(path)
}

func (APIImplementor) Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (APIImplementor) Rmdir(path string, force bool) error {
	if force {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func (APIImplementor) LinkPath(tgt string, src string) error {
	return os.Symlink(tgt, src)
}

// IsMountPoint - returns true if its a mount point.
// A path is considered a mount point if:
//  - directory exists and
//  - it is a soft link and
//  - the target path of the link exists.
func (APIImplementor) IsMountPoint(tgt string) (bool, error) {
	// This code is similar to k8s.io/kubernetes/pkg/util/mount except the pathExists usage.
	// Also in a remote call environment the os error cannot be passed directly back, hence the callers
	// are expected to perform the isExists check before calling this call in CSI proxy.
	stat, err := os.Lstat(tgt)
	if err != nil {
		return false, err
	}

	// If its a link and it points to an existing file then its a mount point.
	if stat.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(tgt)
		if err != nil {
			return false, fmt.Errorf("readlink error: %v", err)
		}
		exists, err := pathExists(target)
		if err != nil {
			return false, err
		}
		return exists, nil
	}

	return false, nil
}
