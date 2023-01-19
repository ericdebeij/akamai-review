package report

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
)

type AlbReport struct {
	EdgeSession session.Session
	AlbService  *aksv.AlbService
	DnsService  *akutil.Dns
	ExportDC    string
}

type originUsage struct {
	policies []string
}

func (ar AlbReport) Report() {
	var w *csv.Writer
	var f *os.File
	var err error

	origins := make(map[string]originUsage)

	pl, err := ar.AlbService.ListPolicies()
	if err != nil {
		log.Fatalln("failed to retrieve policies for ALB", err)
	}
	for _, pol := range *pl {
		if pol.Deleted {
			continue
		}
		vs := []int{}
		for _, a := range pol.Activations {
			if akutil.FindInt(vs, a.PolicyInfo.Version) >= 0 {
				pv, err := ar.AlbService.PolicyVersion(a.PolicyInfo.PolicyID, a.PolicyInfo.Version)
				if err != nil {
					log.Fatalln("failed to retrieve policyversion for ALB", err)
				}
				for _, pm := range pv.MatchRules {
					oid := pm.ForwardSettings.OriginID
					x, found := origins[oid]
					if !found {
						x = originUsage{}
					}
					if akutil.FindString(x.policies, pol.Name) < 0 {
						x.policies = append(x.policies, pol.Name)
						origins[oid] = x
					}
				}
			}
			vs = append(vs, a.PolicyInfo.Version)
		}

	}

	if ar.ExportDC != "" {
		f, err = os.Create(ar.ExportDC)
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		defer f.Close()

		w = csv.NewWriter(f)
		defer w.Flush()

		r := []string{"albid", "staging", "production", "hostname"}
		w.Write(r)
	}

	x, err := ar.AlbService.ListActivation()
	if err != nil {
		fmt.Println(err)
	}
	for ak, av := range *x {
		var p, s int
		vs := []int{}
		ac, found := av["PRODUCTION"]
		if found {
			p = ac.Version
			vs = append(vs, ac.Version)
		}

		ac, found = av["STAGING"]
		if found {
			s = ac.Version
			if akutil.FindInt(vs, s) < 0 {
				vs = append(vs, ac.Version)
			}
		}

		fmt.Println("ALB Definition:", ak)
		for _, v := range vs {
			activeon := ""
			if v == s {
				activeon = "Staging,"
			}
			if v == p {
				activeon += " Production"

			}
			activeon = strings.Trim(activeon, " ,")
			fmt.Println("- Version:", v, "(active on:", activeon+")")

			vd, err := ar.AlbService.ListAlbVersionDetails(ak, v)
			if err != nil {
				fmt.Println(err)
			}
			ipf := false
			fmt.Println("  DataCenters:")
			for _, dc := range vd.DataCenters {
				ips, _, err := ar.DnsService.DnsInfo(dc.Hostname)
				ipinfo := ""
				if err != nil {
					fmt.Println(err)
				}
				if len(ips) == 0 {
					ipinfo = "no-ip"
				} else {
					ipinfo = "ips:" + strings.Join(ips, ",")
					ipf = true
				}
				ls := "no-liveness"
				if vd.LivenessSettings != nil {
					ls = "liveness:" + strconv.Itoa(vd.LivenessSettings.Interval) + "s"
				}
				fmt.Println("  -", dc.Hostname, "("+ls+", "+ipinfo+")")
				if ar.ExportDC != "" {
					w.Write([]string{ak, strconv.Itoa(p), strconv.Itoa(s), dc.Hostname})
				}
			}
			if !ipf {
				fmt.Println("  ** No valid datacenters **")
			}

			fmt.Println("  Policies:")
			op, of := origins[ak]
			if of {
				for _, p := range op.policies {
					fmt.Println("  -", p)
				}
			} else {
				fmt.Println("  ** No active policies **")
			}

		}
	}
}
