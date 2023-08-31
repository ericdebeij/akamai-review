package securityreport

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/exportx"
	"github.com/ericdebeij/akamai-review/v3/services"
	"github.com/hako/durafmt"
)

type SecHostReport struct {
	Export      string
	WarningDays int
	Match       string
	Skip        string
	HttpTest    bool
}

func (hr *SecHostReport) Report() {
	log.Infof("hosts-certificate %+v", hr)
	srvs := services.Services
	// skip failover hosts and only .com hosts: `^(?!fail).*\.com$`
	var matchre, skipre *regexp.Regexp
	if hr.Match != "" {
		matchre = regexp.MustCompile(hr.Match)
	}
	if hr.Skip != "" {
		skipre = regexp.MustCompile(hr.Skip)
	}

	coverageRequest := appsec.GetApiHostnameCoverageRequest{}
	x, err := srvs.SecClient.GetApiHostnameCoverage(context.Background(), coverageRequest)

	if err != nil {
		log.Fatalf("api hostnamecoverage %v", err)
	}

	csvx, err := exportx.Create(hr.Export)
	if err != nil {
		log.Fatalf("export file %v", err)
		os.Exit(1)
	}
	defer csvx.Close()
	if hr.HttpTest {
		csvx.Header("host", "cdn", "sec-config", "sec-policies", "IPs", "subject-cn", "issuer-cn", "expires", "expire-days", "http-https")
	} else {
		csvx.Header("host", "cdn", "sec-config", "sec-policies", "IPs", "subject-cn", "issuer-cn", "expires", "expire-days")
	}
	nu := time.Now()
	for _, ch := range x.HostnameCoverage {
		hn := strings.ToLower(ch.Hostname)
		//fmt.Print(hn)
		if (skipre != nil && skipre.MatchString(hn)) || (matchre != nil && !matchre.MatchString(hn)) {
			//fmt.Println("..skip")
			continue
		}

		//fmt.Print("..test")
		testresult := srvs.ClientTest.Testhost(hn)
		expiredays := 0
		httptest := ""

		if testresult.Subject != "" {
			expiredays = int(testresult.Expire.Sub(nu).Hours() / 24)
		}

		secconfig := ""
		if ch.Configuration != nil {
			secconfig = ch.Configuration.Name
		}

		if hr.HttpTest {
			if len(testresult.Ips) > 0 {
				//fmt.Print("..httptest..")
				httptest = srvs.ClientTest.TestHttp("http://" + hn + "/")
			}
			csvx.Write(testresult.Hostname, testresult.Cdn, secconfig, ch.PolicyNames, testresult.Ips, testresult.Subject, testresult.Issuer, testresult.Expire, expiredays, httptest)
		} else {
			csvx.Write(testresult.Hostname, testresult.Cdn, secconfig, ch.PolicyNames, testresult.Ips, testresult.Subject, testresult.Issuer, testresult.Expire, expiredays)
		}
		if testresult.Err == "" && testresult.Cdn == "akamai" && nu.After(testresult.Expire.AddDate(0, 0, 0-hr.WarningDays)) {
			fmt.Println("Host       :", testresult.Hostname)
			fmt.Println("Expire date:", testresult.Expire)
			fmt.Println("Subject    :", testresult.Subject)
			fmt.Println("Issuer     :", testresult.Issuer)
			diff := testresult.Expire.Sub(nu)
			dura := durafmt.Parse(diff)
			fmt.Println("Time left:", dura)
		}
		//fmt.Println("ready")
	}
}
