package healthcheck

import (
	"time"

	"github.com/gorift/gorift/pkg/server"
)

type Func func(Request) (Report, error)

type Request struct {
	Address server.Address
	Port    server.Port
}

type Report struct {
	Available bool
	LastCheck time.Time
}
