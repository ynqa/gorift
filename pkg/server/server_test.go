package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostString(t *testing.T) {
	testCases := []struct {
		host     Host
		expected string
	}{
		{
			host:     Host("localhost"),
			expected: "localhost",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.host.String(), tc.expected)
	}
}

func TestAddressString(t *testing.T) {
	testCases := []struct {
		address  Address
		expected string
	}{
		{
			address:  Address("127.0.0.1"),
			expected: "127.0.0.1",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.address.String(), tc.expected)
	}
}

func TestPortString(t *testing.T) {
	testCases := []struct {
		port     Port
		expected string
	}{
		{
			port:     Port(8080),
			expected: "8080",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.port.String(), tc.expected)
	}
}
