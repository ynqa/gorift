package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorift/gorift/pkg/metrics"
	"github.com/gorift/gorift/pkg/server"
)

func TestP2C(t *testing.T) {
	testCases := []struct {
		members []*server.Member
		isErr   bool
	}{
		{
			members: []*server.Member{
				server.NewMember(server.Server{Host: server.Host("h1"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, metrics.NewMetricsRepository()),
				server.NewMember(server.Server{Host: server.Host("h2"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, metrics.NewMetricsRepository()),
				server.NewMember(server.Server{Host: server.Host("h3"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, metrics.NewMetricsRepository()),
			},
			isErr: false,
		},
		{
			members: []*server.Member{
				server.NewMember(server.Server{Host: server.Host("h1"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, metrics.NewMetricsRepository()),
				server.NewMember(server.Server{Host: server.Host("h2"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, metrics.NewMetricsRepository()),
			},
			isErr: false,
		},
		{
			members: []*server.Member{
				server.NewMember(server.Server{Host: server.Host("h1"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, nil),
			},
			isErr: false,
		},
		{
			members: []*server.Member{
				server.NewMember(server.Server{Host: server.Host("h1"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, nil),
				server.NewMember(server.Server{Host: server.Host("h2"), Port: server.Port(8080)}, server.Address(""), server.HealthStatus{}, nil),
			},
			isErr: true,
		},
		{
			members: []*server.Member{},
			isErr:   true,
		},
		{
			members: nil,
			isErr:   true,
		},
	}

	for _, tc := range testCases {
		algo := NewP2C()
		_, err := algo.Pick(tc.members)
		assert.Equal(t, tc.isErr, err != nil)
	}
}
