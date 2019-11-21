package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricsRepository(t *testing.T) {
	testCases := []struct {
		entries     []Entry
		expectedLen int
	}{
		{
			expectedLen: 1,
		},
		{
			entries: []Entry{
				{
					Label:  Label("fake"),
					Metric: &intMetric{},
				},
			},
			expectedLen: 2,
		},
	}

	for _, tc := range testCases {
		repo := NewMetricsRepository(tc.entries...)
		assert.Equal(t, len(repo), tc.expectedLen)
	}
}

func TestIntMetrics(t *testing.T) {
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
