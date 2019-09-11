package integrationtests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/api/filesystem/v1alpha1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/client/filesystem/v1alpha1"
	filesystem "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
)

func TestFilesystemAPIGroup(t *testing.T) {
	defer startServer(t, &filesystem.Server{})()

	t.Run("it works", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer close(t, client)

		path := "/dummy/path"
		request := &v1alpha1.PathExistsRequest{
			Path: path,
		}
		response, err := client.PathExists(context.Background(), request)
		if assert.Nil(t, err) {
			assert.False(t, response.Success)

			if assert.NotNil(t, response.CmdletError) {
				assert.Equal(t, "dummy", response.CmdletError.CmdletName)
				assert.Equal(t, uint32(12), response.CmdletError.Code)
				assert.Equal(t, "hey there "+path, response.CmdletError.Message)
			}
		}
	})
}
