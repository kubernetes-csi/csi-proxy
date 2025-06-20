package volume

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/pkg/cim"
	"github.com/kubernetes-csi/csi-proxy/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"k8s.io/klog/v2"
)

const (
	minimumResizeSize = 100 * 1024 * 1024
)

// API exposes the internal volume operations available in the server
type API interface {
	// ListVolumesOnDisk lists volumes on a disk identified by a `diskNumber` and optionally a partition identified by `partitionNumber`.
	ListVolumesOnDisk(diskNumber uint32, partitionNumber uint32) (volumeIDs []string, err error)
	// MountVolume mounts the volume at the requested global staging target path.
	MountVolume(volumeID, targetPath string) error
	// UnmountVolume gracefully dismounts a volume.
	UnmountVolume(volumeID, targetPath string) error
	// IsVolumeFormatted checks if a volume is formatted with NTFS.
	IsVolumeFormatted(volumeID string) (bool, error)
	// FormatVolume formats a volume with the NTFS format.
	FormatVolume(volumeID string) error
	// ResizeVolume performs resizing of the partition and file system for a block based volume.
	ResizeVolume(volumeID string, sizeBytes int64) error
	// GetVolumeStats gets the volume information.
	GetVolumeStats(volumeID string) (int64, int64, error)
	// GetDiskNumberFromVolumeID returns the disk number for a given volumeID.
	GetDiskNumberFromVolumeID(volumeID string) (uint32, error)
	// GetVolumeIDFromTargetPath returns the volume id of a given target path.
	GetVolumeIDFromTargetPath(targetPath string) (string, error)
	// WriteVolumeCache writes the volume `volumeID`'s cache to disk.
	WriteVolumeCache(volumeID string) error
	// GetClosestVolumeIDFromTargetPath returns the volume id of a given target path.
	GetClosestVolumeIDFromTargetPath(targetPath string) (string, error)
}

// VolumeAPI implements the internal Volume APIs
type VolumeAPI struct{}

// verifies that the API is implemented
var _ API = &VolumeAPI{}

var (
	// VolumeRegexp matches a Windows Volume
	// example: Volume{452e318a-5cde-421e-9831-b9853c521012}
	//
	// The field UniqueId has an additional prefix which is NOT included in the regex
	// however the regex can match UniqueId too
	// PS C:\disks> (Get-Disk -Number 1 | Get-Partition | Get-Volume).UniqueId
	// \\?\Volume{452e318a-5cde-421e-9831-b9853c521012}\
	VolumeRegexp = regexp.MustCompile(`Volume\{[\w-]*\}`)

	notMountedFolder = errors.New("not a mounted folder")
)

// New - Construct a new Volume API Implementation.
func New() VolumeAPI {
	return VolumeAPI{}
}

// ListVolumesOnDisk - returns back list of volumes(volumeIDs) in a disk and a partition.
func (VolumeAPI) ListVolumesOnDisk(diskNumber uint32, partitionNumber uint32) (volumeIDs []string, err error) {
	partitions, err := cim.ListPartitionsOnDisk(diskNumber, partitionNumber, cim.PartitionSelectorListObjectID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list partition on disk %d", diskNumber)
	}

	volumes, err := cim.FindVolumesByPartition(partitions)
	if cim.IgnoreNotFound(err) != nil {
		return nil, errors.Wrapf(err, "failed to list volumes on disk %d", diskNumber)
	}

	for _, volume := range volumes {
		uniqueID, err := cim.GetVolumeUniqueID(volume)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get unique ID for volume %v", volume)
		}
		volumeIDs = append(volumeIDs, uniqueID)
	}

	return volumeIDs, nil
}

