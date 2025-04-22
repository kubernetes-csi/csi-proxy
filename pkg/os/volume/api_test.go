package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
