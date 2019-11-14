package server

import (
	"strconv"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/metrics"
)

type Server struct {
	Host Host
	Port Port
}

type Host string

func (h Host) String() string {
	return string(h)
}

type Address string

func (a Address) String() string {
	return string(a)
}

type Port int

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

type Member struct {
	Host    Host
	Address Address
	Port    Port

	HealthStatus HealthStatus

	mu                sync.RWMutex
	metricsRepository metrics.MetricsRepository
}

type HealthStatus struct {
	Available bool
	LastCheck time.Time
}

func NewMember(
	host Host,
	address Address,
	port Port,
	healthStatus HealthStatus,
	metricsRepository metrics.MetricsRepository,
) *Member {
	return &Member{
		Host:    host,
		Address: address,
		Port:    port,

		HealthStatus: healthStatus,

		metricsRepository: metricsRepository,
	}
}

func (m *Member) GetMetrics(label metrics.MetricsLabel) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metric, ok := m.metricsRepository[label]
	if !ok {
		return nil, xerrors.Errorf("%s is not found in metrics repository", label)
	}
	return metric.Get(), nil
}

func (m *Member) AddMetrics(label metrics.MetricsLabel, val interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.metricsRepository[label]
	if !ok {
		return xerrors.Errorf("%s is not found in metrics repository", label)
	}
	return metric.Add(val)
}
