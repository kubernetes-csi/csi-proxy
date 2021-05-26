package integrationtests

import (
	"testing"
)

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
// Skip on Github Actions as it is expected to fail
func TestDiskAPIGroup(t *testing.T) {
	t.Run("v1beta3Tests", func(t *testing.T) {
		skipTestOnCondition(t, isRunningOnGhActions())
		v1beta3DiskTests(t)
	})
}
