package apiversion

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVersion(t *testing.T) {
	testCases := []struct {
		name     string
		expected *Version
	}{
		{
			name: "v1",
			expected: &Version{
				major:     1,
				qualifier: stable,
			},
		},
		{
			name: "v2alpha7",
			expected: &Version{
				major:            2,
				qualifier:        alpha,
				qualifierVersion: 7,
			},
		},
		{
			name: "v975beta986654",
			expected: &Version{
				major:            975,
				qualifier:        beta,
				qualifierVersion: 986654,
			},
		},
		{
			name: "1",
		},
		{
			name: "v0",
		},
		{
			name: "v2alpha0",
		},
		{
			name: "whatever",
		},
	}

	for _, testCase := range testCases {
		t.Run("with input "+testCase.name, func(t *testing.T) {
			actual, err := NewVersion(testCase.name)

			if testCase.expected == nil {
				assert.NotNil(t, err)
			} else {
				testCase.expected.rawName = testCase.name
				assert.Equal(t, *testCase.expected, actual)
			}
		})
	}
}

func TestNewVersionOrPanic(t *testing.T) {
	t.Run("with a valid input, it passes the result along", func(t *testing.T) {
		version := NewVersionOrPanic("v8")
		assert.Equal(t, uint(8), version.major)
	})

	t.Run("with an invalid input, it panics", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()

		NewVersionOrPanic("notvalid")
	})
}

func TestVersionComparison(t *testing.T) {
	testCases := []struct {
		v1 string
		v2 string
		// relation is the expected output of v1.Compare(v2)
		relation Comparison
	}{
		{
			v1:       "v1",
			v2:       "v1",
			relation: Equal,
		},
		{
			v1:       "v1alpha1",
			v2:       "v1",
			relation: Lesser,
		},
		{
			v1:       "v2beta8",
			v2:       "v2alpha7",
			relation: Greater,
		},
		{
			v1:       "v8alpha1",
			v2:       "v7",
			relation: Greater,
		},
		{
			v1:       "v7beta12",
			v2:       "v7",
			relation: Lesser,
		},
		{
			v1:       "v2alpha98",
			v2:       "v2alpha98",
			relation: Equal,
		},
		{
			v1:       "v2alpha98",
			v2:       "v2alpha97",
			relation: Greater,
		},
		{
			v1:       "v3beta5",
			v2:       "v3beta5",
			relation: Equal,
		},
		{
			v1:       "v3beta4",
			v2:       "v3beta5",
			relation: Lesser,
		},
		{
			v1:       "v1",
			v2:       "v3",
			relation: Lesser,
		},
	}

	runTestCase := func(v1, v2 Version, relation Comparison) {
		name := v1.rawName + " is "
		switch relation {
		case Lesser:
			name += "lesser than"
		case Equal:
			name += "equal to"
		case Greater:
			name += "greater than"
		}
		name += " " + v2.rawName

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, relation, v1.Compare(v2))
		})
	}

	for _, testCase := range testCases {
		v1, err := NewVersion(testCase.v1)
		require.Nil(t, err)
		v2, err := NewVersion(testCase.v2)
		require.Nil(t, err)

		runTestCase(v1, v2, testCase.relation)

		if testCase.relation == Equal {
			assert.Equal(t, testCase.v1, testCase.v2)
		} else {
			runTestCase(v2, v1, -testCase.relation)
		}
	}
}

func TestVersionToString(t *testing.T) {
	raw := "v8alpha12"
	v, err := NewVersion(raw)
	require.Nil(t, err)

	assert.Equal(t, raw, fmt.Sprintf("%v", v))
}
