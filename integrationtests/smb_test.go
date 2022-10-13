package integrationtests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fs "github.com/kubernetes-csi/csi-proxy/pkg/filesystem"
	fsapi "github.com/kubernetes-csi/csi-proxy/pkg/filesystem/api"
	"github.com/kubernetes-csi/csi-proxy/pkg/smb"
	smbapi "github.com/kubernetes-csi/csi-proxy/pkg/smb/api"
)

func TestSmbAPIGroup(t *testing.T) {
	t.Run("v1alpha1SmbTests", func(t *testing.T) {
		v1alpha1SmbTests(t)
	})
	t.Run("v1beta1SmbTests", func(t *testing.T) {
		v1beta1SmbTests(t)
	})
	t.Run("v1beta2SmbTests", func(t *testing.T) {
		v1beta2SmbTests(t)
	})
	t.Run("v1SmbTests", func(t *testing.T) {
		v1SmbTests(t)
	})
}

func TestSmb(t *testing.T) {
	fsClient, err := fs.New(fsapi.New())
	require.Nil(t, err)
	client, err := smb.New(smbapi.New(), fsClient)
	require.Nil(t, err)

	username := randomString(5)
	password := randomString(10) + "!"
	sharePath := fmt.Sprintf("C:\\smbshare%s", randomString(5))
	smbShare := randomString(6)

	localPath := fmt.Sprintf("C:\\localpath%s", randomString(5))

	if err = setupUser(username, password); err != nil {
		t.Fatalf("TestSmbAPIGroup %v", err)
	}
	defer removeUser(t, username)

	if err = setupSmbShare(smbShare, sharePath, username); err != nil {
		t.Fatalf("TestSmbAPIGroup %v", err)
	}
	defer removeSmbShare(t, smbShare)

	hostname, err := os.Hostname()
	assert.Nil(t, err)

	username = "domain\\" + username
	remotePath := "\\\\" + hostname + "\\" + smbShare
	// simulate Mount SMB operations around staging a volume on a node
	mountSmbShareReq := &smb.NewSmbGlobalMappingRequest{
		RemotePath: remotePath,
		Username:   username,
		Password:   password,
	}
	_, err = client.NewSmbGlobalMapping(context.Background(), mountSmbShareReq)
	if err != nil {
		t.Fatalf("TestSmbAPIGroup %v", err)
	}

	err = getSmbGlobalMapping(remotePath)
	assert.Nil(t, err)

	err = writeReadFile(remotePath)
	assert.Nil(t, err)

	unmountSmbShareReq := &smb.RemoveSmbGlobalMappingRequest{
		RemotePath: remotePath,
	}
	_, err = client.RemoveSmbGlobalMapping(context.Background(), unmountSmbShareReq)
	if err != nil {
		t.Fatalf("TestSmbAPIGroup %v", err)
	}
	err = getSmbGlobalMapping(remotePath)
	assert.NotNil(t, err)
	err = writeReadFile(localPath)
	assert.NotNil(t, err)
}
