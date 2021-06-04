package integrationtests

import (
	"os"
	"testing"
)

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

func TestFilesystemAPIGroup(t *testing.T) {
	t.Run("v1beta2FilesystemTests", func(t *testing.T) {
		v1beta2FilesystemTests(t)
	})
	t.Run("v1beta1FilesystemTests", func(t *testing.T) {
		v1beta1FilesystemTests(t)
	})
	t.Run("v1alpha1FilesystemTests", func(t *testing.T) {
		v1alpha1FilesystemTests(t)
	})
}
