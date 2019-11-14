package resolve

import (
	"time"

	"github.com/gorift/gorift/pkg/server"
)

type Resolver interface {
	Lookup(ResolveRequest) (ResolveReport, error)
}

type ResolveRequest struct {
	Host server.Host
}

type ResolveReport struct {
	Addresses []server.Address
	LastCheck time.Time
}
