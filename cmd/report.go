package cmd

import (
	"strings"

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
		case "product":
			cmd.Flags().StringVar(&defaultReport.Product, "product", "", "product identification")
		case "group":
			cmd.Flags().StringVar(&defaultReport.Group, "group", "", "group identification")
		case "export":
			cmd.Flags().StringVar(&defaultReport.Export, "export", "", "export filename")
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
		repdef.Export = defaultValue(repdef.Export, defaultReport.Export, repname+".csv")
		repdef.Contract = defaultValue(repdef.Contract, defaultReport.Contract)
		repdef.Product = defaultValue(repdef.Product, defaultReport.Product)
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
		defaultReport.Export = defaultValue(defaultReport.Export, reportType+".csv")
		runareport(&defaultReport)
	}
	return
}
