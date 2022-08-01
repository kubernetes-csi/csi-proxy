package smb

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/pkg/os/filesystem"
	"github.com/kubernetes-csi/csi-proxy/pkg/os/smb"
	fsserver "github.com/kubernetes-csi/csi-proxy/pkg/server/filesystem"
	internal "github.com/kubernetes-csi/csi-proxy/pkg/server/smb/impl"
)

type fakeSmbAPI struct{}

var _ smb.API = &fakeSmbAPI{}

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
func (fakeFileSystemAPI) RmdirContents(path string) error {
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
			remote:      "\\\\hostname\\path",
			username:    "",
			password:    "",
			version:     v1,
			expectError: false,
		},
	}
	fsSrv, err := fsserver.NewServer([]string{`C:\var\lib\kubelet`}, &fakeFileSystemAPI{})
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
		_, err := srv.NewSmbGlobalMapping(context.TODO(), req, tc.version)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but NewSmbGlobalMapping returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but NewSmbGlobalMapping returned error: %v", err)
		}
	}
}

func TestGetRootMappingPath(t *testing.T) {
	testCases := []struct {
		remote       string
		expectResult string
		expectError  bool
	}{
		{
			remote:       "",
			expectResult: "",
			expectError:  true,
		},
		{
			remote:       "hostname",
			expectResult: "",
			expectError:  true,
		},
		{
			remote:       "\\\\hostname\\path",
			expectResult: "\\\\hostname\\path",
			expectError:  false,
		},
		{
			remote:       "\\\\hostname\\path\\",
			expectResult: "\\\\hostname\\path",
			expectError:  false,
		},
		{
			remote:       "\\\\hostname\\path\\subpath",
			expectResult: "\\\\hostname\\path",
			expectError:  false,
		},
	}
	for _, tc := range testCases {
		result, err := getRootMappingPath(tc.remote)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but getRootMappingPath returned a nil error")
		}
		if !tc.expectError {
			if err != nil {
				t.Errorf("Expected no errors but getRootMappingPath returned error: %v", err)
			}
			if expectResult != result {
				t.Errorf("Expected (%s) but getRootMappingPath returned (%s)", expectResult, result)
			} 
		}
	}
}