package filter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorift/gorift/pkg/server"
)

func TestAvailables(t *testing.T) {
	testCases := []struct {
		members  []*server.Member
		expected []*server.Member
	}{
		{
			members: []*server.Member{
				server.NewMember("h1", server.Address(""), server.Port(8080), server.HealthStatus{Available: true}, nil),
				server.NewMember("h2", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
			},
			expected: []*server.Member{
				server.NewMember("h1", server.Address(""), server.Port(8080), server.HealthStatus{Available: true}, nil),
			},
		},
		{
			members: []*server.Member{
				server.NewMember("h1", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
				server.NewMember("h2", server.Address(""), server.Port(8080), server.HealthStatus{}, nil),
			},
			expected: []*server.Member{},
		},
		{
			members:  []*server.Member{},
			expected: []*server.Member{},
		},
		{
			members:  nil,
			expected: []*server.Member{},
		},
	}

	for _, tc := range testCases {
		filtered := Availables()(tc.members)
		assert.True(t, reflect.DeepEqual(tc.expected, filtered))
	}
}
