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
var propertiesCmd = &cobra.Command{
	Use:   "properties",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		properties()
	},
}

func init() {
	rootCmd.AddCommand(propertiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// propertiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// propertiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func properties() {
	papiClient := papi.Client(akamaiSession)
	propreport := report.PropReport{
		EdgeSession: akamaiSession,
		DnsService:  akutil.NewDnsService(viper.GetString("resolver")),
		DiagService: aksv.NewDiagnosticsService(akamaiSession, viper.GetString("akamai.cache")),
		PropService: aksv.NewPropertyService(papiClient, viper.GetString("akamai.cache")),
		Export:      viper.GetString("export"),
		Group:       viper.GetString("report.group"),
	}
	propreport.Report()
}
