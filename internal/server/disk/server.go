package disk

import (
	"context"
	"fmt"

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
	ListDiskIDs() (map[string]shared.DiskIDs, error)
	DiskStats(diskID string) (int64, error)
	SetAttachState(diskID string, isOnline bool) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) ListDiskLocations(context context.Context, request *internal.ListDiskLocationsRequest, version apiversion.Version) (*internal.ListDiskLocationsResponse, error) {
	klog.V(4).Infof("calling ListDiskLocations")
	response := &internal.ListDiskLocationsResponse{}
	m, err := s.hostAPI.ListDiskLocations()
	if err != nil {
		klog.Errorf("failed ListDiskLocations: %v", err)
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
	klog.V(4).Infof("calling PartitionDisk with diskID %q", request.DiskID)
	response := &internal.PartitionDiskResponse{}
	diskID := request.DiskID

	initialized, err := s.hostAPI.IsDiskInitialized(diskID)
	if err != nil {
		klog.Errorf("failed check IsDiskInitialized %v", err)
		return response, err
	}
	if !initialized {
		klog.V(4).Infof("Initializing disk %s", diskID)
		err = s.hostAPI.InitializeDisk(diskID)
		if err != nil {
			klog.Errorf("failed InitializeDisk %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %s already initialized", diskID)
	}

	klog.V(4).Infof("Checking if disk %s is partitioned", diskID)
	partitioned, err := s.hostAPI.PartitionsExist(diskID)
	if err != nil {
		klog.Errorf("failed check PartitionsExist %v", err)
		return response, err
	}
	if !partitioned {
		klog.V(4).Infof("Creating partition on disk %s", diskID)
		err = s.hostAPI.CreatePartition(diskID)
		if err != nil {
			klog.Errorf("failed CreatePartition %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %s already partitioned", diskID)
	}
	return response, nil
}

func (s *Server) Rescan(context context.Context, request *internal.RescanRequest, version apiversion.Version) (*internal.RescanResponse, error) {
	klog.V(4).Infof("calling PartitionDisk")
	response := &internal.RescanResponse{}
	err := s.hostAPI.Rescan()
	if err != nil {
		klog.Errorf("failed Rescan %v", err)
		return nil, err
	}
	return response, nil
}

func (s *Server) GetDiskNumberByName(context context.Context, request *internal.GetDiskNumberByNameRequest, version apiversion.Version) (*internal.GetDiskNumberByNameResponse, error) {
	klog.V(4).Infof("calling GetDiskNumberByName with diskName %q", request.DiskName)
	response := &internal.GetDiskNumberByNameResponse{}
	diskName := request.DiskName
	number, err := s.hostAPI.GetDiskNumberByName(diskName)
	if err != nil {
		klog.Errorf("failed GetDiskNumberByName %v", err)
		return nil, err
	}
	response.DiskNumber = number
	return response, nil
}

func (s *Server) ListDiskIDs(context context.Context, request *internal.ListDiskIDsRequest, version apiversion.Version) (*internal.ListDiskIDsResponse, error) {
	klog.V(4).Infof("calling ListDiskIDs")
	minimumVersion := apiversion.NewVersionOrPanic("v1beta1")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("ListDiskIDs requires CSI-Proxy API version v1beta1 or greater")
	}

	response := &internal.ListDiskIDsResponse{}
	diskIDs, err := s.hostAPI.ListDiskIDs()

	if err != nil {
		klog.Errorf("failed ListDiskIDs %v", err)
		return nil, err
	}
	responseDiskIDs := make(map[string]*internal.DiskIDs)

	// Convert from shared to internal type
	for k, v := range diskIDs {
		diskIDs := internal.DiskIDs{Identifiers: make(map[string]string)}
		for k1, v1 := range v.Identifiers {
			diskIDs.Identifiers[k1] = v1
		}
		responseDiskIDs[k] = &diskIDs
	}
	response.DiskIDs = responseDiskIDs
	return response, nil
}

func (s *Server) DiskStats(context context.Context, request *internal.DiskStatsRequest, version apiversion.Version) (*internal.DiskStatsResponse, error) {
	klog.V(4).Infof("calling GetDiskStats with diskID %q", request.DiskID)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta1")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("GetDiskStats requires CSI-Proxy API version v1beta1 or greater")
	}
	diskID := request.DiskID

	diskSize, err := s.hostAPI.DiskStats(diskID)
	if err != nil {
		klog.Errorf("failed GetDiskStats %v", err)
		return nil, err
	}

	response := &internal.DiskStatsResponse{
		DiskSize: diskSize,
	}

	return response, nil
}

func (s *Server) SetAttachState(_ context.Context, request *internal.SetAttachStateRequest, version apiversion.Version) (*internal.SetAttachStateResponse, error) {
	klog.V(4).Infof("calling SetDiskState with diskID %q and isOnline %v", request.DiskID, request.IsOnline)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("SetDiskState requires CSI-Proxy API version v1beta2 or greater")
	}

	err := s.hostAPI.SetAttachState(request.DiskID, request.IsOnline)
	if err != nil {
		klog.Errorf("failed SetDiskState %v", err)
		return nil, err
	}

	response := &internal.SetAttachStateResponse{}

	return response, nil
}
