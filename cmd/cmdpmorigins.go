/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/propreport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// originlistCmd represents the originlist command
var pmoriginsCmd = &cobra.Command{
	Use:   "pm-origins",
	Short: "An overview of the origins",
	Long: `An overview of the properties and the origins used within the property. In order to find this information the property manager rules are downloaded (and stored in a cache).
At a high level the property match criteria are extracted from the config file. (Only a limited number of conditions are shown). 
For an origin the type of origin (web,ns), forward host header and ip-address is shown (not possible for variable origins)`,
	Run: func(cmd *cobra.Command, args []string) {
		or := &propreport.OriginReport{
			Export: viper.GetString("pm-origins.export"),
			Group:  viperAlias("pm-origins", "group"),
		}
		or.Report()
	},
}

func init() {
	param(pmoriginsCmd, "export", "pm-origins.export", "pm-origins.csv", "name of the exportfile")
	param(pmoriginsCmd, "group", "pm-origins.group", "", "filter for the group")
	rootCmd.AddCommand(pmoriginsCmd)
}
