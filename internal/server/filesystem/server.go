package filesystem

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/filesystem/internal"
	"github.com/kubernetes-csi/csi-proxy/internal/utils"
	"k8s.io/klog"
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
	PathValid(path string) (bool, error)
	Mkdir(path string) error
	Rmdir(path string, force bool) error
	LinkPath(tgt string, src string) error
	IsMountPoint(path string) (bool, error)
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
	return absPathRegexWindows.MatchString(path)
}

// ValidatePluginPath - Validates the path is compatible with the 'plugin path'
// restrictions.
// Note: The reason why we cannot reuse the validatePathWindows directly
// from other parts of the library is that it seems internal.PLUGIN was not
// usable from outside the internal path tree.
func (s *Server) ValidatePluginPath(path string) error {
	return s.validatePathWindows(internal.PLUGIN, path)
}

func (s *Server) validatePathWindows(pathCtx internal.PathContext, path string) error {
	prefix := ""
	if pathCtx == internal.PLUGIN {
		prefix = s.kubeletCSIPluginsPath
	} else if pathCtx == internal.POD {
		prefix = s.kubeletPodPath
	} else {
		return fmt.Errorf("invalid PathContext: %v", pathCtx)
	}

	pathlen := len(path)

	if pathlen > utils.MaxPathLengthWindows {
		return fmt.Errorf("path length %d exceeds maximum characters: %d", pathlen, utils.MaxPathLengthWindows)
	}

	if pathlen > 0 && (path[0] == '\\') {
		return fmt.Errorf("invalid character \\ at beginning of path: %s", path)
	}

	if isUNCPathWindows(path) {
		return fmt.Errorf("unsupported UNC path prefix: %s", path)
	}

	if containsInvalidCharactersWindows(path) {
		return fmt.Errorf("path contains invalid characters: %s", path)
	}

	if !isAbsWindows(path) {
		return fmt.Errorf("not an absolute Windows path: %s", path)
	}

	if !strings.HasPrefix(strings.ToLower(path), strings.ToLower(prefix)) {
		return fmt.Errorf("path: %s is not within context path: %s", path, prefix)
	}

	return nil
}

// PathExists checks if the given path exists on the host.
func (s *Server) PathExists(ctx context.Context, request *internal.PathExistsRequest, version apiversion.Version) (*internal.PathExistsResponse, error) {
	klog.V(4).Infof("calling PathExists with path %q", request.Path)
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return &internal.PathExistsResponse{
			Error: err.Error(),
		}, err
	}
	exists, err := s.hostAPI.PathExists(request.Path)
	if err != nil {
		klog.Errorf("failed check PathExists %v", err)
		return &internal.PathExistsResponse{
			Error: err.Error(),
		}, err
	}
	return &internal.PathExistsResponse{
		Error:  "",
		Exists: exists,
	}, err
}

// PathValid checks if the given path is accessiable.
func (s *Server) PathValid(ctx context.Context, path string) (bool, error) {
	klog.V(4).Infof("calling PathValid with path %q", path)
	return s.hostAPI.PathValid(path)
}

func (s *Server) Mkdir(ctx context.Context, request *internal.MkdirRequest, version apiversion.Version) (*internal.MkdirResponse, error) {
	klog.V(4).Infof("calling Mkdir with path %q", request.Path)
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, err
	}
	err = s.hostAPI.Mkdir(request.Path)
	if err != nil {
		klog.Errorf("failed Mkdir %v", err)
		return &internal.MkdirResponse{
			Error: err.Error(),
		}, err
	}

	return &internal.MkdirResponse{
		Error: "",
	}, err
}

func (s *Server) Rmdir(ctx context.Context, request *internal.RmdirRequest, version apiversion.Version) (*internal.RmdirResponse, error) {
	klog.V(2).Infof("calling Rmdir with path %q", request.Path)
	err := s.validatePathWindows(request.Context, request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return &internal.RmdirResponse{
			Error: err.Error(),
		}, err
	}
	err = s.hostAPI.Rmdir(request.Path, request.Force)
	if err != nil {
		klog.Errorf("failed Rmdir %v", err)
		return &internal.RmdirResponse{
			Error: err.Error(),
		}, err
	}
	return &internal.RmdirResponse{
		Error: "",
	}, err
}

func (s *Server) LinkPath(ctx context.Context, request *internal.LinkPathRequest, version apiversion.Version) (*internal.LinkPathResponse, error) {
	klog.V(4).Infof("calling LinkPath with targetPath %q sourcePath %q", request.TargetPath, request.SourcePath)
	err := s.validatePathWindows(internal.POD, request.TargetPath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for target path %v", err)
		return &internal.LinkPathResponse{
			Error: err.Error(),
		}, err
	}
	err = s.validatePathWindows(internal.PLUGIN, request.SourcePath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for source path %v", err)
		return &internal.LinkPathResponse{
			Error: err.Error(),
		}, err
	}
	err = s.hostAPI.LinkPath(request.SourcePath, request.TargetPath)
	errString := ""
	if err != nil {
		klog.Errorf("failed LinkPath %v", err)
		errString = err.Error()
	}
	return &internal.LinkPathResponse{
		Error: errString,
	}, err
}

func (s *Server) IsMountPoint(ctx context.Context, request *internal.IsMountPointRequest, version apiversion.Version) (*internal.IsMountPointResponse, error) {
	klog.V(4).Infof("calling IsMountPoint with path %q", request.Path)
	isMount, err := s.hostAPI.IsMountPoint(request.Path)
	if err != nil {
		klog.Errorf("failed IsMountPoint %v", err)
		return &internal.IsMountPointResponse{
			IsMountPoint: false,
			Error:        err.Error(),
		}, err
	}
	return &internal.IsMountPointResponse{
		Error:        "",
		IsMountPoint: isMount,
	}, err
}
