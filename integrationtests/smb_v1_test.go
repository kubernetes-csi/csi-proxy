package integrationtests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-proxy/client/api/smb/v1"
	client "github.com/kubernetes-csi/csi-proxy/client/groups/smb/v1"

	"github.com/stretchr/testify/assert"
)

func v1SmbTests(t *testing.T) {
	client, err := client.NewClient()
	if err != nil {
		t.Fatalf("Fail to get smb API group client %v", err)
	}
	defer client.Close()

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
	mountSmbShareReq := &v1.NewSmbGlobalMappingRequest{
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

	unmountSmbShareReq := &v1.RemoveSmbGlobalMappingRequest{
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
