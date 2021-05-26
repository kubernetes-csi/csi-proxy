package v1beta2

import (
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta2"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_internal_ListDiskLocationsResponse_To_v1beta2_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1beta2.ListDiskLocationsResponse) error {
	// conversion function added in v1beta3, the disk_locations map key was changed from string to uint32
	// v1beta2 is expecting a string so convert the uint32 to string
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[string]*v1beta2.DiskLocation, len(*in))
		for key, val := range *in {
			newVal := new(v1beta2.DiskLocation)
			if err := Convert_internal_DiskLocation_To_v1beta2_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[strconv.FormatUint(uint64(key), 10)] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}

func Convert_v1beta2_ListDiskLocationsResponse_To_internal_ListDiskLocationsResponse(in *v1beta2.ListDiskLocationsResponse, out *internal.ListDiskLocationsResponse) error {
	// there's no need to implement this function because it's never used (we only convert an internal response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_internal_ListDiskIDsResponse_To_v1beta2_ListDiskIDsResponse(in *internal.ListDiskIDsResponse, out *v1beta2.ListDiskIDsResponse) error {
	// conversion function added in v1beta3, diskIDs was renamed to disk_ids, also the disk_ids map key was converted from string to uint32
	// v1beta2 is expecting diskIDs to have a key of the type string instead of uint32
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[string]*v1beta2.DiskIDs, len(*in))
		for key, val := range *in {
			newVal := new(v1beta2.DiskIDs)

			// copy internal.DiskIDs struct to a map by known keys (page83 and serialNumber)
			newVal.Identifiers = make(map[string]string)
			newVal.Identifiers["page83"] = val.Page83
			newVal.Identifiers["serialNumber"] = val.SerialNumber
			(*out)[strconv.FormatUint(uint64(key), 10)] = newVal
		}
	} else {
		out.DiskIDs = nil
	}
	return nil
}

func Convert_v1beta2_ListDiskIDsResponse_To_internal_ListDiskIDsResponse(in *v1beta2.ListDiskIDsResponse, out *internal.ListDiskIDsResponse) error {
	// there's no need to implement this function because it's never used (we only convert an internal response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_v1beta2_PartitionDiskRequest_To_internal_PartitionDiskRequest(in *v1beta2.PartitionDiskRequest, out *internal.PartitionDiskRequest) error {
	diskNumber, err := strconv.ParseUint(in.DiskID, 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to convert v1beta2.PartitionDiskRequest.DiskID to uint32: %w", err)
	}
	out.DiskNumber = uint32(diskNumber)
	return nil
}

// no need to override PartitionDiskResponse because the request doesn't return anything
