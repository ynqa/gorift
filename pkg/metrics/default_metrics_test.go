package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64Metrics(t *testing.T) {
	testCases := []struct {
		input         interface{}
		isErrOnAdd    bool
		expectedOnGet interface{}
	}{
		{
			input:         int(1),
			isErrOnAdd:    false,
			expectedOnGet: int(1),
		},
		{
			input:      nil,
			isErrOnAdd: true,
		},
		{
			input:      "str",
			isErrOnAdd: true,
		},
	}

	for _, tc := range testCases {
		metric := &intMetric{}
		err := metric.Add(tc.input)
		if tc.isErrOnAdd {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tc.expectedOnGet, metric.Get())
		}
	}
}
