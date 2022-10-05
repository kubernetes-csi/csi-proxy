package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathValid(t *testing.T) {
	tests := []struct {
		remotepath     string
		expectedResult bool
		expectError    bool
	}{
		{
			"c:",
			true,
			false,
		},
		{
			"invalid-path",
			false,
			false,
		},
	}

	for _, test := range tests {
		result, err := pathValid(test.remotepath)
		assert.Equal(t, result, test.expectedResult, "Expect result not equal with pathValid(%s) return: %q, expected: %q, error: %v",
			test.remotepath, result, test.expectedResult, err)
		if test.expectError {
			assert.NotNil(t, err, "Expect error during pathValid(%s)", test.remotepath)
		} else {
			assert.Nil(t, err, "Expect error is nil during pathValid(%s)", test.remotepath)
		}
	}
}
