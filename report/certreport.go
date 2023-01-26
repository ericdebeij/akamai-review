package report

import (
	"context"
	"encoding/csv"
	"fmt"

	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/apex/log"
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
	hostname  string
	status    string
	err       string
	cdn       string
	subject   string
	issuer    string
	expire    time.Time
	iscovered string
}

func (t *testclass) testhost(hostname string) (info *hostinfo) {
	info = t.hosts[hostname]

	ips, _, err := t.dnsService.DnsInfo(hostname)
	if err != nil {
		log.Errorf("dns error %w", err)
	}

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
	EdgeSession    session.Session
	DnsService     *akutil.Dns
	DiagService    *aksv.DiagnosticsService
	Export         string
	UseCoverage    bool
	UseHostnames   []string
	SkipHostnames  []string
	MatchHostnames []string
	WarningDays    int
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
		log.Fatalf("failed to open file %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	secClient := appsec.Client(tst.edgeSession)

	if certreport.UseCoverage {
		coverageRequest := appsec.GetApiHostnameCoverageRequest{}
		x, err := secClient.GetApiHostnameCoverage(context.Background(), coverageRequest)

		if err != nil {
			log.Fatalf("api hostnamecoverage %w", err)
			os.Exit(1)
		}
		for _, ch := range x.HostnameCoverage {
			hn := strings.ToLower(ch.Hostname)

			tst.hosts[hn] = &hostinfo{
				hostname:  hn,
				iscovered: ch.Status,
			}
		}
	}

	for _, hn := range certreport.UseHostnames {
		hn = strings.ToLower(hn)
		_, found := tst.hosts[hn]
		if !found {
			tst.hosts[hn] = &hostinfo{
				hostname:  hn,
				iscovered: "unknown",
			}
		}
	}

	skipre := make([]*regexp.Regexp, 0, 10)
	for _, hn := range certreport.SkipHostnames {
		if hn[0] == '^' {
			re, err := regexp.Compile(hn)
			if err != nil {
				log.Errorf("compile regex %w", err)
			} else {
				skipre = append(skipre, re)
			}
		}

		_, found := tst.hosts[hn]
		if found {
			tst.hosts[hn].status = "skip"
		}
	}
	matchre := make([]*regexp.Regexp, 0, 10)
	matchexact := make([]string, 0, 10)
	for _, hn := range certreport.MatchHostnames {
		if hn[0] == '^' {
			re, err := regexp.Compile(hn)
			if err != nil {
				log.Fatalf("regex compile %w", err)
			} else {
				matchre = append(matchre, re)
			}
		} else {
			matchexact = append(matchexact, hn)
		}
	}
	n := time.Now()
	fmt.Printf("Checking %d hosts", len(tst.hosts))
	r := []string{"hostname", "cdn", "subject", "issuer", "expire", "error", "covered"}
	w.Write(r)
	for hn, hi := range tst.hosts {
		for _, re := range skipre {
			if re.MatchString(hn) {
				hi.status = "skip"
			}
		}

		found := false
		for _, re := range matchre {
			if re.MatchString(hn) {
				found = true
				break
			}
		}
		if !found {
			for _, hx := range matchexact {
				if hx == hn {
					found = true
					break
				}
			}
		}
		if !found && len(certreport.MatchHostnames) > 0 {
			hi.status = "skip"
		}

		if hi.status == "skip" {
			continue
		}

		i := tst.testhost(hn)
		w.Write([]string{i.hostname, i.cdn, i.subject, i.issuer, i.expire.String(), i.err, i.iscovered})
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
