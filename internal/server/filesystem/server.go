package filesystem

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/filesystem/internal"
	"github.com/kubernetes-csi/csi-proxy/internal/utils"
)

type Server struct {
	kubeletCSIPluginsPath string
	kubeletPodPath        string
	hostAPI               API
}

var invalidPathCharsRegexWindows = regexp.MustCompile(`["/\:\?\*|]`)
var absPathRegexWindows = regexp.MustCompile(`^[a-zA-Z]:\\`)

type API interface {
	PathExists(path string) (bool, error)
	Mkdir(path string) error
	Rmdir(path string, force bool) error
	LinkPath(tgt string, src string) error
}

func NewServer(kubeletCSIPluginsPath string, kubeletPodPath string, hostAPI API) (*Server, error) {
	return &Server{
		kubeletCSIPluginsPath: kubeletCSIPluginsPath,
		kubeletPodPath:        kubeletPodPath,
		hostAPI:               hostAPI,
	}, nil
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
	if absPathRegexWindows.MatchString(path) {
		return true
	}
	return false
}

func (s *Server) validatePathWindows(pathCtx internal.PathContext, path string) error {
	prefix := ""
	if pathCtx == internal.PLUGIN {
		prefix = s.kubeletCSIPluginsPath
	} else if pathCtx == internal.POD {
		prefix = s.kubeletPodPath
	} else {
		return fmt.Errorf("Invalid PathContext: %v", pathCtx)
	}

	pathlen := len(path)

	if pathlen > utils.MAX_PATH_LENGTH_WINDOWS {
		return fmt.Errorf("Path length %d exceeds maximum characters: %d", pathlen, utils.MAX_PATH_LENGTH_WINDOWS)
	}

	if pathlen > 0 && (path[0] == '\\') {
		return fmt.Errorf("Invalid character \\ at begining of path: %s", path)
	}

	if isUNCPathWindows(path) {
		return fmt.Errorf("Unsupported UNC path prefix: %s", path)
	}

	if containsInvalidCharactersWindows(path) {
		return fmt.Errorf("Path contains invalid characters: %s", path)
	}

	if !isAbsWindows(path) {
		return fmt.Errorf("Not an absolute Windows path: %s", path)
	}

	if !strings.HasPrefix(path, prefix) {
		return fmt.Errorf("Path: %s is not within context path: %s", path, prefix)
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
	} else if pathCtx == internal.POD {
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
	exists, err := s.hostAPI.PathExists(request.Path)
	if err != nil {
		return &internal.PathExistsResponse{
			Error: err.Error(),
		}, nil
	}
	return &internal.PathExistsResponse{
		Error:  "",
		Exists: exists,
	}, nil
}

func (s *Server) Mkdir(ctx context.Context, request *internal.MkdirRequest, version apiversion.Version) (*internal.MkdirResponse, error) {
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, err
	}
	err = s.hostAPI.Mkdir(request.Path)
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
	err = s.hostAPI.Rmdir(request.Path, request.Force)
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
	err := s.validatePathWindows(internal.POD, request.SourcePath)
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
