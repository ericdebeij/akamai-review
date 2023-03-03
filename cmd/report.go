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

var defaultReport ReportFields

func defaultValue(v ...string) string {
	for _, x := range v {
		if x != "" {
			return x
		}
	}
	return ""
}

func ReportParameters(cmd *cobra.Command, param ...string) {
	for _, p := range param {
		switch p {
		case "period":
			cmd.Flags().StringVar(&defaultReport.Period, "period", "", "period, format YYYY-MM")
		case "contract":
			cmd.Flags().StringVar(&defaultReport.Contract, "contract", "", "contract identification")
		case "group":
			cmd.Flags().StringVar(&defaultReport.Group, "group", "", "group identification")
		case "export":
			defaultExport := cmd.Use + ".csv"
			if akutil.FindString(param, "period") >= 0 {
				defaultExport = cmd.Use + "-PERIOD.csv"
			}
			cmd.Flags().StringVar(&defaultReport.Export, "export", defaultExport, "export filename")
		}
	}
	//reportCmd.Flags().StringVar(&reportName, "name", "", "limit to specific report")
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Run reports",
	Long:  `This overall command is used to run multiple reports`,
	Run: func(cmd *cobra.Command, args []string) {
		runreport("")
	},
}

var reportName string

func init() {
	reportCmd.Flags().StringVar(&reportName, "name", "", "limit to specific report")
	rootCmd.AddCommand(reportCmd)
}

func runareport(repdef *ReportFields) {
	if strings.HasPrefix(repdef.Type, "usage-cpcode") {
		usagereport(repdef)
	}

	if strings.HasPrefix(repdef.Type, "properties-") {
		propreport(repdef)
	}
}

func runreport(reportType string) (runned int) {
	reports := make(map[string]*ReportFields)
	viper.UnmarshalKey("reports", &reports)
	for repname, repdef := range reports {
		repdef.Period = defaultValue(repdef.Period, defaultReport.Period)
		repdef.Export = defaultValue(repdef.Export, defaultReport.Export)
		repdef.Contract = defaultValue(repdef.Contract, defaultReport.Contract)
		repdef.Group = defaultValue(repdef.Group, defaultReport.Group)

		if reportName != "" && repname != reportName {
			continue
		}
		if reportType != "" && repdef.Type != reportType {
			continue
		}

		runareport(repdef)
		runned += 1
	}
	if runned == 0 && reportType != "" {
		defaultReport.Type = reportType
		runareport(&defaultReport)
	}
	return
}
