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
	ReportName  string
	Export      string
	Group       string
}

func (pr PropReport) Report() {

	if pr.ReportName != "origin" {
		log.Fatalf("not (yet) supported report")
	}

	log.Infof("property report %v", pr.Export)

	f, err := os.Create(pr.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	groupResponse, err := pr.PropService.PapiClient.GetGroups(context.Background())
	if err != nil {
		log.Fatalf("get groups %w", err)
		return
	}

	r := []string{"group", "property", "origintype", "origin", "forward", "hostmatch", "pathmatch", "siteshield"}
	w.Write(r)
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

						origins, f := pb.Behaviors["origin"]
						if f {
							for _, o := range origins {
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
								w.Write([]string{grp.GroupName,
									x.PropertyName,
									otype,
									ohostname,
									ohostheader,
									hostmatch,
									pathmatch,
									siteshield,
								})
							}
						}
					}
				}
			}
		}
	}
}
