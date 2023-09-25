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
	Property    string
	WarningDays int
	HttpTest    bool
}

type OriginReport struct {
	Export   string
	Group    string
	Property string
}

type BehaviorReport struct {
	Export   string
	Group    string
	Property string
	Behavior string
	Criteria bool
}

type PropertyInfo struct {
	Contract     string
	Groupname    string
	Propertyname string
	Siteshield   string
	Hosts        []*Hostinfo
	Origins      []*OriginInfo
	Summary      *properties.PropSum
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

func (br BehaviorReport) Report() {
	log.Infof("pm-behavior %+v", br)
	csvx, err := exportx.Create(br.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer csvx.Close()

	properties := Build(br.Group, br.Property)

	if br.Behavior == "" {
		csvx.Header("group", "property", "behaviors")
		for _, p := range properties {
			ba := make([]string, 0, 100)
			for k, b := range p.Summary.Behaviors {
				ki := fmt.Sprintf("%s(%d)", k, len(b))
				ba = append(ba, ki)
			}

			csvx.Write(p.Groupname, p.Propertyname, ba)
		}
	} else {

		if br.Criteria {
			csvx.Header("group", "property", br.Behavior, "Criteria")

		} else {
			csvx.Header("group", "property", br.Behavior)

		}

		for _, p := range properties {
			b := p.Summary.Behaviors[br.Behavior]
			if len(b) == 0 {
				csvx.Write(p.Groupname, p.Propertyname, "not used")
			}
			for _, bi := range b {
				if br.Criteria {
					csvx.Write(p.Groupname, p.Propertyname, bi.Behavior.Options, bi.Criteria)
				} else {
					csvx.Write(p.Groupname, p.Propertyname, bi.Behavior.Options)
				}
			}
		}
	}
}
func (or OriginReport) Report() {
	log.Infof("pm-origins %+v", or)
	csvx, err := exportx.Create(or.Export)
	if err != nil {
		log.Fatalf("failed to open file %w", err)
	}
	defer csvx.Close()

	properties := Build(or.Group, or.Property)

	csvx.Header("contract", "group", "property", "origin", "origintype", "forward", "hostmatch", "pathmatch", "siteshield", "ips")

	for _, p := range properties {
		for _, o := range p.Origins {
			csvx.Write(
				p.Contract,
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
	log.Infof("pm-hosts %+v", hr)
	srvs := services.Services
	csvx, err := exportx.Create(hr.Export)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvx.Close()

	properties := Build(hr.Group, hr.Property)

	if hr.HttpTest {
		csvx.Header("contract", "group", "property", "host", "edgehost", "cdn", "ips", "cert-subject", "cert-issuer", "cert-expire", "httptest")
	} else {
		csvx.Header("contract", "group", "property", "host", "edgehost", "cdn", "ips", "cert-subject", "cert-issuer", "cert-expire")
	}

	n := time.Now()
	for _, p := range properties {
		for _, h := range p.Hosts {
			ci := h.Clientinfo

			httptest := ""
			if hr.HttpTest && len(h.Clientinfo.Ips) > 0 {
				httptest = srvs.ClientTest.TestHttp("http://" + h.Hostname + "/")
			}

			csvx.Write(
				p.Contract,
				p.Groupname,
				p.Propertyname,
				h.Hostname,
				h.Edgehost,
				ci.Cdn,
				strings.Join(h.Clientinfo.Ips, " "),
				ci.Subject,
				ci.Issuer,
				ci.Expire,
				httptest,
			)

			if hr.WarningDays != 0 && ci.Err == "" && ci.Cdn == "akamai" && n.After(ci.Expire.AddDate(0, 0, 0-hr.WarningDays)) {
				log.Warnf("Host       : %v", ci.Hostname)
				log.Warnf("Expire date: %v", ci.Expire)
				log.Warnf("Subject    : %v", ci.Subject)
				log.Warnf("Issuer     : %v", ci.Issuer)
				diff := ci.Expire.Sub(n)
				dura := durafmt.Parse(diff)
				log.Warnf("Time left: %v", dura)
			}
		}

	}
}

func Build(group string, property string) (properties []*PropertyInfo) {
	log.Infof("buildup property info, filter: group: %v property: %v", group, property)
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
						if property != "" && x.PropertyName != property {
							continue
						}
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
						if pt == nil {
							log.Fatalf("No rule tree for %+v", prtrq)
						}

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
											x := critmatch.Options["values"]
											switch x := x.(type) {
											case string:
												pathmatch += " " + x
											case []interface{}:
												for _, tv := range x {
													pathmatch += " " + tv.(string)
												}
											default:
												log.Warnf("expected type is []interface or string, received type %T value %+v", x, x)
											}
										}
									}
								}
								pathmatch = strings.Trim(pathmatch, " ")

								ips, _, err := srvs.Dns.DnsInfo(ohostname)
								if err != nil {
									log.Errorf("dns %s: %v", ohostname, err)
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
							Contract:     strings.Replace(x.ContractID, "ctr_", "", 1),
							Groupname:    grp.GroupName,
							Propertyname: x.PropertyName,
							Siteshield:   siteshield,
							Origins:      origins,
							Hosts:        hosts,
							Summary:      pb,
						}
						properties = append(properties, propinfo)
					}
				}
			}
		}
	}
	return
}
