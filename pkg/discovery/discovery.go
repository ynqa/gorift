package discovery

import (
	"time"

	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/metrics"
	"github.com/gorift/gorift/pkg/monitor"
	"github.com/gorift/gorift/pkg/monitor/healthcheck"
	"github.com/gorift/gorift/pkg/resolve"
	"github.com/gorift/gorift/pkg/server"
)

type Option struct {
	Interval time.Duration
	Resolver resolve.Resolver
}

func (opt Option) Validate() error {
	if opt.Interval <= 0 {
		return xerrors.New("non-positive interval for ticker")
	}
	if opt.Resolver == nil {
		return xerrors.New("no Resolver")
	}
	return nil
}

type Discovery interface {
	GetMembers() []*server.Member
	Shutdown()
}

func New(
	srv server.Server,
	maybeDiscoveryOption *Option,
	maybeHealthcheckMonitorOption *healthcheck.Option,
	metricsEntries []metrics.Entry,
) Discovery {
	discovery := newNopDiscovery(
		srv,
		maybeHealthcheckMonitorOption,
		metricsEntries,
	)
	if maybeDiscoveryOption != nil {
		if err := maybeDiscoveryOption.Validate(); err == nil {
			discovery = newDefaultDiscovery(
				srv,
				*maybeDiscoveryOption,
				maybeHealthcheckMonitorOption,
				metricsEntries,
			)
		}
	}
	return discovery
}

type nopDiscovery struct {
	srv server.Server

	monitor *monitor.Monitor
}

func newNopDiscovery(
	srv server.Server,
	maybeHealthcheckMonitorOption *healthcheck.Option,
	metricsEntries []metrics.Entry,
) Discovery {
	m := monitor.New(
		server.Address(srv.Host),
		srv.Port,
		maybeHealthcheckMonitorOption,
		metricsEntries,
	)
	return &nopDiscovery{
		srv: srv,

		monitor: m,
	}
}

func (d *nopDiscovery) GetMembers() []*server.Member {
	return []*server.Member{
		server.NewMember(
			d.srv,
			server.Address(d.srv.Host),
			d.monitor.GetHealthStatus(),
			d.monitor.GetMetricsRepository(),
		),
	}
}

func (d *nopDiscovery) Shutdown() {
	d.monitor.Shutdown()
}
