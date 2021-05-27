package integrationtests

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/internal/server"
	srvtypes "github.com/kubernetes-csi/csi-proxy/internal/server/types"
)

// startServer starts the proxy's GRPC servers, and returns a function to shut them down when done with testing
func startServer(t *testing.T, apiGroups ...srvtypes.APIGroup) func() {
	s := server.NewServer(apiGroups...)

	listeningChan := make(chan interface{})
	go func() {
		assert.Nil(t, s.Start(listeningChan))
	}()

	select {
	case <-listeningChan:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for GRPC servers to start listening")
	}

	return func() {
		assert.Nil(t, s.Stop())
	}
}

func close(t *testing.T, closer io.Closer) {
	assert.Nil(t, closer.Close())
}

// recursiveDiff ensures that dir1 and dir2 contain the same files, with the same contents.
// fileSuffixesToRemove will be removed from file names, if found.
func recursiveDiff(t *testing.T, dir1, dir2 string, fileSuffixesToRemove ...string) {
	hashesDir1, err := fileHashes(dir1, fileSuffixesToRemove...)
	require.Nil(t, err, "unable to get file hashes for directory %q", dir1)
	hashesDir2, err := fileHashes(dir2, fileSuffixesToRemove...)
	require.Nil(t, err, "unable to get file hashes for directory %q", dir2)

	t.Logf("Hashes for dir1: %+v", hashesDir1)
	t.Logf("Hashes for dir2: %+v", hashesDir2)

	for filePath, hash1 := range hashesDir1 {
		if hash2, present := hashesDir2[filePath]; assert.True(t, present, "%q present in %q but not in %q", filePath, dir1, dir2) {
			if hash1 != hash2 {
				contents1 := readFile(t, filepath.Join(dir1, filePath))
				contents2 := readFile(t, filepath.Join(dir2, filePath))

				differ := diffmatchpatch.New()
				diffs := differ.DiffMain(contents1, contents2, true)
				assert.Fail(t, fmt.Sprintf("File %q differs in %q and %q:\n", filePath, dir1, dir2), "Diff:\n%s", differ.DiffPrettyText(diffs))
			}
			delete(hashesDir2, filePath)
		}
	}

	for filePath := range hashesDir2 {
		assert.Fail(t, fmt.Sprintf("%q present in %q but not in %q", filePath, dir2, dir1))
	}
}

// fileHashes walks through dir, and returns a map mapping file paths to MD5 hashes
func fileHashes(dir string, fileSuffixesToRemove ...string) (map[string]string, error) {
	dir = strings.ReplaceAll(dir, "/", string(os.PathSeparator))
	hashes := make(map[string]string)

	if walkErr := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "unable to descend into %q", filePath)
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "unable to open %q", filePath)
		}
		defer file.Close()

		hasher := md5.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return errors.Wrapf(err, "unable to read %q", filePath)
		}

		hashBytes := hasher.Sum(nil)[:16]

		relativePath := strings.TrimPrefix(strings.TrimPrefix(filePath, dir), "/")
		for _, suffix := range fileSuffixesToRemove {
			relativePath = strings.TrimSuffix(relativePath, suffix)
		}

		hashes[relativePath] = hex.EncodeToString(hashBytes)
		return nil
	}); walkErr != nil {
		return nil, walkErr
	}

	return hashes, nil
}

func readFile(t *testing.T, filePath string) string {
	contents, err := ioutil.ReadFile(filePath)
	require.Nil(t, err, "unable to read %q", filePath)
	return string(contents)
}

// GetWorkDirPath returns the path to the current working directory
// to be used anytime the filepath is required to be within context of csi-proxy
func getWorkDirPath(dir string, t *testing.T) string {
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}
	return filepath.Join(path, "testdir", dir)
}

// returns true if CSI_PROXY_GH_ACTIONS is set to "TRUE"
func isRunningOnGhActions() bool {
	return os.Getenv("CSI_PROXY_GH_ACTIONS") == "TRUE"
}

// returns true if underlying os is windows
func isRunningWindows() bool {
	return runtime.GOOS == "windows"
}

