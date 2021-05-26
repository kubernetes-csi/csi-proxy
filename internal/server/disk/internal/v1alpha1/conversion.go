package v1alpha1

import (
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1alpha1"
	internal "github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_internal_ListDiskLocationsResponse_To_v1alpha1_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1alpha1.ListDiskLocationsResponse) error {
	// conversion function added in v1beta3, the disk_locations map key was changed from string to uint32
	// v1alpha1 is expecting a string so convert the uint32 to string
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[string]*v1alpha1.DiskLocation, len(*in))
		for key, val := range *in {
			newVal := new(v1alpha1.DiskLocation)
			if err := Convert_internal_DiskLocation_To_v1alpha1_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[strconv.FormatUint(uint64(key), 10)] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}

func Convert_v1alpha1_ListDiskLocationsResponse_To_internal_ListDiskLocationsResponse(in *v1alpha1.ListDiskLocationsResponse, out *internal.ListDiskLocationsResponse) error {
	// there's no need to implement this function because it's never used (we only convert an internal response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_v1alpha1_PartitionDiskRequest_To_internal_PartitionDiskRequest(in *v1alpha1.PartitionDiskRequest, out *internal.PartitionDiskRequest) error {
	diskNumber, err := strconv.ParseUint(in.DiskID, 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to convert v1alpha1.PartitionDiskRequest.DiskID to uint32: %w", err)
	}
	out.DiskNumber = uint32(diskNumber)
	return nil
}

// no need to override PartitionDiskResponse because the request doesn't return anything

func Convert_v1alpha1_GetDiskNumberByNameResponse_To_internal_GetDiskNumberByNameResponse(in *v1alpha1.GetDiskNumberByNameResponse, out *internal.GetDiskNumberByNameResponse) error {
	// there's no need to implement this function because it's never used (we only convert an internal response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_internal_GetDiskNumberByNameResponse_To_v1alpha1_GetDiskNumberByNameResponse(in *internal.GetDiskNumberByNameResponse, out *v1alpha1.GetDiskNumberByNameResponse) error {
	out.DiskNumber = strconv.FormatUint(uint64(in.DiskNumber), 10)
	return nil
}
