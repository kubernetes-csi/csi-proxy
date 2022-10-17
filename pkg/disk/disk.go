package disk

import (
	"context"
	"fmt"
	"strconv"

	diskapi "github.com/kubernetes-csi/csi-proxy/pkg/disk/api"
	"k8s.io/klog/v2"
)

type Disk struct {
	hostAPI diskapi.API
}

type Interface interface {
	DiskStats(context.Context, *DiskStatsRequest) (*DiskStatsResponse, error)
	GetAttachState(context.Context, *GetAttachStateRequest) (*GetAttachStateResponse, error)
	GetDiskNumberByName(context.Context, *GetDiskNumberByNameRequest) (*GetDiskNumberByNameResponse, error)
	// GetDiskState gets the offline/online state of a disk.
	GetDiskState(context.Context, *GetDiskStateRequest) (*GetDiskStateResponse, error)

	// GetDiskStats returns the stats of a disk (currently it returns the disk size).
	GetDiskStats(context.Context, *GetDiskStatsRequest) (*GetDiskStatsResponse, error)

	// ListDiskIDs returns a map of DiskID objects where the key is the disk number.
	ListDiskIDs(context.Context, *ListDiskIDsRequest) (*ListDiskIDsResponse, error)

	// ListDiskLocations returns locations <Adapter, Bus, Target, LUN ID> of all
	// disk devices enumerated by the host.
	ListDiskLocations(context.Context, *ListDiskLocationsRequest) (*ListDiskLocationsResponse, error)

	// PartitionDisk initializes and partitions a disk device with the GPT partition style
	// (if the disk has not been partitioned already) and returns the resulting volume device ID.
	PartitionDisk(context.Context, *PartitionDiskRequest) (*PartitionDiskResponse, error)

	// Rescan refreshes the host's storage cache.
	Rescan(context.Context, *RescanRequest) (*RescanResponse, error)
	SetAttachState(context.Context, *SetAttachStateRequest) (*SetAttachStateResponse, error)

	// SetDiskState sets the offline/online state of a disk.
	SetDiskState(context.Context, *SetDiskStateRequest) (*SetDiskStateResponse, error)
}

// check that Disk implements Interface
var _ Interface = &Disk{}

func New(hostAPI diskapi.API) (*Disk, error) {
	return &Disk{
		hostAPI: hostAPI,
	}, nil
}

func (d *Disk) ListDiskLocations(context context.Context, request *ListDiskLocationsRequest) (*ListDiskLocationsResponse, error) {
	klog.V(2).Infof("Request: ListDiskLocations: %+v", request)
	response := &ListDiskLocationsResponse{}
	m, err := d.hostAPI.ListDiskLocations()
	if err != nil {
		klog.Errorf("ListDiskLocations failed: %v", err)
		return response, err
	}

	response.DiskLocations = make(map[uint32]*DiskLocation)
	for k, v := range m {
		d := &DiskLocation{}
		d.Adapter = v.Adapter
		d.Bus = v.Bus
		d.Target = v.Target
		d.LUNID = v.LUNID
		response.DiskLocations[k] = d
	}
	return response, nil
}

