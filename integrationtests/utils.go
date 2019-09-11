package integrationtests

import (
	"io"
	"testing"
	"time"

	"github.com/kubernetes-csi/csi-proxy/internal/server"
	"github.com/stretchr/testify/assert"
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
