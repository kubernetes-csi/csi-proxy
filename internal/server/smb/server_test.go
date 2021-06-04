package smb

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/os/filesystem"
	fsserver "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
)

type fakeSmbAPI struct{}

func (fakeSmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	return nil
}

func (fakeSmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	return nil
}

func (fakeSmbAPI) IsSmbMapped(remotePath string) (bool, error) {
	return false, nil
}

func (fakeSmbAPI) NewSmbLink(remotePath, localPath string) error {
	return nil
}

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
func (fakeFileSystemAPI) CreateSymlink(tgt string, src string) error {
	return nil
}

func (fakeFileSystemAPI) IsSymlink(path string) (bool, error) {
	return true, nil
}

func TestNewSmbGlobalMapping(t *testing.T) {
	v1, err := apiversion.NewVersion("v1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}
	testCases := []struct {
		remote      string
		local       string
		username    string
		password    string
		version     apiversion.Version
		expectError bool
	}{
		{
			remote:      "",
			username:    "",
			password:    "",
			version:     v1,
			expectError: true,
		},
		{
			remote:      "\\test\\path",
			username:    "",
			password:    "",
			version:     v1,
			expectError: false,
		},
	}
	fsSrv, err := fsserver.NewServer(`C:\var\lib\kubelet\plugins`, `C:\var\lib\kubelet\pods`, &fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem Server could not be initialized for testing: %v", err)
	}

	srv, err := NewServer(&fakeSmbAPI{}, fsSrv)
	if err != nil {
		t.Fatalf("Smb Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		req := &internal.NewSmbGlobalMappingRequest{
			LocalPath:  tc.local,
			RemotePath: tc.remote,
			Username:   tc.username,
			Password:   tc.password,
		}
		response, err := srv.NewSmbGlobalMapping(context.TODO(), req, tc.version)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but NewSmbGlobalMapping returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but NewSmbGlobalMapping returned error: %s", response.Error)
		}
	}
}
