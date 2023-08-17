package cpsreport

import (
	"context"

	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/exportx"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/ericdebeij/akamai-review/v3/services"
)

type CertReport struct {
	Contract string
	Export   string
}

func (cr *CertReport) Report() {
	srvs := services.Services
	listenrollreq := cps.ListEnrollmentsRequest{
		ContractID: cr.Contract,
	}
	listenrollresp, err := srvs.AkamaiCps.ListEnrollments(context.Background(), listenrollreq)
	if err != nil {
		log.Fatalf("list enrollments %w", err)
	}

	csvx, err := exportx.Create(cr.Export)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvx.Close()

	csvx.Header("cert-cn", "san", "cdn")

	for _, rl := range listenrollresp.Enrollments {
		hostinfo := srvs.ClientTest.Testhost(rl.CSR.CN)
		csvx.Write(rl.CSR.CN, rl.CSR.CN, hostinfo.Cdn)

		sans := rl.CSR.SANS
		if rl.NetworkConfiguration.DNSNameSettings != nil {
			if len(rl.NetworkConfiguration.DNSNameSettings.DNSNames) > 0 {
				sans = rl.NetworkConfiguration.DNSNameSettings.DNSNames
			}
		}

		for _, san := range sans {
			if san != rl.CSR.CN {
				hostinfo := srvs.ClientTest.Testhost(san)
				csvx.Write(rl.CSR.CN, san, hostinfo.Cdn)
			}
		}
	}
}
