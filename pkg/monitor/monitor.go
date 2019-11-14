package monitor

import (
	"github.com/gorift/gorift/pkg/metrics"
	healthcheckmonitor "github.com/gorift/gorift/pkg/monitor/healthcheck"
	"github.com/gorift/gorift/pkg/server"
)

type Monitor struct {
	healthcheckMonitor healthcheckmonitor.Monitor
	metricsRepository  metrics.MetricsRepository

	doneCh chan struct{}
}

func New(
	address server.Address,
	port server.Port,
	maybeHealthcheckMonitorOption *healthcheckmonitor.Option,
	metricsEntries []metrics.MetricEntry,
) *Monitor {
	healthcheckMonitor := healthcheckmonitor.NewNopMonitor()
	if maybeHealthcheckMonitorOption != nil {
		if err := maybeHealthcheckMonitorOption.Validate(); err == nil {
			healthcheckMonitor = healthcheckmonitor.NewDefaultMonitor(
				address, port, *maybeHealthcheckMonitorOption)
		}
	}
	return newMonitor(
		healthcheckMonitor,
		metricsEntries,
	)
}

func newMonitor(
	healthcheckMonitor healthcheckmonitor.Monitor,
	metricsEntries []metrics.MetricEntry,
) *Monitor {
	monitor := &Monitor{
		healthcheckMonitor: healthcheckMonitor,
		metricsRepository:  metrics.NewMetricsRepository(metricsEntries...),

		doneCh: make(chan struct{}),
	}
	go monitor.exec()
	return monitor
}

func (m *Monitor) GetHealthStatus() server.HealthStatus {
	return m.healthcheckMonitor.GetHealthStatus()
}

func (m *Monitor) GetMetricsRepository() metrics.MetricsRepository {
	return m.metricsRepository
}

func (m *Monitor) exec() {
	for {
		select {
		case <-m.doneCh:
			m.healthcheckMonitor.Shutdown()
		}
	}
}

func (m *Monitor) Shutdown() {
	close(m.doneCh)
}