// FormatVolume - Formats a volume with the NTFS format.
func (VolumeAPI) FormatVolume(volumeID string) (err error) {
	volume, err := cim.QueryVolumeByUniqueID(volumeID, nil)
	if err != nil {
		return fmt.Errorf("error formatting volume (%s). error: %v", volumeID, err)
	}

	result, err := cim.FormatVolume(volume,
		"NTFS", // Format,
		"",     // FileSystemLabel,
		nil,    // AllocationUnitSize,
		false,  // Full,
		true,   // Force
		nil,    // Compress,
		nil,    // ShortFileNameSupport,
		nil,    // SetIntegrityStreams,
		nil,    // UseLargeFRS,
		nil,    // DisableHeatGathering,
	)
	if result != 0 || err != nil {
		return fmt.Errorf("error formatting volume (%s). result: %d, error: %v", volumeID, result, err)
	}
	return nil
}

// WriteVolumeCache - Writes the file system cache to disk with the given volume id
func (VolumeAPI) WriteVolumeCache(volumeID string) (err error) {
	return writeCache(volumeID)
}

// IsVolumeFormatted - Check if the volume is formatted with the pre specified filesystem(typically ntfs).
func (VolumeAPI) IsVolumeFormatted(volumeID string) (bool, error) {
	volume, err := cim.QueryVolumeByUniqueID(volumeID, cim.VolumeSelectorListForFileSystemType)
	if err != nil {
		return false, fmt.Errorf("error checking if volume (%s) is formatted. error: %v", volumeID, err)
	}

	fsType, err := cim.GetVolumeFileSystemType(volume)
	if err != nil {
		return false, fmt.Errorf("failed to query volume file system type (%s): %w", volumeID, err)
	}

	return fsType != cim.FileSystemUnknown, nil
}

// MountVolume - mounts a volume to a path. This is done using Win32 API SetVolumeMountPoint for presenting the volume via a path.
func (VolumeAPI) MountVolume(volumeID, path string) error {
	mountPoint := path
	if !strings.HasSuffix(mountPoint, "\\") {
		mountPoint += "\\"
	}
	utf16MountPath, _ := windows.UTF16PtrFromString(mountPoint)
	utf16VolumeID, _ := windows.UTF16PtrFromString(volumeID)
	err := windows.SetVolumeMountPoint(utf16MountPath, utf16VolumeID)
	if err != nil {
		if errors.Is(windows.GetLastError(), windows.ERROR_DIR_NOT_EMPTY) {
			targetVolumeID, err := getTarget(path)
			if err != nil {
				return fmt.Errorf("error get target volume (%s) to path %s. error: %v", volumeID, path, err)
			}

			if volumeID == targetVolumeID {
				return nil
			}
		}

		return fmt.Errorf("error mount volume (%s) to path %s. error: %v", volumeID, path, err)
	}

	return nil
}

// UnmountVolume - unmounts the volume path by removing the partition access path
func (VolumeAPI) UnmountVolume(volumeID, path string) error {
	if err := writeCache(volumeID); err != nil {
		return err
	}

	mountPoint := path
	if !strings.HasSuffix(mountPoint, "\\") {
		mountPoint += "\\"
	}
	utf16MountPath, _ := windows.UTF16PtrFromString(mountPoint)
	err := windows.DeleteVolumeMountPoint(utf16MountPath)
	if err != nil {
		return fmt.Errorf("error umount volume (%s) from path %s. error: %v", volumeID, path, err)
	}
	return nil
}

