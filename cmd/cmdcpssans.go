/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/cpsreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cpssansCmd represents the cpsList command
var cpssansCmd = &cobra.Command{
	Use:     "cps-sans",
	Aliases: []string{"cs", "cps-s"},
	Short:   "List SANS/certificates as defined in cps and the usage of the SANs",
	Long:    `List of the certificates, the SAN in the certificates. Additional information is provided to check whether the CN or SAN entry is actually served via Akamai`,
	Run: func(cmd *cobra.Command, args []string) {
		cr := &cpsreport.CertSanReport{
			Contract: viperAlias("cps-sans", "contract"),
			Export:   viper.GetString("cps-sans.export"),
		}

		cr.Report()
	},
}

func init() {
	param(cpssansCmd, "export", "cps-sans.export", "cps-sans.csv", "contract to be used")
	param(cpssansCmd, "contract", "cps-sans.contract", "", "contract to be used")
	RootCmd.AddCommand(cpssansCmd)
}
