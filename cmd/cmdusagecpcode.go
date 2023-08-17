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
		um := &usagereport.UsageMonth{
			Contract: viper.GetString("default.contract"),
			Product:  viper.GetString("default.product"),
			Period:   viper.GetString("Period"),
			Export:   viper.GetString("usage.month"),
		}
		um.Report()
	},
}

func init() {
	param(usageCmd, "export", "usage.month", "usage_PERIOD.csv", "name of the exportfile")
	param(usageCmd, "contract", "default.contract", "", "contract to be used")
	param(usageCmd, "period", "period", "", "period to be investigated")
	param(usageCmd, "product", "default.product", "", "product code to be used")
	rootCmd.AddCommand(usageCmd)
}
