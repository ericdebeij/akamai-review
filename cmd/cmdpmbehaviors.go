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
var pmbehaviorsCmd = &cobra.Command{
	Use:     "pm-behaviors",
	Aliases: []string{"pmb"},
	Short:   "An overview of the behaviors in a propery",
	Long:    `An overview of the properties and the behaviors implemented or details about a specific behavior`,
	Run: func(cmd *cobra.Command, args []string) {
		or := &propreport.BehaviorReport{
			Export:   viper.GetString("pm-behaviors.export"),
			Group:    viperAlias("pm-behaviors", "group"),
			Property: viper.GetString("pm-behaviors.property"),
			Behavior: viper.GetString("pm-behaviors.behavior"),
		}
		or.Report()
	},
}

func init() {
	param(pmbehaviorsCmd, "export", "pm-behaviors.export", "pm-behaviors.csv", "name of the exportfile")
	param(pmbehaviorsCmd, "group", "pm-behaviors.group", "", "filter for the group")
	param(pmbehaviorsCmd, "property", "pm-behaviors.property", "", "filter for the property")
	param(pmbehaviorsCmd, "behavior", "pm-behaviors.behavior", "", "provide details about the specific behavior")
	RootCmd.AddCommand(pmbehaviorsCmd)
}
