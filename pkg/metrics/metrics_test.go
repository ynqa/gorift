package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricsRepository(t *testing.T) {
	testCases := []struct {
		entries     []MetricEntry
		expectedLen int
	}{
		{
			expectedLen: 1,
		},
		{
			entries: []MetricEntry{
				{
					Label:  MetricsLabel("fake"),
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
