package dnstun

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/pkg/errors"
)

func init() { plugin.Register("dnstun", setup) }

func setup(c *caddy.Controller) error {
	p, err := NewDNSTun(Options{})
	if err != nil {
		return plugin.Error("dnstun", errors.Wrapf(err, "failed to create plugin"))
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return chainPlugin{"dnstun", next, p.ServeDNS}
	})
	return nil
}
