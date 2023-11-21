package dnsutil

import "net"

type dnsinfo struct {
	ips    []string
	cnames []string
	err    error
}
type Dns struct {
	Resolver string
	cache    map[string]dnsinfo
}

func NewDnsService(resolver string) (d *Dns) {
	d = &Dns{
		Resolver: resolver,
		cache:    make(map[string]dnsinfo, 100),
	}
	return
}
func (d *Dns) DnsInfo(hostname string) (ips, cnames []string, err error) {
	var cn string
	cnames = make([]string, 0)

	testhost := hostname
	di, dif := d.cache[testhost]
	if dif {
		ips = di.ips
		cnames = di.cnames
		err = di.err
		return
	}
	for {

		cn, err = net.LookupCNAME(testhost)
		if err != nil || cn == "" || cn == testhost || cn == testhost+"." {
			break
		}
		cnames = append(cnames, cn)
		testhost = cn
	}

	ips, err = net.LookupHost(hostname)
	if dnsErr, ok := err.(*net.DNSError); ok {
		// Check if it's a "no such host" error
		if dnsErr.IsNotFound {
			err = nil
		}
	}
	return
}
