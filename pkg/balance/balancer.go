package balance

import (
	"time"

	"go.uber.org/zap"

	"github.com/gorift/gorift/pkg/balance/algorithm"
	"github.com/gorift/gorift/pkg/balance/middleware/filter"
	"github.com/gorift/gorift/pkg/discovery"
	"github.com/gorift/gorift/pkg/errors"
	"github.com/gorift/gorift/pkg/errors/global"
	"github.com/gorift/gorift/pkg/healthcheck"
	"github.com/gorift/gorift/pkg/metrics"
	healthcheckmonitor "github.com/gorift/gorift/pkg/monitor/healthcheck"
	"github.com/gorift/gorift/pkg/resolve"
	"github.com/gorift/gorift/pkg/server"
)

var (
	defaultLogger       = zap.NewNop()
	defaultAlgorithm    = algorithm.NewRandom()
	defaultFilterFnList = []filter.FilterFn{filter.Availables()}
)

type option struct {
	logger       *zap.Logger
	algorithm    algorithm.Algorithm
	filterFnList []filter.FilterFn

	maybeDiscoveryOption          *discovery.Option
	maybeHealthcheckMonitorOption *healthcheckmonitor.Option
	metricsEntries                []metrics.MetricEntry
}

type BalancerOption func(*option)

func WithZapLogger(logger *zap.Logger) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.logger = logger
	})
}

func WithBalancerAlgorithm(algorithm algorithm.Algorithm) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.algorithm = algorithm
	})
}

func WithFilterFnList(funcs ...filter.FilterFn) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.filterFnList = funcs
	})
}

func EnableDiscovery(
	interval time.Duration,
	resolver resolve.Resolver,
) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.maybeDiscoveryOption = &discovery.Option{
			Interval: interval,
			Resolver: resolver,
		}
	})
}

func EnableHealthcheck(
	interval time.Duration,
	fn healthcheck.HealthcheckFn,
) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.maybeHealthcheckMonitorOption = &healthcheckmonitor.Option{
			Interval: interval,
			Fn:       fn,
		}
	})
}

func AddCustomMetrics(
	entries ...metrics.MetricEntry,
) BalancerOption {
	return BalancerOption(func(opt *option) {
		opt.metricsEntries = entries
	})
}

type Balancer struct {
	logger       *zap.Logger
	algorithm    algorithm.Algorithm
	filterFnList []filter.FilterFn

	maybeDiscoveryOption          *discovery.Option
	maybeHealthcheckMonitorOption *healthcheckmonitor.Option
	metricsEntries                []metrics.MetricEntry

	multiDiscovery *discovery.MultiDiscovery
}

func New(opts ...BalancerOption) (*Balancer, error) {
	opt := option{
		logger:       defaultLogger,
		algorithm:    defaultAlgorithm,
		filterFnList: defaultFilterFnList,
	}

	for _, fn := range opts {
		fn(&opt)
	}

	balancer := &Balancer{
		logger:       opt.logger,
		algorithm:    opt.algorithm,
		filterFnList: opt.filterFnList,

		maybeDiscoveryOption:          opt.maybeDiscoveryOption,
		maybeHealthcheckMonitorOption: opt.maybeHealthcheckMonitorOption,
		metricsEntries:                opt.metricsEntries,

		multiDiscovery: discovery.NewMultiDiscovery(),
	}
	if err := balancer.Validate(); err != nil {
		return nil, err
	}

	go global.LogError(balancer.logger)

	return balancer, nil
}

func (b *Balancer) Validate() error {
	var merged errors.MergedError
	if b.maybeDiscoveryOption != nil {
		merged.Add(b.maybeDiscoveryOption.Validate())
	}
	if b.maybeHealthcheckMonitorOption != nil {
		merged.Add(b.maybeHealthcheckMonitorOption.Validate())
	}
	if merged.Len() <= 0 {
		return nil
	}
	return merged
}

func (b *Balancer) Register(servers ...server.Server) {
	for _, srv := range servers {
		d := discovery.New(
			srv,
			b.maybeDiscoveryOption,
			b.maybeHealthcheckMonitorOption,
			b.metricsEntries,
		)
		b.multiDiscovery.Register(srv, d)
	}
}

func (b *Balancer) GetMembers() []*server.Member {
	return b.multiDiscovery.GetMembers()
}

func (b *Balancer) Pick() (*server.Member, error) {
	candidate := b.multiDiscovery.GetMembers()
	for _, fn := range b.filterFnList {
		candidate = fn(candidate)
	}
	picked, err := b.algorithm.Pick(candidate)
	postPick(picked, err)
	return picked, err
}

func postPick(picked *server.Member, errOnPick error) {
	if errOnPick == nil && picked != nil {
		picked.AddMetrics(metrics.TotalPickedLabel, 1)
	}
}

func (b *Balancer) Shutdown() {
	b.multiDiscovery.Shutdown()
	global.Close()
}
