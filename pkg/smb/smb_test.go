package smb

import (
	"context"
	"testing"

	fs "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem"
	fsapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem/hostapi"
	smbapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/smb/hostapi"
)

type fakeSMBAPI struct{}

var _ smbapi.HostAPI = &fakeSMBAPI{}

func (fakeSMBAPI) NewSMBGlobalMapping(remotePath, username, password string) error {
	return nil
}

func (fakeSMBAPI) RemoveSMBGlobalMapping(remotePath string) error {
	return nil
}

func (fakeSMBAPI) IsSMBMapped(remotePath string) (bool, error) {
	return false, nil
}

func (fakeSMBAPI) NewSMBLink(remotePath, localPath string) error {
	return nil
}

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

func TestNewSMBGlobalMapping(t *testing.T) {
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

	client, err := New(&fakeSMBAPI{}, fsClient)
	if err != nil {
		t.Fatalf("SMB client could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		req := &NewSMBGlobalMappingRequest{
			LocalPath:  tc.local,
			RemotePath: tc.remote,
			Username:   tc.username,
			Password:   tc.password,
		}
		_, err := client.NewSMBGlobalMapping(context.TODO(), req)
		if tc.expectError && err == nil {
			t.Errorf("Expected error but NewSMBGlobalMapping returned a nil error")
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no errors but NewSMBGlobalMapping returned error: %v", err)
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
