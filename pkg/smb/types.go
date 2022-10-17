package smb

type NewSMBGlobalMappingRequest struct {
	RemotePath string
	LocalPath  string
	Username   string
	Password   string
}

type NewSMBGlobalMappingResponse struct {
	// Intentionally empty.
}

type RemoveSMBGlobalMappingRequest struct {
	RemotePath string
}

type RemoveSMBGlobalMappingResponse struct {
	// Intentionally empty.
}
