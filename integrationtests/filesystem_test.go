package integrationtests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
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
	t.Run("v2alpha1FilesystemTests", func(t *testing.T) {
		v2alpha1FilesystemTests(t)
		testMetrics(t)
	})
	t.Run("v1FilesystemTests", func(t *testing.T) {
		v1FilesystemTests(t)
		testMetrics(t)
	})
	t.Run("v1beta2FilesystemTests", func(t *testing.T) {
		v1beta2FilesystemTests(t)
		testMetrics(t)
	})
	t.Run("v1beta1FilesystemTests", func(t *testing.T) {
		v1beta1FilesystemTests(t)
		testMetrics(t)
	})
	t.Run("v1alpha1FilesystemTests", func(t *testing.T) {
		v1alpha1FilesystemTests(t)
		testMetrics(t)
	})
}

func testMetrics(t *testing.T) {
	metricsAddress := "localhost:8888"
	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", metricsAddress))
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check if the response body contains 'grpc_server_handling_seconds_bucket' metrics for CreateSymlink
	if !bytes.Contains(body, []byte("grpc_server_handling_seconds_bucket{grpc_method=\"CreateSymlink\"")) {
		t.Fatalf("Response did not contain 'grpc_server_handling_seconds_bucket' metrics. Response: %s", body)
	}
}
