package propreport

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/exportx"
	"github.com/ericdebeij/akamai-review/v3/service/clienttest"
	"github.com/ericdebeij/akamai-review/v3/service/properties"
	"github.com/ericdebeij/akamai-review/v3/services"
	"github.com/hako/durafmt"
)

type HostReport struct {
	Export      string
	Group       string
	WarningDays int
}

type OriginReport struct {
	Export string
	Group  string
}

type PropertyReport struct {
	Export string
	Group  string
}

type PropertyInfo struct {
	Groupname    string
	Propertyname string
	Siteshield   string
	Hosts        []*Hostinfo
	Origins      []*OriginInfo
	Behaviors    *properties.PropSum
}

type Hostinfo struct {
	Hostname   string
	Edgehost   string
	Clientinfo *clienttest.ClientInfo
}

type OriginInfo struct {
	Origin     string
	Hostheader string
	Type       string
	Hostmatch  string
	Pathmatch  string
	Ips        []string
}

func (or OriginReport) Report() {
	csvx, err := exportx.Create(or.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer csvx.Close()

	properties := Build(or.Group)

	csvx.Header("group", "property", "origin", "origintype", "forward", "hostmatch", "pathmatch", "siteshield", "ips")

	for _, p := range properties {
		for _, o := range p.Origins {
			csvx.Write(
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
			)
		}
	}
}
func (hr HostReport) Report() {
	csvx, err := exportx.Create(hr.Export)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvx.Close()

	properties := Build(hr.Group)

	csvx.Header("group", "property", "host", "edgehost", "cdn", "ips", "cert-subject", "cert-issuer", "cert-expire")

	n := time.Now()
	for _, p := range properties {
		for _, h := range p.Hosts {
			ci := h.Clientinfo
			csvx.Write(
				p.Groupname,
				p.Propertyname,
				h.Hostname,
				h.Edgehost,
				ci.Cdn,
				strings.Join(h.Clientinfo.Ips, " "),
				ci.Subject,
				ci.Issuer,
				ci.Expire,
			)

			if hr.WarningDays != 0 && ci.Err == "" && ci.Cdn == "akamai" && n.After(ci.Expire.AddDate(0, 0, 0-hr.WarningDays)) {
				fmt.Println("Host       :", ci.Hostname)
				fmt.Println("Expire date:", ci.Expire)
				fmt.Println("Subject    :", ci.Subject)
				fmt.Println("Issuer     :", ci.Issuer)
				diff := ci.Expire.Sub(n)
				dura := durafmt.Parse(diff)
				fmt.Println("Time left:", dura)
			}
		}

	}
}

func Build(group string) (properties []*PropertyInfo) {
	srvs := services.Services
	properties = make([]*PropertyInfo, 0, 1000)
	groupResponse, err := srvs.Properties.PapiClient.GetGroups(context.Background())
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
			if group == "" || group == grp.GroupName {

				pl, err2 := srvs.Properties.GetProperties(context.Background(), plrq)
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
						hl, _ := srvs.Properties.GetPropertyVersionHostnames(prhrq)

						hll := len(hl.Hostnames.Items)
						hostnames := make([]string, hll, hll)
						hosts := make([]*Hostinfo, 0, 10)
						for hii, hiv := range hl.Hostnames.Items {
							hostnames[hii] = hiv.CnameFrom

							hostinfo := &Hostinfo{
								Hostname:   hiv.CnameFrom,
								Edgehost:   hiv.CnameTo,
								Clientinfo: srvs.ClientTest.Testhost(hiv.CnameFrom),
							}

							hosts = append(hosts, hostinfo)
						}

						prtrq := papi.GetRuleTreeRequest{
							PropertyID:      x.PropertyID,
							PropertyVersion: pv,
							ContractID:      x.ContractID,
							GroupID:         x.GroupID,
							ValidateRules:   false,
						}
						pt := srvs.Properties.GetRuleTree(prtrq)

						pb := srvs.Properties.FindBehaviors(&pt.Rules)

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

								ips, _, err := srvs.Dns.DnsInfo(ohostname)
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
							Hosts:        hosts,
							Behaviors:    pb,
						}
						properties = append(properties, propinfo)
					}
				}
			}
		}
	}
	return
}
