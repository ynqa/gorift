package metrics

import (
	"golang.org/x/xerrors"
)

const (
	TotalPickedLabel MetricsLabel = "Gorift_TotalPicked"
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
