package internal

type ComputeDoubleRequest struct {
	Input64 int64
}

type ComputeDoubleResponse struct {
	Response int64
}

type TellMeAPoemRequest struct {
	IWantATitle bool
}

type TellMeAPoemResponse struct {
	Title string
	Lines []string
}
