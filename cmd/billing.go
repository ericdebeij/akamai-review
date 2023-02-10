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
	"github.com/spf13/viper"
)

// reportCmd represents the usage command
var reportCmd = &cobra.Command{
	Use:   "billing",
	Short: "Billing-usage report traffic / cpcode",
	Long: `This report extracts traffic information for the requested period (default the previous period)
	and compares this with the period before`,
	Run: func(cmd *cobra.Command, args []string) {
		billingreport()
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

type ReportFields struct {
	Type     string
	Period   string
	Contract string
	Product  string
	Export   string
}

func billingreport() {
	reports := make(map[string]*ReportFields)
	viper.UnmarshalKey("reports", &reports)
	for repname, repdef := range reports {
		if repdef.Period == "" {
			repdef.Period = yearmonth.Add(yearmonth.FromTime(time.Now()), -1)
		}
		repdef.Export = strings.NewReplacer("PERIOD", repdef.Period).Replace(repdef.Export)

		log.Infof("report %s, period %s, export %s", repname, repdef.Period, repdef.Export)
		if repdef.Type == "usagecpcode" {
			ur := report.UsageCpcode{
				EdgeSession:    akamaiSession,
				BillingService: aksv.NewBillingService(akamaiSession),
				CpCodeService:  aksv.NewCpcodeService(akamaiSession),
				Contract:       repdef.Contract,
				Product:        repdef.Product,
				Period:         repdef.Period,
				Export:         repdef.Export,
			}
			ur.Report()
		}
	}
}
