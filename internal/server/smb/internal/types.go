package internal

type NewSmbGlobalMappingRequest struct {
	RemotePath string
	Username   string
	Password   string
}

type NewSmbGlobalMappingResponse struct {
	Error string
}

type RemoveSmbGlobalMappingRequest struct {
	RemotePath string
}

type RemoveSmbGlobalMappingResponse struct {
	Error string
}
