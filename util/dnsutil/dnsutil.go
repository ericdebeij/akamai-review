package dnsutil

import "net"

type Dns struct {
	Resolver string
}

func NewDnsService(resolver string) (d *Dns) {
	d = &Dns{
		Resolver: resolver,
	}
	return
}
func (d *Dns) DnsInfo(hostname string) (ips, cnames []string, err error) {
	var cn string
	cnames = make([]string, 0)

	testhost := hostname
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
