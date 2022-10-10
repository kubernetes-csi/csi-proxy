package filesystem

import (
	"context"

	filesystemapi "github.com/kubernetes-csi/csi-proxy/pkg/filesystem/api"
	"k8s.io/klog/v2"
)

type Filesystem struct {
	hostAPI filesystemapi.API
}

type Interface interface {
	CreateSymlink(context.Context, *CreateSymlinkRequest) (*CreateSymlinkResponse, error)
	IsMountPoint(context.Context, *IsMountPointRequest) (*IsMountPointResponse, error)
	IsSymlink(context.Context, *IsSymlinkRequest) (*IsSymlinkResponse, error)
	LinkPath(context.Context, *LinkPathRequest) (*LinkPathResponse, error)
	Mkdir(context.Context, *MkdirRequest) (*MkdirResponse, error)
	PathExists(context.Context, *PathExistsRequest) (*PathExistsResponse, error)
	PathValid(context.Context, *PathValidRequest) (*PathValidResponse, error)
	Rmdir(context.Context, *RmdirRequest) (*RmdirResponse, error)
	RmdirContents(context.Context, *RmdirContentsRequest) (*RmdirContentsResponse, error)
}

// check that Filesystem implements Interface
var _ Interface = &Filesystem{}

func New(hostAPI filesystemapi.API) (*Filesystem, error) {
	return &Filesystem{
		hostAPI: hostAPI,
	}, nil
}

// PathExists checks if the given path exists on the host.
func (f *Filesystem) PathExists(ctx context.Context, request *PathExistsRequest) (*PathExistsResponse, error) {
	klog.V(2).Infof("Request: PathExists with path=%q", request.Path)
	err := ValidatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	exists, err := f.hostAPI.PathExists(request.Path)
	if err != nil {
		klog.Errorf("failed check PathExists %v", err)
		return nil, err
	}
	return &PathExistsResponse{
		Exists: exists,
	}, err
}

// PathValid checks if the given path is accessible.
func (f *Filesystem) PathValid(ctx context.Context, request *PathValidRequest) (*PathValidResponse, error) {
	klog.V(2).Infof("Request: PathValid with path %q", request.Path)
	valid, err := f.hostAPI.PathValid(request.Path)
	return &PathValidResponse{
		Valid: valid,
	}, err
}

func (f *Filesystem) Mkdir(ctx context.Context, request *MkdirRequest) (*MkdirResponse, error) {
	klog.V(2).Infof("Request: Mkdir with path=%q", request.Path)
	err := ValidatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	err = f.hostAPI.Mkdir(request.Path)
	if err != nil {
		klog.Errorf("failed Mkdir %v", err)
		return nil, err
	}

	return &MkdirResponse{}, err
}

func (f *Filesystem) Rmdir(ctx context.Context, request *RmdirRequest) (*RmdirResponse, error) {
	klog.V(2).Infof("Request: Rmdir with path=%q", request.Path)
	err := ValidatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	err = f.hostAPI.Rmdir(request.Path, request.Force)
	if err != nil {
		klog.Errorf("failed Rmdir %v", err)
		return nil, err
	}
	return nil, err
}

func (f *Filesystem) RmdirContents(ctx context.Context, request *RmdirContentsRequest) (*RmdirContentsResponse, error) {
	klog.V(2).Infof("Request: RmdirContents with path=%q", request.Path)
	err := ValidatePathWindows(request.Path)
	if err != nil {
		klog.Errorf("failed validatePathWindows %v", err)
		return nil, err
	}
	err = f.hostAPI.RmdirContents(request.Path)
	if err != nil {
		klog.Errorf("failed RmdirContents %v", err)
		return nil, err
	}
	return nil, err
}

func (f *Filesystem) LinkPath(ctx context.Context, request *LinkPathRequest) (*LinkPathResponse, error) {
	klog.V(2).Infof("Request: LinkPath with targetPath=%q sourcePath=%q", request.TargetPath, request.SourcePath)
	createSymlinkRequest := &CreateSymlinkRequest{
		SourcePath: request.SourcePath,
		TargetPath: request.TargetPath,
	}
	if _, err := f.CreateSymlink(ctx, createSymlinkRequest); err != nil {
		klog.Errorf("Failed to forward to CreateSymlink: %v", err)
		return nil, err
	}
	return &LinkPathResponse{}, nil
}

func (f *Filesystem) CreateSymlink(ctx context.Context, request *CreateSymlinkRequest) (*CreateSymlinkResponse, error) {
	klog.V(2).Infof("Request: CreateSymlink with targetPath=%q sourcePath=%q", request.TargetPath, request.SourcePath)
	err := ValidatePathWindows(request.TargetPath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for target path %v", err)
		return nil, err
	}
	err = ValidatePathWindows(request.SourcePath)
	if err != nil {
		klog.Errorf("failed validatePathWindows for source path %v", err)
		return nil, err
	}
	err = f.hostAPI.CreateSymlink(request.SourcePath, request.TargetPath)
	if err != nil {
		klog.Errorf("failed CreateSymlink: %v", err)
		return nil, err
	}
	return &CreateSymlinkResponse{}, nil
}

func (f *Filesystem) IsMountPoint(ctx context.Context, request *IsMountPointRequest) (*IsMountPointResponse, error) {
	klog.V(2).Infof("Request: IsMountPoint with path=%q", request.Path)
	isSymlinkRequest := &IsSymlinkRequest{
		Path: request.Path,
	}
	isSymlinkResponse, err := f.IsSymlink(ctx, isSymlinkRequest)
	if err != nil {
		klog.Errorf("Failed to forward to IsSymlink: %v", err)
		return nil, err
	}
	return &IsMountPointResponse{
		IsMountPoint: isSymlinkResponse.IsSymlink,
	}, nil
}

func (f *Filesystem) IsSymlink(ctx context.Context, request *IsSymlinkRequest) (*IsSymlinkResponse, error) {
	klog.V(2).Infof("Request: IsSymlink with path=%q", request.Path)
	isSymlink, err := f.hostAPI.IsSymlink(request.Path)
	if err != nil {
		klog.Errorf("failed IsSymlink %v", err)
		return nil, err
	}
	return &IsSymlinkResponse{
		IsSymlink: isSymlink,
	}, nil
}
