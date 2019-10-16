package integrationtests

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kubernetes-csi/csi-proxy/internal/server"
)

// startServer starts the proxy's GRPC servers, and returns a function to shut them down when done with testing
func startServer(t *testing.T, apiGroups ...server.APIGroup) func() {
	s := server.NewServer(apiGroups...)

	listeningChan := make(chan interface{})
	go func() {
		assert.Nil(t, s.Start(listeningChan))
	}()

	select {
	case <-listeningChan:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for GRPC servers to start listening")
	}

	return func() {
		assert.Nil(t, s.Stop())
	}
}

func close(t *testing.T, closer io.Closer) {
	assert.Nil(t, closer.Close())
}

// recursiveDiff ensures that dir1 and dir2 contain the same files, with the same contents.
// fileSuffixesToRemove will be removed from file names, if found.
func recursiveDiff(t *testing.T, dir1, dir2 string, fileSuffixesToRemove ...string) {
	hashesDir1, err := fileHashes(dir1, fileSuffixesToRemove...)
	require.Nil(t, err, "unable to get file hashes for directory %q", dir1)
	hashesDir2, err := fileHashes(dir2, fileSuffixesToRemove...)
	require.Nil(t, err, "unable to get file hashes for directory %q", dir2)

	for filePath, hash1 := range hashesDir1 {
		if hash2, present := hashesDir2[filePath]; assert.True(t, present, "%q present in %q but not in %q", filePath, dir1, dir2) {
			if hash1 != hash2 {
				contents1 := readFile(t, filepath.Join(dir1, filePath))
				contents2 := readFile(t, filepath.Join(dir2, filePath))

				differ := diffmatchpatch.New()
				diffs := differ.DiffMain(string(contents1), string(contents2), true)
				assert.Fail(t, fmt.Sprintf("File %q differs in %q and %q:\n", filePath, dir1, dir2), "Diff:\n%s", differ.DiffPrettyText(diffs))
			}
			delete(hashesDir2, filePath)
		}
	}

	for filePath := range hashesDir2 {
		assert.Fail(t, fmt.Sprintf("%q present in %q but not in %q", filePath, dir2, dir1))
	}

	return
}

// fileHashes walks through dir, and returns a map mapping file paths to MD5 hashes
func fileHashes(dir string, fileSuffixesToRemove ...string) (map[string]string, error) {
	dir = strings.ReplaceAll(dir, "/", string(os.PathSeparator))
	hashes := make(map[string]string)

	if walkErr := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "unable to descend into %q", filePath)
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "unable to open %q", filePath)
		}
		defer file.Close()

		hasher := md5.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return errors.Wrapf(err, "unable to read %q", filePath)
		}

		hashBytes := hasher.Sum(nil)[:16]

		relativePath := strings.TrimPrefix(strings.TrimPrefix(filePath, dir), "/")
		for _, suffix := range fileSuffixesToRemove {
			relativePath = strings.TrimSuffix(relativePath, suffix)
		}

		hashes[relativePath] = hex.EncodeToString(hashBytes)
		return nil
	}); walkErr != nil {
		return nil, walkErr
	}

	return hashes, nil
}

func readFile(t *testing.T, filePath string) string {
	contents, err := ioutil.ReadFile(filePath)
	require.Nil(t, err, "unable to read %q", filePath)
	return string(contents)
}
