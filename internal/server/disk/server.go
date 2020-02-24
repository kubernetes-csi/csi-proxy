package disk

import (
	"context"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
	shared "github.com/kubernetes-csi/csi-proxy/internal/shared/disk"
	"k8s.io/klog"
)

type Server struct {
	hostAPI API
}

type API interface {
	ListDiskLocations() (map[string]shared.DiskLocation, error)
	IsDiskInitialized(diskID string) (bool, error)
	InitializeDisk(diskID string) error
	PartitionsExist(diskID string) (bool, error)
	CreatePartition(diskID string) error
	Rescan() error
	GetDiskNumberByName(diskName string) (string, error) 
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) GetDiskNumberByName(context context.Context, request *internal.GetDiskNumberByNameRequest, version apiversion.Version) (*internal.GetDiskNumberByNameResponse, error) {
	response := &internal.GetDiskNumberByNameResponse{}
	diskName := request.DiskName
	number, err := s.hostAPI.GetDiskNumberByName(diskName)
	if err != nil {
		return nil, err
	}

	response.DiskNumber = number
	return response, nil
}

func (s *Server) ListDiskLocations(context context.Context, request *internal.ListDiskLocationsRequest, version apiversion.Version) (*internal.ListDiskLocationsResponse, error) {
	response := &internal.ListDiskLocationsResponse{}
	m, err := s.hostAPI.ListDiskLocations()
	if err != nil {
		return response, err
	}

	response.DiskLocations = make(map[string]*internal.DiskLocation)
	for k, v := range m {
		d := &internal.DiskLocation{}
		d.Adapter = v.Adapter
		d.Bus = v.Bus
		d.Target = v.Target
		d.LUNID = v.LUNID
		response.DiskLocations[k] = d
	}

	return response, nil
}

func (s *Server) PartitionDisk(context context.Context, request *internal.PartitionDiskRequest, version apiversion.Version) (*internal.PartitionDiskResponse, error) {
	response := &internal.PartitionDiskResponse{}
	diskID := request.DiskID

	klog.Infof("Checking if disk %s is initialized", diskID)
	initialized, err := s.hostAPI.IsDiskInitialized(diskID)
	if err != nil {
		return response, err
	}
	if !initialized {
		klog.Infof("Initializing disk %s", diskID)
		err = s.hostAPI.InitializeDisk(diskID)
		if err != nil {
			return response, err
		}
	} else {
		klog.Infof("Disk %s already initialized", diskID)
	}

	klog.Infof("Checking if disk %s is partitioned", diskID)
	paritioned, err := s.hostAPI.PartitionsExist(diskID)
	if err != nil {
		return response, err
	}
	if !paritioned {
		klog.Infof("Creating partition on disk %s", diskID)
		err = s.hostAPI.CreatePartition(diskID)
		if err != nil {
			return response, err
		}
	} else {
		klog.Infof("Disk %s already partitioned", diskID)
	}

	return response, nil
}

func (s *Server) Rescan(context context.Context, request *internal.RescanRequest, version apiversion.Version) (*internal.RescanResponse, error) {
	response := &internal.RescanResponse{}
	err := s.hostAPI.Rescan()
	if err != nil {
		return nil, err
	}
	return response, nil
}