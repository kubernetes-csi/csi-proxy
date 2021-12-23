package impl

// PathExistsRequest is the internal representation of requests to the PathExists endpoint.
type PathExistsRequest struct {
	// The path whose existence we want to check in the host's filesystem
	Path string
}

// PathExistsResponse is the internal representation of responses from the PathExists endpoint.
type PathExistsResponse struct {
	// Indicates whether the path in PathExistsRequest exists in the host's filesystem
	Exists bool
}

type MkdirRequest struct {
	// The path to create in the host's filesystem.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	// Non-existent parent directories in the path will be automatically created.
	// Directories will be created with Read and Write privileges of the Windows
	// User account under which csi-proxy is started (typically LocalSystem).
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// Depending on the context parameter of this function, the path prefix needs
	// to match the paths specified either as kubelet-csi-plugins-path
	// or as kubelet-pod-path parameters of csi-proxy.
	// The path parameter cannot already exist on host filesystem.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// Maximum path length will be capped to 260 characters.
	Path string
}

type MkdirResponse struct {
}

type RmdirRequest struct {
	// The path to remove in the host's filesystem.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// Depending on the context parameter of this function, the path prefix needs
	// to match the paths specified either as kubelet-csi-plugins-path
	// or as kubelet-pod-path parameters of csi-proxy.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// Path cannot be a file of type symlink.
	// Maximum path length will be capped to 260 characters.
	Path string
	// Force remove all contents under path (if any).
	Force bool
}

type RmdirResponse struct {
}

type LsdirRequest struct {
	// The directory to list files from.
	Path string
}

type LsdirResponse struct {
	// The files under the directory.
	Files []string
}

type CreateSymlinkRequest struct {
	// The path of the existing directory to be linked.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// The path prefix needs needs to match the paths specified as
	// kubelet-csi-plugins-path parameter of csi-proxy.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// source_path cannot already exist in the host filesystem.
	// Maximum path length will be capped to 260 characters.
	SourcePath string
	// Target path is the location of the new directory entry to be created in the host's filesystem.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// The path prefix needs to match the paths specified as
	// kubelet-pod-path parameter of csi-proxy.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// target_path needs to exist as a directory in the host that is empty.
	// target_path cannot be a symbolic link.
	// Maximum path length will be capped to 260 characters.
	TargetPath string
}

type CreateSymlinkResponse struct {
}

type IsSymlinkRequest struct {
	Path string
}

type IsSymlinkResponse struct {
	IsSymlink bool
}

// Compatibility for pre v1beta2 APIs

type LinkPathRequest struct {
	// The path where the symlink is created in the host's filesystem.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// source_path cannot already exist in the host filesystem.
	// Maximum path length will be capped to 260 characters.
	SourcePath string
	// Target path in the host's filesystem used for the symlink creation.
	// All special characters allowed by Windows in path names will be allowed
	// except for restrictions noted below. For details, please check:
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	//
	// Restrictions:
	// Only absolute path (indicated by a drive letter prefix: e.g. "C:\") is accepted.
	// UNC paths of the form "\\server\share\path\file" are not allowed.
	// All directory separators need to be backslash character: "\".
	// Characters: .. / : | ? * in the path are not allowed.
	// target_path needs to exist as a directory in the host that is empty.
	// target_path cannot be a symbolic link.
	// Maximum path length will be capped to 260 characters.
	TargetPath string
}

type LinkPathResponse struct {
}

type IsMountPointRequest struct {
	Path string
}

type IsMountPointResponse struct {
	IsMountPoint bool
}
