package integrationtests

import (
	"testing"
)

// This test is meant to run on GCE where the page83 ID of the first disk contains
// the host name
// Skip on Github Actions as it is expected to fail
func TestDiskAPIGroup(t *testing.T) {
	t.Run("v1beta3Tests", func(t *testing.T) {
		v1beta3DiskTests(t)
	})
	t.Run("v1beta2Tests", func(t *testing.T) {
		v1beta2DiskTests(t)
	})
	// t.Run("v1beta1Tests", func(t *testing.T) {
	// 	v1beta1DiskTests(t)
	// })
	// t.Run("v1alpha1Tests", func(t *testing.T) {
	// 	v1alpha1DiskTests(t)
	// })
}
