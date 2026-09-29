package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const gib = 1024 * 1024 * 1024

func TestUnallocatedAfter(t *testing.T) {
	tests := []struct {
		name                   string
		offset, partSize, disk uint64
		want                   uint64
		skipsGetSupportedSize  bool
	}{
		{
			// Values from a 128 GiB data disk whose partition already fills the disk:
			// only the GPT backup header area is left.
			name: "partition fills disk", offset: 16777216, partSize: 137421127680, disk: 128 * gib,
			want: 128*gib - 16777216 - 137421127680, skipsGetSupportedSize: true,
		},
		{
			name: "disk expanded from 128 to 256 GiB", offset: 16777216, partSize: 137421127680, disk: 256 * gib,
			want: 256*gib - 16777216 - 137421127680, skipsGetSupportedSize: false,
		},
		{
			name: "just under the minimum resize size", offset: 0, partSize: 10 * gib, disk: 10*gib + minimumResizeSize - 1,
			want: minimumResizeSize - 1, skipsGetSupportedSize: true,
		},
		{
			name: "exactly the minimum resize size", offset: 0, partSize: 10 * gib, disk: 10*gib + minimumResizeSize,
			want: minimumResizeSize, skipsGetSupportedSize: false,
		},
		{
			name: "partition end past reported disk size", offset: 1024, partSize: 10 * gib, disk: 10 * gib,
			want: 0, skipsGetSupportedSize: true,
		},
	}
	for _, test := range tests {
		got := unallocatedAfter(test.offset, test.partSize, test.disk)
		if got != test.want {
			t.Errorf("%s: unallocatedAfter() = %d, want %d", test.name, got, test.want)
		}
		if skips := got < minimumResizeSize; skips != test.skipsGetSupportedSize {
			t.Errorf("%s: skips GetSupportedSize = %v, want %v", test.name, skips, test.skipsGetSupportedSize)
		}
	}
}

func TestGetTarget(t *testing.T) {
	tests := []struct {
		mountpath      string
		expectedResult string
		expectError    bool
	}{
		{
			"c:\\",
			"",
			true,
		},
	}
	for _, test := range tests {
		target, err := getTarget(test.mountpath)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during getTarget(%s)", test.mountpath)
		} else {
			assert.Nil(t, err, "Expect error is nil during getTarget(%s)", test.mountpath)
		}
		assert.Equal(t, target, test.expectedResult, "Expect result not equal with getTarget(%s) return: %q, expected: %s, error: %v",
			test.mountpath, target, test.expectedResult, err)
	}
}
