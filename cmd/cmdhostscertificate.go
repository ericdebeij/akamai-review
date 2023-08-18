/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/securityreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pmhostsCmd represents the hostlist command
var hostscertCmd = &cobra.Command{
	Use:   "hosts-certificate",
	Short: "List of all hostnames in your account per property with dns and certificate information",
	Long: `An overview of the properties and the hostnames associated within the property. In order to find this information the property manager hostnames are downloaded (and stored in a cache).
The related edgehost is shown and the host is checked to see if it is actually served by Akamai, resolves in a proper IP-address and information regarding the certificate being used`,
	Run: func(cmd *cobra.Command, args []string) {
		hr := &securityreport.SecHostReport{
			Export:      viper.GetString("hosts-certificate.export"),
			WarningDays: viper.GetInt("warningdays"),
			Match:       viper.GetString("hosts-certificate.hostmatch"),
			Skip:        viper.GetString("hosts-certificate.hostskip"),
			HttpTest:    viper.GetBool("hosts-certificate.httptest"),
		}
		hr.Report()
	},
}

func init() {
	param(hostscertCmd, "export", "hosts-certificate.export", "hosts-certificate.csv", "name of the exportfile")
	param(hostscertCmd, "group", "hosts-certificate.group", "", "filter for the group")
	param(hostscertCmd, "match", "hosts-certificate.hostmatch", "", "regular expression for hostmatch")
	param(hostscertCmd, "skip", "hosts-certificate.hostskip", "", "regular expression for hostskip")
	param(hostscertCmd, "httptest", "hosts-certificate.httptest", false, "run an http test to check if http->https redirect is implemented")
	rootCmd.AddCommand(hostscertCmd)
}
