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
	t.Run("IsMount", func(t *testing.T) {
		client, err := v1alpha1client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		rand1 := r1.Intn(100)
		rand2 := r1.Intn(100)

		testDir := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io", rand1)
		err = os.MkdirAll(testDir, os.ModeDir)
		require.Nil(t, err)
		defer os.RemoveAll(testDir)

		// 1. Check the isMount on a path which does not exist. Failure scenario.
		stagepath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\volume%d", rand1, rand2)
		isMountRequest := &v1alpha1.IsMountPointRequest{
			Path: stagepath,
		}
		isMountResponse, err := client.IsMountPoint(context.Background(), isMountRequest)
		require.Nil(t, err)
		require.Equal(t, isMountResponse.IsMountPoint, false)

		// 2. Create the directory. This time its not a mount point. Failure scenario.
		err = os.Mkdir(stagepath, os.ModeDir)
		require.Nil(t, err)
		defer os.Remove(stagepath)
		isMountRequest = &v1alpha1.IsMountPointRequest{
			Path: stagepath,
		}
		isMountResponse, err = client.IsMountPoint(context.Background(), isMountRequest)
		require.Nil(t, err)
		require.Equal(t, isMountResponse.IsMountPoint, false)

		err = os.Remove(stagepath)
		require.Nil(t, err)
		targetStagePath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\volume%d-tgt", rand1, rand2)
		lnTargetStagePath := fmt.Sprintf("C:\\var\\lib\\kubelet\\plugins\\testplugin-%d.csi.io\\volume%d-tgt-ln", rand1, rand2)

		// 3. Create soft link to the directory and make sure target exists. Success scenario.
		os.Mkdir(targetStagePath, os.ModeDir)
		require.Nil(t, err)
		defer os.Remove(targetStagePath)
		// Create a sym link
		err = os.Symlink(targetStagePath, lnTargetStagePath)
		require.Nil(t, err)
		defer os.Remove(lnTargetStagePath)

		isMountRequest = &v1alpha1.IsMountPointRequest{
			Path: lnTargetStagePath,
		}
		isMountResponse, err = client.IsMountPoint(context.Background(), isMountRequest)
		require.Nil(t, err)
		require.Equal(t, isMountResponse.IsMountPoint, true)

		// 4. Remove the path. Failure scenario.
		err = os.Remove(targetStagePath)
		require.Nil(t, err)
		isMountRequest = &v1alpha1.IsMountPointRequest{
			Path: lnTargetStagePath,
		}
		isMountResponse, err = client.IsMountPoint(context.Background(), isMountRequest)
		require.Nil(t, err)
		require.Equal(t, isMountResponse.IsMountPoint, false)
	})
}
