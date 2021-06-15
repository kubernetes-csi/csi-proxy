package impl

// ComputeDoubleRequest is the internal representation of requests to the ComputeDouble endpoint.
type ComputeDoubleRequest struct {
	Input64 int64
}

// ComputeDoubleResponse is the internal representation of responses from the ComputeDouble endpoint.
type ComputeDoubleResponse struct {
	Response int64
}

// TellMeAPoemRequest is the internal representation of requests to the TellMeAPoem endpoint.
type TellMeAPoemRequest struct {
	IWantATitle bool
}

// TellMeAPoemResponse is the internal representation of responses from the TellMeAPoem endpoint.
type TellMeAPoemResponse struct {
	Title string
	Lines []string
}