// ResizeVolume - resizes a volume with the given size, if size == 0 then max supported size is used
func (VolumeAPI) ResizeVolume(volumeID string, size int64) error {
	var err error
	var finalSize int64
	part, err := cim.GetPartitionByVolumeUniqueID(volumeID)
	if err != nil {
		return err
	}

	// If size is 0 then we will resize to the maximum size possible, otherwise just resize to size
	if size == 0 {
		var result int
		var status string
		result, _, finalSize, status, err = cim.GetPartitionSupportedSize(part)
		if result != 0 || err != nil {
			return fmt.Errorf("error getting sizeMin, sizeMax from volume (%s). result: %d, status: %s, error: %v", volumeID, result, status, err)
		}

	} else {
		finalSize = size
	}

	currentSize, err := cim.GetPartitionSize(part)
	if err != nil {
		return fmt.Errorf("error getting the current size of volume (%s) with error (%v)", volumeID, err)
	}

	// only resize if finalSize - currentSize is greater than 100MB
	if finalSize-currentSize < minimumResizeSize {
		klog.V(2).Infof("minimum resize difference (100MB) not met, skipping resize. volumeID=%s currentSize=%d finalSize=%d", volumeID, currentSize, finalSize)
		return nil
	}

	//if the partition's size is already the size we want this is a noop, just return
	if currentSize >= finalSize {
		klog.V(2).Infof("Attempted to resize volume (%s) to a lower size, from currentBytes=%d wantedBytes=%d", volumeID, currentSize, finalSize)
		return nil
	}

	result, _, err := cim.ResizePartition(part, finalSize)
	if result != 0 || err != nil {
		return fmt.Errorf("error resizing volume (%s). size:%v, finalSize %v, error: %v", volumeID, size, finalSize, err)
	}

	diskNumber, err := cim.GetPartitionDiskNumber(part)
	if err != nil {
		return fmt.Errorf("error parsing disk number of volume (%s). error: %v", volumeID, err)
	}

	disk, err := cim.QueryDiskByNumber(diskNumber, nil)
	if err != nil {
		return fmt.Errorf("error query disk of volume (%s). error: %v", volumeID, err)
	}

	result, _, err = cim.RefreshDisk(disk)
	if result != 0 || err != nil {
		return fmt.Errorf("error rescan disk (%d). result %d, error: %v", diskNumber, result, err)
	}

	return nil
}

// GetVolumeStats - retrieves the volume stats for a given volume
func (VolumeAPI) GetVolumeStats(volumeID string) (int64, int64, error) {
	volume, err := cim.QueryVolumeByUniqueID(volumeID, cim.VolumeSelectorListForStats)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting capacity and used size of volume (%s). error: %v", volumeID, err)
	}

	volumeSize, err := cim.GetVolumeSize(volume)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to query volume size (%s): %w", volumeID, err)
	}

	volumeSizeRemaining, err := cim.GetVolumeSizeRemaining(volume)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to query volume remaining size (%s): %w", volumeID, err)
	}

	volumeUsedSize := volumeSize - volumeSizeRemaining
	return volumeSize, volumeUsedSize, nil
}

// GetDiskNumberFromVolumeID - gets the disk number where the volume is.
func (VolumeAPI) GetDiskNumberFromVolumeID(volumeID string) (uint32, error) {
	// get the size and sizeRemaining for the volume
	part, err := cim.GetPartitionByVolumeUniqueID(volumeID)
	if err != nil {
		return 0, err
	}

	diskNumber, err := cim.GetPartitionDiskNumber(part)
	if err != nil {
		return 0, fmt.Errorf("error query disk number of volume (%s). error: %v", volumeID, err)
	}

	return diskNumber, nil
}

// GetVolumeIDFromTargetPath - gets the volume ID given a mount point, the function is recursive until it find a volume or errors out
func (VolumeAPI) GetVolumeIDFromTargetPath(mount string) (string, error) {
	volumeString, err := getTarget(mount)
	if err != nil {
		return "", err
	}

	return volumeString, nil
}

func getTarget(mount string) (string, error) {
	mountedFolder, err := utils.IsMountedFolder(mount)
	if err != nil {
		return "", err
	}

	if !mountedFolder {
		return "", notMountedFolder
	}

	utf16FullMountPath, _ := windows.UTF16PtrFromString(mount)
	outPathBuffer := make([]uint16, windows.MAX_LONG_PATH)
	err = windows.GetVolumePathName(utf16FullMountPath, &outPathBuffer[0], uint32(len(outPathBuffer)))
	if err != nil {
		return "", err
	}
	targetPath := utils.EnsureLongPath(windows.UTF16PtrToString(&outPathBuffer[0]))
	if !strings.HasSuffix(targetPath, "\\") {
		targetPath += "\\"
	}
	utf16TargetPath, _ := windows.UTF16PtrFromString(targetPath)
	outPathBuffer = make([]uint16, windows.MAX_LONG_PATH)
	err = windows.GetVolumeNameForVolumeMountPoint(utf16TargetPath, &outPathBuffer[0], uint32(len(outPathBuffer)))
	if err != nil {
		return "", err
	}
	return windows.UTF16PtrToString(&outPathBuffer[0]), nil
}

