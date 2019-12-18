package smb

import (
	"context"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/smb/internal"
)

type fakeSmbAPI struct{}

func (fakeSmbAPI) NewSmbGlobalMapping(target, source, username, password string, readOnly bool) error {
	return nil
}

func (fakeSmbAPI) RemoveSmbGlobalMapping(target, source string) error {
	return nil
}

func TestNewSmbGlobalMapping(t *testing.T) {
	v1alpha1, err := apiversion.NewVersion("v1alpha1")
	testCases := []struct {
		remote		string
		local 		string
		username	string
		password	string
		readOnly	bool
		version 	apiversion.Version
		expectError bool
	} {
		{
			remote: "",
			local: "",
			username: "",
			password: "",
			readOnly: true,
			version: v1alpha1,
			expectError: true,
		},
		{
			remote: "\\test\\path",
			local: "10.1.1.1\\share",
			username: "",
			password: "",
			readOnly: true,
			version: v1alpha1,
			expectError: false,
		},
	}
	srv, err := NewServer(&fakeSmbAPI{})
	if err != nil {
		t.Fatalf("Smb Server could not be initialized for testing: %v", err)
	}
	for _, tc := range testCases {
		req := &internal.NewSmbGlobalMappingRequest{
			RemotePath: 	tc.remote,
			LocalPath:		tc.local,
			Username:		tc.username,
			Password:		tc.password,
			ReadOnly:		tc.readOnly,
		}
		response, _ := srv.NewSmbGlobalMapping(context.TODO(), req, tc.version)
		if tc.expectError && response.Error == "" {
			t.Errorf("Expected error but MountSmbShare returned a nil error")
		}
		if !tc.expectError && response.Error != "" {
			t.Errorf("Expected no errors but MountSmbShare returned error: %s", response.Error)
		}
	}
}