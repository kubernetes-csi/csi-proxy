package smb

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
)

type fakeSmbAPI struct{}

func (fakeSmbAPI) NewSmbGlobalMapping(remotePath, username, password string) error {
	return nil
}

func (fakeSmbAPI) RemoveSmbGlobalMapping(remotePath string) error {
	return nil
}

func TestNewSmbGlobalMapping(t *testing.T) {
	v1alpha1, err := apiversion.NewVersion("v1alpha1")
	if err != nil {
		t.Fatalf("New version error: %v", err)
	}
	testCases := []struct {
		remote      string
		username    string
		password    string
		version     apiversion.Version
		expectError bool
	}{
		{
			remote:      "",
			username:    "",
			password:    "",
			version:     v1alpha1,
			expectError: true,
		},
		{
			remote:      "\\test\\path",
			username:    "",
			password:    "",
			version:     v1alpha1,
			expectError: false,
		},
	}
	srv, err := NewServer(&fakeSmbAPI{})
	if err != nil {
		t.Fatalf("Smb Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		req := &internal.NewSmbGlobalMappingRequest{
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
