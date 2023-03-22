package aksv

import (
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
)

type ClientTester struct {
	EdgeSession session.Session
	DnsService  *akutil.Dns
	DiagService *DiagnosticsService
	//Hosts       map[string]*ClientTest
}

type ClientInfo struct {
	Hostname string
	Ips      []string
	Status   string
	Err      string
	Cdn      string
	Subject  string
	Issuer   string
	Expire   time.Time
	//Iscovered string
}

func (t *ClientTester) Testhost(hostname string) (info *ClientInfo) {
	info = &ClientInfo{
		Hostname: hostname,
	}

	hn := hostname
	if strings.HasPrefix(hn, "*.") {
		hn = "wildcard" + hn[1:]
	}
	ips, _, err := t.DnsService.DnsInfo(hn)
	if err != nil {
		log.Errorf("dns error %w", err)
	}

	info.Ips = ips

	if len(ips) == 0 {
		info.Cdn = "no-ip"
		return
	}
	_, akamaized, e := t.DiagService.IsAkamaiIp(ips)

	if e != nil {
		info.Err = e.Error()
		return
	}

	if akamaized > 0 {
		info.Cdn = "akamai"
	} else {
		info.Cdn = "other"
	}

	certs, err := akutil.Loadcerts(hostname)
	if err != nil {
		info.Err = err.Error()
		return
	}

	info.Subject = certs[0].Subject.ToRDNSequence().String()
	info.Issuer = certs[0].Issuer.ToRDNSequence().String()
	info.Expire = certs[0].NotAfter

	return
}
