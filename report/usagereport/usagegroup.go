package usagereport

import (
	"fmt"

	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/exportx"
	"github.com/ericdebeij/akamai-review/v3/services"
)

type UsageRG struct {
	Contract       string
	Product        string
	ReportingGroup string
	FromMonth      string
	ToMonth        string
	Export         string
}

func (ur *UsageRG) GroupCsv() {
	//how to export the usage group as csv in golang?

}

func (ur *UsageRG) ByGroup() {
	srvs := services.Services

	csvx, errx := exportx.Create(ur.Export)
	if errx != nil {
		log.Fatal(errx.Error())
		return
	}
	defer csvx.Close()
	csvx.Header("month", "GB")

	if ur.Contract == "" || ur.Product == "" {
		log.Fatalf("Contract and Product are mandatory parameters for usage-cpcode")
	}

	x, err := srvs.AkamaiBilling.GetUsageCpcode(ur.Contract, ur.Product, ur.FromMonth, ur.ToMonth)
	if err != nil {
		log.Fatalf("getusagecpcode error: %w", err)
		return
	}

	rginfos, err := srvs.AkamaiCpcodes.GetRepgroups("")
	if err != nil {
		log.Fatalf("reportinggroup error: %w", err)
		return
	}
	cprgmap := rginfos.MapCpcodeRepgroup()
	fmt.Println(cprgmap)
	rg := rginfos.FindByName(ur.ReportingGroup)

	for _, p := range x.UsagePeriods {
		gb := 0.0

		for _, s := range p.CpCodeStats {
			if rg.FindCpcode(s.CpCode) {
				for _, v := range s.Stats {
					if v.StatType == "Bytes" {
						gb += v.Value
					}
				}
			}
		}
		csvx.Write(p.Month, gb)
	}
}
