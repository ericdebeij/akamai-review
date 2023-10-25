/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/cpsreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cpsoverviewCmd represents the cpsList command
var cpsoverviewCmd = &cobra.Command{
	Use:     "cps-overview",
	Aliases: []string{"co", "cps-overview"},
	Short:   "Overview of all certificates as defined in cps",
	Long:    `Overview of all certificates, ciphers, tls settings, mtls usage`,
	Run: func(cmd *cobra.Command, args []string) {
		cr := &cpsreport.CertOverviewReport{
			Contract: viperAlias("cps-overview", "contract"),
			Export:   viper.GetString("cps-overview.export"),
		}

		cr.Report()
	},
}

func init() {
	param(cpsoverviewCmd, "export", "cps-overview.export", "cps-overview.csv", "contract to be used")
	param(cpsoverviewCmd, "contract", "cps-overview.contract", "", "contract to be used")
	RootCmd.AddCommand(cpsoverviewCmd)
}
