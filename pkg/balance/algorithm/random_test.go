package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorift/gorift/pkg/server"
)

func TestRandom(t *testing.T) {
	testCases := []struct {
		members []*server.Member
		isErr   bool
	}{
		{
			members: []*server.Member{
				server.NewMember("h1", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
				server.NewMember("h2", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
			},
			isErr: false,
		},
		{
			members: []*server.Member{
				server.NewMember("h1", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
			},
			isErr: false,
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
		algo := NewRandom()
		_, err := algo.Pick(tc.members)
		assert.Equal(t, tc.isErr, err != nil)
	}
}
