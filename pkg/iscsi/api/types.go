package api

// TargetPortal is an address and port pair for a specific iSCSI storage
// target.
// JSON field names are the WMI MSFT_iSCSITargetPortal field names.
type TargetPortal struct {
	Address string `json:"TargetPortalAddress"`
	Port    uint32 `json:"TargetPortalPortNumber"`
}
