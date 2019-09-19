package dummy

import (
	"context"
	"fmt"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	internal "github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/server/dummy/internal"
)

type Server struct{}

func (s *Server) ComputeDouble(ctx context.Context, request *internal.ComputeDoubleRequest, version apiversion.Version) (*internal.ComputeDoubleResponse, error) {
	in := request.Input64
	out := 2 * in

	if sign(in) != sign(out) {
		// overflow
		return nil, fmt.Errorf("int64 overflow with input: %d", in)
	}

	return &internal.ComputeDoubleResponse{
		Response: out,
	}, nil
}

func sign(x int64) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func (s *Server) TellMeAPoem(ctx context.Context, request *internal.TellMeAPoemRequest, version apiversion.Version) (*internal.TellMeAPoemResponse, error) {
	lines := []string{
		`Not like the brazen giant of Greek fame,`,
		`With conquering limbs astride from land to land;`,
		`Here at our sea-washed, sunset gates shall stand`,
		`A mighty woman with a torch, whose flame`,
		`Is the imprisoned lightning, and her name`,
		`Mother of Exiles. From her beacon-hand`,
		`Glows world-wide welcome; her mild eyes command`,
		`The air-bridged harbor that twin cities frame.`,
		`"Keep, ancient lands, your storied pomp!" cries she`,
		`With silent lips. "Give me your tired, your poor,`,
		`Your huddled masses yearning to breathe free,`,
		`The wretched refuse of your teeming shore.`,
		`Send these, the homeless, tempest-tost to me,`,
		`I lift my lamp beside the golden door!"`,
	}

	response := &internal.TellMeAPoemResponse{
		Lines: lines,
	}

	if request.IWantATitle {
		response.Title = "The New Colossus"
	}

	return response, nil
}
