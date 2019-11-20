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

func (r *nopResolver) Lookup(req resolve.Request) (resolve.Report, error) {
	return resolve.Report{
		Addresses: []server.Address{
			server.Address(req.Host),
		},
		LastCheck: time.Now(),
	}, nil
}

func nopHealthcheckFn() healthcheck.Func {
	return healthcheck.Func(func(req healthcheck.Request) (healthcheck.Report, error) {
		return healthcheck.Report{
			Available: true,
			LastCheck: time.Now(),
		}, nil
	})
}

func nopFilterFn() filter.Func {
	return filter.Func(func(members []*server.Member) []*server.Member {
		return members
	})
}

const (
	fakeMetricLabel metrics.Label = "fake"
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
		balance.WithFilterFuncs(nopFilterFn()),
		balance.EnableDiscovery(time.Second, &nopResolver{}),
		balance.EnableHealthcheck(time.Second, nopHealthcheckFn()),
		balance.AddCustomMetrics(metrics.Entry{
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
