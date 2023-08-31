package clienttest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/service/diagnostics"
	"github.com/ericdebeij/akamai-review/v3/util/certutil"
	"github.com/ericdebeij/akamai-review/v3/util/dnsutil"
)

type ClientTester struct {
	EdgeSession session.Session
	DnsService  *dnsutil.Dns
	DiagService *diagnostics.DiagnosticsService
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

func (t *ClientTester) TestHttp(url string) (info string) {

	client := &http.Client{
		Timeout: time.Second * 4,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Debugf("http test", url, err)
		info = "error"
		return
	}

	loc5 := ""
	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		loc5 = location[:5]
	}

	bywho := ""
	srv := resp.Header.Get("Server")
	if srv == "AkamaiGHost" {
		bywho = "akamai"
	} else if srv != "" {
		bywho = "origin"
	}
	info = fmt.Sprintf("%v-%v-%v", resp.StatusCode, loc5, bywho)
	log.Debugf("http test", url, resp.StatusCode, resp.Header.Get("Location"), resp.Header.Get("Server"))
	return
}

func (t *ClientTester) Testhost(hostname string) (info *ClientInfo) {
	log.Debugf("Test host: %s", hostname)
	info = &ClientInfo{
		Hostname: hostname,
	}

	hn := hostname
	if strings.HasPrefix(hn, "*.") {
		hn = "wildcard" + hn[1:]
	}

	ips, _, err := t.DnsService.DnsInfo(hn)
	if err != nil {
		log.Debugf("dns error, host %v, %v, sleep 5s, retry", hn, err)
		time.Sleep(5 * time.Second)
		ips, _, err = t.DnsService.DnsInfo(hn)
		if err != nil {
			log.Errorf("dns (retry) error, host %v, %v", hn, err)
		}
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

	certs, err := certutil.Loadcerts(hostname)
	log.Debugf("Certificates loaded for hostname %s: %v", hostname, err)
	if err != nil {
		info.Err = err.Error()
		return
	}

	//info.Subject = certs[0].Subject.ToRDNSequence().String()
	info.Subject = certs[0].Subject.CommonName
	//info.Issuer = certs[0].Issuer.ToRDNSequence().String()
	info.Issuer = certs[0].Issuer.CommonName
	info.Expire = certs[0].NotAfter

	return
}
