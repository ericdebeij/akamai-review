/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/yearmonth"
	"github.com/ericdebeij/akamai-review/v2/report"
	"github.com/spf13/cobra"
)

func init() {
	bycpcodeCmd := &cobra.Command{
		Use:   "usage-cpcode",
		Short: "Reports based on usage as part of the billing data",
		Long: `This report extracts traffic information for the requested period (default the previous period)
		and compares this with the period before`,
		Run: func(cmd *cobra.Command, args []string) {
			defaultParam.Type = "usage-cpcode"
			usagereport(&defaultParam)
		},
	}
	ReportParameters(bycpcodeCmd, "contract", "product", "period", "export")
	rootCmd.AddCommand(bycpcodeCmd)

	byrpgroupCmd := &cobra.Command{
		Use:   "usage-reportinggroup",
		Short: "Reports based on usage as part of the billing data",
		Long: `This report extracts traffic information for the requested period (default the previous period)
		and compares this with the period before`,
		Run: func(cmd *cobra.Command, args []string) {
			defaultParam.Type = "usage-reportinggroup"
			usagereport(&defaultParam)
		},
	}
	ReportParameters(byrpgroupCmd, "contract", "product", "period", "export")
	rootCmd.AddCommand(byrpgroupCmd)
}

func usagereport(rp *ReportFields) {
	if rp.Period == "" {
		rp.Period = yearmonth.Add(yearmonth.FromTime(time.Now()), -1)
	}
	rp.Export = strings.NewReplacer("PERIOD", rp.Period).Replace(rp.Export)

	log.Infof("report %s, period %s, export %s", rp.Type, rp.Period, rp.Export)
	ur := report.BillingUsage{
		EdgeSession:    akamaiSession,
		BillingService: aksv.NewBillingService(akamaiSession),
		CpCodeService:  aksv.NewCpcodeService(akamaiSession),
		ReportType:     rp.Type,
		Contract:       rp.Contract,
		Product:        rp.Product,
		Period:         rp.Period,
		Export:         rp.Export,
	}
	ur.Report()
}
