/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/propreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pmhostsCmd represents the hostlist command
var pmhostsCmd = &cobra.Command{
	Use:   "pm-hosts",
	Short: "List of all hostnames in your account per property with dns and certificate information",
	Long: `An overview of the properties and the hostnames associated within the property. In order to find this information the property manager hostnames are downloaded (and stored in a cache).
The related edgehost is shown and the host is checked to see if it is actually served by Akamai, resolves in a proper IP-address and information regarding the certificate being used`,
	Run: func(cmd *cobra.Command, args []string) {
		hr := &propreport.HostReport{
			Export:      viper.GetString("properties.hosts"),
			Group:       viper.GetString("group"),
			WarningDays: viper.GetInt("warningdays"),
		}
		hr.Report()
	},
}

func init() {
	param(pmhostsCmd, "export", "properties.hosts", "pmhosts.csv", "name of the exportfile")
	param(pmhostsCmd, "group", "group", "", "filter for the group")
	rootCmd.AddCommand(pmhostsCmd)
}
