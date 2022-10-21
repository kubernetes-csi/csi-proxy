package integrationtests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	disk "github.com/kubernetes-csi/csi-proxy/pkg/disk"
	diskapi "github.com/kubernetes-csi/csi-proxy/pkg/disk/hostapi"
	iscsi "github.com/kubernetes-csi/csi-proxy/pkg/iscsi"
	iscsiapi "github.com/kubernetes-csi/csi-proxy/pkg/iscsi/hostapi"
	system "github.com/kubernetes-csi/csi-proxy/pkg/system"
	systemapi "github.com/kubernetes-csi/csi-proxy/pkg/system/hostapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultISCSIPort = 3260
const defaultProtoPort = 0 // default value when port is not set

func TestISCSI(t *testing.T) {
	skipTestOnCondition(t, !shouldRunISCSITests())

	err := installISCSITarget()
	require.NoError(t, err, "Failed installing iSCSI target")

	t.Run("List/Add/Remove TargetPortal (Port=3260)", func(t *testing.T) {
		targetPortalTest(t, defaultISCSIPort)
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
		e2eTest(t)
	})

}

func e2eTest(t *testing.T) {
	config, err := setupEnv("e2e")
	require.NoError(t, err)

	defer requireCleanup(t)

	iSCSIClient, err := iscsi.New(iscsiapi.New())
	require.Nil(t, err)

	diskClient, err := disk.New(diskapi.New())
	require.Nil(t, err)

	systemClient, err := system.New(systemapi.New())
	require.Nil(t, err)

	startReq := &system.StartServiceRequest{Name: "MSiSCSI"}
	_, err = systemClient.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultISCSIPort,
	}

	addTpReq := &iscsi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iSCSIClient.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := iSCSIClient.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.IQNs, config.IQN)
	}

	connectReq := &iscsi.ConnectTargetRequest{TargetPortal: tp, IQN: config.IQN}
	_, err = iSCSIClient.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	tgtDisksReq := &iscsi.GetTargetDisksRequest{TargetPortal: tp, IQN: config.IQN}
	tgtDisksResp, err := iSCSIClient.GetTargetDisks(context.TODO(), tgtDisksReq)
	require.Nil(t, err)
	require.Len(t, tgtDisksResp.DiskIDs, 1)

	diskId := tgtDisksResp.DiskIDs[0]
	diskNumber, err := strconv.ParseUint(diskId, 10, 64)
	require.NoError(t, err)

	attachReq := &disk.SetDiskStateRequest{DiskNumber: uint32(diskNumber), IsOnline: true}
	_, err = diskClient.SetDiskState(context.TODO(), attachReq)
	require.Nil(t, err)

	partReq := &disk.PartitionDiskRequest{DiskNumber: uint32(diskNumber)}
	_, err = diskClient.PartitionDisk(context.TODO(), partReq)
	assert.Nil(t, err)

	detachReq := &disk.SetDiskStateRequest{DiskNumber: uint32(diskNumber), IsOnline: false}
	_, err = diskClient.SetDiskState(context.TODO(), detachReq)
	assert.Nil(t, err)
}

func targetTest(t *testing.T) {
	config, err := setupEnv("target")
	require.NoError(t, err)

	defer requireCleanup(t)

	iSCSIClient, err := iscsi.New(iscsiapi.New())
	require.Nil(t, err)

	systemClient, err := system.New(systemapi.New())
	require.Nil(t, err)

	startReq := &system.StartServiceRequest{Name: "MSiSCSI"}
	_, err = systemClient.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultISCSIPort,
	}

	addTpReq := &iscsi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iSCSIClient.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := iSCSIClient.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.IQNs, config.IQN)
	}

	connectReq := &iscsi.ConnectTargetRequest{TargetPortal: tp, IQN: config.IQN}
	_, err = iSCSIClient.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsi.DisconnectTargetRequest{TargetPortal: tp, IQN: config.IQN}
	_, err = iSCSIClient.DisconnectTarget(context.TODO(), disconReq)
	assert.Nil(t, err)
}

