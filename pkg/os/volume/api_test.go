package volume

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTarget(t *testing.T) {
	tests := []struct {
		mountpath      string
		expectedResult string
		expectError    bool
		counter        int
	}{
		{
			"c:\\",
			"",
			true,
			1,
		},
	}
	for _, test := range tests {
		target, err := getTarget(test.mountpath, test.counter)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during getTarget(%s)", test.mountpath)
		} else {
			assert.Nil(t, err, "Expect error is nil during getTarget(%s)", test.mountpath)
		}
		assert.Equal(t, target, test.expectedResult, "Expect result not equal with getTarget(%s) return: %q, expected: %s, error: %v",
			test.mountpath, target, test.expectedResult, err)
	}
}
