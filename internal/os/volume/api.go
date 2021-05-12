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

// API exposes the internal volume operations available in the server
type API interface {
	ListVolumesOnDisk(diskNumber int64) (volumeIDs []string, err error)
	// MountVolume mounts the volume at the requested global staging target path
	MountVolume(volumeID, targetPath string) error
	// UnmountVolume gracefully dismounts a volume
	UnmountVolume(volumeID, targetPath string) error
	// IsVolumeFormatted checks if a volume is formatted with NTFS
	IsVolumeFormatted(volumeID string) (bool, error)
	// FormatVolume formats a volume with the provided file system
	FormatVolume(volumeID string) error
	// ResizeVolume performs resizing of the partition and file system for a block based volume
	ResizeVolume(volumeID string, sizeBytes int64) error
	// GetVolumeStats gets the volume information
	GetVolumeStats(volumeID string) (int64, int64, error)
	// GetDiskNumberFromVolumeID returns the disk number for where the volume is at
	GetDiskNumberFromVolumeID(volumeID string) (int64, error)
	// GetVolumeIDFromTargetPath returns the volume id of a given target path
	GetVolumeIDFromTargetPath(targetPath string) (string, error)
	// WriteVolumeCache writes volume cache to disk
	WriteVolumeCache(volumeID string) error
}

// VolumeAPI implements the internal Volume APIs
type VolumeAPI struct{}

// verifies that the API is implemented
var _ API = &VolumeAPI{}

// New - Construct a new Volume API Implementation.
func New() VolumeAPI {
	return VolumeAPI{}
}

func runExec(cmd string) ([]byte, error) {
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	return out, err
}

func getVolumeSize(volumeID string) (int64, error) {
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" | Get-partition).Size", volumeID)
	out, err := runExec(cmd)

	if err != nil || len(out) == 0 {
		return -1, fmt.Errorf("error getting size of the partition from mount. cmd %s, output: %s, error: %v", cmd, string(out), err)
	}

	outString := strings.TrimSpace(string(out))
	volumeSize, err := strconv.ParseInt(outString, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("error parsing size of volume %s received %v trimmed to %v err %v", volumeID, out, outString, err)
	}

	return volumeSize, nil
}