func targetChapTest(t *testing.T) {
	const targetName = "chapTarget"
	const username = "someuser"
	const password = "verysecretpass"

	config, err := setupEnv(targetName)
	require.NoError(t, err)

	defer requireCleanup(t)

	err = setChap(targetName, username, password)
	require.NoError(t, err)

	iSCSIClient, err := iscsi.New(iscsiapi.New())
	require.Nil(t, err)

	systemClient, err := system.New(systemapi.New())
	require.Nil(t, err)

	startReq := &system.StartServiceRequest{Name: "MSiSCSI"}
	_, err = systemClient.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultISCSIPort,
	}

	addTpReq := &iscsi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iSCSIClient.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := iSCSIClient.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.IQNs, config.IQN)
	}

	connectReq := &iscsi.ConnectTargetRequest{
		TargetPortal: tp,
		IQN:          config.IQN,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsi.ONE_WAY_CHAP,
	}
	_, err = iSCSIClient.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsi.DisconnectTargetRequest{TargetPortal: tp, IQN: config.IQN}
	_, err = iSCSIClient.DisconnectTarget(context.TODO(), disconReq)
	assert.Nil(t, err)
}

func targetMutualChapTest(t *testing.T) {
	const targetName = "mutualChapTarget"
	const username = "anotheruser"
	const password = "averylongsecret"
	const reverse_password = "reversssssssse"

	config, err := setupEnv(targetName)
	require.NoError(t, err)

	defer requireCleanup(t)

	err = setChap(targetName, username, password)
	require.NoError(t, err)

	err = setReverseChap(targetName, reverse_password)
	require.NoError(t, err)

	iSCSIClient, err := iscsi.New(iscsiapi.New())
	require.Nil(t, err)

	systemClient, err := system.New(systemapi.New())
	require.Nil(t, err)

	{
		req := &system.StartServiceRequest{Name: "MSiSCSI"}
		resp, err := systemClient.StartService(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	tp := &iscsi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultISCSIPort,
	}

	{
		req := &iscsi.AddTargetPortalRequest{
			TargetPortal: tp,
		}
		resp, err := iSCSIClient.AddTargetPortal(context.Background(), req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}

	{
		req := &iscsi.DiscoverTargetPortalRequest{TargetPortal: tp}
		resp, err := iSCSIClient.DiscoverTargetPortal(context.TODO(), req)
		if assert.Nil(t, err) && assert.NotNil(t, resp) {
			assert.Contains(t, resp.IQNs, config.IQN)
		}
	}

	{
		// Try using a wrong initiator password and expect error on connection
		req := &iscsi.SetMutualChapSecretRequest{MutualChapSecret: "made-up-pass"}
		resp, err := iSCSIClient.SetMutualChapSecret(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	connectReq := &iscsi.ConnectTargetRequest{
		TargetPortal: tp,
		IQN:          config.IQN,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsi.MUTUAL_CHAP,
	}

	_, err = iSCSIClient.ConnectTarget(context.TODO(), connectReq)
	assert.NotNil(t, err)

	{
		req := &iscsi.SetMutualChapSecretRequest{MutualChapSecret: reverse_password}
		resp, err := iSCSIClient.SetMutualChapSecret(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	_, err = iSCSIClient.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	{
		req := &iscsi.DisconnectTargetRequest{TargetPortal: tp, IQN: config.IQN}
		resp, err := iSCSIClient.DisconnectTarget(context.TODO(), req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}
}

func targetPortalTest(t *testing.T, port uint32) {
	config, err := setupEnv(fmt.Sprintf("targetportal-%d", port))
	require.NoError(t, err)

	defer requireCleanup(t)

	iSCSIClient, err := iscsi.New(iscsiapi.New())
	require.Nil(t, err)

	systemClient, err := system.New(systemapi.New())
	require.Nil(t, err)

	startReq := &system.StartServiceRequest{Name: "MSiSCSI"}
	_, err = systemClient.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    port,
	}

	listReq := &iscsi.ListTargetPortalsRequest{}

	listResp, err := iSCSIClient.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, listResp.TargetPortals, 0,
			"Expect no registered target portals")
	}

	addTpReq := &iscsi.AddTargetPortalRequest{TargetPortal: tp}
	_, err = iSCSIClient.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	// Port 0 (unset) is handled as the default iSCSI port
	expectedPort := port
	if expectedPort == 0 {
		expectedPort = defaultISCSIPort
	}

	gotListResp, err := iSCSIClient.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, gotListResp.TargetPortals, 1)
		assert.Equal(t, gotListResp.TargetPortals[0].TargetPort, expectedPort)
		assert.Equal(t, gotListResp.TargetPortals[0].TargetAddress, tp.TargetAddress)
	}

	remReq := &iscsi.RemoveTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iSCSIClient.RemoveTargetPortal(context.Background(), remReq)
	assert.Nil(t, err)

	listResp, err = iSCSIClient.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, listResp.TargetPortals, 0,
			"Expect no registered target portals after delete")
	}
}
