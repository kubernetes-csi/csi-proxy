package integrationtests

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/api/dummy/v1"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/api/dummy/v1alpha1"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/api/dummy/v1alpha2"
	v1client "github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/client/dummy/v1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/client/dummy/v1alpha1"
	v1alpha2client "github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/client/dummy/v1alpha2"
	"github.com/kubernetes-csi/csi-proxy/integrationtests/apigroups/server/dummy"
)

// This tests the general API structure; it uses a test dummy API group.

func TestAPIGroups(t *testing.T) {
	defer startServer(t, &dummy.Server{})()

	t.Run("happy path with v1alpha1", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		request := &v1alpha1.ComputeDoubleRequest{
			Input32: 28,
		}
		response, err := client.ComputeDouble(context.Background(), request)
		if assert.Nil(t, err) {
			assert.Equal(t, int32(56), response.Response32)
		}
	})

	t.Run("overflow with v1alpha1", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		request := &v1alpha1.ComputeDoubleRequest{
			Input32: math.MaxInt32/2 + 1,
		}
		response, err := client.ComputeDouble(context.Background(), request)
		assert.Nil(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "int32 overflow")
		}
	})

	t.Run("happy path with v1alpha2", func(t *testing.T) {
		client, err := v1alpha2client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		request := &v1alpha2.ComputeDoubleRequest{
			Input64: math.MaxInt32/2 + 1,
		}
		response, err := client.ComputeDouble(context.Background(), request)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(math.MaxInt32+1), response.Response)
		}
	})

	t.Run("overflow with v1alpha2", func(t *testing.T) {
		client, err := v1alpha2client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		request := &v1alpha2.ComputeDoubleRequest{
			Input64: math.MinInt64/2 - 1,
		}
		response, err := client.ComputeDouble(context.Background(), request)
		assert.Nil(t, response)
		if assert.NotNil(t, err) {
			assert.Contains(t, err.Error(), "int64 overflow")
		}
	})

	t.Run("with v1", func(t *testing.T) {
		client, err := v1client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		request := &v1.TellMeAPoemRequest{
			IWantATitle: true,
		}
		response, err := client.TellMeAPoem(context.Background(), request)
		if assert.Nil(t, err) {
			assert.Equal(t, "The New Colossus", response.Title)
		}
	})
}
