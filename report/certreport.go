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
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/hako/durafmt"
)

type CertReport struct {
	EdgeSession    session.Session
	DnsService     *akutil.Dns
	DiagService    *aksv.DiagnosticsService
	Export         string
	UseCoverage    bool
	UseHostnames   []string
	UseCps         bool
	SkipHostnames  []string
	MatchHostnames []string
	WarningDays    int
	Contracts      []string
}

type hostinfo struct {
	hostname    string
	status      string
	iscovered   string
	enrollments []string
	//clientinfo *ClientTest
}

func (certreport CertReport) Report() {
	tst := &aksv.ClientTester{}
	hostlist := make(map[string]*hostinfo, 500)

	tst.EdgeSession = certreport.EdgeSession
	tst.DiagService = certreport.DiagService
	defer tst.DiagService.FlushCache()
	tst.DnsService = certreport.DnsService

	f, err := os.Create(certreport.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	secClient := appsec.Client(tst.EdgeSession)

	if certreport.UseCoverage {
		coverageRequest := appsec.GetApiHostnameCoverageRequest{}
		x, err := secClient.GetApiHostnameCoverage(context.Background(), coverageRequest)

		if err != nil {
			log.Fatalf("api hostnamecoverage %w", err)
			os.Exit(1)
		}
		for _, ch := range x.HostnameCoverage {
			hn := strings.ToLower(ch.Hostname)

			hostlist[hn] = &hostinfo{
				hostname:  hn,
				iscovered: ch.Status,
			}
		}
	}

	if certreport.UseCps {
		cpsClient := cps.Client(certreport.EdgeSession)

		for _, contract := range certreport.Contracts {
			listenrollreq := cps.ListEnrollmentsRequest{
				ContractID: contract,
			}
			listenrollresp, err := cpsClient.ListEnrollments(context.Background(), listenrollreq)
			if err != nil {
				log.Fatalf("list enrollments %w", err)
				return
			}
			for _, rl := range listenrollresp.Enrollments {
				hn := strings.ToLower(rl.CSR.CN)
				he, found := hostlist[hn]
				if !found {
					he = &hostinfo{
						hostname:    hn,
						enrollments: []string{},
					}
				}
				//idcode := fmt.Sprintf("%s(%d)", rl.CSR.CN, rl.e  )
				he.enrollments = akutil.UpsertString(he.enrollments, rl.CSR.CN)
				hostlist[hn] = he
				sans := rl.CSR.SANS
				if rl.NetworkConfiguration.DNSNameSettings != nil {
					if len(rl.NetworkConfiguration.DNSNameSettings.DNSNames) > 0 {
						sans = rl.NetworkConfiguration.DNSNameSettings.DNSNames
					}
				}

				for _, san := range sans {
					hn := strings.ToLower(san)
					he, found := hostlist[hn]
					if !found {
						he = &hostinfo{
							hostname:    hn,
							enrollments: []string{},
						}
					}
					//idcode := fmt.Sprintf("%s(%d)", rl.CSR.CN, rl.e  )
					he.enrollments = akutil.UpsertString(he.enrollments, rl.CSR.CN)
					hostlist[hn] = he
				}
			}
		}
	}

	for _, hn := range certreport.UseHostnames {
		hn = strings.ToLower(hn)
		_, found := hostlist[hn]
		if !found {
			hostlist[hn] = &hostinfo{
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

		_, found := hostlist[hn]
		if found {
			hostlist[hn].status = "skip"
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
	fmt.Printf("Checking %d hosts\n", len(hostlist))
	r := []string{"hostname", "cdn", "subject", "issuer", "expire", "error", "covered"}
	w.Write(r)
	for hn, hi := range hostlist {
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

		i := tst.Testhost(hn)
		w.Write([]string{i.Hostname, i.Cdn, i.Subject, i.Issuer, i.Expire.String(), i.Err, hi.iscovered})
		if i.Err == "" && i.Cdn == "akamai" && n.After(i.Expire.AddDate(0, 0, 0-certreport.WarningDays)) {
			fmt.Println("Host       :", i.Hostname)
			fmt.Println("Expire date:", i.Expire)
			fmt.Println("Subject    :", i.Subject)
			fmt.Println("Issuer     :", i.Issuer)
			diff := i.Expire.Sub(n)
			dura := durafmt.Parse(diff)
			fmt.Println("Time left:", dura)
		}
	}

}
