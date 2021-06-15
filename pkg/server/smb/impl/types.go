package impl

type NewSmbGlobalMappingRequest struct {
	RemotePath string
	LocalPath  string
	Username   string
	Password   string
}

type NewSmbGlobalMappingResponse struct {
	// Intentionally empty.
}

type RemoveSmbGlobalMappingRequest struct {
	RemotePath string
}

type RemoveSmbGlobalMappingResponse struct {
	// Intentionally empty.
}
