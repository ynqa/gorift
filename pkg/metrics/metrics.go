package metrics

import (
	"github.com/mohae/deepcopy"
	"golang.org/x/xerrors"
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

const (
	TotalPickedLabel Label = "Gorift_TotalPicked"
)

var (
	TotalPickedMetric Metric = &intMetric{}
)

type intMetric struct {
	val int
}

func (m *intMetric) Add(val interface{}) error {
	var t int
	switch val.(type) {
	case int, int8, int16, int32, int64:
		t = val.(int)
	default:
		return xerrors.Errorf("expected int, but got %v", val)
	}
	m.val += t
	return nil
}

func (m *intMetric) Get() interface{} {
	return m.val
}
