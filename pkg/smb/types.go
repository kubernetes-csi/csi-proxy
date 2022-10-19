package smb

type NewSMBGlobalMappingRequest struct {
	// A remote SMB share to mount
	// All unicode characters allowed in SMB server name specifications are
	// permitted except for restrictions below
	//
	// Restrictions:
	// SMB remote path specified in the format: \\server-name\sharename, \\server.fqdn\sharename or \\a.b.c.d\sharename
	// If not an IP address, share name has to be a valid DNS name.
	// UNC specifications to local paths or prefix: \\?\ is not allowed.
	// Characters: + [ ] " / : ; | < > , ? * = $ are not allowed
	RemotePath string

	// Optional local path to mount the SMB on
	LocalPath string

	// Username credential associated with the share
	Username string

	// Password credential associated with the share
	Password string
}

type NewSMBGlobalMappingResponse struct {
	// Intentionally empty
}

type RemoveSMBGlobalMappingRequest struct {
	// A remote SMB share mapping to remove
	// All unicode characters allowed in SMB server name specifications are
	// permitted except for restrictions below
	//
	// Restrictions:
	// SMB share specified in the format: \\server-name\sharename, \\server.fqdn\sharename or \\a.b.c.d\sharename
	// If not an IP address, share name has to be a valid DNS name.
	// UNC specifications to local paths or prefix: \\?\ is not allowed.
	// Characters: + [ ] " / : ; | < > , ? * = $ are not allowed
	RemotePath string
}

type RemoveSMBGlobalMappingResponse struct {
	// Intentionally empty
}
