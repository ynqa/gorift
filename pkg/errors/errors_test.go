package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergedError(t *testing.T) {
	testCases := []struct {
		errs     []error
		expected string
	}{
		{
			errs:     []error{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		var merged MergedError
		for _, err := range tc.errs {
			merged.Add(err)
		}
		assert.Equal(t, tc.expected, merged.Error())
	}
}
