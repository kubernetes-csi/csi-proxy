package filesystem

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/filesystem/internal"
)

type Server struct {
	kubeletCSIPluginsPath string
	kubeletPodPath        string
	hostAPI               API
}

const MAX_PATH_LENGTH_WINDOWS = 260

var invalidPathCharsRegexWindows = regexp.MustCompile(`["/\:\?\*|]`)

type API interface {
	PathExists(path string) (bool, error)
	Mkdir(path string) error
	Rmdir(path string) error
	LinkPath(tgt string, src string) error
}

func NewServer(hostOS string, kubeletCSIPluginsPath string, kubeletPodPath string, hostAPI API) (*Server, error) {
	if hostOS != "windows" {
		return nil, fmt.Errorf("Unsupported OS for FileSystem API server: %s", hostOS)
	}
	return &Server{
		kubeletCSIPluginsPath: kubeletCSIPluginsPath,
		kubeletPodPath:        kubeletPodPath,
		hostAPI:               hostAPI,
	}, nil
}

func isDriveLetterWindows(c uint8) bool {
	return ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z')
}

func containsInvalidCharactersWindows(path string) bool {
	if isAbsWindows(path) {
		path = path[3:]
	}
	if invalidPathCharsRegexWindows.MatchString(path) {
		return true
	}
	if strings.Contains(path, `..`) {
		return true
	}
	return false
}

func isUNCPathWindows(path string) bool {
	// check for UNC/pipe prefixes like "\\"
	if len(path) < 2 {
		return false
	}
	if path[0] == '\\' && path[1] == '\\' {
		return true
	}
	return false
}

func isAbsWindows(path string) bool {
	// for Windows check for C:\\.. prefix only
	// UNC prefixes of the form \\ are not considered
	// absolute in the context of CSI proxy
	if len(path) < 3 {
		return false
	}
	c := path[0]
	if !isDriveLetterWindows(c) {
		return false
	}
	if path[1] != ':' {
		return false
	}
	if path[2] != '\\' && path[2] != '/' {
		return false
	}
	return true
}

func (s *Server) validatePathWindows(pathCtx internal.PathContext, path string) error {
	prefix := ""
	if pathCtx == internal.PLUGIN {
		prefix = s.kubeletCSIPluginsPath
	} else if pathCtx == internal.CONTAINER {
		prefix = s.kubeletPodPath
	} else {
		return fmt.Errorf("Invalid PathContext: %v", pathCtx)
	}

	pathlen := len(path)

	if pathlen > MAX_PATH_LENGTH_WINDOWS {
		return fmt.Errorf("Path length %d exceeds maximum characters: %d", pathlen, MAX_PATH_LENGTH_WINDOWS)
	}

	if len(path) > 0 && (path[0] == '\\') {
		return fmt.Errorf("Invalid character \\ at begining of path: %s", path)
	}

	if isUNCPathWindows(path) {
		return fmt.Errorf("Unsupported UNC path prefix: %s", path)
	}

	if containsInvalidCharactersWindows(path) {
		return fmt.Errorf("Path contains invalid characters: %s", path)
	}

	if isAbsWindows(path) && !strings.HasPrefix(path, prefix) {
		return fmt.Errorf("Absolute path: %s is not within context path: %s", path, prefix)
	}

	return nil
}

func (s *Server) abs(pathCtx internal.PathContext, path string) (string, error) {
	if isAbsWindows(path) {
		return path, nil
	}
	prefix := ""
	if pathCtx == internal.PLUGIN {
		prefix = s.kubeletCSIPluginsPath
	} else if pathCtx == internal.CONTAINER {
		prefix = s.kubeletPodPath
	} else {
		return "", fmt.Errorf("Invalid PathContext: %v", pathCtx)
	}
	return prefix + "\\" + path, nil
}

// PathExists checks if the given path exists on the host.
func (s *Server) PathExists(ctx context.Context, request *internal.PathExistsRequest, version apiversion.Version) (*internal.PathExistsResponse, error) {
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		return &internal.PathExistsResponse{
			Error: err.Error(),
		}, nil
	}

	path, err := s.abs(request.Context, request.Path)
	if err != nil {
		return &internal.PathExistsResponse{
			Error: err.Error(),
		}, nil
	}

	exists, err := s.hostAPI.PathExists(path)
	return &internal.PathExistsResponse{
		Error:  err.Error(),
		Exists: exists,
	}, nil
}

func (s *Server) Mkdir(ctx context.Context, request *internal.MkdirRequest, version apiversion.Version) (*internal.MkdirResponse, error) {
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, nil
	}

	path, err := s.abs(request.Context, request.Path)
	if err != nil {
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, nil
	}

	err = s.hostAPI.Mkdir(path)
	if err != nil {
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, nil
	}

	return &internal.MkdirResponse{
		Error: "",
	}, nil
}

func (s *Server) Rmdir(ctx context.Context, request *internal.RmdirRequest, version apiversion.Version) (*internal.RmdirResponse, error) {
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		return &internal.RmdirResponse{
			Error: err.Error(),
		}, nil
	}
	err = s.hostAPI.Rmdir(request.Path)
	if err != nil {
		return &internal.RmdirResponse{
			Error: err.Error(),
		}, nil
	}
	return &internal.RmdirResponse{
		Error: "",
	}, nil
}

func (s *Server) LinkPath(ctx context.Context, request *internal.LinkPathRequest, version apiversion.Version) (*internal.LinkPathResponse, error) {
	err := s.validatePathWindows(internal.CONTAINER, request.SourcePath)
	if err != nil {
		return &internal.LinkPathResponse{
			Error: err.Error(),
		}, nil
	}
	err = s.validatePathWindows(internal.PLUGIN, request.TargetPath)
	if err != nil {
		return &internal.LinkPathResponse{
			Error: err.Error(),
		}, nil
	}
	err = s.hostAPI.LinkPath(request.TargetPath, request.SourcePath)
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	return &internal.LinkPathResponse{
		Error: errString,
	}, nil
}
