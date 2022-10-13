package smb

import (
	"context"
	"testing"

	fs "github.com/kubernetes-csi/csi-proxy/pkg/filesystem"
	fsapi "github.com/kubernetes-csi/csi-proxy/pkg/filesystem/api"
	smbapi "github.com/kubernetes-csi/csi-proxy/pkg/smb/api"
)

type fakeSmbAPI struct{}

var _ smbapi.API = &fakeSmbAPI{}

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

var _ fsapi.API = &fakeFileSystemAPI{}

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
	testCases := []struct {
		remote      string
		local       string
		username    string
		password    string
		expectError bool
	}{
		{
			remote:      "",
			username:    "",
			password:    "",
			expectError: true,
		},
		{
			remote:      "\\\\hostname\\path",
			username:    "",
			password:    "",
			expectError: false,
		},
	}
	fsClient, err := fs.New(&fakeFileSystemAPI{})
	if err != nil {
		t.Fatalf("FileSystem client could not be initialized for testing: %v", err)
	}

	client, err := New(&fakeSmbAPI{}, fsClient)
	if err != nil {
		t.Fatalf("Smb client could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		req := &NewSmbGlobalMappingRequest{
			LocalPath:  tc.local,
			RemotePath: tc.remote,
			Username:   tc.username,
			Password:   tc.password,
		}
		_, err := client.NewSmbGlobalMapping(context.TODO(), req)
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
			if tc.expectResult != result {
				t.Errorf("Expected (%s) but getRootMappingPath returned (%s)", tc.expectResult, result)
			}
		}
	}
}
