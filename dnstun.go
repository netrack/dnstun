package dnstun

//go:generate protoc --proto_path=. --go_out=plugins=grpc:. resolver.proto

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"google.golang.org/grpc"
)

type Options struct {
	Host string
}

// DNSTun is a plugin to block DNS tunneling queries.
type DNSTun struct {
	resolver ResolverClient
}

func NewDNSTun(opts Options) (*DNSTun, error) {
	conn, err := grpc.Dial(opts.Host)
	if err != nil {
		return nil, err
	}
	return &DNSTun{
		resolver: NewResolverClient(conn),
	}, nil
}

func (dt *DNSTun) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	req := ResolveRequest{Name: state.QName()}

	// When the remote resolver is not available, simply forward request
	// processing to the next handler as there is no error.
	resp, err := dt.resolver.Resolve(ctx, &req)
	if err != nil {
		return dns.RcodeSuccess, nil
	}

	switch resp.Action {
	case Action_REJECT:
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeRefused)
		w.WriteMsg(m)

		return dns.RcodeRefused, nil
	}

	// Pass the control to the next plugin.
	return dns.RcodeSuccess, nil
}

type chainPlugin struct {
	name     string
	next     plugin.Handler
	serveDNS func(context.Context, dns.ResponseWriter, *dns.Msg) (int, error)
}

// Name returns the name of the plugin. This method implements plugin.Handler
// interface.
func (p chainPlugin) Name() string {
	return p.name
}

func (p chainPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	rcode, err := p.serveDNS(ctx, w, r)
	if rcode != dns.RcodeSuccess {
		return rcode, err
	}

	state := request.Request{W: w, Req: r}
	return plugin.NextOrFailure(state.Name(), p.next, ctx, w, r)
}
