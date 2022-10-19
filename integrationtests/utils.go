package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/pkg/volume"
)

// getKubeletPathForTest returns the path to the current working directory
// to be used anytime the filepath is required to be within context of csi-proxy
func getKubeletPathForTest(dir string, t *testing.T) string {
	return filepath.Join("C:\\var\\lib\\kubelet", "testdir", dir)
}

// returns true if CSI_PROXY_GH_ACTIONS is set to "TRUE"
func isRunningOnGhActions() bool {
	return os.Getenv("CSI_PROXY_GH_ACTIONS") == "TRUE"
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
func shouldRunISCSITests() bool {
	return os.Getenv("ENABLE_ISCSI_TESTS") == "TRUE"
}

func runPowershellCmd(t *testing.T, command string) (string, error) {
	cmd := exec.Command("powershell", "/c", fmt.Sprintf("& { $global:ProgressPreference = 'SilentlyContinue'; %s }", command))
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

func getTestPluginPath() (string, int) {
	s1 := rand.NewSource(time.Now().UTC().UnixNano())
	r1 := rand.New(s1)

	testId := r1.Intn(1000000)
	return fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\", testId), testId
}

func diskInit(t *testing.T) (*VirtualHardDisk, func()) {
	testPluginPath, testId := getTestPluginPath()
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
	diskNumUnparsed = strings.TrimSpace(diskNumUnparsed)
	if diskNum, err = strconv.ParseUint(diskNumUnparsed, 10, 32); err != nil {
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

// sizeIsAround returns true if the actual size is around the expected size
// (considers the fact that some bytes were lost)
func sizeIsAround(t *testing.T, actualSize, expectedSize int64) bool {
	// An upper bound on the number of bytes that are lost when creating or resizing a partition
	var volumeSizeBytesLoss int64 = (20 * 1024 * 1024)
	var lowerBound = expectedSize - volumeSizeBytesLoss
	var upperBound = expectedSize
	t.Logf("Checking that the size is inside the bounds: %d < (actual) %d < %d", lowerBound, actualSize, upperBound)
	return lowerBound <= actualSize && actualSize <= upperBound
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// volumeInit initializes a volume, it creates a VHD, initializes it,
// creates a partition with the max size and formats the volume corresponding to that partition
func volumeInit(volumeClient volume.Interface, t *testing.T) (*VirtualHardDisk, string, func()) {
	vhd, vhdCleanup := diskInit(t)

	listRequest := &volume.ListVolumesOnDiskRequest{
		DiskNumber: vhd.DiskNumber,
	}
	listResponse, err := volumeClient.ListVolumesOnDisk(context.TODO(), listRequest)
	if err != nil {
		t.Fatalf("List response: %v", err)
	}

	volumeIDsLen := len(listResponse.VolumeIds)
	if volumeIDsLen != 1 {
		t.Fatalf("Number of volumes not equal to 1: %d", volumeIDsLen)
	}
	volumeID := listResponse.VolumeIds[0]
	t.Logf("VolumeId %v", volumeID)

	isVolumeFormattedRequest := &volume.IsVolumeFormattedRequest{
		VolumeId: volumeID,
	}
	isVolumeFormattedResponse, err := volumeClient.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err != nil {
		t.Fatalf("Is volume formatted request error: %v", err)
	}
	if isVolumeFormattedResponse.Formatted {
		t.Fatal("Volume formatted. Unexpected !!")
	}

	formatVolumeRequest := &volume.FormatVolumeRequest{
		VolumeId: volumeID,
	}
	_, err = volumeClient.FormatVolume(context.TODO(), formatVolumeRequest)
	if err != nil {
		t.Fatalf("Volume format failed. Error: %v", err)
	}

	isVolumeFormattedResponse, err = volumeClient.IsVolumeFormatted(context.TODO(), isVolumeFormattedRequest)
	if err != nil {
		t.Fatalf("Is volume formatted request error: %v", err)
	}
	if !isVolumeFormattedResponse.Formatted {
		t.Fatal("Volume should be formatted. Unexpected !!")
	}
	return vhd, volumeID, vhdCleanup
}
