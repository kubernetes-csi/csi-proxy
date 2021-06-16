package v1beta1

import (
	"fmt"
	"strconv"

	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta1"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/disk/impl"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_impl_ListDiskLocationsResponse_To_v1beta1_ListDiskLocationsResponse(in *impl.ListDiskLocationsResponse, out *v1beta1.ListDiskLocationsResponse) error {
	// conversion function added in v1beta3, the disk_locations map key was changed from string to uint32
	// v1beta1 is expecting a string so convert the uint32 to string
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[string]*v1beta1.DiskLocation, len(*in))
		for key, val := range *in {
			newVal := new(v1beta1.DiskLocation)
			if err := Convert_impl_DiskLocation_To_v1beta1_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[strconv.FormatUint(uint64(key), 10)] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}

func Convert_v1beta1_ListDiskLocationsResponse_To_impl_ListDiskLocationsResponse(in *v1beta1.ListDiskLocationsResponse, out *impl.ListDiskLocationsResponse) error {
	// there's no need to implement this function because it's never used (we only convert an impl response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_impl_ListDiskIDsResponse_To_v1beta1_ListDiskIDsResponse(in *impl.ListDiskIDsResponse, out *v1beta1.ListDiskIDsResponse) error {
	// conversion function added in v1beta3, diskIDs was renamed to disk_ids, also the disk_ids map key was converted from string to uint32
	// v1beta1 is expecting diskIDs to have a key of the type string instead of uint32
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[string]*v1beta1.DiskIDs, len(*in))
		for key, val := range *in {
			newVal := new(v1beta1.DiskIDs)

			// copy impl.DiskIDs struct to a map by known keys (page83 and serialNumber)
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

func Convert_v1beta1_ListDiskIDsResponse_To_impl_ListDiskIDsResponse(in *v1beta1.ListDiskIDsResponse, out *impl.ListDiskIDsResponse) error {
	// there's no need to implement this function because it's never used (we only convert an impl response to a client response)
	// however we need to override it so that the generator doesn't generate the wrong implementation
	// see https://kubernetes.slack.com/archives/CN5JCCW31/p1621979489011400
	return nil
}

func Convert_v1beta1_PartitionDiskRequest_To_impl_PartitionDiskRequest(in *v1beta1.PartitionDiskRequest, out *impl.PartitionDiskRequest) error {
	diskNumber, err := strconv.ParseUint(in.DiskID, 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to convert v1beta1.PartitionDiskRequest.DiskID to uint32: %w", err)
	}
	out.DiskNumber = uint32(diskNumber)
	return nil
}

// no need to override PartitionDiskResponse because the request doesn't return anything
