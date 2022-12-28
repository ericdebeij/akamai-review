/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/ericdebeij/akamai-review/v2/report"
	"github.com/spf13/viper"
)

// certificatesCmd represents the certificates command
var certificatesCmd = &cobra.Command{
	Use:   "certificates",
	Short: "Check your certificates",
	Long: `Collect information regarding the status of your certificates,
	collect hosts from certificate locations in the account and checks
	the related certificate status`,
	Run: func(cmd *cobra.Command, args []string) {
		certificates()
	},
}

func init() {
	rootCmd.AddCommand(certificatesCmd)
	viper.SetDefault("resolver", "8.8.8.8:53")
	viper.SetDefault("export", "export.csv")
	viper.SetDefault("warning.days", 14)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// certificatesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// certificatesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func certificates() {
	certreport := report.CertReport{
		EdgeSession:   akamaiSession,
		DnsService:    akutil.NewDnsService(viper.GetString("resolver")),
		DiagService:   aksv.NewDiagnosticsService(akamaiSession, viper.GetString("akamai.cache")),
		Export:        viper.GetString("export"),
		UseCoverage:   viper.GetBool("input.appsec"),
		UseHostnames:  viper.GetStringSlice("input.hostnames"),
		SkipHostnames: viper.GetStringSlice("input.skiphosts"),
		WarningDays:   viper.GetInt("warning.days"),
	}
	certreport.Report()
}
