package integrationtests

import (
	"context"
	"fmt"
	"testing"

	diskApi "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1beta3"
	iscsiApi "github.com/kubernetes-csi/csi-proxy/client/api/iscsi/v1alpha2"
	systemApi "github.com/kubernetes-csi/csi-proxy/client/api/system/v1alpha1"
	diskClient "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1beta3"
	iscsiClient "github.com/kubernetes-csi/csi-proxy/client/groups/iscsi/v1alpha2"
	systemClient "github.com/kubernetes-csi/csi-proxy/client/groups/system/v1alpha1"

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
		e2eTest(t)
	})

}

func e2eTest(t *testing.T) {
	config, err := setupEnv("e2e")
	require.NoError(t, err)

	defer requireCleanup(t)

	iscsi, err := iscsiClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, iscsi.Close()) }()

	disk, err := diskClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, disk.Close()) }()

	system, err := systemClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &systemApi.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsiApi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsiApi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = iscsi.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsiApi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := iscsi.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsiApi.ConnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = iscsi.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	tgtDisksReq := &iscsiApi.GetTargetDisksRequest{TargetPortal: tp, Iqn: config.Iqn}
	tgtDisksResp, err := iscsi.GetTargetDisks(context.TODO(), tgtDisksReq)
	require.Nil(t, err)
	require.Len(t, tgtDisksResp.DiskIDs, 1)

	diskId := tgtDisksResp.DiskIDs[0]

	attachReq := &diskApi.SetAttachStateRequest{DiskID: diskId, IsOnline: true}
	_, err = disk.SetAttachState(context.TODO(), attachReq)
	require.Nil(t, err)

	partReq := &diskApi.PartitionDiskRequest{DiskID: diskId}
	_, err = disk.PartitionDisk(context.TODO(), partReq)
	assert.Nil(t, err)

	detachReq := &diskApi.SetAttachStateRequest{DiskID: diskId, IsOnline: false}
	_, err = disk.SetAttachState(context.TODO(), detachReq)
	assert.Nil(t, err)
}

func targetTest(t *testing.T) {
	config, err := setupEnv("target")
	require.NoError(t, err)

	defer requireCleanup(t)

	client, err := iscsiClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := systemClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &systemApi.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsiApi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsiApi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsiApi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := client.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsiApi.ConnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsiApi.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.DisconnectTarget(context.TODO(), disconReq)
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

	client, err := iscsiClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := systemClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &systemApi.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsiApi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	addTpReq := &iscsiApi.AddTargetPortalRequest{
		TargetPortal: tp,
	}
	_, err = client.AddTargetPortal(context.Background(), addTpReq)
	assert.Nil(t, err)

	discReq := &iscsiApi.DiscoverTargetPortalRequest{TargetPortal: tp}
	discResp, err := client.DiscoverTargetPortal(context.TODO(), discReq)
	if assert.Nil(t, err) {
		assert.Contains(t, discResp.Iqns, config.Iqn)
	}

	connectReq := &iscsiApi.ConnectTargetRequest{
		TargetPortal: tp,
		Iqn:          config.Iqn,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsiApi.AuthenticationType_ONE_WAY_CHAP,
	}
	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	disconReq := &iscsiApi.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
	_, err = client.DisconnectTarget(context.TODO(), disconReq)
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

	client, err := iscsiClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := systemClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	{
		req := &systemApi.StartServiceRequest{Name: "MSiSCSI"}
		resp, err := system.StartService(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	tp := &iscsiApi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    defaultIscsiPort,
	}

	{
		req := &iscsiApi.AddTargetPortalRequest{
			TargetPortal: tp,
		}
		resp, err := client.AddTargetPortal(context.Background(), req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}

	{
		req := &iscsiApi.DiscoverTargetPortalRequest{TargetPortal: tp}
		resp, err := client.DiscoverTargetPortal(context.TODO(), req)
		if assert.Nil(t, err) && assert.NotNil(t, resp) {
			assert.Contains(t, resp.Iqns, config.Iqn)
		}
	}

	{
		// Try using a wrong initiator password and expect error on connection
		req := &iscsiApi.SetMutualChapSecretRequest{MutualChapSecret: "made-up-pass"}
		resp, err := client.SetMutualChapSecret(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	connectReq := &iscsiApi.ConnectTargetRequest{
		TargetPortal: tp,
		Iqn:          config.Iqn,
		ChapUsername: username,
		ChapSecret:   password,
		AuthType:     iscsiApi.AuthenticationType_MUTUAL_CHAP,
	}

	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.NotNil(t, err)

	{
		req := &iscsiApi.SetMutualChapSecretRequest{MutualChapSecret: reverse_password}
		resp, err := client.SetMutualChapSecret(context.TODO(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	_, err = client.ConnectTarget(context.TODO(), connectReq)
	assert.Nil(t, err)

	{
		req := &iscsiApi.DisconnectTargetRequest{TargetPortal: tp, Iqn: config.Iqn}
		resp, err := client.DisconnectTarget(context.TODO(), req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}
}

func targetPortalTest(t *testing.T, port uint32) {
	config, err := setupEnv(fmt.Sprintf("targetportal-%d", port))
	require.NoError(t, err)

	defer requireCleanup(t)

	client, err := iscsiClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, client.Close()) }()

	system, err := systemClient.NewClient()
	require.Nil(t, err)

	defer func() { assert.NoError(t, system.Close()) }()

	startReq := &systemApi.StartServiceRequest{Name: "MSiSCSI"}
	_, err = system.StartService(context.TODO(), startReq)
	require.NoError(t, err)

	tp := &iscsiApi.TargetPortal{
		TargetAddress: config.Ip,
		TargetPort:    port,
	}

	listReq := &iscsiApi.ListTargetPortalsRequest{}

	listResp, err := client.ListTargetPortals(context.Background(), listReq)
	if assert.Nil(t, err) {
		assert.Len(t, listResp.TargetPortals, 0,
			"Expect no registered target portals")
	}

	addTpReq := &iscsiApi.AddTargetPortalRequest{TargetPortal: tp}
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

	remReq := &iscsiApi.RemoveTargetPortalRequest{
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
