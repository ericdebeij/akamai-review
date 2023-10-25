/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/cpsreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cpscertificatesCmd represents the cpsList command
var cpscertificatesCmd = &cobra.Command{
	Use:     "cps-certificates",
	Aliases: []string{"cc", "cps-c"},
	Short:   "List certificates as defined in cps and the usage of the SANs",
	Long:    `List of the certificates, the SAN in the certificates. Additional information is provided to check whether the CN or SAN entry is actually served via Akamai`,
	Run: func(cmd *cobra.Command, args []string) {
		cr := &cpsreport.CertSanReport{
			Contract: viperAlias("cps-certificates", "contract"),
			Export:   viper.GetString("cps-certificates.export"),
		}

		cr.Report()
	},
}

func init() {
	param(cpscertificatesCmd, "export", "cps-certificates.export", "cps-certificates.csv", "contract to be used")
	param(cpscertificatesCmd, "contract", "cps-certificates.contract", "", "contract to be used")
	RootCmd.AddCommand(cpscertificatesCmd)
}
