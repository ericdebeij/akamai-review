/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ericdebeij/akamai-review/v3/report/usagereport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// originlistCmd represents the originlist command
var usageCmd = &cobra.Command{
	Use:   "usage-cpcode",
	Short: "An overview of the usage for a month per cpcode and a comparison with the previous month",
	Long:  `Uses the billing API to get an overview of the usage for a specific month and compares this with the previous month, both bytes and hits`,
	Run: func(cmd *cobra.Command, args []string) {
		um := &usagereport.UsageCpcode{
			Contract: viperAlias("usage-cpcode", "contract"),
			Product:  viperAlias("usage-cpcode", "product"),
			Period:   viper.GetString("usage-cpcode.period"),
			Export:   viper.GetString("usage-cpcode.export"),
		}
		um.Report()
	},
}

func init() {
	param(usageCmd, "export", "usage-cpcode.export", "usage_PERIOD.csv", "name of the exportfile")
	param(usageCmd, "contract", "usage-cpcode.contract", "", "contract to be used")
	param(usageCmd, "period", "usage-cpcode.period", "", "period to be investigated")
	param(usageCmd, "product", "usage-cpcode.product", "", "product code to be used")
	rootCmd.AddCommand(usageCmd)
}
