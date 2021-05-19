package v1alpha1

// Add manual conversion functions here to override automatic conversion functions

// import (
// 	"github.com/kubernetes-csi/csi-proxy/client/api/disk/v1alpha1"
// 	"github.com/kubernetes-csi/csi-proxy/internal/server/disk/internal"
// )

// func Convert_internal_ListDiskLocationsResponse_To_v1alpha1_ListDiskLocationsResponse(in *internal.ListDiskLocationsResponse, out *v1alpha1.ListDiskLocationsResponse) error {
// 	if in.DiskLocations != nil {
// 		in, out := &in.DiskLocations, &out.DiskLocations
// 		*out = make(map[string]*v1alpha1.DiskLocation, len(*in))
// 		for key, val := range *in {
// 			newVal := new(v1alpha1.DiskLocation)
// 			if err := Convert_internal_DiskLocation_To_v1alpha1_DiskLocation(val, newVal); err != nil {
// 				return err
// 			}
// 			(*out)[key] = newVal
// 		}
// 	} else {
// 		out.DiskLocations = nil
// 	}
// 	return nil
// }
