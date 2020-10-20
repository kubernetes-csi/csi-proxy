package internal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/internal/server/iscsi/internal"
	v1alpha1_internal "github.com/kubernetes-csi/csi-proxy/internal/server/iscsi/internal/v1alpha1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestListTargetPortals_Conversion(t *testing.T) {
	testCases := []struct {
		in      *internal.ListTargetPortalsResponse
		wantOut *v1alpha1.ListTargetPortalsResponse
		wantErr bool
	}{
		{
			in: &internal.ListTargetPortalsResponse{
				TargetPortals: []*internal.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantOut: &v1alpha1.ListTargetPortalsResponse{
				TargetPortals: []*v1alpha1.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantErr: false,
		},
		{
			in: &internal.ListTargetPortalsResponse{
				TargetPortals: []*internal.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"},
					{TargetPort: 1000, TargetAddress: "test.iqn.2"},
				},
			},
			wantOut: &v1alpha1.ListTargetPortalsResponse{
				TargetPortals: []*v1alpha1.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"},
					{TargetPort: 1000, TargetAddress: "test.iqn.2"},
				},
			},
			wantErr: false,
		},
		{
			in:      &internal.ListTargetPortalsResponse{},
			wantOut: &v1alpha1.ListTargetPortalsResponse{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		got := v1alpha1.ListTargetPortalsResponse{}
		err := v1alpha1_internal.Convert_internal_ListTargetPortalsResponse_To_v1alpha1_ListTargetPortalsResponse(tc.in, &got)
		if tc.wantErr && err == nil {
			t.Errorf("Expected error but returned a nil error")
		}
		if !tc.wantErr && err != nil {
			t.Errorf("Expected no errors but returned error: %s", err)
		}
		if diff := cmp.Diff(tc.wantOut, got, protocmp.Transform()); diff != "" {
			t.Errorf("Returned unexpected difference between conversion (-want +got):\n%s", diff)
		}
	}
}
