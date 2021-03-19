package v1

// Add manual conversion functions here to override automatic conversion functions

import (
	v1 "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1"
	internal "github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
)

func Convert_internal_ListDiskLocationsResponse_To_v1_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1.ListDiskLocationsResponse) error {
	if in.DiskLocations != nil {
		in, out := &in.DiskLocations, &out.DiskLocations
		*out = make(map[string]*v1.DiskLocation, len(*in))
		for key, val := range *in {
			newVal := new(v1.DiskLocation)
			if err := Convert_internal_DiskLocation_To_v1_DiskLocation(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal
		}
	} else {
		out.DiskLocations = nil
	}
	return nil
}

func Convert_internal_ListDiskIDsResponse_To_v1_ListDiskIDsResponse(in *internal.ListDiskIDsResponse, out *v1.ListDiskIDsResponse) error {
	if in.DiskIDs != nil {
		in, out := &in.DiskIDs, &out.DiskIDs
		*out = make(map[string]*v1.DiskIDs, len(*in))
		for key, val := range *in {
			newVal := new(v1.DiskIDs)
			if err := Convert_internal_DiskIDs_To_v1_DiskIDs(val, newVal); err != nil {
				return err
			}
			(*out)[key] = newVal
		}
	} else {
		out.DiskIDs = nil
	}
	return nil
}

func Convert_internal_DiskStatsResponse_To_v1_DiskStatsResponse(in *internal.DiskStatsResponse, out *v1.DiskStatsResponse) error {
	out.DiskSize = in.DiskSize
	return nil
}

func Convert_v1_DiskStatsRequest_To_internal_DiskStatsRequest(in *v1.DiskStatsRequest, out *internal.DiskStatsRequest) error {
	out.DiskID = in.DiskID
	return nil
}
