/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ericdebeij/akamai-review/v3/report/certreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cpscertificatesCmd represents the cpsList command
var cpscertificatesCmd = &cobra.Command{
	Use:   "cpscertificates",
	Short: "List certificates as defined in cps",
	Long:  `List of the certificates, the SAN in the certificates. Additional information is provided to check whether the CN or SAN entry is actually served via Akamai`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("contract"))
		cr := &certreport.CertReport{
			Contract: viper.GetString("default.contract"),
			Export:   viper.GetString("cps.certificates"),
		}

		cr.Report()
	},
}

func init() {
	param(cpscertificatesCmd, "export", "cps.certificates", "cpscertificates.csv", "contract to be used")
	param(cpscertificatesCmd, "contract", "default.contract", "", "contract to be used")
	rootCmd.AddCommand(cpscertificatesCmd)
}
