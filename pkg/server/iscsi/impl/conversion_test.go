package impl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha2"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/iscsi/impl"
	v1alpha1_impl "github.com/kubernetes-csi/csi-proxy/pkg/server/iscsi/impl/v1alpha1"
	v1alpha2_impl "github.com/kubernetes-csi/csi-proxy/pkg/server/iscsi/impl/v1alpha2"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestListTargetPortals_Conversion_v1alpha1(t *testing.T) {
	testCases := []struct {
		in      *impl.ListTargetPortalsResponse
		wantOut *v1alpha1.ListTargetPortalsResponse
		wantErr bool
	}{
		{
			in: &impl.ListTargetPortalsResponse{
				TargetPortals: []*impl.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantOut: &v1alpha1.ListTargetPortalsResponse{
				TargetPortals: []*v1alpha1.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantErr: false,
		},
		{
			in: &impl.ListTargetPortalsResponse{
				TargetPortals: []*impl.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"},
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
			in:      &impl.ListTargetPortalsResponse{},
			wantOut: &v1alpha1.ListTargetPortalsResponse{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		got := v1alpha1.ListTargetPortalsResponse{}
		err := v1alpha1_impl.Convert_impl_ListTargetPortalsResponse_To_v1alpha1_ListTargetPortalsResponse(tc.in, &got)
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

func TestListTargetPortals_Conversion_v1alpha2(t *testing.T) {
	testCases := []struct {
		in      *impl.ListTargetPortalsResponse
		wantOut *v1alpha2.ListTargetPortalsResponse
		wantErr bool
	}{
		{
			in: &impl.ListTargetPortalsResponse{
				TargetPortals: []*impl.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantOut: &v1alpha2.ListTargetPortalsResponse{
				TargetPortals: []*v1alpha2.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"}},
			},
			wantErr: false,
		},
		{
			in: &impl.ListTargetPortalsResponse{
				TargetPortals: []*impl.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"},
					{TargetPort: 1000, TargetAddress: "test.iqn.2"},
				},
			},
			wantOut: &v1alpha2.ListTargetPortalsResponse{
				TargetPortals: []*v1alpha2.TargetPortal{{TargetPort: 3260, TargetAddress: "test.iqn"},
					{TargetPort: 1000, TargetAddress: "test.iqn.2"},
				},
			},
			wantErr: false,
		},
		{
			in:      &impl.ListTargetPortalsResponse{},
			wantOut: &v1alpha2.ListTargetPortalsResponse{},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		got := v1alpha2.ListTargetPortalsResponse{}
		err := v1alpha2_impl.Convert_impl_ListTargetPortalsResponse_To_v1alpha2_ListTargetPortalsResponse(tc.in, &got)
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
