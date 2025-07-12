package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureLongPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "AlreadyHasPrefix",
			input:    `\\?\C:\Some\Path`,
			expected: `\\?\C:\Some\Path`,
		},
		{
			name:     "MissingPrefix",
			input:    `C:\Some\Path`,
			expected: `\\?\C:\Some\Path`,
		},
		{
			name:     "EmptyPath",
			input:    ``,
			expected: `\\?\`,
		},
		{
			name:     "UNCPathWithoutPrefix",
			input:    `\\Server\Share`,
			expected: `\\?\\\Server\Share`, // depends on how you want to treat UNC paths
		},
		{
			name:     "PrefixOnlyOnce",
			input:    `\\?\C:\Some\Path`,
			expected: `\\?\C:\Some\Path`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnsureLongPath(tt.input)
			if result != tt.expected {
				t.Errorf("EnsureLongPath(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPathValid(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "valid-file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpDir, err := os.MkdirTemp("", "valid-dir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	nonexistent := filepath.Join(os.TempDir(), "does-not-exist-"+filepath.Base(tmpFile.Name()))

	invalid := string([]byte{0x00}) // illegal null character

	tests := []struct {
		name             string
		path             string
		expectValid      bool
		expectErrMessage string
	}{
		{
			name:             "ValidFile",
			path:             tmpFile.Name(),
			expectValid:      true,
			expectErrMessage: "",
		},
		{
			name:             "ValidDirectory",
			path:             tmpDir,
			expectValid:      true,
			expectErrMessage: "",
		},
		{
			name:             "NonExistentPath",
			path:             nonexistent,
			expectValid:      false,
			expectErrMessage: "",
		},
		{
			name:             "InvalidPath",
			path:             invalid,
			expectValid:      false,
			expectErrMessage: "invalid path: invalid argument",
		},
		{
			name:             "Drive C",
			path:             "c:",
			expectValid:      true,
			expectErrMessage: "",
		},
		{
			name:             "InvalidRemotePath",
			path:             "invalid-remote-path",
			expectValid:      false,
			expectErrMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := IsPathValid(tt.path)
			if valid != tt.expectValid {
				t.Errorf("Expected valid = %v, got %v", tt.expectValid, valid)
			}
			if err == nil && tt.expectErrMessage != "" {
				t.Errorf("Expected error message = %s, got no error", tt.expectErrMessage)
			}
			if err != nil {
				if tt.expectErrMessage != "" && err.Error() != tt.expectErrMessage {
					t.Errorf("Expected error message = %s, got error = %s", tt.expectErrMessage, err.Error())
				} else if tt.expectErrMessage == "" {
					t.Errorf("Expected no error, got error = %s", err.Error())
				}
			}
		})
	}
}

func runPowershellCmd(t *testing.T, command string) (string, error) {
	cmd := exec.Command("powershell", "/c", fmt.Sprintf("& { $global:ProgressPreference = 'SilentlyContinue'; %s }", command))
	t.Logf("Executing command: %q", cmd.String())
	result, err := cmd.CombinedOutput()
	return string(result), err
}

func createMountedFolder(t *testing.T, vhdxPath, mountedPath string, initialSize int) {
	cmd := fmt.Sprintf("New-VHD -Path %s -SizeBytes %d", vhdxPath, initialSize)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s.", err, cmd, out)
	}
	cmd = fmt.Sprintf("Mount-VHD -Path %s", vhdxPath)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
	}
	cmd = fmt.Sprintf("Mount-VHD -Path %s", vhdxPath)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %q. Out: %s", err, cmd, out)
	}
	cmd = fmt.Sprintf("(Get-VHD -Path %s).DiskNumber", vhdxPath)
	diskNumUnparsed, err := runPowershellCmd(t, cmd)
	if err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}
	diskNumUnparsed = strings.TrimSpace(diskNumUnparsed)
	cmd = fmt.Sprintf("Initialize-Disk -Number %s -PartitionStyle GPT", diskNumUnparsed)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error initializing disk: %v. Command: %q. Out: %s", err, cmd, out)
	}
	// Create a new partition using all available space
	cmd = fmt.Sprintf("New-Partition -DiskNumber %s -UseMaximumSize", diskNumUnparsed)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error creating partition: %v. Command: %q. Out: %s", err, cmd, out)
	}
	// Format the partition with NTFS
	cmd = fmt.Sprintf("(Get-Disk -Number %s | Get-Partition | Get-Volume) | Format-Volume -FileSystem NTFS -Confirm:$false", diskNumUnparsed)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error formatting volume: %v. Command: %q. Out: %s", err, cmd, out)
	}
	cmd = fmt.Sprintf(`(Get-Disk -Number %s | Get-Partition ) | Add-PartitionAccessPath -AccessPath %s`, diskNumUnparsed, mountedPath)
	if _, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}
}

func unmountFolder(t *testing.T, vhdxPath, mountedPath string) {
	cmd := fmt.Sprintf("(Get-VHD -Path %s).DiskNumber", vhdxPath)
	diskNumUnparsed, err := runPowershellCmd(t, cmd)
	if err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}
	diskNumUnparsed = strings.TrimSpace(diskNumUnparsed)
	cmd = fmt.Sprintf(`Get-Disk -Number %s | Get-Partition | Remove-PartitionAccessPath -AccessPath %s`, diskNumUnparsed, mountedPath)
	if _, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error: %v. Command: %s", err, cmd)
	}
	cmd = fmt.Sprintf("Dismount-VHD -Path %s", vhdxPath)
	if out, err := runPowershellCmd(t, cmd); err != nil {
		t.Fatalf("Error unmounting VHD: %v. Command: %q. Out: %s", err, cmd, out)
	}
}

func TestIsMountedFolder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-dir")
	require.NoError(t, err, "Failed to create temporary directory.")

	tests := []struct {
		name           string
		path           string
		setup          func()
		cleanup        func()
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "Non-existent path",
			path:           filepath.Join(tempDir, "nonexistent"),
			expectedResult: false,
			expectedError:  errors.New("The system cannot find the file specified."),
		},
		{
			name: "Regular directory",
			path: filepath.Join(tempDir, "regular_dir"),
			setup: func() {
				err := os.MkdirAll(filepath.Join(tempDir, "regular_dir"), 0644)
				require.NoError(t, err, "Failed to create regular_dir directory.")
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Mounted folder",
			path: filepath.Join(tempDir, "mounted_folder"),
			setup: func() {
				err := os.MkdirAll(filepath.Join(tempDir, "mounted_folder"), 0644)
				require.NoError(t, err, "Failed to create regular_dir directory.")

				createMountedFolder(t, filepath.Join(tempDir, "test.vhdx"), filepath.Join(tempDir, "mounted_folder"), 1024*1024*1024)
			},
			cleanup: func() {
				unmountFolder(t, filepath.Join(tempDir, "test.vhdx"), filepath.Join(tempDir, "mounted_folder"))
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "Regular file",
			path: filepath.Join(tempDir, "regular_file"),
			setup: func() {
				err := os.WriteFile(filepath.Join(tempDir, "regular_file"), []byte("just_a_test"), 0644)
				require.NoError(t, err, "Failed to create regular_file.")
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "Regular symlink",
			path: filepath.Join(tempDir, "regular_symlink"),
			setup: func() {
				err := os.WriteFile(filepath.Join(tempDir, "regular_file"), []byte("just_a_test"), 0644)
				require.NoError(t, err, "Failed to create regular_file.")

				err = os.Symlink(filepath.Join(tempDir, "regular_file"), filepath.Join(tempDir, "regular_symlink"))
				require.NoError(t, err, "Failed to create regular_file.")
			},
			cleanup: func() {
				err := os.RemoveAll(filepath.Join(tempDir, "regular_file"))
				require.NoError(t, err, "Failed to delete regular_file.")

				err = os.RemoveAll(filepath.Join(tempDir, "regular_symlink"))
				require.NoError(t, err, "Failed to delete regular_symlink.")
			},
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			// Run test
			result, err := IsMountedFolder(tt.path)

			if tt.cleanup != nil {
				tt.cleanup()
			}

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}

	err = os.RemoveAll(tempDir)
	require.NoError(t, err, "Failed to remove directory.")
}

func TestIsPathSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a regular file
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Create a symlink to the file
	symlinkPath := filepath.Join(tmpDir, "file_symlink")
	if err := os.Symlink(filePath, symlinkPath); err != nil {
		t.Skipf("Symlinks not supported on this platform or permission denied: %v", err)
	}

	// Create a directory
	dirPath := filepath.Join(tmpDir, "dir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create a symlink to the directory
	dirSymlinkPath := filepath.Join(tmpDir, "dir_symlink")
	if err := os.Symlink(dirPath, dirSymlinkPath); err != nil {
		t.Skipf("Directory symlinks not supported or permission denied: %v", err)
	}

	// Non-existent path
	nonExistent := filepath.Join(tmpDir, "not_exists")

	tests := []struct {
		name             string
		path             string
		expectLink       bool
		expectErrMessage string
	}{
		{
			name:             "RegularFile",
			path:             filePath,
			expectLink:       false,
			expectErrMessage: "",
		},
		{
			name:             "FileSymlink",
			path:             symlinkPath,
			expectLink:       true,
			expectErrMessage: "",
		},
		{
			name:             "Directory",
			path:             dirPath,
			expectLink:       false,
			expectErrMessage: "",
		},
		{
			name:             "DirectorySymlink",
			path:             dirSymlinkPath,
			expectLink:       true,
			expectErrMessage: "",
		},
		{
			name:             "NonExistent",
			path:             nonExistent,
			expectLink:       false,
			expectErrMessage: "The system cannot find the file specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isLink, err := IsPathSymlink(tt.path)
			if isLink != tt.expectLink {
				t.Errorf("Expected isLink=%v, got %v", tt.expectLink, isLink)
			}
			if err == nil && tt.expectErrMessage != "" {
				t.Errorf("Expected error message = %s, got no error", tt.expectErrMessage)
			}
			if err != nil {
				if tt.expectErrMessage != "" && !strings.Contains(err.Error(), tt.expectErrMessage) {
					t.Errorf("Expected error message = %s, got error = %s", tt.expectErrMessage, err.Error())
				} else if tt.expectErrMessage == "" {
					t.Errorf("Expected no error, got error = %s", err.Error())
				}
			}
		})
	}
}

func TestCreateSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a target file
	targetFile := filepath.Join(tmpDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}
	fileLink := filepath.Join(tmpDir, "file_link.txt")

	// Create a target directory
	targetDir := filepath.Join(tmpDir, "target_dir")
	if err := os.Mkdir(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}
	dirLink := filepath.Join(tmpDir, "dir_link")

	tests := []struct {
		name     string
		link     string
		target   string
		isDir    bool
		wantErr  bool
		validate func(t *testing.T, linkPath string)
	}{
		{
			name:    "FileSymlink",
			link:    fileLink,
			target:  targetFile,
			isDir:   false,
			wantErr: false,
			validate: func(t *testing.T, linkPath string) {
				fi, err := os.Lstat(linkPath)
				if err != nil {
					t.Fatalf("Symlink not created: %v", err)
				}
				if fi.Mode()&os.ModeSymlink == 0 {
					t.Errorf("Expected symlink, got mode: %v", fi.Mode())
				}
			},
		},
		{
			name:    "DirSymlink",
			link:    dirLink,
			target:  targetDir,
			isDir:   true,
			wantErr: false,
			validate: func(t *testing.T, linkPath string) {
				fi, err := os.Lstat(linkPath)
				if err != nil {
					t.Fatalf("Symlink not created: %v", err)
				}
				if fi.Mode()&os.ModeSymlink == 0 {
					t.Errorf("Expected symlink, got mode: %v", fi.Mode())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateSymlink(tt.link, tt.target, tt.isDir)
			if err != nil && !tt.wantErr {
				t.Fatalf("CreateSymlink() error = %v, wantErr = %v", err, tt.wantErr)
			} else if err == nil && tt.wantErr {
				t.Fatalf("CreateSymlink() no error, wantErr = %v", tt.wantErr)
			}
			if tt.validate != nil {
				tt.validate(t, tt.link)
			}
		})
	}
}

func TestPathExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	nonexistentPath := filepath.Join(os.TempDir(), "definitely_does_not_exist.test")
	_ = os.Remove(nonexistentPath) // ensure it doesn't exist

	invalidPath := string([]byte{0x00}) // causes an error on most systems

	tests := []struct {
		name             string
		path             string
		expectExist      bool
		expectErrMessage string
	}{
		{
			name:             "ExistingFile",
			path:             tmpFile.Name(),
			expectExist:      true,
			expectErrMessage: "",
		},
		{
			name:             "ExistingDirectory",
			path:             tmpDir,
			expectExist:      true,
			expectErrMessage: "",
		},
		{
			name:             "NonExistentPath",
			path:             nonexistentPath,
			expectExist:      false,
			expectErrMessage: "",
		},
		{
			name:             "InvalidPath",
			path:             invalidPath,
			expectExist:      false,
			expectErrMessage: "invalid argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := PathExists(tt.path)
			if exists != tt.expectExist {
				t.Errorf("Expected exists = %v, got %v", tt.expectExist, exists)
			}
			if err == nil && tt.expectErrMessage != "" {
				t.Errorf("Expected error message = %s, got no error", tt.expectErrMessage)
			}
			if err != nil {
				if tt.expectErrMessage != "" && !strings.Contains(err.Error(), tt.expectErrMessage) {
					t.Errorf("Expected error message = %s, got error = %s", tt.expectErrMessage, err.Error())
				} else if tt.expectErrMessage == "" {
					t.Errorf("Expected no error, got error = %s", err.Error())
				}
			}
		})
	}
}
