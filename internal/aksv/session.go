package aksv

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
)

type EdgeConfig struct {
	Edgerc    string
	Section   string
	AccountID string
}

func NewSession(param *EdgeConfig) (s session.Session, err error) {
	edgercOps := []edgegrid.Option{edgegrid.WithEnv(true)}
	edgercOps = append(edgercOps, edgegrid.WithFile(param.Edgerc))
	edgercOps = append(edgercOps, edgegrid.WithSection(param.Section))

	edgerc, err := edgegrid.New(edgercOps...)
	edgerc.AccountKey = param.AccountID

	if err != nil {
		log.Fatalf("edgerc error: %w")
		return
	}

	//fmt.Println(edgerc)
	s, err = session.New(
		session.WithSigner(edgerc),
		session.WithHTTPTracing(true),
		session.WithLog(log.Log),
	)
	return
}
