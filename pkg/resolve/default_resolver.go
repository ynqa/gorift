package resolve

import (
	"net"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/xerrors"

	"github.com/gorift/gorift/pkg/server"
)

const (
	defaultResolvConfigPath = "/etc/resolv.conf"
)

type DefaultResolver struct {
	cfg    *dns.ClientConfig
	client *dns.Client
}

type option struct {
	resolveConfigPath string
}

type DefaultResolverOption func(*option)

func WithResolvConfigPath(path string) DefaultResolverOption {
	return DefaultResolverOption(func(opt *option) {
		opt.resolveConfigPath = path
	})
}

func NewDefaultResolver(opts ...DefaultResolverOption) (Resolver, error) {
	opt := option{
		resolveConfigPath: defaultResolvConfigPath,
	}
	for _, fn := range opts {
		fn(&opt)
	}

	cfg, err := dns.ClientConfigFromFile(opt.resolveConfigPath)
	if err != nil {
		return nil, xerrors.Errorf("failed to create resolver: %w", err)
	}

	return &DefaultResolver{
		cfg:    cfg,
		client: &dns.Client{},
	}, nil
}

func (r *DefaultResolver) Lookup(req Request) (Report, error) {
	m4 := &dns.Msg{}
	m4.SetQuestion(dns.Fqdn(req.Host.String()), dns.TypeA)

	resp, _, err := r.client.Exchange(m4, selectServer(r.cfg))
	if err != nil {
		return Report{}, err
	}

	addresses := make([]server.Address, 0)
	if resp.Rcode == dns.RcodeSuccess {
		for _, ans := range resp.Answer {
			record := ans.(*dns.A)
			addresses = append(addresses, server.Address(record.A.String()))
		}
	}
	return Report{
		Addresses: addresses,
		LastCheck: time.Now(),
	}, nil
}

func selectServer(cfg *dns.ClientConfig) string {
	return net.JoinHostPort(cfg.Servers[0], cfg.Port)
}
