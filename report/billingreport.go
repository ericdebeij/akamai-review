package report

import (
	"encoding/csv"
	"fmt"
	"strings"

	"os"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/yearmonth"
)

type UsageCpcode struct {
	EdgeSession    session.Session
	BillingService *aksv.BillingService
	CpCodeService  *aksv.CpcodeService

	Contract string
	Product  string
	Period   string
	Export   string
}

type r struct {
	prevBytes    float64
	currentBytes float64
	prevHits     float64
	currentHits  float64
}

func (ur UsageCpcode) Report() {

	tm := yearmonth.Add(ur.Period, 1)
	fm := yearmonth.Add(tm, -2)
	x, err := ur.BillingService.GetUsageCpcode(ur.Contract, ur.Product, fm, tm)
	sum := make(map[int]r, 5000)
	if err != nil {
		log.Fatalf("usage by cpcode: %w", err)
		return
	}
	for _, p := range x.UsagePeriods {
		for _, s := range p.CpCodeStats {
			for _, v := range s.Stats {
				a, f := sum[s.CpCode]
				if !f {
					a = r{}
				}
				if p.Month == ur.Period {
					if v.StatType == "Bytes" {
						a.currentBytes = v.Value
					}
					if v.StatType == "Hits" {
						a.currentHits = v.Value
					}
				}
				if p.Month == fm {
					if v.StatType == "Bytes" {
						a.prevBytes = v.Value
					}
					if v.StatType == "Hits" {
						a.prevHits = v.Value
					}
				}
				sum[s.CpCode] = a
			}
		}
	}

	cpinfos, err := ur.CpCodeService.GetCpcodes()
	if err != nil {
		log.Fatalf("cpcode error: %w", err)
		return
	}

	rginfos, err := ur.CpCodeService.GetRepgroups("")
	if err != nil {
		log.Fatalf("reportinggroup error: %w", err)
		return
	}

	f, err := os.Create(ur.Export)
	if err != nil {
		log.Fatalf("failed to open file: %w", err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	r := []string{"cpcode", "cpname", "repgrp", ur.Period + "(GB)", fm + "(GB)", "diff(GB)", ur.Period + "(Hits)", fm + "(Hits)", "diff(Hits)"}
	w.Write(r)

	for cpcode, vs := range sum {
		cpinfo := cpinfos.FindCpcode(cpcode)
		rginfo := rginfos.FindByCpcode(cpcode)
		rgroup := make([]string, 0, 3)
		for _, rg := range rginfo.Groups {
			rgroup = append(rgroup, rg.ReportingGroupName)
		}
		w.Write([]string{strconv.Itoa(cpcode), cpinfo.CpcodeName,
			strings.Join(rgroup, ","),
			fmt.Sprintf("%f", vs.currentBytes),
			fmt.Sprintf("%f", vs.prevBytes),
			fmt.Sprintf("%f", vs.currentBytes-vs.prevBytes),
			fmt.Sprintf("%f", vs.currentHits),
			fmt.Sprintf("%f", vs.prevHits),
			fmt.Sprintf("%f", vs.currentHits-vs.prevHits),
		})
	}
}
