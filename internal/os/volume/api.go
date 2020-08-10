package volume

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const formatFilesystem = "ntfs"

// VolAPIImplementor - struct for implementing the internal Volume APIs
type VolAPIImplementor struct{}

// New - Construct a new Volume API Implementation.
func New() VolAPIImplementor {
	return VolAPIImplementor{}
}

func runExec(cmd string) ([]byte, error) {
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	return out, err
}

// ListVolumesOnDisk - returns back list of volumes(volumeIDs) in the disk (requested in diskID).
func (VolAPIImplementor) ListVolumesOnDisk(diskID string) (volumeIDs []string, err error) {
	cmd := fmt.Sprintf("(Get-Disk -DeviceId %s |Get-Partition | Get-Volume).UniqueId", diskID)
	out, err := runExec(cmd)
	if err != nil {
		return []string{}, fmt.Errorf("error list volumes on disk. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	volumeIds := strings.Split(strings.TrimSpace(string(out)), "\r\n")
	return volumeIds, nil
}

// FormatVolume - Format a volume with a pre specified filesystem (typically ntfs)
func (VolAPIImplementor) FormatVolume(volumeID string) (err error) {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Format-Volume  -FileSystem %s -Confirm:$false", volumeID, formatFilesystem)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error formatting volume. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	// TODO: Do we need to handle anything for len(out) == 0
	return nil
}

// IsVolumeFormatted - Check if the volume is formatted with the pre specified filesystem(typically ntfs).
func (VolAPIImplementor) IsVolumeFormatted(volumeID string) (bool, error) {
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" -ErrorAction Stop).FileSystemType", volumeID)
	out, err := runExec(cmd)
	if err != nil {
		return false, fmt.Errorf("error checking if volume is formatted. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	if len(out) == 0 || !strings.EqualFold(strings.TrimSpace(string(out)), formatFilesystem) {
		return false, nil
	}
	return true, nil
}

// MountVolume - mounts a volume to a path. This is done using the Add-PartitionAccessPath for presenting the volume via a path.
func (VolAPIImplementor) MountVolume(volumeID, path string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Add-PartitionAccessPath -AccessPath %s", volumeID, path)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error mount volume to path. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	return nil
}

// DismountVolume - unmounts the volume path by removing the patition access path
func (VolAPIImplementor) DismountVolume(volumeID, path string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Remove-PartitionAccessPath -AccessPath %s", volumeID, path)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error getting driver letter to mount volume. cmd: %s, output: %s,error: %v", cmd, string(out), err)
	}
	return nil
}

// ResizeVolume - resize the volume to the size specified as parameter.
func (VolAPIImplementor) ResizeVolume(volumeID string, size int64) error {
	// TODO: Check the size of the resize
	// TODO: We have to get the right partition.
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-partition | Resize-Partition -Size %d", volumeID, size)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error resizing volume. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	return nil
}

// VolumeStats - resize the volume to the size specified as parameter.
func (VolAPIImplementor) VolumeStats(volumeID string) (int64, int64, int64, error) {

	//first get the disk size where the volume resides
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" | Get-partition | Get-Disk).Size", volumeID)
	out, err := runExec(cmd)
	if err != nil || len(out) == 0 {
		return -1, -1, -1, fmt.Errorf("error getting size of disk. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return -1, -1, -1, fmt.Errorf("error compiling regex. err: %v", err)
	}
	diskSizeOutput := reg.ReplaceAllString(string(out), "")

	diskSize, err := strconv.ParseInt(diskSizeOutput, 10, 64)
	if err != nil {
		return -1, -1, -1, fmt.Errorf("error parsing size of disk. cmd: %s, output: %s, error: %v", cmd, diskSizeOutput, err)
	}

	// secondly, get the size and sizeRemaining for the volume
	cmd = fmt.Sprintf("(Get-Volume -UniqueId \"%s\" | Select SizeRemaining,Size) | ConvertTo-Json", volumeID)
	out, err = runExec(cmd)

	if err != nil {
		return -1, -1, -1, fmt.Errorf("error getting capacity and used size of volume. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	var getVolume map[string]int64
	outString := string(out)
	err = json.Unmarshal([]byte(outString), &getVolume)
	if err != nil {
		return -1, -1, -1, fmt.Errorf("out %v outstring %v err %v", out, outString, err)
	}
	var volumeSizeRemaining int64
	var volumeSize int64

	volumeSize = getVolume["Size"]
	volumeSizeRemaining = getVolume["SizeRemaining"]

	volumeUsedSize := volumeSize - volumeSizeRemaining
	return diskSize, volumeSize, volumeUsedSize, nil
}
