package volume

import (
	"context"
	"fmt"

	volumeapi "github.com/kubernetes-csi/csi-proxy/pkg/volume/api"
	"k8s.io/klog/v2"
)

// Volume wraps the host API and implements the interface
type Volume struct {
	hostAPI volumeapi.API
}

type Interface interface {
	DismountVolume(context.Context, *DismountVolumeRequest) (*DismountVolumeResponse, error)
	FormatVolume(context.Context, *FormatVolumeRequest) (*FormatVolumeResponse, error)
	GetClosestVolumeIDFromTargetPath(context.Context, *GetClosestVolumeIDFromTargetPathRequest) (*GetClosestVolumeIDFromTargetPathResponse, error)
	GetDiskNumberFromVolumeID(context.Context, *GetDiskNumberFromVolumeIDRequest) (*GetDiskNumberFromVolumeIDResponse, error)
	GetVolumeDiskNumber(context.Context, *VolumeDiskNumberRequest) (*VolumeDiskNumberResponse, error)
	GetVolumeIDFromMount(context.Context, *VolumeIDFromMountRequest) (*VolumeIDFromMountResponse, error)
	GetVolumeIDFromTargetPath(context.Context, *GetVolumeIDFromTargetPathRequest) (*GetVolumeIDFromTargetPathResponse, error)
	GetVolumeStats(context.Context, *GetVolumeStatsRequest) (*GetVolumeStatsResponse, error)
	IsVolumeFormatted(context.Context, *IsVolumeFormattedRequest) (*IsVolumeFormattedResponse, error)
	ListVolumesOnDisk(context.Context, *ListVolumesOnDiskRequest) (*ListVolumesOnDiskResponse, error)
	MountVolume(context.Context, *MountVolumeRequest) (*MountVolumeResponse, error)
	ResizeVolume(context.Context, *ResizeVolumeRequest) (*ResizeVolumeResponse, error)
	UnmountVolume(context.Context, *UnmountVolumeRequest) (*UnmountVolumeResponse, error)
	VolumeStats(context.Context, *VolumeStatsRequest) (*VolumeStatsResponse, error)
	WriteVolumeCache(context.Context, *WriteVolumeCacheRequest) (*WriteVolumeCacheResponse, error)
}

var _ Interface = &Volume{}

func New(hostAPI volumeapi.API) (*Volume, error) {
	return &Volume{
		hostAPI: hostAPI,
	}, nil
}

func (v *Volume) ListVolumesOnDisk(context context.Context, request *ListVolumesOnDiskRequest) (*ListVolumesOnDiskResponse, error) {
	klog.V(2).Infof("ListVolumesOnDisk: Request: %+v", request)
	response := &ListVolumesOnDiskResponse{}

	volumeIDs, err := v.hostAPI.ListVolumesOnDisk(request.DiskNumber, request.PartitionNumber)
	if err != nil {
		klog.Errorf("failed ListVolumeOnDisk %v", err)
		return response, err
	}

	response.VolumeIds = volumeIDs
	return response, nil
}

func (v *Volume) MountVolume(context context.Context, request *MountVolumeRequest) (*MountVolumeResponse, error) {
	klog.V(2).Infof("MountVolume: Request: %+v", request)
	response := &MountVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("MountVolumeRequest.VolumeId is empty")
	}
	targetPath := request.TargetPath
	if targetPath == "" {
		klog.Errorf("targetPath empty")
		return response, fmt.Errorf("MountVolumeRequest.TargetPath is empty")
	}

	err := v.hostAPI.MountVolume(volumeID, targetPath)
	if err != nil {
		klog.Errorf("failed MountVolume %v", err)
		return response, err
	}
	return response, nil
}

func (v *Volume) DismountVolume(context context.Context, request *DismountVolumeRequest) (*DismountVolumeResponse, error) {
	unmountVolumeRequest := &UnmountVolumeRequest{
		VolumeId:   request.VolumeId,
		TargetPath: request.Path,
	}
	_, err := v.UnmountVolume(context, unmountVolumeRequest)
	if err != nil {
		return nil, fmt.Errorf("Forward to UnmountVolume failed, err=%+v", err)
	}
	dismountVolumeResponse := &DismountVolumeResponse{}
	return dismountVolumeResponse, nil
}

