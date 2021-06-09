package filesystem

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/os/filesystem"
	"github.com/kubernetes-csi/csi-proxy/internal/server/filesystem/internal"
	"github.com/kubernetes-csi/csi-proxy/internal/utils"
	"k8s.io/klog/v2"
)

type Server struct {
	kubeletPath string
	hostAPI     filesystem.API
}

// check that Server fulfills internal.ServerInterface
var _ internal.ServerInterface = &Server{}

var invalidPathCharsRegexWindows = regexp.MustCompile(`["/\:\?\*|]`)
var absPathRegexWindows = regexp.MustCompile(`^[a-zA-Z]:\\`)

func NewServer(kubeletPath string, hostAPI filesystem.API) (*Server, error) {
	return &Server{
		kubeletPath: kubeletPath,
		hostAPI:     hostAPI,
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
	return s.validatePathWindows(path)
}

func (s *Server) validatePathWindows(path string) error {
	prefix := s.kubeletPath

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
	klog.V(2).Infof("Request: PathExists with path=%q", request.Path)
	err := s.validatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	exists, err := s.hostAPI.PathExists(request.Path)
	if err != nil {
		klog.Errorf("failed check PathExists %v", err)
		return nil, err
	}
	return &internal.PathExistsResponse{
		Exists: exists,
	}, err
}

// PathValid checks if the given path is accessiable.
func (s *Server) PathValid(ctx context.Context, path string) (bool, error) {
	klog.V(2).Infof("Request: PathValid with path %q", path)
	return s.hostAPI.PathValid(path)
}

func (s *Server) Mkdir(ctx context.Context, request *internal.MkdirRequest, version apiversion.Version) (*internal.MkdirResponse, error) {
	klog.V(2).Infof("Request: Mkdir with path=%q", request.Path)
	err := s.validatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	err = s.hostAPI.Mkdir(request.Path)
	if err != nil {
		klog.Errorf("failed Mkdir %v", err)
		return nil, err
	}

	return &internal.MkdirResponse{}, err
}

func (s *Server) Rmdir(ctx context.Context, request *internal.RmdirRequest, version apiversion.Version) (*internal.RmdirResponse, error) {
	klog.V(2).Infof("Request: Rmdir with path=%q", request.Path)
	err := s.validatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	err = s.hostAPI.Rmdir(request.Path, request.Force)
	if err != nil {
		klog.Errorf("failed Rmdir %v", err)
		return nil, err
	}
	return nil, err
}
func (s *Server) LinkPath(ctx context.Context, request *internal.LinkPathRequest, version apiversion.Version) (*internal.LinkPathResponse, error) {
	klog.V(2).Infof("Request: LinkPath with targetPath=%q sourcePath=%q", request.TargetPath, request.SourcePath)
	createSymlinkRequest := &internal.CreateSymlinkRequest{
		SourcePath: request.SourcePath,
		TargetPath: request.TargetPath,
	}
	if _, err := s.CreateSymlink(ctx, createSymlinkRequest, version); err != nil {
		klog.Errorf("Failed to forward to CreateSymlink: %v", err)
		return nil, err
	}
	return &internal.LinkPathResponse{}, nil
}

func (s *Server) CreateSymlink(ctx context.Context, request *internal.CreateSymlinkRequest, version apiversion.Version) (*internal.CreateSymlinkResponse, error) {
	klog.V(2).Infof("Request: CreateSymlink with targetPath=%q sourcePath=%q", request.TargetPath, request.SourcePath)
	err := s.validatePathWindows(request.TargetPath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for target path %v", err)
		return nil, err
	}
	err = s.validatePathWindows(request.SourcePath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for source path %v", err)
		return nil, err
	}
	err = s.hostAPI.CreateSymlink(request.SourcePath, request.TargetPath)
	if err != nil {
		klog.Errorf("failed CreateSymlink: %v", err)
		return nil, err
	}
	return &internal.CreateSymlinkResponse{}, nil
}

func (s *Server) IsMountPoint(ctx context.Context, request *internal.IsMountPointRequest, version apiversion.Version) (*internal.IsMountPointResponse, error) {
	klog.V(2).Infof("Request: IsMountPoint with path=%q", request.Path)
	isSymlinkRequest := &internal.IsSymlinkRequest{
		Path: request.Path,
	}
	isSymlinkResponse, err := s.IsSymlink(ctx, isSymlinkRequest, version)
	if err != nil {
		klog.Errorf("Failed to forward to IsSymlink: %v", err)
		return nil, err
	}
	return &internal.IsMountPointResponse{
		IsMountPoint: isSymlinkResponse.IsSymlink,
	}, nil
}

func (s *Server) IsSymlink(ctx context.Context, request *internal.IsSymlinkRequest, version apiversion.Version) (*internal.IsSymlinkResponse, error) {
	klog.V(2).Infof("Request: IsSymlink with path=%q", request.Path)
	isSymlink, err := s.hostAPI.IsSymlink(request.Path)
	if err != nil {
		klog.Errorf("failed IsSymlink %v", err)
		return nil, err
	}
	return &internal.IsSymlinkResponse{
		IsSymlink: isSymlink,
	}, nil
}
