package diagnostics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/util/osutil"
)

type DiagnosticsService struct {
	Session  session.Session
	Response *http.Response
	CacheDir string
	IpCache  map[string]bool
	dirty    int
}

func NewDiagnosticsService(s session.Session, cacheFile string) (ds *DiagnosticsService) {
	ds = &DiagnosticsService{
		Session: s,
	}

	ds.CacheDir = osutil.ExpandPath(cacheFile)
	ds.ReadCache()

	return
}

type IsAkamaiIpResponse struct {
	IsCdnIP bool `json:"isCdnIp"`
}

func (ds *DiagnosticsService) ReadCache() {
	jsonFile, err := os.Open(ds.CacheDir + "/diagiscdn.json")

	if err == nil {
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &ds.IpCache)
		log.Debug(fmt.Sprint("IP Cache loaded:", len(ds.IpCache)))
	} else {
		ds.IpCache = make(map[string]bool, 100)
	}
	ds.dirty = 0
}

func (ds *DiagnosticsService) FlushCache() {

	if ds.dirty > 0 {
		byteblob, err := json.Marshal(ds.IpCache)
		if err != nil {
			fmt.Println(err)
		}
		os.MkdirAll(filepath.Dir(ds.CacheDir+"/diagiscdn.json"), 0750)
		err = os.WriteFile(ds.CacheDir+"/diagiscdn.json", byteblob, 0644)
		if err != nil {
			fmt.Println(err)
		}
		ds.dirty = 0
	}
}

type RequestVerifyIP struct {
	IPAddresses []string `json:"ipAddresses"`
}

type ResponseVerifyIP struct {
	Request struct {
		IPAddresses []string `json:"ipAddresses"`
	} `json:"request"`
	CreatedBy       string `json:"createdBy"`
	CreatedTime     string `json:"createdTime"`
	CompletedTime   string `json:"completedTime"`
	ExecutionStatus string `json:"executionStatus"`
	Results         []struct {
		ExecutionStatus string `json:"executionStatus"`
		IPAddress       string `json:"ipAddress"`
		IsEdgeIP        bool   `json:"isEdgeIp"`
	} `json:"results"`
}

func ipCidr(theIp string) string {
	// Replace this with your IPv4 address
	ip := net.ParseIP(theIp)
	cidrPrefixLength := 24
	cidrBits := 32
	if ip.To4() == nil {
		cidrPrefixLength = 64
		cidrBits = 128
	}
	// Get the network address by masking the IP with CIDR mask
	networkAddress := ip.Mask(net.CIDRMask(cidrPrefixLength, cidrBits))

	// Create an IPNet structure representing the CIDR range
	ipNet := &net.IPNet{
		IP:   networkAddress,
		Mask: net.CIDRMask(cidrPrefixLength, cidrBits),
	}
	return ipNet.String()
}

func (ds *DiagnosticsService) IsAkamaiIp(ips []string) (ismap map[string]bool, is int, err error) {
	ismap = make(map[string]bool, len(ips))
	allfound := true
	for _, ip := range ips {
		cidr := ipCidr(ip)
		b, f := ds.IpCache[cidr]

		if f {
			ismap[ip] = b
			if b {
				is = is + 1
			}
		} else {
			allfound = false
			break
		}
	}
	if allfound {
		return
	}
	is = 0
	//"/diagnostic-tools/v2/ip-addresses/{ipAddress}/is-cdn-ip"
	url := "/edge-diagnostics/v1/verify-edge-ip"

	rqb := RequestVerifyIP{
		IPAddresses: ips,
	}

	rsb := ResponseVerifyIP{}

	attempt := 0
	for {
		body, e := json.Marshal(rqb)
		if e != nil {
			err = e
			return
		}
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-type", "application/json")
		ds.Response, err = ds.Session.Exec(req, &rsb)

		if err != nil {
			log.Errorf("isakamaiip - request, %v, %+v", err, ips)
			return
		}

		if ds.Response.StatusCode == 200 {
			for _, rip := range rsb.Results {
				cidr := ipCidr(rip.IPAddress)
				ds.IpCache[cidr] = rip.IsEdgeIP
				ismap[cidr] = rip.IsEdgeIP
				if rip.IsEdgeIP {
					is = is + 1
				}
				ds.dirty += 1
			}
			if ds.dirty > 5 {
				ds.FlushCache()
			}
			return
		}
		if ds.Response.StatusCode != 429 {
			data, e := io.ReadAll(ds.Response.Body)
			if e != nil {
				log.Debugf("isakamaiip, status %v, data error %v", ds.Response.StatusCode, e)
			}
			log.Debugf("isakamaiip, status %v, data %v", ds.Response.StatusCode, string(data))

			err = fmt.Errorf("status code: %d", ds.Response.StatusCode)
			attempt += 1
			if attempt >= 2 {
				return
			}
		}
		log.Infof("running into rate control (%v), wait 60 seconds, (%v)", ds.Response.StatusCode, ips)
		ds.FlushCache()
		time.Sleep(time.Minute)
	}
}