func (v *Volume) UnmountVolume(context context.Context, request *UnmountVolumeRequest) (*UnmountVolumeResponse, error) {
	klog.V(2).Infof("UnmountVolume: Request: %+v", request)
	response := &UnmountVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("volume id empty")
	}
	targetPath := request.TargetPath
	if targetPath == "" {
		klog.Errorf("target path empty")
		return response, fmt.Errorf("target path empty")
	}
	err := v.hostAPI.UnmountVolume(volumeID, targetPath)
	if err != nil {
		klog.Errorf("failed UnmountVolume %v", err)
		return response, err
	}
	return response, nil
}

func (v *Volume) IsVolumeFormatted(context context.Context, request *IsVolumeFormattedRequest) (*IsVolumeFormattedResponse, error) {
	klog.V(2).Infof("IsVolumeFormatted: Request: %+v", request)
	response := &IsVolumeFormattedResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("volume id empty")
	}
	isFormatted, err := v.hostAPI.IsVolumeFormatted(volumeID)
	if err != nil {
		klog.Errorf("failed IsVolumeFormatted %v", err)
		return response, err
	}
	klog.V(5).Infof("IsVolumeFormatted: return: %v", isFormatted)
	response.Formatted = isFormatted
	return response, nil
}

func (v *Volume) FormatVolume(context context.Context, request *FormatVolumeRequest) (*FormatVolumeResponse, error) {
	klog.V(2).Infof("FormatVolume: Request: %+v", request)
	response := &FormatVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("volume id empty")
	}

	err := v.hostAPI.FormatVolume(volumeID)
	if err != nil {
		klog.Errorf("failed FormatVolume %v", err)
		return response, err
	}
	return response, nil
}

func (v *Volume) WriteVolumeCache(context context.Context, request *WriteVolumeCacheRequest) (*WriteVolumeCacheResponse, error) {
	klog.V(2).Infof("WriteVolumeCache: Request: %+v", request)
	response := &WriteVolumeCacheResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("volume id empty")
	}

	err := v.hostAPI.WriteVolumeCache(volumeID)
	if err != nil {
		klog.Errorf("failed WriteVolumeCache %v", err)
		return response, err
	}
	return response, nil
}

func (v *Volume) ResizeVolume(context context.Context, request *ResizeVolumeRequest) (*ResizeVolumeResponse, error) {
	klog.V(2).Infof("ResizeVolume: Request: %+v", request)
	response := &ResizeVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		klog.Errorf("volume id empty")
		return response, fmt.Errorf("volume id empty")
	}
	sizeBytes := request.SizeBytes
	// TODO : Validate size param

	err := v.hostAPI.ResizeVolume(volumeID, sizeBytes)
	if err != nil {
		klog.Errorf("failed ResizeVolume %v", err)
		return response, err
	}
	return response, nil
}

func (v *Volume) VolumeStats(context context.Context, request *VolumeStatsRequest) (*VolumeStatsResponse, error) {
	getVolumeStatsRequest := &GetVolumeStatsRequest{
		VolumeId: request.VolumeId,
	}
	getVolumeStatsResponse, err := v.GetVolumeStats(context, getVolumeStatsRequest)
	if err != nil {
		return nil, fmt.Errorf("Forward to GetVolumeStats failed, err=%+v", err)
	}
	volumeStatsResponse := &VolumeStatsResponse{
		VolumeSize:     getVolumeStatsResponse.TotalBytes,
		VolumeUsedSize: getVolumeStatsResponse.UsedBytes,
	}
	return volumeStatsResponse, nil
}

func (v *Volume) GetVolumeStats(context context.Context, request *GetVolumeStatsRequest) (*GetVolumeStatsResponse, error) {
	klog.V(2).Infof("GetVolumeStats: Request: %+v", request)
	volumeID := request.VolumeId
	if volumeID == "" {
		return nil, fmt.Errorf("volume id empty")
	}

	totalBytes, usedBytes, err := v.hostAPI.GetVolumeStats(volumeID)
	if err != nil {
		klog.Errorf("failed GetVolumeStats %v", err)
		return nil, err
	}

	klog.V(2).Infof("VolumeStats: returned: Capacity %v Used %v", totalBytes, usedBytes)

	response := &GetVolumeStatsResponse{
		TotalBytes: totalBytes,
		UsedBytes:  usedBytes,
	}

	return response, nil
}

