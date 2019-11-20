package resolve

import (
	"time"

	"github.com/gorift/gorift/pkg/server"
)

type Resolver interface {
	Lookup(Request) (Report, error)
}

type Request struct {
	Host server.Host
}

type Report struct {
	Addresses []server.Address
	LastCheck time.Time
}
