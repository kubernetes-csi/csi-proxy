package v1

import (
	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1"
	impl "github.com/kubernetes-csi/csi-proxy/pkg/server/disk/impl"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_impl_ListDiskIDsResponse_To_v1_ListDiskIDsResponse(in *impl.ListDiskIDsResponse, out *v1.ListDiskIDsResponse) error {
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[uint32]*v1.DiskIDs, len(*in))
		for key, val := range *in {

			// This function is almost generated correctly, it has an issue in the type of arguments sent
			// e.g.  if err := Convert_impl_DiskIDs_To_v1_DiskIDs(*&val, *newVal); err != nil {
			newVal := new(v1.DiskIDs)
			if err := Convert_impl_DiskIDs_To_v1_DiskIDs(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal

		}
	} else {
		out.DiskIDs = nil
	}
	return nil
}

func Convert_impl_ListDiskLocationsResponse_To_v1_ListDiskLocationsResponse(in *impl.ListDiskLocationsResponse, out *v1.ListDiskLocationsResponse) error {
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[uint32]*v1.DiskLocation, len(*in))
		for key, val := range *in {

			// This function is almost generated correctly, it has an issue in the type of arguments sent
			// e.g.  if err := Convert_impl_DiskLocation_To_v1_DiskLocation(*&val, *newVal); err != nil {
			newVal := new(v1.DiskLocation)
			if err := Convert_impl_DiskLocation_To_v1_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}
