package filesystem

import (
	"context"
	"testing"

	fsapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem/hostapi"
)

type fakeFileSystemAPI struct{}

var _ fsapi.HostAPI = &fakeFileSystemAPI{}

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
func (fakeFileSystemAPI) RmdirContents(path string) error {
	return nil
}
func (fakeFileSystemAPI) CreateSymlink(tgt string, src string) error {
	return nil
}

func (fakeFileSystemAPI) IsSymlink(path string) (bool, error) {
	return true, nil
}

func TestMkdirWindows(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "path inside pod context with pod context set",
			path:        `C:\var\lib\kubelet\pods\pv1`,
			expectError: false,
		},
		{
			name:        "path inside plugin context with plugin context set",
			path:        `C:\var\lib\kubelet\plugins\pv1`,
			expectError: false,
		},
		{
			name:        "path with invalid character `:` beyond drive letter prefix",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1:foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `/`",
			path:        `C:\var\lib\kubelet\pods\pv1/foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `*`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1*foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `?`",
			path:        `C:\var\lib\kubelet\pods\pv1?foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `|`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin|pv1\foo`,
			expectError: true,
		},
		{
			name:        "path with invalid characters `..`",
			path:        `C:\var\lib\kubelet\pods\pv1\..\..\..\system32`,
			expectError: true,
		},
		{
			name:        "path with invalid prefix `\\`",
			path:        `\\csi-plugin\..\..\..\system32`,
			expectError: true,
		},
		{
			name:        "relative path",
			path:        `pv1\foo`,
			expectError: true,
		},
	}
	client, err := New(&fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		req := &MkdirRequest{
			Path: tc.path,
		}
		_, err := client.Mkdir(context.TODO(), req)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but Mkdir returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but Mkdir returned error: %v", err)
		}
	}
}

func TestRmdirWindows(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		expectError bool
		force       bool
	}{
		{
			name:        "path inside pod context with pod context set",
			path:        `C:\var\lib\kubelet\pods\pv1`,
			expectError: false,
		},
		{
			name:        "path inside plugin context with plugin context set",
			path:        `C:\var\lib\kubelet\plugins\pv1`,
			expectError: false,
		},
		{
			name:        "path with invalid character `:` beyond drive letter prefix",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1:foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `/`",
			path:        `C:\var\lib\kubelet\pods\pv1/foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `*`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin\pv1*foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `?`",
			path:        `C:\var\lib\kubelet\pods\pv1?foo`,
			expectError: true,
		},
		{
			name:        "path with invalid character `|`",
			path:        `C:\var\lib\kubelet\plugins\csi-plugin|pv1\foo`,
			expectError: true,
		},
		{
			name:        "path with invalid characters `..`",
			path:        `C:\var\lib\kubelet\pods\pv1\..\..\..\system32`,
			expectError: true,
		},
		{
			name:        "path with invalid prefix `\\`",
			path:        `\\csi-plugin\..\..\..\system32`,
			expectError: true,
		},
		{
			name:        "relative path",
			path:        `pv1\foo`,
			expectError: true,
		},
	}
	client, err := New(&fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		t.Logf("test case: %s", tc.name)
		req := &RmdirRequest{
			Path:  tc.path,
			Force: tc.force,
		}
		_, err := client.Rmdir(context.TODO(), req)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but Rmdir returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but Rmdir returned error: %v", err)
		}
	}
}