// GetClosestVolumeIDFromTargetPath returns the volume id of a given target path.
func (VolumeAPI) GetClosestVolumeIDFromTargetPath(targetPath string) (string, error) {
	volumeString, err := findClosestVolume(targetPath)

	if err != nil {
		return "", fmt.Errorf("error getting the closest volume for the path=%s, err=%v", targetPath, err)
	}

	return volumeString, nil
}

// findClosestVolume finds the closest volume id for a given target path
// by following symlinks and moving up in the filesystem. if after moving up in the filesystem
// we get to a DriveLetter then the volume corresponding to this drive letter is returned instead.
func findClosestVolume(path string) (string, error) {
	candidatePath := path

	// Run in a bounded loop to avoid doing an infinite loop
	// while trying to follow symlinks
	//
	// The maximum path length in Windows is 260, it could be possible to end
	// up in a scenario where we do more than 256 iterations (e.g. by following symlinks from
	// a place high in the hierarchy to a nested sibling location many times)
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#:~:text=In%20editions%20of%20Windows%20before,required%20to%20remove%20the%20limit.
	//
	// The number of iterations is 256, which is similar to the number of iterations in filepath-securejoin
	// https://github.com/cyphar/filepath-securejoin/blob/64536a8a66ae59588c981e2199f1dcf410508e07/join.go#L51
	for i := 0; i < 256; i += 1 {
		isSymlink, err := utils.IsPathSymlink(candidatePath)
		if err != nil {
			return "", err
		}

		// mounted folder created by SetVolumeMountPoint may still report ModeSymlink == 0
		mountedFolder, err := utils.IsMountedFolder(candidatePath)
		if err != nil {
			return "", err
		}

		if isSymlink && mountedFolder {
			target, err := getTarget(candidatePath)
			if err != nil && !errors.Is(err, notMountedFolder) {
				return "", err
			}
			// if it has the form Volume{volumeid} then it's a volume
			if target != "" && VolumeRegexp.Match([]byte(target)) {
				// symlinks that are pointing to Volumes don't have this prefix
				return target, nil
			}
			// otherwise follow the symlink
			candidatePath = target
		} else {
			// if it's not a symlink move one level up
			previousPath := candidatePath
			candidatePath = filepath.Dir(candidatePath)

			// if the new path is the same as the previous path then we reached the root path
			if previousPath == candidatePath {
				// find the volume for the root path (assuming that it's a DriveLetter)
				target, err := getVolumeForDriveLetter(candidatePath[0:1])
				if err != nil {
					return "", err
				}
				return target, nil
			}
		}
	}

	return "", fmt.Errorf("failed to find the closest volume for path=%s", path)
}

// getVolumeForDriveLetter gets a volume from a drive letter (e.g. C:/).
func getVolumeForDriveLetter(path string) (string, error) {
	if len(path) != 1 {
		return "", fmt.Errorf("the path %s is not a valid drive letter", path)
	}

	volume, err := cim.GetVolumeByDriveLetter(path, cim.VolumeSelectorListUniqueID)
	if err != nil {
		return "", nil
	}

	uniqueID, err := cim.GetVolumeUniqueID(volume)
	if err != nil {
		return "", fmt.Errorf("error query unique ID of volume (%v). error: %v", volume, err)
	}

	return uniqueID, nil
}

func writeCache(volumeID string) error {
	volume, err := cim.QueryVolumeByUniqueID(volumeID, nil)
	if err != nil {
		return fmt.Errorf("error writing volume (%s) cache. error: %v", volumeID, err)
	}

	result, err := cim.FlushVolume(volume)
	if result != 0 || err != nil {
		return fmt.Errorf("error writing volume (%s) cache. result: %d, error: %v", volumeID, result, err)
	}
	return nil
}
