package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kr/pretty"
	"go.uber.org/zap"

	"github.com/gorift/gorift/pkg/balance"
	"github.com/gorift/gorift/pkg/balance/algorithm"
	"github.com/gorift/gorift/pkg/balance/middleware/filter"
	"github.com/gorift/gorift/pkg/healthcheck"
	"github.com/gorift/gorift/pkg/metrics"
	"github.com/gorift/gorift/pkg/resolve"
	"github.com/gorift/gorift/pkg/server"
)

type nopResolver struct{}

func (r *nopResolver) Lookup(req resolve.ResolveRequest) (resolve.ResolveReport, error) {
	return resolve.ResolveReport{
		Addresses: []server.Address{
			server.Address(req.Host),
		},
		LastCheck: time.Now(),
	}, nil
}

func nopHealthcheckFn() healthcheck.HealthcheckFn {
	return healthcheck.HealthcheckFn(func(req healthcheck.HealthcheckRequest) (healthcheck.HealthcheckReport, error) {
		return healthcheck.HealthcheckReport{
			Available: true,
			LastCheck: time.Now(),
		}, nil
	})
}

func nopFilterFn() filter.FilterFn {
	return filter.FilterFn(func(members []*server.Member) []*server.Member {
		return members
	})
}

const (
	fakeMetricLabel metrics.MetricsLabel = "fake"
)

type fakeMetric struct {
	val int
}

func (f *fakeMetric) Add(val interface{}) error {
	f.val += val.(int)
	return nil
}
func (f *fakeMetric) Get() interface{} {
	return f.val
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	balancer, err := balance.New(
		balance.WithZapLogger(logger),
		balance.WithBalancerAlgorithm(algorithm.NewRandom()),
		balance.WithFilterFnList(nopFilterFn()),
		balance.EnableDiscovery(time.Second, &nopResolver{}),
		balance.EnableHealthcheck(time.Second, nopHealthcheckFn()),
		balance.AddCustomMetrics(metrics.MetricEntry{
			Label:  fakeMetricLabel,
			Metric: &fakeMetric{},
		}),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	balancer.Register(
		server.Server{
			Host: server.Host("host1"),
			Port: server.Port(8080),
		},
		server.Server{
			Host: server.Host("host2"),
			Port: server.Port(8080),
		},
	)

	for i := 0; i < 10; i++ {
		member, err := balancer.Pick()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		member.AddMetrics(fakeMetricLabel, 1)
	}

	members := balancer.GetMembers()
	for _, member := range members {
		fmt.Printf("%# v", pretty.Formatter(member))
	}
}
