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
