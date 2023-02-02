/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"
	"time"

	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/yearmonth"
	"github.com/ericdebeij/akamai-review/v2/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usagetrafficCmd represents the usage command
var usagetrafficCmd = &cobra.Command{
	Use:   "usagetraffic",
	Short: "Billing-usage report traffic / cpcode",
	Long: `This report extracts traffic information for the requested period (default the previous period)
	and compares this with the period before`,
	Run: func(cmd *cobra.Command, args []string) {
		usagetraffic()
	},
}

func init() {
	rootCmd.AddCommand(usagetrafficCmd)
}

func usagetraffic() {
	period := viper.GetString("report.period")
	if period == "" {
		period = yearmonth.Add(yearmonth.FromTime(time.Now()), -1)
	}

	ur := report.UsageReport{
		EdgeSession:    akamaiSession,
		BillingService: aksv.NewBillingService(akamaiSession),
		CpCodeService:  aksv.NewCpcodeService(akamaiSession),
		Contract:       viper.GetString("report.contract"),
		Product:        viper.GetString("report.product"),
		Period:         period,
		Export:         strings.NewReplacer("PERIOD", period).Replace(viper.GetString("report.export")),
	}
	ur.Report()
}
