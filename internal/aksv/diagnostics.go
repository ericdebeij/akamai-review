package aksv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
)

type DiagnosticsService struct {
	Session  session.Session
	Response *http.Response
	CacheDir string
	IpCache  map[string]bool
	dirty    bool
}

func NewDiagnosticsService(s session.Session, cacheFile string) (ds *DiagnosticsService) {
	ds = &DiagnosticsService{
		Session: s,
	}

	ds.CacheDir = akutil.ExpandPath(cacheFile)
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
	ds.dirty = false
}

func (ds *DiagnosticsService) FlushCache() {
	if ds.dirty {
		byteblob, err := json.Marshal(ds.IpCache)
		if err != nil {
			fmt.Println(err)
		}
		os.MkdirAll(filepath.Dir(ds.CacheDir+"/diagiscdn.json"), 0750)
		err = os.WriteFile(ds.CacheDir+"/diagiscdn.json", byteblob, 0644)
		if err != nil {
			fmt.Println(err)
		}

		ds.dirty = false
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

func (ds *DiagnosticsService) IsAkamaiIp(ips []string) (ismap map[string]bool, is int, err error) {
	ismap = make(map[string]bool, len(ips))
	allfound := true
	for _, ip := range ips {
		b, f := ds.IpCache[ip]

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
			fmt.Println(err)
			return
		}

		if ds.Response.StatusCode == 200 {
			for _, rip := range rsb.Results {
				ds.IpCache[rip.IPAddress] = rip.IsEdgeIP
				ismap[rip.IPAddress] = rip.IsEdgeIP
				if rip.IsEdgeIP {
					is = is + 1
				}
				ds.dirty = true
			}
			return
		}
		if ds.Response.StatusCode != 429 {
			data, e := io.ReadAll(ds.Response.Body)
			if e != nil {
				fmt.Println(e)
			}
			fmt.Println(string(data))

			err = fmt.Errorf("status code: %d", ds.Response.StatusCode)
			return
		}
		fmt.Println("running into rate control, wait 60 seconds")
		ds.FlushCache()
		time.Sleep(time.Minute)
	}
}
