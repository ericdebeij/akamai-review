package dnsutil

import (
	"context"
	"net"

	"github.com/apex/log"
)

type dnsinfo struct {
	ips    []string
	cnames []string
	err    error
}
type Dns struct {
	Address  string
	Resolver *net.Resolver
	cache    map[string]dnsinfo
}

func NewDnsService(resolverAddress string) (d *Dns) {
	d = &Dns{
		Address: resolverAddress,
		cache:   make(map[string]dnsinfo, 100),
	}

	log.Infof("using resolver: %v", resolverAddress)
	if resolverAddress == "" || resolverAddress == "default" {
		d.Resolver = net.DefaultResolver
	} else {
		r := &net.Resolver{
			PreferGo:     true,
			StrictErrors: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				c, e := d.DialContext(ctx, "udp", resolverAddress)
				return c, e
			},
		}

		d.Resolver = r
	}
	return
}

func (d *Dns) DnsInfo(hostname string) (ips, cnames []string, err error) {
	var cn string
	cnames = make([]string, 0)

	ip := net.ParseIP(hostname)
	if ip != nil {
		ips = make([]string, 1)
		ips[0] = hostname
		return
	}

	testhost := hostname
	di, dif := d.cache[testhost]
	if dif {
		ips = di.ips
		cnames = di.cnames
		err = di.err
		return
	}
	for {
		cn, err = d.Resolver.LookupCNAME(context.Background(), testhost)
		if err != nil || cn == "" || cn == testhost || cn == testhost+"." {
			break
		}
		cnames = append(cnames, cn)
		testhost = cn
	}

	ips, err = d.Resolver.LookupHost(context.Background(), hostname)

	if dnsErr, ok := err.(*net.DNSError); ok {
		// Check if it's a "no such host" error
		if dnsErr.IsNotFound {
			err = nil
		}
	}
	return
}
