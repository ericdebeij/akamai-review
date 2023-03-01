/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/ericdebeij/akamai-review/v2/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// propertiesCmd represents the properties command

func init() {

	propcmd := &cobra.Command{
		Use:   "properties",
		Short: "report on properties in the account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			properties("properties")
		},
	}
	ReportParameters(propcmd, "export", "group")
	rootCmd.AddCommand(propcmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "properties-origin",
		Short: "report on origins used in the properties in the account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			properties("origin")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "properties-host",
		Short: "report on hosts used in the properties in the account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			properties("host")
		},
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// propertiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// propertiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//propertiesCmd.Flags().StringVar(&preportname, "report", "", "name of the report, options: origin")
}

func properties(reportname string) {
	defaultParam.Type = "properties-" + reportname
	propreport(&defaultParam)
}

func propreport(rp *ReportFields) {
	papiClient := papi.Client(akamaiSession)
	propreport := report.PropReport{
		EdgeSession: akamaiSession,
		DnsService:  akutil.NewDnsService(viper.GetString("resolver")),
		DiagService: aksv.NewDiagnosticsService(akamaiSession, viper.GetString("akamai.cache")),
		PropService: aksv.NewPropertyService(papiClient, viper.GetString("akamai.cache")),
		Export:      rp.Export,
		Group:       rp.Group,
		ReportType:  rp.Type,
	}
	propreport.Report()

}