func (v *Volume) GetVolumeDiskNumber(context context.Context, request *VolumeDiskNumberRequest) (*VolumeDiskNumberResponse, error) {
	getDiskNumberFromVolumeIDRequest := &GetDiskNumberFromVolumeIDRequest{
		VolumeId: request.VolumeId,
	}
	getDiskNumberFromVolumeIDResponse, err := v.GetDiskNumberFromVolumeID(context, getDiskNumberFromVolumeIDRequest)
	if err != nil {
		return nil, fmt.Errorf("Forward to GetDiskNumberFromVolumeID failed, err=%+v", err)
	}
	volumeStatsResponse := &VolumeDiskNumberResponse{
		DiskNumber: int64(getDiskNumberFromVolumeIDResponse.DiskNumber),
	}
	return volumeStatsResponse, nil
}

func (v *Volume) GetDiskNumberFromVolumeID(context context.Context, request *GetDiskNumberFromVolumeIDRequest) (*GetDiskNumberFromVolumeIDResponse, error) {
	klog.V(2).Infof("GetDiskNumberFromVolumeID: Request: %+v", request)

	volumeId := request.VolumeId
	if volumeId == "" {
		return nil, fmt.Errorf("volume id empty")
	}

	diskNumber, err := v.hostAPI.GetDiskNumberFromVolumeID(volumeId)
	if err != nil {
		klog.Errorf("failed GetDiskNumberFromVolumeID %v", err)
		return nil, err
	}

	response := &GetDiskNumberFromVolumeIDResponse{
		DiskNumber: diskNumber,
	}

	return response, nil
}

func (v *Volume) GetVolumeIDFromMount(context context.Context, request *VolumeIDFromMountRequest) (*VolumeIDFromMountResponse, error) {
	getVolumeIDFromTargetPathRequest := &GetVolumeIDFromTargetPathRequest{
		TargetPath: request.Mount,
	}
	getVolumeIDFromTargetPathResponse, err := v.GetVolumeIDFromTargetPath(context, getVolumeIDFromTargetPathRequest)
	if err != nil {
		return nil, fmt.Errorf("Forward to GetVolumeIDFromTargetPath failed, err=%+v", err)
	}
	volumeIDFromMountResponse := &VolumeIDFromMountResponse{
		VolumeId: getVolumeIDFromTargetPathResponse.VolumeId,
	}
	return volumeIDFromMountResponse, nil
}

func (v *Volume) GetVolumeIDFromTargetPath(context context.Context, request *GetVolumeIDFromTargetPathRequest) (*GetVolumeIDFromTargetPathResponse, error) {
	klog.V(2).Infof("GetVolumeIDFromTargetPath: Request: %+v", request)

	targetPath := request.TargetPath
	if targetPath == "" {
		return nil, fmt.Errorf("target path is empty")
	}

	volume, err := v.hostAPI.GetVolumeIDFromTargetPath(targetPath)
	if err != nil {
		klog.Errorf("failed GetVolumeIDFromTargetPath: %v", err)
		return nil, err
	}

	response := &GetVolumeIDFromTargetPathResponse{
		VolumeId: volume,
	}

	return response, nil
}

func (v *Volume) GetClosestVolumeIDFromTargetPath(context context.Context, request *GetClosestVolumeIDFromTargetPathRequest) (*GetClosestVolumeIDFromTargetPathResponse, error) {
	klog.V(2).Infof("GetClosestVolumeIDFromTargetPath: Request: %+v", request)

	targetPath := request.TargetPath
	if targetPath == "" {
		return nil, fmt.Errorf("target path is empty")
	}

	volume, err := v.hostAPI.GetClosestVolumeIDFromTargetPath(targetPath)
	if err != nil {
		klog.Errorf("failed GetClosestVolumeIDFromTargetPath: %v", err)
		return nil, err
	}

	response := &GetClosestVolumeIDFromTargetPathResponse{
		VolumeId: volume,
	}

	return response, nil
}
