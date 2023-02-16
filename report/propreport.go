package report

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
)

type PropReport struct {
	EdgeSession session.Session
	DnsService  *akutil.Dns
	DiagService *aksv.DiagnosticsService
	PropService *aksv.Propsv
	ReportType  string
	Export      string
	Group       string
	LoadRules   bool
	LoadHosts   bool
}

type PropertyInfo struct {
	Groupname    string
	Propertyname string
	Siteshield   string
	Hosts        []*Hostinfo
	Origins      []*OriginInfo
}

type Hostinfo struct {
	Hostname string
}

type OriginInfo struct {
	Origin     string
	Hostheader string
	Type       string
	Hostmatch  string
	Pathmatch  string
	Ips        []string
}

func (pr PropReport) Report() {

	if pr.ReportType != "origin" && pr.ReportType != "host" {
		log.Fatalf("not (yet) supported report", pr.ReportType)
	}

	log.Infof("property report %v", pr.Export)

	f, err := os.Create(pr.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	r := []string{"group", "property", "origin", "origintype", "forward", "hostmatch", "pathmatch", "siteshield", "ips"}
	w.Write(r)

	if pr.ReportType == "origin" {
		pr.LoadRules = true
		properties := pr.Build()
		for _, p := range properties {
			for _, o := range p.Origins {
				w.Write([]string{
					p.Groupname,
					p.Propertyname,
					//strings.Join(hostnames, " "),
					o.Origin,
					o.Type,
					o.Hostheader,
					o.Hostmatch,
					o.Pathmatch,
					p.Siteshield,
					strings.Join(o.Ips, " "),
				})
			}
		}
	}
}

func (pr PropReport) Build() (properties []*PropertyInfo) {
	properties = make([]*PropertyInfo, 0, 1000)
	groupResponse, err := pr.PropService.PapiClient.GetGroups(context.Background())
	if err != nil {
		log.Fatalf("get groups %w", err)
		return
	}

	for _, grp := range groupResponse.Groups.Items {
		for _, contractId := range grp.ContractIDs {
			plrq := papi.GetPropertiesRequest{
				ContractID: contractId,
				GroupID:    grp.GroupID,
			}
			if pr.Group == "" || pr.Group == grp.GroupName {

				pl, err2 := pr.PropService.GetProperties(context.Background(), plrq)
				if err2 != nil {
					log.Errorf("get properties for group %s - %w", grp, err2)
					continue
				}

				for _, x := range pl.Properties.Items {
					pv := 0
					if x.ProductionVersion != nil {

						pv = *x.ProductionVersion

						prhrq := papi.GetPropertyVersionHostnamesRequest{
							PropertyID:        x.PropertyID,
							PropertyVersion:   pv,
							ContractID:        x.ContractID,
							GroupID:           x.GroupID,
							ValidateHostnames: false,
							IncludeCertStatus: true,
						}
						hl, _ := pr.PropService.GetPropertyVersionHostnames(prhrq)

						hll := len(hl.Hostnames.Items)
						hostnames := make([]string, hll, hll)
						for hii, hiv := range hl.Hostnames.Items {
							hostnames[hii] = hiv.CnameFrom
						}

						prtrq := papi.GetRuleTreeRequest{
							PropertyID:      x.PropertyID,
							PropertyVersion: pv,
							ContractID:      x.ContractID,
							GroupID:         x.GroupID,
							ValidateRules:   false,
						}
						pt := pr.PropService.GetRuleTree(prtrq)

						pb := pr.PropService.FindBehaviors(&pt.Rules)

						siteshields, f := pb.Behaviors["siteShield"]
						siteshield := ""
						if f {
							for _, ss := range siteshields {
								siteshield += " " + ss.Behavior.Options["ssmap"].(map[string]interface{})["value"].(string)
							}
						}

						pmorigins, f := pb.Behaviors["origin"]
						var origins []*OriginInfo
						if f {
							origins = make([]*OriginInfo, len(pmorigins))
							for oi, o := range pmorigins {

								otype := o.Behavior.Options["originType"].(string)
								ohostname := ""
								ohostheader := ""
								if otype == "NET_STORAGE" {
									n := o.Behavior.Options["netStorage"].(map[string]interface{})["cpCode"].(float64)
									ohostheader = fmt.Sprintf("cpcode:%v", n)
									otype = "ns"
									ohostname = fmt.Sprintf("%v", o.Behavior.Options["netStorage"].(map[string]interface{})["downloadDomainName"])
								} else {
									if otype == "CUSTOMER" {
										ohostname = fmt.Sprint(o.Behavior.Options["hostname"])
										otype = "web"
										ohostheader = fmt.Sprint(o.Behavior.Options["forwardHostHeader"])
										if ohostheader == "CUSTOM" {
											ohostheader = fmt.Sprint(o.Behavior.Options["customForwardHostHeader"])
										}
									}
								}
								hostmatch := ""
								pathmatch := ""
								for _, critlist := range o.Criteria {
									for _, critmatch := range critlist {
										if critmatch.Name == "hostname" {
											hostmatch += fmt.Sprint(critmatch.Options["values"])
											//hostmatch += " " + strings.Join(critmatch.Options["values"].([]string), ",")
										}

										if critmatch.Name == "path" {
											tt := critmatch.Options["values"].([]interface{})
											for _, tv := range tt {
												pathmatch += " " + tv.(string)
											}
										}
									}
								}
								pathmatch = strings.Trim(pathmatch, " ")

								ips, _, err := pr.DnsService.DnsInfo(ohostname)
								if err != nil {
									log.Infof("dns %s: %w")
								}

								origin := &OriginInfo{
									Origin:     ohostname,
									Hostheader: ohostheader,
									Type:       otype,
									Hostmatch:  hostmatch,
									Pathmatch:  pathmatch,
									Ips:        ips,
								}
								origins[oi] = origin

							}
						}
						propinfo := &PropertyInfo{
							Groupname:    grp.GroupName,
							Propertyname: x.PropertyName,
							Siteshield:   siteshield,
							Origins:      origins,
						}
						properties = append(properties, propinfo)
					}
				}
			}
		}
	}

	return
}
