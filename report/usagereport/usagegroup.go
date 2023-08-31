package usagereport

import (
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/conv/yearmonth"
	"github.com/ericdebeij/akamai-review/v3/exportx"
	"github.com/ericdebeij/akamai-review/v3/services"
	"golang.org/x/exp/slices"
)

type UsageRG struct {
	Contract        string
	Product         string
	ReportingGroups []string
	FromMonth       string
	ToMonth         string
	Export          string
	StatType        string
	Unit            string
}

func (ur *UsageRG) GroupCsv() {
	//how to export the usage group as csv in golang?

}

func (ur *UsageRG) Report() {
	srvs := services.Services

	// Mandatory parameters
	if ur.Contract == "" || ur.Product == "" || ur.FromMonth == "" || ur.ToMonth == "" {
		log.Fatalf("Contract, Product, FromMonth, ToMonth are mandatory parameters for usage-repgroup")
	}
	if ur.StatType == "" {
		ur.StatType = "Bytes"
	}
	if ur.StatType == "Bytes" && ur.Unit == "" {
		ur.Unit = "GB"
	}

	tomonth := yearmonth.Add(ur.ToMonth, +1)

	// yes, it is weird, the unit return from the bulling API for Bytes is MB, but we can do the calculation
	factor := 1.0
	if ur.StatType == "Bytes" {
		switch ur.Unit {
		case "B":
			factor = 1e9
		case "KB":
			factor = 1e6
		case "MB":
			factor = 1e3
		case "GB":

		case "TB":
			factor = 1e-3
		case "PB":
			factor = 1e-6
		default:
			log.Errorf("usage-repgroup: unknown unit for Bytes: %v", ur.Unit)
		}
	}

	log.Infof("usage-repgroup tomonth %v, factor %v, %+v", tomonth, factor, ur)

	x, err := srvs.AkamaiBilling.GetUsageCpcode(ur.Contract, ur.Product, ur.FromMonth, tomonth)
	if err != nil {
		log.Fatalf("getusagecpcode error: %w", err)
		return
	}

	rginfos, err := srvs.AkamaiCpcodes.GetRepgroups("")
	if err != nil {
		log.Fatalf("reportinggroup error: %w", err)
		return
	}

	// By default use all reportinggroups
	reportingGroups := ur.ReportingGroups
	if len(reportingGroups) == 0 {
		for _, rg := range rginfos.Groups {
			for _, contract := range rg.Contracts {
				if contract.ContractID == ur.Contract {
					reportingGroups = append(reportingGroups, rg.ReportingGroupName)
				}
			}
		}
	}
	log.Infof("reporting groups selected %v", reportingGroups)
	numGroups := len(reportingGroups)

	cpcodes, errc := srvs.AkamaiCpcodes.GetCpcodes()
	if errc != nil {
		log.Fatalf("usage-repgroup: get cpcode %v", errc)
	}

	// Prepare CSV file export
	csvx, errx := exportx.Create(ur.Export)
	if errx != nil {
		log.Fatal(errx.Error())
		return
	}

	defer csvx.Close()
	hdrs := []string{"month"}
	hdrs = append(hdrs, reportingGroups...)
	hdrs = append(hdrs, "other")
	csvx.Header(hdrs...)

	// Create a map with a column indication for every CPCode
	cpcode2colmap := make(map[int]int, 500)
	for rgi, rgn := range reportingGroups {
		rg := rginfos.FindByName(rgn)
		if rg == nil {
			log.Errorf("reporting group %s not available", rgn)
			continue
		}

		for _, rgc := range rg.Contracts {
			for _, cpv := range rgc.Cpcodes {
				colid, found := cpcode2colmap[cpv.CpcodeID]
				if found {
					log.Warnf("cpcode %v (%v) found in multiple reporting groups (%s, %s), first group used", cpv.CpcodeID, cpv.CpcodeName, reportingGroups[colid], rg.ReportingGroupName)
				} else {
					cpcode2colmap[cpv.CpcodeID] = rgi
				}
			}
		}
	}

	var restcps []int

	// Now run thru the data and find the Bytes
	for _, p := range x.UsagePeriods {
		values := make([]float64, numGroups+1)
		for _, s := range p.CpCodeStats {
			col, found := cpcode2colmap[s.CpCode]
			if !found {
				col = numGroups
				if !slices.Contains(restcps, s.CpCode) {
					cpcode := cpcodes.FindCpcode(s.CpCode)
					log.Warnf("cpcode %v (%v) used without reportinggroup", s.CpCode, cpcode.CpcodeName)
				}
			}
			for _, v := range s.Stats {
				if v.StatType == ur.StatType {
					values[col] += v.Value * factor
				}
			}
		}

		vs := []interface{}{p.Month}
		for _, v := range values {
			vs = append(vs, v)
		}
		csvx.Write(vs...)
	}
}
