/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/ericdebeij/akamai-review/v2/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// albCmd represents the alb command
var albCmd = &cobra.Command{
	Use:   "alb",
	Short: "Priving an overview of ALB configuration",
	Long:  `Providing an overview of ALB, functionality will be added once needed`,
	Run: func(cmd *cobra.Command, args []string) {
		alb()
	},
}

var albExportDC string

func init() {
	albCmd.Flags().StringVar(&albExportDC, "exportdc", "", "export datacenters used in alb definitions")
	rootCmd.AddCommand(albCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// albCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// albCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func alb() {
	if albExportDC == "" {
		albExportDC = viper.GetString("alb.exportdc")
	}
	albreport := report.AlbReport{
		EdgeSession: akamaiSession,
		AlbService:  aksv.NewAlbService(akamaiSession),
		DnsService:  akutil.NewDnsService(viper.GetString("resolver")),
		ExportDC:    albExportDC,
	}
	albreport.Report()
}
