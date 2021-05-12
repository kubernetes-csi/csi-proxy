package v1beta3

// Add manual conversion functions here to override automatic conversion functions

import (
	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta3"
	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

func Convert_internal_ListDiskLocationsResponse_To_v1beta3_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1beta3.ListDiskLocationsResponse) error {
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[string]*v1beta3.DiskLocation, len(*in))
		for key, val := range *in {
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

func Convert_internal_ListDiskIDsResponse_To_v1beta3_ListDiskIDsResponse(in *internal.ListDiskIDsResponse, out *v1beta3.ListDiskIDsResponse) error {
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[string]*v1beta3.DiskIDs, len(*in))
		for key, val := range *in {
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

func Convert_internal_DiskStatsResponse_To_v1beta3_DiskStatsResponse(in *internal.DiskStatsResponse, out *v1beta3.DiskStatsResponse) error {
	out.DiskSize = in.DiskSize
	return nil
}

func Convert_v1beta3_DiskStatsRequest_To_internal_DiskStatsRequest(in *v1beta3.DiskStatsRequest, out *internal.DiskStatsRequest) error {
	out.DiskID = in.DiskID
	return nil
}
