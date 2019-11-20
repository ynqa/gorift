package discovery

import (
	"sync"
	"time"

	"github.com/gorift/gorift/pkg/errors/global"
	"github.com/gorift/gorift/pkg/metrics"
	"github.com/gorift/gorift/pkg/monitor"
	"github.com/gorift/gorift/pkg/monitor/healthcheck"
	"github.com/gorift/gorift/pkg/resolve"
	"github.com/gorift/gorift/pkg/server"
)

type defaultDiscovery struct {
	srv                      server.Server
	option                   Option
	healthcheckMonitorOption *healthcheck.Option
	metricsEntries           []metrics.Entry

	mu       sync.RWMutex
	marks    map[server.Address]bool
	monitors map[server.Address]*monitor.Monitor

	doneCh chan struct{}
}

func newDefaultDiscovery(
	srv server.Server,
	option Option,
	maybeHealthcheckMonitorOption *healthcheck.Option,
	metricsEntries []metrics.Entry,
) Discovery {
	marks := make(map[server.Address]bool)
	monitors := make(map[server.Address]*monitor.Monitor)

	// [TODO] initial status: whether it registers host for monitor or not.
	hostAsAddress := server.Address(srv.Host)
	marks[hostAsAddress] = true
	monitors[hostAsAddress] = monitor.New(
		hostAsAddress,
		srv.Port,
		maybeHealthcheckMonitorOption,
		metricsEntries,
	)

	d := &defaultDiscovery{
		srv:                      srv,
		option:                   option,
		healthcheckMonitorOption: maybeHealthcheckMonitorOption,
		metricsEntries:           metricsEntries,

		marks:    marks,
		monitors: monitors,

		doneCh: make(chan struct{}),
	}
	go d.exec()
	return d
}

func (d *defaultDiscovery) GetMembers() []*server.Member {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var members []*server.Member
	for address, monitor := range d.monitors {
		members = append(members, server.NewMember(
			d.srv,
			address,
			monitor.GetHealthStatus(),
			monitor.GetMetricsRepository(),
		))
	}
	return members
}

func (d *defaultDiscovery) exec() {
	ticker := time.NewTicker(d.option.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-d.doneCh:
			go d.shutdownMembers()
			return
		case <-ticker.C:
			go d.handle()
		}
	}
}

func (d *defaultDiscovery) Shutdown() {
	close(d.doneCh)
}

func (d *defaultDiscovery) shutdownMembers() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.monitors {
		v.Shutdown()
	}
}

func (d *defaultDiscovery) handle() {
	report, err := d.option.Resolver.Lookup(
		resolve.Request{
			Host: d.srv.Host,
		},
	)
	if err != nil {
		global.SendError(err)
		return
	}
	go d.update(report)
}

func (d *defaultDiscovery) update(report resolve.Report) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for address := range d.marks {
		d.marks[address] = false
	}

	for _, address := range report.Addresses {
		d.marks[address] = true
	}

	for address, marked := range d.marks {
		if marked {
			if _, ok := d.monitors[address]; !ok {
				// marked and existed
				d.monitors[address] = monitor.New(
					address,
					d.srv.Port,
					d.healthcheckMonitorOption,
					d.metricsEntries,
				)
			}
		} else {
			// not marked
			delete(d.monitors, address)
		}
	}
}
