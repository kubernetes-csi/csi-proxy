package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1alpha1"
	v1alpha1client "github.com/kubernetes-csi/csi-proxy/client/groups/filesystem/v1alpha1"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TestFilesystemAPIGroup(t *testing.T) {
	t.Run("PathExists positive", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		// simulate FS operations around staging a volume on a node
		stagepath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\volume%d", r1.Intn(100), r1.Intn(100))
		mkdirReq := &v1alpha1.MkdirRequest{
			Path:    stagepath,
			Context: v1alpha1.PathContext_PLUGIN,
		}
		mkdirRsp, err := client.Mkdir(context.Background(), mkdirReq)
		if assert.Nil(t, err) {
			assert.Equal(t, "", mkdirRsp.Error)
		}

		exists, err := pathExists(stagepath)
		assert.True(t, exists, err)

		// simulate operations around publishing a volume to a pod
		podpath := fmt.Sprintf("C:\\var\\lib\\kubelet\\pods\\test-pod-id\\volumes\\kubernetes.io~csi\\pvc-test%d", r1.Intn(100))
		mkdirReq = &v1alpha1.MkdirRequest{
			Path:    podpath,
			Context: v1alpha1.PathContext_POD,
		}
		mkdirRsp, err = client.Mkdir(context.Background(), mkdirReq)
		if assert.Nil(t, err) {
			assert.Equal(t, "", mkdirRsp.Error)
		}
		linkReq := &v1alpha1.LinkPathRequest{
			SourcePath: podpath + "\\rootvol",
			TargetPath: stagepath,
		}
		linkRsp, err := client.LinkPath(context.Background(), linkReq)
		if assert.Nil(t, err) {
			assert.Equal(t, "", linkRsp.Error)
		}

		exists, err = pathExists(podpath + "\\rootvol")
		assert.True(t, exists, err)

		// cleanup pvpath
		rmdirReq := &v1alpha1.RmdirRequest{
			Path:    podpath,
			Context: v1alpha1.PathContext_POD,
			Force:   true,
		}
		rmdirRsp, err := client.Rmdir(context.Background(), rmdirReq)
		if assert.Nil(t, err) {
			assert.Equal(t, "", rmdirRsp.Error)
		}

		exists, err = pathExists(podpath)
		assert.False(t, exists, err)

		// cleanup plugin path
		rmdirReq = &v1alpha1.RmdirRequest{
			Path:    stagepath,
			Context: v1alpha1.PathContext_PLUGIN,
			Force:   true,
		}
		rmdirRsp, err = client.Rmdir(context.Background(), rmdirReq)
		if assert.Nil(t, err) {
			assert.Equal(t, "", rmdirRsp.Error)
		}

		exists, err = pathExists(stagepath)
		assert.False(t, exists, err)

	})
}
