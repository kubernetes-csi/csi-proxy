package filesystem

import (
	// "fmt"
	"os"
	// "os/exec"
	// "runtime"
)

// Implements the Filesystem OS API calls. All code here should be very simple
// pass-through to the OS APIs. Any logic around the APIs should go in
// internal/server/filesystem/server.go so that logic can be easily unit-tested
// without requiring specific OS environments.

type APIImplementor struct{}

func New() APIImplementor {
	return APIImplementor{}
}

func (APIImplementor) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
