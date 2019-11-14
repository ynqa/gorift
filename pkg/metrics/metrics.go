package metrics

import (
	"github.com/mohae/deepcopy"
)

type MetricsLabel string

type Metric interface {
	Add(interface{}) error
	Get() interface{}
}

type MetricEntry struct {
	Label  MetricsLabel
	Metric Metric
}

type MetricsRepository map[MetricsLabel]Metric

func NewMetricsRepository(entries ...MetricEntry) MetricsRepository {
	repository := make(MetricsRepository)

	// default metrics for balancer
	repository[TotalPickedLabel] = deepcopy.Copy(TotalPickedMetric).(Metric)

	for _, entry := range entries {
		repository[entry.Label] = deepcopy.Copy(entry.Metric).(Metric)
	}
	return repository
}