func (d *Disk) PartitionDisk(context context.Context, request *PartitionDiskRequest) (*PartitionDiskResponse, error) {
	klog.V(2).Infof("Request: PartitionDisk with diskNumber=%d", request.DiskNumber)
	response := &PartitionDiskResponse{}
	diskNumber := request.DiskNumber

	initialized, err := d.hostAPI.IsDiskInitialized(diskNumber)
	if err != nil {
		klog.Errorf("IsDiskInitialized failed: %v", err)
		return response, err
	}
	if !initialized {
		klog.V(4).Infof("Initializing disk %d", diskNumber)
		err = d.hostAPI.InitializeDisk(diskNumber)
		if err != nil {
			klog.Errorf("failed InitializeDisk %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %d already initialized", diskNumber)
	}

	klog.V(4).Infof("Checking if disk %d has basic partitions", diskNumber)
	partitioned, err := d.hostAPI.BasicPartitionsExist(diskNumber)
	if err != nil {
		klog.Errorf("failed check BasicPartitionsExist %v", err)
		return response, err
	}
	if !partitioned {
		klog.V(4).Infof("Creating basic partition on disk %d", diskNumber)
		err = d.hostAPI.CreateBasicPartition(diskNumber)
		if err != nil {
			klog.Errorf("failed CreateBasicPartition %v", err)
			return response, err
		}
	} else {
		klog.V(4).Infof("Disk %d already partitioned", diskNumber)
	}
	return response, nil
}

func (d *Disk) Rescan(context context.Context, request *RescanRequest) (*RescanResponse, error) {
	klog.V(2).Infof("Request: Rescan")
	response := &RescanResponse{}
	err := d.hostAPI.Rescan()
	if err != nil {
		klog.Errorf("Rescan failed %v", err)
		return nil, err
	}
	return response, nil
}

func (d *Disk) GetDiskNumberByName(context context.Context, request *GetDiskNumberByNameRequest) (*GetDiskNumberByNameResponse, error) {
	klog.V(4).Infof("Request: GetDiskNumberByName with diskName %q", request.DiskName)
	response := &GetDiskNumberByNameResponse{}
	diskName := request.DiskName
	number, err := d.hostAPI.GetDiskNumberByName(diskName)
	if err != nil {
		klog.Errorf("GetDiskNumberByName failed: %v", err)
		return nil, err
	}
	response.DiskNumber = number
	return response, nil
}

func (d *Disk) ListDiskIDs(context context.Context, request *ListDiskIDsRequest) (*ListDiskIDsResponse, error) {
	klog.V(4).Infof("Request: ListDiskIDs")

	diskIDs, err := d.hostAPI.ListDiskIDs()
	if err != nil {
		klog.Errorf("ListDiskIDs failed: %v", err)
		return nil, err
	}

	// Convert from shared to internal type
	responseDiskIDs := make(map[uint32]*DiskIDs)
	for k, v := range diskIDs {
		responseDiskIDs[k] = &DiskIDs{
			Page83:       v.Page83,
			SerialNumber: v.SerialNumber,
		}
	}
	response := &ListDiskIDsResponse{DiskIDs: responseDiskIDs}
	klog.V(5).Infof("Response=%v", response)
	return response, nil
}

func (d *Disk) DiskStats(context context.Context, request *DiskStatsRequest) (*DiskStatsResponse, error) {
	klog.V(2).Infof("Request: DiskStats: diskNumber=%d", request.DiskID)
	// forward to GetDiskStats
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format DiskStatsRequest.DiskID with err: %w", err)
	}
	getDiskStatsRequest := &GetDiskStatsRequest{
		DiskNumber: uint32(diskNumber),
	}
	getDiskStatsResponse, err := d.GetDiskStats(context, getDiskStatsRequest)
	if err != nil {
		klog.Errorf("Forward to GetDiskStats failed: %+v", err)
		return nil, err
	}
	return &DiskStatsResponse{
		DiskSize: getDiskStatsResponse.TotalBytes,
	}, nil
}

func (d *Disk) GetDiskStats(context context.Context, request *GetDiskStatsRequest) (*GetDiskStatsResponse, error) {
	klog.V(2).Infof("Request: GetDiskStats: diskNumber=%d", request.DiskNumber)
	diskNumber := request.DiskNumber
	totalBytes, err := d.hostAPI.GetDiskStats(diskNumber)
	if err != nil {
		klog.Errorf("GetDiskStats failed: %v", err)
		return nil, err
	}
	return &GetDiskStatsResponse{
		TotalBytes: totalBytes,
	}, nil
}

func (d *Disk) SetAttachState(context context.Context, request *SetAttachStateRequest) (*SetAttachStateResponse, error) {
	klog.V(2).Infof("Request: SetAttachState: %+v", request)

	// forward to SetDiskState
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format SetAttachStateRequest.DiskID with err: %w", err)
	}
	setDiskStateRequest := &SetDiskStateRequest{
		DiskNumber: uint32(diskNumber),
		IsOnline:   request.IsOnline,
	}
	_, err = d.SetDiskState(context, setDiskStateRequest)
	if err != nil {
		klog.Errorf("Forward to SetDiskState failed with: %+v", err)
		return nil, err
	}
	return &SetAttachStateResponse{}, nil
}

func (d *Disk) SetDiskState(context context.Context, request *SetDiskStateRequest) (*SetDiskStateResponse, error) {
	klog.V(2).Infof("Request: SetDiskState with diskNumber=%d and isOnline=%v", request.DiskNumber, request.IsOnline)
	err := d.hostAPI.SetDiskState(request.DiskNumber, request.IsOnline)
	if err != nil {
		klog.Errorf("SetDiskState failed: %v", err)
		return nil, err
	}
	return &SetDiskStateResponse{}, nil
}

func (d *Disk) GetAttachState(context context.Context, request *GetAttachStateRequest) (*GetAttachStateResponse, error) {
	klog.V(2).Infof("Request: GetAttachState: %+v", request)

	// forward to GetDiskState
	diskNumber, err := strconv.ParseUint(request.DiskID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to format GetAttachStateRequest.DiskID with err: %w", err)
	}
	getDiskStateRequest := &GetDiskStateRequest{
		DiskNumber: uint32(diskNumber),
	}
	getDiskStateResponse, err := d.GetDiskState(context, getDiskStateRequest)
	if err != nil {
		klog.Errorf("Forward to GetDiskState failed with: %+v", err)
		return nil, err
	}
	return &GetAttachStateResponse{
		IsOnline: getDiskStateResponse.IsOnline,
	}, nil
}

func (d *Disk) GetDiskState(context context.Context, request *GetDiskStateRequest) (*GetDiskStateResponse, error) {
	klog.V(4).Infof("Request: GetDiskState with diskNumber=%d", request.DiskNumber)
	isOnline, err := d.hostAPI.GetDiskState(request.DiskNumber)
	if err != nil {
		klog.Errorf("GetDiskState failed with: %v", err)
		return nil, err
	}
	return &GetDiskStateResponse{IsOnline: isOnline}, nil
}
