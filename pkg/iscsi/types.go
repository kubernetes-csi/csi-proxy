package iscsi

type AddTargetPortalRequest struct {
	// iSCSI Target Portal to register in the initiator
	TargetPortal *TargetPortal
}

type AddTargetPortalResponse struct {
	// Intentionally empty
}

type AuthenticationType uint32

const (
	// No authentication is used
	NONE = 0

	// One way CHAP authentication. The target authenticates the initiator.
	ONE_WAY_CHAP = 1

	// Mutual CHAP authentication. The target and initiator authenticate each
	// other.
	MUTUAL_CHAP = 2
)

type ConnectTargetRequest struct {
	// Target portal to which the initiator will connect.
	TargetPortal *TargetPortal

	// IQN of the iSCSI Target
	Iqn string

	// Connection authentication type, None by default
	//
	// One Way Chap uses the chap_username and chap_secret
	// fields mentioned below to authenticate the initiator.
	//
	// Mutual Chap uses both the user/secret mentioned below
	// and the Initiator Chap Secret to authenticate the target and initiator.
	AuthType AuthenticationType

	// CHAP Username used to authenticate the initiator
	ChapUsername string

	// CHAP password used to authenticate the initiator
	ChapSecret string
}

type ConnectTargetResponse struct {
	// Intentionally empty
}

type DisconnectTargetRequest struct {
	// Target portal from which initiator will disconnect
	TargetPortal *TargetPortal
	// IQN of the iSCSI Target
	Iqn string
}

type DisconnectTargetResponse struct {
	// Intentionally empty
}

type DiscoverTargetPortalRequest struct {
	// iSCSI Target Portal on which to initiate discovery
	TargetPortal *TargetPortal
}

type DiscoverTargetPortalResponse struct {
	// List of discovered IQN addresses
	// follows IQN format: iqn.yyyy-mm.naming-authority:unique-name
	Iqns []string
}

type GetTargetDisksRequest struct {
	// Target portal whose disks will be queried
	TargetPortal *TargetPortal
	// IQN of the iSCSI Target
	Iqn string
}

type GetTargetDisksResponse struct {
	// List composed of disk ids (numbers) that are associated with the
	// iSCSI target
	DiskIDs []string
}

type ListTargetPortalsRequest struct {
}

type ListTargetPortalsResponse struct {
	// A list of Target Portals currently registered in the initiator
	TargetPortals []*TargetPortal
}

type RemoveTargetPortalRequest struct {
	// iSCSI Target Portal
	TargetPortal *TargetPortal
}

type RemoveTargetPortalResponse struct {
	// Intentionally empty
}

type TargetPortal struct {
	// iSCSI Target (server) address
	TargetAddress string
	// iSCSI Target port (default iSCSI port is 3260)
	TargetPort uint32
}

type SetMutualChapSecretRequest struct {
	// the default CHAP secret that all initiators on this machine (node) use to
	// authenticate the target on mutual CHAP authentication.
	// Must be at least 12 byte long for non-Ipsec connections, at least one
	// byte long for Ipsec connections, and at most 16 bytes long.
	MutualChapSecret string
}

type SetMutualChapSecretResponse struct {
	// Intentionally empty
}
