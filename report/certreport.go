package report

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/hako/durafmt"
)

type testclass struct {
	edgeSession session.Session
	dnsService  *akutil.Dns
	diagService *aksv.DiagnosticsService
	hosts       map[string]*hostinfo
}

type hostinfo struct {
	hostname string
	status   string
	err      string
	cdn      string
	subject  string
	issuer   string
	expire   time.Time
}

func (t *testclass) testhost(hostname string) (info *hostinfo) {
	info = t.hosts[hostname]

	ips, _ := t.dnsService.DnsInfo(hostname)
	if len(ips) == 0 {
		info.cdn = "no-ip"
		return
	}
	_, akamaized, e := t.diagService.IsAkamaiIp(ips)

	if e != nil {
		info.err = e.Error()
		return
	}

	if akamaized > 0 {
		info.cdn = "akamai"
	} else {
		info.cdn = "other"
	}

	certs, err := akutil.Loadcerts(hostname)
	if err != nil {
		info.err = err.Error()
		return
	}

	info.subject = certs[0].Subject.ToRDNSequence().String()
	info.issuer = certs[0].Issuer.ToRDNSequence().String()
	info.expire = certs[0].NotAfter

	return
}

type CertReport struct {
	EdgeSession   session.Session
	DnsService    *akutil.Dns
	DiagService   *aksv.DiagnosticsService
	Export        string
	UseCoverage   bool
	UseHostnames  []string
	SkipHostnames []string
	WarningDays   int
}

func (certreport CertReport) Report() {
	tst := &testclass{}
	tst.hosts = make(map[string]*hostinfo, 500)

	tst.edgeSession = certreport.EdgeSession
	tst.diagService = certreport.DiagService
	defer tst.diagService.FlushCache()
	tst.dnsService = certreport.DnsService

	f, err := os.Create(certreport.Export)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	r := []string{"hostname", "cdn", "subject", "issuer", "expire", "error"}
	w.Write(r)

	secClient := appsec.Client(tst.edgeSession)

	if certreport.UseCoverage {
		coverageRequest := appsec.GetApiHostnameCoverageRequest{}
		x, err := secClient.GetApiHostnameCoverage(context.Background(), coverageRequest)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, ch := range x.HostnameCoverage {
			tst.hosts[ch.Hostname] = &hostinfo{
				hostname: ch.Hostname,
			}
		}
	}

	for _, hn := range certreport.UseHostnames {
		tst.hosts[hn] = &hostinfo{
			hostname: hn,
		}
	}

	skiphosts := certreport.SkipHostnames
	skipre := make([]*regexp.Regexp, 0, 10)
	for _, hn := range skiphosts {
		if hn[0] == '^' {
			re, err := regexp.Compile(hn)
			if err != nil {
				log.Print(err)
			} else {
				skipre = append(skipre, re)
			}
		}

		_, found := tst.hosts[hn]
		if found {
			tst.hosts[hn].status = "skip"
		}
	}

	n := time.Now()
	fmt.Printf("Checking %d hosts\n", len(tst.hosts))
	for hn, hi := range tst.hosts {
		for _, re := range skipre {
			if re.MatchString(hn) {
				hi.status = "skip"
			}
		}

		if hi.status == "skip" {
			continue
		}
		i := tst.testhost(hn)
		w.Write([]string{i.hostname, i.cdn, i.subject, i.issuer, i.expire.String(), i.err})
		if i.err == "" && i.cdn == "akamai" && n.After(i.expire.AddDate(0, 0, 0-certreport.WarningDays)) {
			fmt.Println("Host       :", i.hostname)
			fmt.Println("Expire date:", i.expire)
			fmt.Println("Subject    :", i.subject)
			fmt.Println("Issuer     :", i.issuer)
			diff := i.expire.Sub(n)
			dura := durafmt.Parse(diff)
			fmt.Println("Time left:", dura)
		}
	}

}
