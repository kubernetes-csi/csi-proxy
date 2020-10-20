package integrationtests

// Needs to run the test in a VM with a local iSCSI target.
// See iscsi_test_setup.ps1 script

import (
	"context"
	"fmt"
	"testing"

	disk_api "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta2"
	iscsi_api "github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha1"
	system_api "github.com/kubernetes-csi/csi-proxy/client/api/system/v1alpha1"
	disk_client "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1beta2"
	iscsi_client "github.com/kubernetes-csi/csi-proxy/client/groups/iscsi/v1alpha1"
	system_client "github.com/kubernetes-csi/csi-proxy/client/groups/system/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultIscsiPort = 3260
const defaultProtoPort = 0 // default value when port is not set

func TestIscsiAPIGroup(t *testing.T) {
	skipTestOnCondition(t, !shouldRunIscsiTests())

	err := installIscsiTarget()
	require.NoError(t, err, "Failed installing iSCSI target")

	t.Run("List/Add/Remove TargetPortal (Port=3260)", func(t *testing.T) {
		targetPortalTest(t, defaultIscsiPort)
	})

	t.Run("List/Add/Remove TargetPortal (Port not mentioned, effectively 3260)", func(t *testing.T) {
		targetPortalTest(t, defaultProtoPort)
	})

	t.Run("Discover Target and Connect/Disconnect (No CHAP)", func(t *testing.T) {
		targetTest(t)
	})

	t.Run("Discover Target and Connect/Disconnect (CHAP)", func(t *testing.T) {
		targetChapTest(t)
	})

	t.Run("Discover Target and Connect/Disconnect (Mutual CHAP)", func(t *testing.T) {
		targetMutualChapTest(t)
	})

	t.Run("Full flow", func(t *testing.T) {
		e2e_test(t)
	})

}

func e2e_test(t *testing.T) {
	config, err := setupEnv("e2e")
	assert.NoError(t, err)

	defer cleanup(t)

	iscsi, err := iscsi_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, iscsi.Close()) }()

	disk, err := disk_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, disk.Close()) }()

	system, err := system_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &system_api.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi_api.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsi_api.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iscsi.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi_api.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := iscsi.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsi_api.ConnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = iscsi.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	tgtDisksReq := &iscsi_api.GetTargetDisksRequest{TargetPortal: tp, Iqn: config.Iqn}
	tgtDisksResp, err := iscsi.GetTargetDisks(context.TODO(), tgtDisksReq)
	require.Nil(t, err)
	require.Len(t, tgtDisksResp.DiskIds, 1)

	diskId := tgtDisksResp.DiskIds[0]

	attachReq := &disk_api.SetAttachStateRequest{DiskID: diskId, IsOnline: true}
	_, err = disk.SetAttachState(context.TODO(), attachReq)
	require.Nil(t, err)

	partReq := &disk_api.PartitionDiskRequest{DiskID: diskId}
	_, err = disk.PartitionDisk(context.TODO(), partReq)
	assert.Nil(t, err)

	detachReq := &disk_api.SetAttachStateRequest{DiskID: diskId, IsOnline: false}
	_, err = disk.SetAttachState(context.TODO(), detachReq)
	assert.Nil(t, err)
}

func targetTest(t *testing.T) {
	config, err := setupEnv("target")
	assert.NoError(t, err)

	defer cleanup(t)

	client, err := iscsi_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := system_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &system_api.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi_api.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsi_api.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi_api.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := client.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsi_api.ConnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsi_api.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.DisconnectTarget(context.TODO(), disconReq)
	assert.Nil(t, err)
}

func targetChapTest(t *testing.T) {
	const targetName = "chapTarget"
	const username = "someuser"
	const password = "verysecretpass"

	config, err := setupEnv(targetName)
	require.NoError(t, err)

	defer cleanup(t)

	err = setChap(targetName, username, password)
	require.NoError(t, err)

	client, err := iscsi_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := system_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &system_api.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi_api.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsi_api.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi_api.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := client.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsi_api.ConnectTargetRequest{
		TargetPortal: tp,
		Iqn:          config.Iqn,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsi_api.AuthenticationType_ONE_WAY_CHAP,
	}
	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsi_api.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.DisconnectTarget(context.TODO(), disconReq)
	assert.Nil(t, err)
}

func targetMutualChapTest(t *testing.T) {
	const targetName = "mutualChapTarget"
	const username = "anotheruser"
	const password = "broccoli-man"
	const reverse_password = "reversssssssse"

	config, err := setupEnv(targetName)
	require.NoError(t, err)

	defer cleanup(t)

	err = setChap(targetName, username, password)
	require.NoError(t, err)

	err = setReverseChap(targetName, reverse_password)
	require.NoError(t, err)

	client, err := iscsi_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := system_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &system_api.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi_api.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsi_api.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi_api.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := client.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsi_api.ConnectTargetRequest{
		TargetPortal: tp,
		Iqn:          config.Iqn,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsi_api.AuthenticationType_MUTUAL_CHAP,
	}

	// TODO: Replace this with the corresponding csi-proxy API when add it
	// Try using a wrong initiator password and expect error
	out, err := runPowershellCmd(`Set-IscsiChapSecret -ChapSecret "made-up-pass"`)
	require.NoErrorf(t, err, "Cannot set initiator chap out=%v", out)

	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.NotNil(t, err)

	// TODO: Replace this with the corresponding csi-proxy API when add it
	out, err = runPowershellCmd(fmt.Sprintf(`Set-IscsiChapSecret -ChapSecret "%s"`,
		reverse_password))
	require.NoErrorf(t, err, "Cannot set initiator chap out=%v", out)

	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsi_api.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.DisconnectTarget(context.TODO(), disconReq)
	assert.Nil(t, err)
}

func targetPortalTest(t *testing.T, port uint32) {
	config, err := setupEnv(fmt.Sprintf("targetportal-%d", port))
	assert.NoError(t, err)

	defer cleanup(t)

	client, err := iscsi_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := system_client.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &system_api.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi_api.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    port,
	}

	listReq := &iscsi_api.ListTargetPortalsRequest{}

	listResp, err := client.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, listResp.TargetPortals, 0,
			"Expect no registered target portals")
	}

	addTpReq := &iscsi_api.AddTargetPortalRequest{TargetPortal: tp}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	// Port 0 (unset) is handled as the default iSCSI port
	expectedPort := port
	if expectedPort == 0 {
		expectedPort = defaultIscsiPort
	}

	gotListResp, err := client.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, gotListResp.TargetPortals, 1)
		assert.Equal(t, gotListResp.TargetPortals[0].TargetPort, expectedPort)
		assert.Equal(t, gotListResp.TargetPortals[0].TargetAddress, tp.TargetAddress)
	}

	remReq := &iscsi_api.RemoveTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.RemoveTargetPortal(context.Background(), remReq)
	assert.Nil(t, err)

	listResp, err = client.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, listResp.TargetPortals, 0,
			"Expect no registered target portals after delete")
	}
}