func skipTestOnCondition(t *testing.T, condition bool) {
	if condition {
		t.Skip("Skipping test")
	}
}

// returns true if ENABLE_ISCSI_TESTS is set to "TRUE"
// used to skip iSCSI tests as they affect the test machine
// e.g. install an iSCSI target, format a disk, etc.
// Take care to use disposable clean VMs for tests
func shouldRunIscsiTests() bool {
	return os.Getenv("ENABLE_ISCSI_TESTS") == "TRUE"
}

func runPowershellCmd(t *testing.T, command string) (string, error) {
	cmd := exec.Command("powershell", "/c", command)
	t.Logf("Executing command: %q", cmd.String())
	result, err := cmd.CombinedOutput()
	return string(result), err
}

func diskCleanup(t *testing.T, vhdxPath, mountPath, testPluginPath string) {
	if t.Failed() {
		t.Logf("Test failed. Skipping cleanup!")
		t.Logf("Mount path located at %s", mountPath)
		t.Logf("VHDx path located at %s", vhdxPath)
		return
	}
	var cmd, out string
	var err error

	cmd = fmt.Sprintf("Dismount-VHD -Path %s", vhdxPath)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, out)
	}
	cmd = fmt.Sprintf("rm %s", vhdxPath)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, out)
	}
	cmd = fmt.Sprintf("rmdir %s", mountPath)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, out)
	}
	if testPluginPath != "" {
		cmd = fmt.Sprintf("rmdir %s", testPluginPath)
		if out, err = runPowershellCmd(t, cmd); err != nil {
			t.Errorf("Error: %v. Command: %s. Out: %s", err, cmd, out)
		}
	}
}

// VirtualHardDisk represents a VHD used for e2e tests
type VirtualHardDisk struct {
	DiskNumber     uint32
	Path           string
	Mount          string
	TestPluginPath string
	InitialSize    int64
}

func diskInit(t *testing.T) (*VirtualHardDisk, func()) {
	s1 := rand.NewSource(time.Now().UTC().UnixNano())
	r1 := rand.New(s1)

	testId := r1.Intn(1000)
	testPluginPath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\", testId)
	mountPath := fmt.Sprintf("%smount-%d", testPluginPath, testId)
	vhdxPath := fmt.Sprintf("%sdisk-%d.vhdx", testPluginPath, testId)

	var cmd, out string
	var err error
	const initialSize = 1 * 1024 * 1024 * 1024
	const partitionStyle = "GPT"

	cmd = fmt.Sprintf("mkdir %s", mountPath)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
	}

	// Initialize the tests, using powershell directly.
	// Create the new vhdx
	cmd = fmt.Sprintf("New-VHD -Path %s -SizeBytes %d", vhdxPath, initialSize)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s.", err, cmd, out)
	}

	// Mount the vhdx as a disk
	cmd = fmt.Sprintf("Mount-VHD -Path %s", vhdxPath)
	if out, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
	}

	var diskNum uint64
	var diskNumUnparsed string
	cmd = fmt.Sprintf("(Get-VHD -Path %s).DiskNumber", vhdxPath)

	if diskNumUnparsed, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}
	if diskNum, err = strconv.ParseUint(strings.TrimRight(diskNumUnparsed, "\r\n"), 10, 32); err != nil {
		t.Fatalf("Error: %v", err)
	}

	cmd = fmt.Sprintf("Initialize-Disk -Number %d -PartitionStyle %s", diskNum, partitionStyle)
	if _, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}

	cmd = fmt.Sprintf("New-Partition -DiskNumber %d -UseMaximumSize", diskNum)
	if _, err = runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}

	cleanup := func() {
		diskCleanup(t, vhdxPath, mountPath, testPluginPath)
	}

	vhd := &VirtualHardDisk{
		DiskNumber:     uint32(diskNum),
		Path:           vhdxPath,
		Mount:          mountPath,
		TestPluginPath: testPluginPath,
		InitialSize:    initialSize,
	}

	return vhd, cleanup
}
