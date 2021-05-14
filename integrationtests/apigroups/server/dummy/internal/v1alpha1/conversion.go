package v1alpha1

import (
	"fmt"
	"math"

	pb "github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/api/dummy/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/server/dummy/internal"
)

func convert_pb_ComputeDoubleRequest_To_internal_ComputeDoubleRequest(in *pb.ComputeDoubleRequest, out *internal.ComputeDoubleRequest) error {
	out.Input64 = int64(in.Input32)
	return nil
}

func convert_internal_ComputeDoubleResponse_To_pb_ComputeDoubleResponse(in *internal.ComputeDoubleResponse, out *pb.ComputeDoubleResponse) error {
	i := in.Response
	if i > math.MaxInt32 || i < math.MinInt32 {
		return fmt.Errorf("int32 overflow for %d", i)
	}
	out.Response32 = int32(i)
	return nil
}
