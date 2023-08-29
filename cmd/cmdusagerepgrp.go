/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/ericdebeij/akamai-review/v3/conv/yearmonth"
	"github.com/ericdebeij/akamai-review/v3/report/usagereport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// originlistCmd represents the originlist command
var usageRepgrpCmd = &cobra.Command{
	Use:     "usage-repgroup",
	Aliases: []string{"urg"},
	Short:   "An overview of the usage for a month per cpcode and a comparison with the previous month",
	Long:    `Uses the billing API to get an overview of the usage for a specific month and compares this with the previous month, both bytes and hits`,
	Run: func(cmd *cobra.Command, args []string) {
		um := &usagereport.UsageRG{
			Contract:        viperAlias("usage-repgroup", "contract"),
			Product:         viperAlias("usage-repgroup", "product"),
			FromMonth:       viper.GetString("usage-repgroup.from"),
			ToMonth:         viper.GetString("usage-repgroup.to"),
			Export:          viper.GetString("usage-repgroup.export"),
			ReportingGroups: viper.GetStringSlice("usage-repgroup.repgroups"),
			StatType:        viper.GetString("usage-repgroup.type"),
			Unit:            viper.GetString("usage-report.unit"),
		}
		um.Report()
	},
}

func init() {
	param(usageRepgrpCmd, "export", "usage-repgroup.export", "usage-repgroup.csv", "name of the exportfile")
	param(usageRepgrpCmd, "contract", "usage-repgroup.contract", "", "contract to be used")
	param(usageRepgrpCmd, "from", "usage-repgroup.from", yearmonth.Add(yearmonth.FromTime(time.Now()), -1), "from month (format YYYY-MM)")
	param(usageRepgrpCmd, "to", "usage-repgroup.to", yearmonth.Add(yearmonth.FromTime(time.Now()), -1), "to month (format YYYY-MM)")
	param(usageRepgrpCmd, "product", "usage-repgroup.product", "", "product code to be used")
	param(usageRepgrpCmd, "rgroup", "usage-repgroup.repgroups", []string{}, "reporting groups (default all)")
	param(usageRepgrpCmd, "type", "usage-reproup.type", "Bytes", "Statistic type [Bytes|Hits]")
	param(usageRepgrpCmd, "unit", "usage-repgroup.unit", "", "Unit (default Bytes:GB, Hits:Hit)")
	RootCmd.AddCommand(usageRepgrpCmd)
}
