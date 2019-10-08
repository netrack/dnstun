package dnstun

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("dnstun", setup) }

func setup(c *caddy.Controller) error {

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return NewDNSTun()
	})

	return nil
}
