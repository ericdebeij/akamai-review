package akutil

import (
	"github.com/miekg/dns"
)

type Dns struct {
	Resolver string
	Client   *dns.Client
}

func NewDnsService(resolver string) (d *Dns) {
	d = &Dns{
		Resolver: resolver,
	}
	d.Client = new(dns.Client)
	return
}
func (d *Dns) DnsInfo(hostname string) (ips, cnames []string, err error) {

	m := &dns.Msg{
		Question: make([]dns.Question, 1),
	}
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeA)

	in, _, err2 := d.Client.Exchange(m, d.Resolver)
	if err2 != nil {
		err = err2
		return
	}
	ips = make([]string, 0, 2)
	cnames = make([]string, 0, 2)
	for _, x := range in.Answer {
		switch t := x.(type) {
		case *dns.CNAME:
			cnames = append(cnames, t.Target)
		case *dns.A:
			ips = append(ips, t.A.String())
		case *dns.AAAA:
			ips = append(ips, t.AAAA.String())
		}
	}
	return
}
