package integrationtests

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1beta2"
	v1beta2client "github.com/kubernetes-csi/csi-proxy/client/groups/filesystem/v1beta2"
)

func v1beta2FilesystemTests(t *testing.T) {
	t.Run("PathExists positive", func(t *testing.T) {
		client, err := v1beta2client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		// simulate FS operations around staging a volume on a node
		stagepath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io\\volume%d", r1.Intn(100), r1.Intn(100)), t)
		mkdirReq := &v1beta2.MkdirRequest{
			Path: stagepath,
		}
		_, err = client.Mkdir(context.Background(), mkdirReq)
		require.NoError(t, err)

		exists, err := pathExists(stagepath)
		assert.True(t, exists, err)

		// simulate operations around publishing a volume to a pod
		podpath := getKubeletPathForTest(fmt.Sprintf("test-pod-id\\volumes\\kubernetes.io~csi\\pvc-test%d", r1.Intn(100)), t)
		mkdirReq = &v1beta2.MkdirRequest{
			Path: podpath,
		}
		_, err = client.Mkdir(context.Background(), mkdirReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath)
		assert.True(t, exists, err)

		sourcePath := stagepath
		targetPath := filepath.Join(podpath, "rootvol")
		// source <- target
		linkReq := &v1beta2.CreateSymlinkRequest{
			SourcePath: sourcePath,
			TargetPath: targetPath,
		}
		_, err = client.CreateSymlink(context.Background(), linkReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath + "\\rootvol")
		assert.True(t, exists, err)

		// cleanup pvpath
		rmdirReq := &v1beta2.RmdirRequest{
			Path:  podpath,
			Force: true,
		}
		_, err = client.Rmdir(context.Background(), rmdirReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath)
		assert.False(t, exists, err)

		// cleanup plugin path
		rmdirReq = &v1beta2.RmdirRequest{
			Path:  stagepath,
			Force: true,
		}
		_, err = client.Rmdir(context.Background(), rmdirReq)
		require.NoError(t, err)

		exists, err = pathExists(stagepath)
		assert.False(t, exists, err)
	})
	t.Run("IsMount", func(t *testing.T) {
		client, err := v1beta2client.NewClient()
		require.Nil(t, err)
		defer client.Close()

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		rand1 := r1.Intn(100)
		rand2 := r1.Intn(100)

		testDir := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io", rand1), t)
		err = os.MkdirAll(testDir, os.ModeDir)
		require.Nil(t, err)
		defer os.RemoveAll(testDir)

		// 1. Check the isMount on a path which does not exist. Failure scenario.
		stagepath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io\\volume%d", rand1, rand2), t)
		IsSymlinkRequest := &v1beta2.IsSymlinkRequest{
			Path: stagepath,
		}
		isSymlink, err := client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.NotNil(t, err)

		// 2. Create the directory. This time its not a mount point. Failure scenario.
		err = os.Mkdir(stagepath, os.ModeDir)
		require.Nil(t, err)
		defer os.Remove(stagepath)
		IsSymlinkRequest = &v1beta2.IsSymlinkRequest{
			Path: stagepath,
		}
		isSymlink, err = client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.Nil(t, err)
		require.Equal(t, isSymlink.IsSymlink, false)

		err = os.Remove(stagepath)
		require.Nil(t, err)
		targetStagePath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io\\volume%d-tgt", rand1, rand2), t)
		lnTargetStagePath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io\\volume%d-tgt-ln", rand1, rand2), t)

		// 3. Create soft link to the directory and make sure target exists. Success scenario.
		err = os.Mkdir(targetStagePath, os.ModeDir)
		require.Nil(t, err)
		defer os.Remove(targetStagePath)
		// Create a sym link
		err = os.Symlink(targetStagePath, lnTargetStagePath)
		require.Nil(t, err)
		defer os.Remove(lnTargetStagePath)

		IsSymlinkRequest = &v1beta2.IsSymlinkRequest{
			Path: lnTargetStagePath,
		}
		isSymlink, err = client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.Nil(t, err)
		require.Equal(t, isSymlink.IsSymlink, true)

		// 4. Remove the path. Failure scenario.
		err = os.Remove(targetStagePath)
		require.Nil(t, err)
		IsSymlinkRequest = &v1beta2.IsSymlinkRequest{
			Path: lnTargetStagePath,
		}
		isSymlink, err = client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.Nil(t, err)
		require.Equal(t, isSymlink.IsSymlink, false)
	})
}
