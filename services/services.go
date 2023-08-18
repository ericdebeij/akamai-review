package services

import (
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/apex/log"
	svbilling "github.com/ericdebeij/akamai-review/v3/service/billing"
	"github.com/ericdebeij/akamai-review/v3/service/clienttest"
	"github.com/ericdebeij/akamai-review/v3/service/cpcodes"
	"github.com/ericdebeij/akamai-review/v3/service/diagnostics"
	"github.com/ericdebeij/akamai-review/v3/service/properties"
	"github.com/ericdebeij/akamai-review/v3/service/svsession"
	"github.com/ericdebeij/akamai-review/v3/util/dnsutil"
	"github.com/spf13/viper"
)

// Yes I know, this should be all references to interface....
type TheServices struct {
	AkamaiSession     session.Session
	AkamaiDiagnostics *diagnostics.DiagnosticsService
	AkamaiCps         cps.CPS
	AkamaiBilling     *svbilling.BillingService
	AkamaiCpcodes     *cpcodes.CpcodeService
	ClientTest        *clienttest.ClientTester
	Dns               *dnsutil.Dns
	PapiClient        papi.PAPI
	Properties        *properties.Propsv
	SecClient         appsec.APPSEC
}

var Services = TheServices{}

func pstr(p string) (v string) {
	v = viper.GetString(p)
	if v == "" {
		v = Parameters[p].Default.(string)
	}
	log.Debugf("Param: %s, Value %v", p, v)
	return
}

func StartServices() {
	var akamaiConfig svsession.EdgeConfig
	akamaiConfig.Edgerc = pstr("akamai.edgerc")
	akamaiConfig.Section = pstr("akamai.section")
	akamaiConfig.AccountID = pstr("akamai.accountkey")
	akamaiConfig.Logger = log.Log

	sess, err := svsession.NewSession(&akamaiConfig)
	if err != nil {
		log.Errorf("session error %v", err)
		os.Exit(1)
	}
	Services.AkamaiSession = sess

	akdiag := diagnostics.NewDiagnosticsService(sess, pstr("akamai.cache"))
	Services.AkamaiDiagnostics = akdiag

	dns := dnsutil.NewDnsService(pstr("dns.resolver"))
	Services.Dns = dns

	tst := &clienttest.ClientTester{
		EdgeSession: sess,
		DnsService:  dns,
		DiagService: akdiag,
	}
	Services.ClientTest = tst

	Services.AkamaiCps = cps.Client(sess)

	Services.AkamaiBilling = svbilling.NewBillingService(sess)
	Services.AkamaiCpcodes = cpcodes.NewCpcodeService(sess)

	papiClient := papi.Client(sess)
	Services.PapiClient = papiClient

	propsv := properties.NewPropertyService(papiClient, pstr("akamai.cache"))
	Services.Properties = propsv

	secClient := appsec.Client(sess)
	Services.SecClient = secClient
}

type TheParameter struct {
	Flag    string
	Viber   string
	Default interface{}
	Help    string
}

var Parameters = make(map[string]TheParameter, 100)

func addparam(flag, viber string, def interface{}, help string) (p TheParameter) {
	p = TheParameter{
		Flag:    flag,
		Viber:   viber,
		Default: def,
		Help:    help,
	}
	Parameters[viber] = p
	return
}
func init() {
	addparam("edgerc", "akamai.edgerc", "~/.edgerc", "akamai location of the credentials file")
	addparam("section", "akamai.section", "default", "akamai section of the credentials file")
	addparam("accountkey", "akamai.accountkey", "", "akamai account switch key")
	addparam("resolver", "dns.resolver", "8.8.8.8:53", "resolver to be used")
	addparam("cache", "akamai.cache", "~/.akamai-cache", "cache folder")
	addparam("loglevel", "log.level", "FATAL", "logging level")
	addparam("logfile", "log.file", "", "logging output")
	addparam("warningdays", "certificate.warningdays", 14, "warning days for certificate issues")
}
