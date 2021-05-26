package disk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/os/disk"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
	"k8s.io/klog/v2"
)

type Server struct {
	hostAPI disk.API
}

// check that Server implements internal.ServerInterface
var _ internal.ServerInterface = &Server{}

func NewServer(hostAPI disk.API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) ListDiskLocations(context context.Context, request *internal.ListDiskLocationsRequest, version apiversion.Version) (*internal.ListDiskLocationsResponse, error) {
	klog.V(2).Infof("Request: ListDiskLocations: %+v", request)
	response := &internal.ListDiskLocationsResponse{}
	m, err := s.hostAPI.ListDiskLocations()
	if err != nil {
		klog.Errorf("ListDiskLocations failed: %v", err)
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
	klog.V(2).Infof("Request: PartitionDisk with diskNumber=%d", request.DiskNumber)
	response := &internal.PartitionDiskResponse{}
	diskNumber := request.DiskNumber

	initialized, err := s.hostAPI.IsDiskInitialized(diskNumber)
	if err != nil {
		klog.Errorf("IsDiskInitialized failed: %v", err)
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
	klog.V(2).Infof("Request: Rescan")
	response := &internal.RescanResponse{}
	err := s.hostAPI.Rescan()
	if err != nil {
		klog.Errorf("Rescan failed %v", err)
		return nil, err
	}
	return response, nil
}

func (s *Server) GetDiskNumberByName(context context.Context, request *internal.GetDiskNumberByNameRequest, version apiversion.Version) (*internal.GetDiskNumberByNameResponse, error) {
	klog.V(4).Infof("Request: GetDiskNumberByName with diskName %q", request.DiskName)
	response := &internal.GetDiskNumberByNameResponse{}
	diskName := request.DiskName
	number, err := s.hostAPI.GetDiskNumberByName(diskName)
	if err != nil {
		klog.Errorf("GetDiskNumberByName failed: %v", err)
		return nil, err
	}
	response.DiskNumber = number
	return response, nil
}

func (s *Server) ListDiskIDs(context context.Context, request *internal.ListDiskIDsRequest, version apiversion.Version) (*internal.ListDiskIDsResponse, error) {
	klog.V(4).Infof("Request: ListDiskIDs")
	minimumVersion := apiversion.NewVersionOrPanic("v1beta1")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("ListDiskIDs requires CSI-Proxy API version v1beta1 or greater")
	}

	diskIDs, err := s.hostAPI.ListDiskIDs()
	if err != nil {
		klog.Errorf("ListDiskIDs failed: %v", err)
		return nil, err
	}

	// Convert from shared to internal type
	responseDiskIDs := make(map[uint32]*internal.DiskIDs)
	for k, v := range diskIDs {
		responseDiskIDs[k] = &internal.DiskIDs{
			Page83:       v.Page83,
			SerialNumber: v.SerialNumber,
		}
	}
	return &internal.ListDiskIDsResponse{DiskIDs: responseDiskIDs}, nil
}

func (s *Server) DiskStats(context context.Context, request *internal.DiskStatsRequest, version apiversion.Version) (*internal.DiskStatsResponse, error) {
	klog.V(2).Infof("Request: DiskStats: diskNumber=%d", request.DiskID)
	minimumVersion := apiversion.NewVersionOrPanic("v1beta1")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("DiskStats requires CSI-Proxy API version v1beta1 or greater")
	}
	// forward to GetDiskStats
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format DiskStatsRequest.DiskID with err: %w", err)
	}
	getDiskStatsRequest := &internal.GetDiskStatsRequest{
		DiskNumber: uint32(diskNumber),
	}
	getDiskStatsResponse, err := s.GetDiskStats(context, getDiskStatsRequest, version)
	if err != nil {
		klog.Errorf("Forward to GetDiskStats failed: %+v", err)
		return nil, err
	}
	return &internal.DiskStatsResponse{
		DiskSize: getDiskStatsResponse.TotalBytes,
	}, nil
}

func (s *Server) GetDiskStats(context context.Context, request *internal.GetDiskStatsRequest, version apiversion.Version) (*internal.GetDiskStatsResponse, error) {
	klog.V(2).Infof("Request: GetDiskStats: diskNumber=%d", request.DiskNumber)
	diskNumber := request.DiskNumber
	totalBytes, err := s.hostAPI.GetDiskStats(diskNumber)
	if err != nil {
		klog.Errorf("GetDiskStats failed: %v", err)
		return nil, err
	}
	return &internal.GetDiskStatsResponse{
		TotalBytes: totalBytes,
	}, nil
}

func (s *Server) SetAttachState(context context.Context, request *internal.SetAttachStateRequest, version apiversion.Version) (*internal.SetAttachStateResponse, error) {
	klog.V(2).Infof("Request: SetAttachState: %+v", request)

	minimumVersion := apiversion.NewVersionOrPanic("v1beta2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("SetAttachState requires CSI-Proxy API version v1beta2 or greater")
	}

	// forward to SetDiskState
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format SetAttachStateRequest.DiskID with err: %w", err)
	}
	setDiskStateRequest := &internal.SetDiskStateRequest{
		DiskNumber: uint32(diskNumber),
		IsOnline:   request.IsOnline,
	}
	_, err = s.SetDiskState(context, setDiskStateRequest, version)
	if err != nil {
		klog.Errorf("Forward to SetDiskState failed with: %+v", err)
		return nil, err
	}
	return &internal.SetAttachStateResponse{}, nil
}

func (s *Server) SetDiskState(context context.Context, request *internal.SetDiskStateRequest, version apiversion.Version) (*internal.SetDiskStateResponse, error) {
	klog.V(2).Infof("Request: SetDiskState with diskNumber=%q and isOnline=%v", request.DiskNumber, request.IsOnline)
	err := s.hostAPI.SetDiskState(request.DiskNumber, request.IsOnline)
	if err != nil {
		klog.Errorf("SetDiskState failed: %v", err)
		return nil, err
	}
	return &internal.SetDiskStateResponse{}, nil
}

func (s *Server) GetAttachState(context context.Context, request *internal.GetAttachStateRequest, version apiversion.Version) (*internal.GetAttachStateResponse, error) {
	klog.V(2).Infof("Request: GetAttachState: %+v", request)

	minimumVersion := apiversion.NewVersionOrPanic("v1beta2")
	if version.Compare(minimumVersion) < 0 {
		return nil, fmt.Errorf("GetAttachState requires CSI-Proxy API version v1beta2 or greater")
	}

	// forward to GetDiskState
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format GetAttachStateRequest.DiskID with err: %w", err)
	}
	getDiskStateRequest := &internal.GetDiskStateRequest{
		DiskNumber: uint32(diskNumber),
	}
	getDiskStateResponse, err := s.GetDiskState(context, getDiskStateRequest, version)
	if err != nil {
		klog.Errorf("Forward to GetDiskState failed with: %+v", err)
		return nil, err
	}
	return &internal.GetAttachStateResponse{
		IsOnline: getDiskStateResponse.IsOnline,
	}, nil
}

func (s *Server) GetDiskState(context context.Context, request *internal.GetDiskStateRequest, version apiversion.Version) (*internal.GetDiskStateResponse, error) {
	klog.V(4).Infof("Request: GetDiskState with diskNumber=%d", request.DiskNumber)
	isOnline, err := s.hostAPI.GetDiskState(request.DiskNumber)
	if err != nil {
		klog.Errorf("GetDiskState failed with: %v", err)
		return nil, err
	}
	return &internal.GetDiskStateResponse{IsOnline: isOnline}, nil
}
