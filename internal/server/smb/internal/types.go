package internal

type MountSmbShareRequest struct {
    RemotePath string
    LocalPath string
    ReadOnly bool
    Username string
    Password string
}

type MountSmbShareResponse struct {
    Error string
}

type UnmountSmbShareRequest struct {
    RemotePath string
    LocalPath string
}

type UnmountSmbShareResponse struct {
    Error string
}