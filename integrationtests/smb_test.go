package integrationtests

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	fs "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem"
	fsapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/filesystem/api"
	"github.com/kubernetes-csi/csi-proxy/v2/pkg/smb"
	smbapi "github.com/kubernetes-csi/csi-proxy/v2/pkg/smb/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const letterset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomString generates a random string with specified length
func randomString(length int) string {
	return stringWithCharset(length, letterset)
}

func setupUser(username, password string) error {
	cmdLine := fmt.Sprintf(`$PWord = ConvertTo-SecureString $Env:password -AsPlainText -Force` +
		`;New-Localuser -name $Env:username -accountneverexpires -password $PWord`)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("username=%s", username),
		fmt.Sprintf("password=%s", password))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("setupUser failed: %v, output: %q", err, string(output))
	}
	return nil
}

func removeUser(t *testing.T, username string) {
	cmdLine := fmt.Sprintf(`Remove-Localuser -name $Env:username`)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("username=%s", username))
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("setupUser failed: %v, output: %q", err, string(output))
	}
}

func setupSmbShare(shareName, localPath, username string) error {
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("setupSmbShare failed to create local path %q: %v", localPath, err)
	}
	cmdLine := fmt.Sprintf(`New-SMBShare -Name $Env:sharename -Path $Env:path -fullaccess $Env:username`)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("sharename=%s", shareName),
		fmt.Sprintf("path=%s", localPath),
		fmt.Sprintf("username=%s", username))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("setupSmbShare failed: %v, output: %q", err, string(output))
	}

	return nil
}

func removeSmbShare(t *testing.T, shareName string) {
	cmdLine := fmt.Sprintf(`Remove-SMBShare -Name $Env:sharename -Force`)
	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("sharename=%s", shareName))
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("setupSmbShare failed: %v, output: %q", err, string(output))
	}
	return
}

func getSmbGlobalMapping(remotePath string) error {
	// use PowerShell Environment Variables to store user input string to prevent command line injection
	// https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-5.1
	cmdLine := fmt.Sprintf(`(Get-SmbGlobalMapping -RemotePath $Env:smbremotepath).Status`)

	cmd := exec.Command("powershell", "/c", cmdLine)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("smbremotepath=%s", remotePath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Get-SmbGlobalMapping failed: %v, output: %q", err, string(output))
	}
	if !strings.Contains(string(output), "OK") {
		return fmt.Errorf("Get-SmbGlobalMapping return status %q instead of OK", string(output))
	}
	return nil
}

func writeReadFile(path string) error {
	fileName := path + "\\hello.txt"
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file %q failed: %v", fileName, err)
	}
	defer f.Close()
	fileContent := "Hello World"
	if _, err = f.WriteString(fileContent); err != nil {
		return fmt.Errorf("write to file %q failed: %v", fileName, err)
	}
	if err = f.Sync(); err != nil {
		return fmt.Errorf("sync file %q failed: %v", fileName, err)
	}
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file %q failed: %v", fileName, err)
	}
	if fileContent != string(dat) {
		return fmt.Errorf("read content of file %q failed: expected %q, got %q", fileName, fileContent, string(dat))
	}
	return nil
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
