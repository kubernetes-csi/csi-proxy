package disk

import (
	"context"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/os/disk"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
	"k8s.io/klog/v2"
)

type Server struct {
	hostAPI disk.API
}

func NewServer(hostAPI disk.API) (*Server, error) {
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

	response.DiskLocations = make(map[uint32]*internal.DiskLocation)
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
	klog.V(4).Infof("calling PartitionDisk with diskNumber=%d", request.DiskNumber)
	response := &internal.PartitionDiskResponse{}
	diskNumber := request.DiskNumber

	initialized, err := s.hostAPI.IsDiskInitialized(diskNumber)
	if err != nil {
		klog.Errorf("failed check IsDiskInitialized %v", err)
		return response, err
	}
	if !initialized {
		klog.V(4).Infof("Initializing disk %d", diskNumber)
		err = s.hostAPI.InitializeDisk(diskNumber)
		if err != nil {
			klog.Errorf("failed InitializeDisk %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %d already initialized", diskNumber)
	}

	klog.V(4).Infof("Checking if disk %d is partitioned", diskNumber)
	partitioned, err := s.hostAPI.PartitionsExist(diskNumber)
	if err != nil {
		klog.Errorf("failed check PartitionsExist %v", err)
		return response, err
	}
	if !partitioned {
		klog.V(4).Infof("Creating partition on disk %d", diskNumber)
		err = s.hostAPI.CreatePartition(diskNumber)
		if err != nil {
			klog.Errorf("failed CreatePartition %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %d already partitioned", diskNumber)
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
	responseDiskIDs := make(map[uint32]*internal.DiskIDs)

	// Convert from shared to internal type
	for k, v := range diskIDs {
		diskIDs := internal.DiskIDs{Identifiers: make(map[string]string)}
		for k1, v1 := range v.Identifiers {
			diskIDs.Identifiers[k1] = v1
		}
		responseDiskIDs[k] = &diskIDs
	}
	response.DiskNumbers = responseDiskIDs
	return response, nil
}

func (s *Server) GetDiskStats(context context.Context, request *internal.GetDiskStatsRequest, version apiversion.Version) (*internal.GetDiskStatsResponse, error) {
	klog.V(4).Infof("calling GetDiskStats with diskNumber=%d", request.DiskNumber)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta1")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("GetDiskStats requires CSI-Proxy API version v1beta1 or greater")
	}
	diskNumber := request.DiskNumber

	diskSize, err := s.hostAPI.GetDiskStats(diskNumber)
	if err != nil {
		klog.Errorf("failed GetDiskStats %v", err)
		return nil, err
	}

	response := &internal.GetDiskStatsResponse{
		DiskSize: diskSize,
	}

	return response, nil
}

func (s *Server) SetDiskState(_ context.Context, request *internal.SetDiskStateRequest, version apiversion.Version) (*internal.SetDiskStateResponse, error) {
	klog.V(4).Infof("calling SetDiskState with diskNumber %q and isOnline %v", request.DiskNumber, request.IsOnline)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("SetDiskState requires CSI-Proxy API version v1beta2 or greater")
	}

	err := s.hostAPI.SetDiskState(request.DiskNumber, request.IsOnline)
	if err != nil {
		klog.Errorf("failed SetDiskState %v", err)
		return nil, err
	}

	response := &internal.SetDiskStateResponse{}

	return response, nil
}

func (s *Server) GetDiskState(_ context.Context, request *internal.GetDiskStateRequest, version apiversion.Version) (*internal.GetDiskStateResponse, error) {
	klog.V(4).Infof("calling GetDiskState with diskNumber %q", request.DiskNumber)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("GetDiskState requires CSI-Proxy API version v1beta2 or greater")
	}

	isOnline, err := s.hostAPI.GetDiskState(request.DiskNumber)
	if err != nil {
		klog.Errorf("failed GetDiskState %v", err)
		return nil, err
	}

	response := &internal.GetDiskStateResponse{IsOnline: isOnline}

	return response, nil
}
