package dnstun

//go:generate protoc --proto_path=. --go_out=plugins=grpc:. resolver.proto

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// DNSTun is a plugin to block DNS tunneling queries.
type DNSTun struct {
	Next plugin.Handler

	Resolver ResolverClient
}

func NewDNSTun() *DNSTun {
	return &DNSTun{
		Next:     nil,
		Resolver: NewResolverClient(nil),
	}
}

// Name returns the name of the plugin. This method implements plugin.Handler
// interface.
func (dt *DNSTun) Name() string {
	return "dnstun"
}

func (dt *DNSTun) ServeDNS(ctx context.Context, rw dns.ResponseWriter, r *dns.Msg) (int, error) {
	return 0, nil
}