// ListVolumesOnDisk - returns back list of volumes(volumeIDs) in the disk.
func (VolumeAPI) ListVolumesOnDisk(diskNumber int64) (volumeIDs []string, err error) {
	cmd := fmt.Sprintf("(Get-Disk -Number %d |Get-Partition | Get-Volume).UniqueId", diskNumber)
	out, err := runExec(cmd)
	if err != nil {
		return []string{}, fmt.Errorf("error list volumes on disk. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	volumeIds := strings.Split(strings.TrimSpace(string(out)), "\r\n")
	return volumeIds, nil
}

// FormatVolume - Format a volume with a pre specified filesystem (typically ntfs)
func (VolumeAPI) FormatVolume(volumeID string) (err error) {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Format-Volume  -FileSystem %s -Confirm:$false", volumeID, formatFilesystem)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error formatting volume. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	// TODO: Do we need to handle anything for len(out) == 0
	return nil
}

// WriteVolumeCache - Writes the file system cache to disk with the given volume id
func (VolumeAPI) WriteVolumeCache(volumeID string) (err error) {
	return writeCache(volumeID)
}

// IsVolumeFormatted - Check if the volume is formatted with the pre specified filesystem(typically ntfs).
func (VolumeAPI) IsVolumeFormatted(volumeID string) (bool, error) {
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" -ErrorAction Stop).FileSystemType", volumeID)
	out, err := runExec(cmd)
	if err != nil {
		return false, fmt.Errorf("error checking if volume is formatted. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	stringOut := strings.TrimSpace(string(out))
	if len(stringOut) == 0 || strings.EqualFold(stringOut, "Unknown") {
		return false, nil
	}
	return true, nil
}

// MountVolume - mounts a volume to a path. This is done using the Add-PartitionAccessPath for presenting the volume via a path.
func (VolumeAPI) MountVolume(volumeID, path string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Add-PartitionAccessPath -AccessPath %s", volumeID, path)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error mount volume to path. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	return nil
}

// UnmountVolume - unmounts the volume path by removing the partition access path
func (VolumeAPI) UnmountVolume(volumeID, path string) error {
	if err := writeCache(volumeID); err != nil {
		return err
	}
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Remove-PartitionAccessPath -AccessPath %s", volumeID, path)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error getting driver letter to mount volume. cmd: %s, output: %s,error: %v", cmd, string(out), err)
	}
	return nil
}

// ResizeVolume - resizes a volume with the given size, if size == 0 then max supported size is used
func (VolumeAPI) ResizeVolume(volumeID string, size int64) error {
	// If size is 0 then we will resize to the maximum size possible, otherwise just resize to size
	var cmd string
	var out []byte
	var err error
	var finalSize int64
	var outString string
	if size == 0 {
		cmd = fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-partition | Get-PartitionSupportedSize | Select SizeMax | ConvertTo-Json", volumeID)
		out, err = runExec(cmd)

		if err != nil || len(out) == 0 {
			return fmt.Errorf("error getting sizemin,sizemax from mount. cmd: %s, output: %s, error: %v", cmd, string(out), err)
		}

		var getVolumeSizing map[string]int64
		outString = string(out)
		err = json.Unmarshal([]byte(outString), &getVolumeSizing)
		if err != nil {
			return fmt.Errorf("out %v outstring %v err %v", out, outString, err)
		}

		sizeMax := getVolumeSizing["SizeMax"]

		finalSize = sizeMax
	} else {
		finalSize = size
	}

	currentSize, err := getVolumeSize(volumeID)
	if err != nil {
		return fmt.Errorf("error getting the current size of volume (%s) with error (%v)", volumeID, err)
	}

	//if the partition's size is already the size we want this is a noop, just return
	if currentSize >= finalSize {
		return nil
	}

	cmd = fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-partition | Resize-Partition -Size %d", volumeID, finalSize)
	out, err = runExec(cmd)
	if err != nil {
		return fmt.Errorf("error resizing volume. cmd: %s, output: %s size:%v, finalSize %v, error: %v", cmd, string(out), size, finalSize, err)
	}
	return nil
}

// GetVolumeStats - retrieves the volume stats for a given volume
func (VolumeAPI) GetVolumeStats(volumeID string) (int64, int64, error) {
	// get the size and sizeRemaining for the volume
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" | Select SizeRemaining,Size) | ConvertTo-Json", volumeID)
	out, err := runExec(cmd)

	if err != nil {
		return -1, -1, fmt.Errorf("error getting capacity and used size of volume. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	var getVolume map[string]int64
	outString := string(out)
	err = json.Unmarshal([]byte(outString), &getVolume)
	if err != nil {
		return -1, -1, fmt.Errorf("out %v outstring %v err %v", out, outString, err)
	}
	var volumeSizeRemaining int64
	var volumeSize int64

	volumeSize = getVolume["Size"]
	volumeSizeRemaining = getVolume["SizeRemaining"]

	volumeUsedSize := volumeSize - volumeSizeRemaining
	return volumeSizeRemaining, volumeUsedSize, nil
}

// GetDiskNumberFromVolumeID - gets the disk number where the volume is.
func (VolumeAPI) GetDiskNumberFromVolumeID(volumeID string) (int64, error) {
	// get the size and sizeRemaining for the volume
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" | Get-Partition).DiskNumber", volumeID)
	out, err := runExec(cmd)

	if err != nil || len(out) == 0 {
		return -1, fmt.Errorf("error getting disk number. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return -1, fmt.Errorf("error compiling regex. err: %v", err)
	}
	diskNumberOutput := reg.ReplaceAllString(string(out), "")

	diskNumber, err := strconv.ParseInt(diskNumberOutput, 10, 64)

	if err != nil {
		return -1, fmt.Errorf("error parsing disk number. cmd: %s, output: %s, error: %v", cmd, diskNumberOutput, err)
	}

	return diskNumber, nil
}

// GetVolumeIDFromTargetPath - gets the volume ID given a mount point, the function is recursive until it find a volume or errors out
func (VolumeAPI) GetVolumeIDFromTargetPath(mount string) (string, error) {
	volumeString, err := getTarget(mount)

	if err != nil {
		return "", fmt.Errorf("error getting the volume for the mount %s, internal error %v", mount, err)
	}

	return volumeString, nil
}

func getTarget(mount string) (string, error) {
	cmd := fmt.Sprintf("(Get-Item -Path %s).Target", mount)
	out, err := runExec(cmd)
	if err != nil || len(out) == 0 {
		return "", fmt.Errorf("error getting volume from mount. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	volumeString := strings.TrimSpace(string(out))
	if !strings.HasPrefix(volumeString, "Volume") {
		return getTarget(volumeString)
	}

	volumeString = "\\\\?\\" + volumeString

	return volumeString, nil
}

func writeCache(volumeID string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Write-Volumecache", volumeID)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error writing volume cache. cmd: %s, output: %s, error: %v", cmd, string(out), err)
	}
	return nil
}
