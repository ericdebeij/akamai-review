package cmd

import (
	"strings"

	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ReportFields struct {
	Type     string
	Period   string
	Contract string
	Product  string
	Export   string
	Group    string
}

var defaultParam ReportFields

func ReportParameters(cmd *cobra.Command, param ...string) {
	for _, p := range param {
		switch p {
		case "period":
			cmd.Flags().StringVar(&defaultParam.Period, "period", "", "period, format YYYY-MM")
		case "contract":
			cmd.Flags().StringVar(&defaultParam.Contract, "contract", "", "contract identification")
		case "group":
			cmd.Flags().StringVar(&defaultParam.Group, "group", "", "group identification")
		case "export":
			defaultExport := cmd.Use + ".csv"
			if akutil.FindString(param, "period") >= 0 {
				defaultExport = cmd.Use + "-PERIOD.csv"
			}
			cmd.Flags().StringVar(&defaultParam.Export, "export", defaultExport, "export filename")
		}
	}
	//reportCmd.Flags().StringVar(&reportName, "name", "", "limit to specific report")
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Run reports",
	Long:  `This overall command is used to run multiple reports`,
	Run: func(cmd *cobra.Command, args []string) {
		runreport()
	},
}

var reportName string

func init() {
	reportCmd.Flags().StringVar(&reportName, "name", "", "limit to specific report")
	rootCmd.AddCommand(reportCmd)
}

func runreport() {
	reports := make(map[string]*ReportFields)
	viper.UnmarshalKey("reports", &reports)
	for repname, repdef := range reports {
		if reportName != "" && repname != reportName {
			continue
		}
		if strings.HasPrefix(repdef.Type, "usage-cpcode") {
			usagereport(repdef)
		}

		if strings.HasPrefix(repdef.Type, "properties-") {
			propreport(repdef)
		}
	}
}
