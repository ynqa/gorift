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
	host server.Host,
	port server.Port,
	maybeDiscoveryOption *Option,
	maybeHealthcheckMonitorOption *healthcheck.Option,
	metricsEntries []metrics.MetricEntry,
) Discovery {
	discovery := newNopDiscovery(
		host,
		port,
		maybeHealthcheckMonitorOption,
		metricsEntries,
	)
	if maybeDiscoveryOption != nil {
		if err := maybeDiscoveryOption.Validate(); err == nil {
			discovery = newDefaultDiscovery(
				host, port,
				*maybeDiscoveryOption,
				maybeHealthcheckMonitorOption,
				metricsEntries,
			)
		}
	}
	return discovery
}

type nopDiscovery struct {
	host server.Host
	port server.Port

	monitor *monitor.Monitor
}

func newNopDiscovery(
	host server.Host,
	port server.Port,
	maybeHealthcheckMonitorOption *healthcheck.Option,
	metricsEntries []metrics.MetricEntry,
) Discovery {
	m := monitor.New(
		server.Address(host),
		port,
		maybeHealthcheckMonitorOption,
		metricsEntries,
	)
	return &nopDiscovery{
		host: host,
		port: port,

		monitor: m,
	}
}

func (d *nopDiscovery) GetMembers() []*server.Member {
	return []*server.Member{
		server.NewMember(
			d.host,
			server.Address(d.host),
			d.port,
			d.monitor.GetHealthStatus(),
			d.monitor.GetMetricsRepository(),
		),
	}
}

func (d *nopDiscovery) Shutdown() {
	d.monitor.Shutdown()
}
