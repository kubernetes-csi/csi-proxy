package filesystem

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/os/filesystem"
	internal "github.com/kubernetes-csi/csi-proxy/pkg/server/filesystem/impl"
)

type fakeFileSystemAPI struct{}

var _ filesystem.API = &fakeFileSystemAPI{}

func (fakeFileSystemAPI) PathExists(path string) (bool, error) {
	return true, nil
}
func (fakeFileSystemAPI) PathValid(path string) (bool, error) {
	return true, nil
}
func (fakeFileSystemAPI) Mkdir(path string) error {
	return nil
}
func (fakeFileSystemAPI) Rmdir(path string, force bool) error {
	return nil
}
func (fakeFileSystemAPI) Lsdir(path string) ([]string, error) {
	return nil, nil
}
func (fakeFileSystemAPI) CreateSymlink(tgt string, src string) error {
	return nil
}

func (fakeFileSystemAPI) IsSymlink(path string) (bool, error) {
	return true, nil
}

func TestMkdirWindows(t *testing.T) {
	v1, err := apiversion.NewVersion("v1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}
	testCases := []struct {
		name        string
		path        string
		version     apiversion.Version
		expectError bool
	}{
		{
			name:        "path outside of pod context with pod context set",
			path:        `C:\foo\bar`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path inside pod context with pod context set",
			path:        `C:\var\lib\kubelet\pods\pv1`,
			version:     v1,
			expectError: false,
		},
		{
			name:        "path outside of plugin context with plugin context set",
			path:        `C:\foo\bar`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path inside plugin context with plugin context set",
			path:        `C:\var\lib\kubelet\plugins\pv1`,
			version:     v1,
			expectError: false,
		},
		{
			name:        "path with invalid character `:` beyond drive letter prefix",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1:foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `/`",
			path:        `C:\var\lib\kubelet\pods\pv1/foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `*`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1*foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `?`",
			path:        `C:\var\lib\kubelet\pods\pv1?foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `|`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin|pv1\foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid characters `..`",
			path:        `C:\var\lib\kubelet\pods\pv1\..\..\..\system32`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid prefix `\\`",
			path:        `\\csi-plugin\..\..\..\system32`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "relative path",
			path:        `pv1\foo`,
			version:     v1,
			expectError: true,
		},
	}
	srv, err := NewServer(`C:\var\lib\kubelet`, &fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		req := &internal.MkdirRequest{
			Path: tc.path,
		}
		_, err := srv.Mkdir(context.TODO(), req, tc.version)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but Mkdir returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but Mkdir returned error: %v", err)
		}
	}
}

func TestRmdirWindows(t *testing.T) {
	v1, err := apiversion.NewVersion("v1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}
	testCases := []struct {
		name        string
		path        string
		version     apiversion.Version
		expectError bool
		force       bool
	}{
		{
			name:        "path outside of pod context with pod context set",
			path:        `C:\foo\bar`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path inside pod context with pod context set",
			path:        `C:\var\lib\kubelet\pods\pv1`,
			version:     v1,
			expectError: false,
		},
		{
			name:        "path outside of plugin context with plugin context set",
			path:        `C:\foo\bar`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path inside plugin context with plugin context set",
			path:        `C:\var\lib\kubelet\plugins\pv1`,
			version:     v1,
			expectError: false,
		},
		{
			name:        "path with invalid character `:` beyond drive letter prefix",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1:foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `/`",
			path:        `C:\var\lib\kubelet\pods\pv1/foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `*`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1*foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `?`",
			path:        `C:\var\lib\kubelet\pods\pv1?foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid character `|`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin|pv1\foo`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid characters `..`",
			path:        `C:\var\lib\kubelet\pods\pv1\..\..\..\system32`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "path with invalid prefix `\\`",
			path:        `\\csi-plugin\..\..\..\system32`,
			version:     v1,
			expectError: true,
		},
		{
			name:        "relative path",
			path:        `pv1\foo`,
			version:     v1,
			expectError: true,
		},
	}
	srv, err := NewServer(`C:\var\lib\kubelet`, &fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		req := &internal.RmdirRequest{
			Path:  tc.path,
			Force: tc.force,
		}
		_, err := srv.Rmdir(context.TODO(), req, tc.version)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but Rmdir returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but Rmdir returned error: %v", err)
		}
	}
}
