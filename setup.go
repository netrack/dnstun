package tensordns

import (
	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin("tensordns", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	return nil
}
