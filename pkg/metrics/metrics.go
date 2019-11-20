package metrics

import (
	"github.com/mohae/deepcopy"
)

type Label string

type Metric interface {
	Add(interface{}) error
	Get() interface{}
}

type Entry struct {
	Label  Label
	Metric Metric
}

type Repository map[Label]Metric

func NewMetricsRepository(entries ...Entry) Repository {
	repository := make(Repository)

	// default metrics for balancer
	repository[TotalPickedLabel] = deepcopy.Copy(TotalPickedMetric).(Metric)

	for _, entry := range entries {
		repository[entry.Label] = deepcopy.Copy(entry.Metric).(Metric)
	}
	return repository
}
