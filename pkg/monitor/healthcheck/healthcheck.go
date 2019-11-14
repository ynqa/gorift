package healthcheck

import (
	"time"

	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/healthcheck"
	"github.com/gorift/gorift/pkg/server"
)

type Option struct {
	Interval time.Duration
	Fn       healthcheck.HealthcheckFn
}

func (opt Option) Validate() error {
	if opt.Interval <= 0 {
		return xerrors.New("non-positive interval for ticker")
	}
	if opt.Fn == nil {
		return xerrors.New("no HealthcheckFn")
	}
	return nil
}

type Monitor interface {
	GetHealthStatus() server.HealthStatus
	Shutdown()
}

type nopMonitor struct{}

func NewNopMonitor() Monitor {
	return &nopMonitor{}
}

func (m *nopMonitor) GetHealthStatus() server.HealthStatus {
	return server.HealthStatus{
		Available: true,
		LastCheck: time.Now(),
	}
}

func (m *nopMonitor) Shutdown() {}
