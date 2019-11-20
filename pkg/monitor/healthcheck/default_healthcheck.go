package healthcheck

import (
	"sync"
	"time"

	"github.com/gorift/gorift/pkg/errors/global"
	"github.com/gorift/gorift/pkg/healthcheck"
	"github.com/gorift/gorift/pkg/server"
)

type defaultMonitor struct {
	address server.Address
	port    server.Port
	option  Option

	mu     sync.RWMutex
	status *server.HealthStatus

	doneCh chan struct{}
}

func NewDefaultMonitor(address server.Address, port server.Port, option Option) Monitor {
	monitor := &defaultMonitor{
		address: address,
		port:    port,
		option:  option,

		// [TODO] initial status: Available true/false, or whether do checkFn on here.
		status: &server.HealthStatus{
			Available: true,
		},

		doneCh: make(chan struct{}),
	}
	go monitor.exec()
	return monitor
}

func (m *defaultMonitor) GetHealthStatus() server.HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return *m.status
}

func (m *defaultMonitor) exec() {
	ticker := time.NewTicker(m.option.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.doneCh:
			return
		case <-ticker.C:
			go m.handle()
		}
	}
}

func (m *defaultMonitor) Shutdown() {
	close(m.doneCh)
}

func (m *defaultMonitor) handle() {
	request, err := m.option.Fn(healthcheck.Request{
		Address: m.address,
		Port:    m.port,
	})
	if err != nil {
		global.SendError(err)
		return
	}
	go m.update(request)
}

func (m *defaultMonitor) update(
	report healthcheck.Report,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.status.Available = report.Available
	m.status.LastCheck = report.LastCheck
}
