package discovery

import (
	"sync"

	"github.com/gorift/gorift/pkg/server"
)

type MultiDiscovery struct {
	mu          sync.RWMutex
	discoveries map[server.Host]Discovery
}

func NewMultiDiscovery() *MultiDiscovery {
	return &MultiDiscovery{
		discoveries: make(map[server.Host]Discovery),
	}
}

func (d *MultiDiscovery) Register(
	srv server.Server,
	discovery Discovery,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.discoveries[srv.Host]; !ok {
		d.discoveries[srv.Host] = discovery
	}
}

func (d *MultiDiscovery) GetMembers() []*server.Member {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var members []*server.Member
	for _, discover := range d.discoveries {
		members = append(members, discover.GetMembers()...)
	}
	return members
}

func (d *MultiDiscovery) Shutdown() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.discoveries {
		v.Shutdown()
	}
}
