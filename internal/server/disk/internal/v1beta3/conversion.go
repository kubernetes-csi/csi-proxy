package v1beta3

import (
	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta3"
	internal "github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_internal_ListDiskIDsResponse_To_v1beta3_ListDiskIDsResponse(in *internal.ListDiskIDsResponse, out *v1beta3.ListDiskIDsResponse) error {
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[uint32]*v1beta3.DiskIDs, len(*in))
		for key, val := range *in {

			// This function is almost generated correctly, it has an issue in the type of arguments sent
			// e.g.  if err := Convert_internal_DiskIDs_To_v1beta3_DiskIDs(*&val, *newVal); err != nil {
			newVal := new(v1beta3.DiskIDs)
			if err := Convert_internal_DiskIDs_To_v1beta3_DiskIDs(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal

		}
	} else {
		out.DiskIDs = nil
	}
	return nil
}

func Convert_internal_ListDiskLocationsResponse_To_v1beta3_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1beta3.ListDiskLocationsResponse) error {
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[uint32]*v1beta3.DiskLocation, len(*in))
		for key, val := range *in {

			// This function is almost generated correctly, it has an issue in the type of arguments sent
			// e.g.  if err := Convert_internal_DiskLocation_To_v1beta3_DiskLocation(*&val, *newVal); err != nil {
			newVal := new(v1beta3.DiskLocation)
			if err := Convert_internal_DiskLocation_To_v1beta3_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}
