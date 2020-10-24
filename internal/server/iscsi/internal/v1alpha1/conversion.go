package v1alpha1

import (
	"github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/internal/server/iscsi/internal"
)

// Add manual conversion functions here to override automatic conversion functions

func Convert_internal_ListTargetPortalsResponse_To_v1alpha1_ListTargetPortalsResponse(in *internal.ListTargetPortalsResponse, out *v1alpha1.ListTargetPortalsResponse) error {
	if in.TargetPortals != nil {
		in, out := &in.TargetPortals, &out.TargetPortals
		*out = make([]*v1alpha1.TargetPortal, len(*in))
		for i := range *in {
			(*out)[i] = new(v1alpha1.TargetPortal)
			if err := Convert_internal_TargetPortal_To_v1alpha1_TargetPortal(*&(*in)[i], *&(*out)[i]); err != nil {
				return err
			}
		}
	} else {
		out.TargetPortals = nil
	}
	return nil
}
