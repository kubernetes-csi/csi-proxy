package integrationtests

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/pkg/filesystem"
	filesystemapi "github.com/kubernetes-csi/csi-proxy/pkg/filesystem/hostapi"
)

func TestFilesystem(t *testing.T) {
	t.Run("PathExists positive", func(t *testing.T) {
		client, err := filesystem.New(filesystemapi.New())
		require.Nil(t, err)

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		// simulate FS operations around staging a volume on a node
		stagepath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io\\volume%d", r1.Intn(100), r1.Intn(100)), t)
		mkdirReq := &filesystem.MkdirRequest{
			Path: stagepath,
		}
		_, err = client.Mkdir(context.Background(), mkdirReq)
		require.NoError(t, err)

		exists, err := pathExists(stagepath)
		assert.True(t, exists, err)

		// simulate operations around publishing a volume to a pod
		podpath := getKubeletPathForTest(fmt.Sprintf("test-pod-id\\volumes\\kubernetes.io~csi\\pvc-test%d", r1.Intn(100)), t)
		mkdirReq = &filesystem.MkdirRequest{
			Path: podpath,
		}
		_, err = client.Mkdir(context.Background(), mkdirReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath)
		assert.True(t, exists, err)

		sourcePath := stagepath
		targetPath := filepath.Join(podpath, "rootvol")
		// source <- target
		linkReq := &filesystem.CreateSymlinkRequest{
			SourcePath: sourcePath,
			TargetPath: targetPath,
		}
		_, err = client.CreateSymlink(context.Background(), linkReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath + "\\rootvol")
		assert.True(t, exists, err)

		// cleanup pvpath
		rmdirReq := &filesystem.RmdirRequest{
			Path:  podpath,
			Force: true,
		}
		_, err = client.Rmdir(context.Background(), rmdirReq)
		require.NoError(t, err)

		exists, err = pathExists(podpath)
		assert.False(t, exists, err)

		// cleanup plugin path
		rmdirReq = &filesystem.RmdirRequest{
			Path:  stagepath,
			Force: true,
		}
		_, err = client.Rmdir(context.Background(), rmdirReq)
		require.NoError(t, err)

		exists, err = pathExists(stagepath)
		assert.False(t, exists, err)
	})
	t.Run("IsMount", func(t *testing.T) {
		client, err := filesystem.New(filesystemapi.New())
		require.Nil(t, err)

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
		IsSymlinkRequest := &filesystem.IsSymlinkRequest{
			Path: stagepath,
		}
		isSymlink, err := client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.NotNil(t, err)

		// 2. Create the directory. This time its not a mount point. Failure scenario.
		err = os.Mkdir(stagepath, os.ModeDir)
		require.Nil(t, err)
		defer os.Remove(stagepath)
		IsSymlinkRequest = &filesystem.IsSymlinkRequest{
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
		// Create a symlink
		err = os.Symlink(targetStagePath, lnTargetStagePath)
		require.Nil(t, err)
		defer os.Remove(lnTargetStagePath)

		IsSymlinkRequest = &filesystem.IsSymlinkRequest{
			Path: lnTargetStagePath,
		}
		isSymlink, err = client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.Nil(t, err)
		require.Equal(t, isSymlink.IsSymlink, true)

		// 4. Remove the path. Failure scenario.
		err = os.Remove(targetStagePath)
		require.Nil(t, err)
		IsSymlinkRequest = &filesystem.IsSymlinkRequest{
			Path: lnTargetStagePath,
		}
		isSymlink, err = client.IsSymlink(context.Background(), IsSymlinkRequest)
		require.Nil(t, err)
		require.Equal(t, isSymlink.IsSymlink, false)
	})
	t.Run("RmdirContents", func(t *testing.T) {
		client, err := filesystem.New(filesystemapi.New())
		require.Nil(t, err)

		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		rand1 := r1.Intn(100)

		rootPath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io", rand1), t)
		// this line should delete the rootPath because only its content were deleted
		defer os.RemoveAll(rootPath)

		paths := []string{
			filepath.Join(rootPath, "foo/goo/"),
			filepath.Join(rootPath, "foo/bar/baz/"),
			filepath.Join(rootPath, "alpha/beta/gamma/"),
		}
		for _, path := range paths {
			err = os.MkdirAll(path, os.ModeDir)
			require.Nil(t, err)
		}

		rmdirContentsRequest := &filesystem.RmdirContentsRequest{
			Path: rootPath,
		}
		_, err = client.RmdirContents(context.Background(), rmdirContentsRequest)
		require.Nil(t, err)

		// the root path should exist
		exists, err := pathExists(rootPath)
		assert.True(t, exists, err)
		// the root path children shouldn't exist
		for _, path := range paths {
			exists, err = pathExists(path)
			assert.False(t, exists, err)
		}
	})

	t.Run("RmdirContentsNoFollowSymlink", func(t *testing.T) {
		// RmdirContents should not delete the target of a symlink, only the symlink
		client, err := filesystem.New(filesystemapi.New())
		require.Nil(t, err)

		r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
		rand1 := r1.Intn(100)

		rootPath := getKubeletPathForTest(fmt.Sprintf("testplugin-%d.csi.io", rand1), t)
		// this line should delete the rootPath because only its content were deleted
		defer os.RemoveAll(rootPath)

		insidePath := filepath.Join(rootPath, "inside/")
		outsidePath := filepath.Join(rootPath, "outside/")
		paths := []string{
			filepath.Join(insidePath, "foo/goo/"),
			filepath.Join(insidePath, "foo/bar/baz/"),
			filepath.Join(insidePath, "foo/beta/gamma/"),
			outsidePath,
		}
		for _, path := range paths {
			err = os.MkdirAll(path, os.ModeDir)
			require.Nil(t, err)
		}

		// create a temp file on the outside and make a symlink from the inside to the outside
		outsideFile := filepath.Join(outsidePath, "target")
		insideFile := filepath.Join(insidePath, "source")

		file, err := os.Create(outsideFile)
		require.Nil(t, err)
		defer file.Close()
		err = os.Symlink(outsideFile, insideFile)
		require.Nil(t, err)

		rmdirContentsRequest := &filesystem.RmdirContentsRequest{
			Path: insidePath,
		}
		_, err = client.RmdirContents(context.Background(), rmdirContentsRequest)
		require.Nil(t, err)

		// the inside path should exist
		exists, err := pathExists(insidePath)
		require.Nil(t, err)
		assert.True(t, exists, "The path shouldn't exist")
		// it should have no children
		children, err := ioutil.ReadDir(insidePath)
		require.Nil(t, err)
		assert.True(t, len(children) == 0, "The RmdirContents path to delete shouldn't have children")
		// the symlink target should exist
		_, err = os.Open(outsideFile)
		if errors.Is(err, os.ErrNotExist) {
			// the file should exist but it was deleted!
			t.Fatalf("File outsideFile=%s doesn't exist", outsideFile)
		}
	})
}
