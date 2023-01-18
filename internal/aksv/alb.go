package aksv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

type AlbService struct {
	Session  session.Session
	Response *http.Response
}

func NewAlbService(s session.Session) (alb *AlbService) {
	alb = &AlbService{
		Session: s,
	}
	return
}

type AlbActivation struct {
	ActivatedBy   string    `json:"activatedBy"`
	ActivatedDate time.Time `json:"activatedDate"`
	Network       string    `json:"network"`
	OriginID      string    `json:"originId"`
	Status        string    `json:"status"`
	Version       int       `json:"version"`
}
type LoadBalancers map[string]AlbActivation
type AlbActivations map[string]LoadBalancers

func (alb *AlbService) ListActivation() (activations *AlbActivations, err error) {
	activations = &AlbActivations{}
	url := "/cloudlets/api/v2/origins/currentActivations"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	alb.Response, err = alb.Session.Exec(req, &activations)
	return
}

type LoadBalancingVersion struct {
	BalancingType    string            `json:"balancingType"`
	CreatedBy        string            `json:"createdBy"`
	CreatedDate      time.Time         `json:"createdDate"`
	Deleted          bool              `json:"deleted"`
	Description      string            `json:"description"`
	Immutable        bool              `json:"immutable"`
	LastModifiedBy   string            `json:"lastModifiedBy"`
	LastModifiedDate time.Time         `json:"lastModifiedDate"`
	OriginID         string            `json:"originId"`
	Version          int               `json:"version"`
	LivenessSettings *LivenessSettings `json:"livenessSettings"`
	DataCenters      []DataCenters     `json:"dataCenters"`
}
type LivenessSettings struct {
	HostHeader                  string  `json:"hostHeader"`
	Interval                    int     `json:"interval"`
	Path                        string  `json:"path"`
	PeerCertificateVerification bool    `json:"peerCertificateVerification"`
	Port                        int     `json:"port"`
	Protocol                    string  `json:"protocol"`
	Status3XxFailure            bool    `json:"status3xxFailure"`
	Status4XxFailure            bool    `json:"status4xxFailure"`
	Status5XxFailure            bool    `json:"status5xxFailure"`
	Timeout                     float32 `json:"timeout"`
}
type DataCenters struct {
	Hostname                      string   `json:"hostname"`
	CloudServerHostHeaderOverride bool     `json:"cloudServerHostHeaderOverride"`
	CloudService                  bool     `json:"cloudService"`
	Continent                     string   `json:"continent"`
	Country                       string   `json:"country"`
	Latitude                      float64  `json:"latitude"`
	Longitude                     float64  `json:"longitude"`
	OriginID                      string   `json:"originId"`
	Percent                       float32  `json:"percent"`
	LivenessHosts                 []string `json:"livenessHosts"`
}

func (alb *AlbService) ListAlbVersionDetails(albid string, ver int) (lbv *LoadBalancingVersion, err error) {
	lbv = &LoadBalancingVersion{}
	url := fmt.Sprintf("/cloudlets/api/v2/origins/%s/versions/%d", albid, ver)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	alb.Response, err = alb.Session.Exec(req, &lbv)
	return
}

type PolicyList []Policy

type Policy struct {
	Location         string              `json:"location"`
	ServiceVersion   interface{}         `json:"serviceVersion"`
	PolicyID         int                 `json:"policyId"`
	GroupID          int                 `json:"groupId"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	PropertyName     interface{}         `json:"propertyName"`
	CreatedBy        string              `json:"createdBy"`
	CreateDate       int64               `json:"createDate"`
	LastModifiedBy   string              `json:"lastModifiedBy"`
	LastModifiedDate int64               `json:"lastModifiedDate"`
	Activations      []PolicyActivations `json:"activations"`
	CloudletCode     string              `json:"cloudletCode"`
	CloudletID       int                 `json:"cloudletId"`
	APIVersion       string              `json:"apiVersion"`
	Deleted          bool                `json:"deleted"`
}
type PolicyInfo struct {
	PolicyID       int    `json:"policyId"`
	Name           string `json:"name"`
	Version        int    `json:"version"`
	Status         string `json:"status"`
	ActivatedBy    string `json:"activatedBy"`
	ActivationDate int64  `json:"activationDate"`
}
type PropertyInfo struct {
	Name           string `json:"name"`
	Version        int    `json:"version"`
	GroupID        int    `json:"groupId"`
	Status         string `json:"status"`
	ActivatedBy    string `json:"activatedBy"`
	ActivationDate int64  `json:"activationDate"`
	ID             int    `json:"id"`
}
type PolicyActivations struct {
	ServiceVersion interface{}  `json:"serviceVersion"`
	Network        string       `json:"network"`
	PolicyInfo     PolicyInfo   `json:"policyInfo"`
	PropertyInfo   PropertyInfo `json:"propertyInfo"`
	APIVersion     string       `json:"apiVersion"`
}

func (alb *AlbService) ListPolicies() (policylist *PolicyList, err error) {
	policylist = &PolicyList{}
	url := "/cloudlets/api/v2/policies?cloudletId=9"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	alb.Response, err = alb.Session.Exec(req, &policylist)
	return
}

type PolicyVersion struct {
	Location         string              `json:"location"`
	RevisionID       int                 `json:"revisionId"`
	PolicyID         int                 `json:"policyId"`
	Version          int                 `json:"version"`
	Description      string              `json:"description"`
	CreatedBy        string              `json:"createdBy"`
	CreateDate       int64               `json:"createDate"`
	LastModifiedBy   string              `json:"lastModifiedBy"`
	LastModifiedDate int64               `json:"lastModifiedDate"`
	Activations      []PolicyActivations `json:"activations"`
	MatchRules       []MatchRules        `json:"matchRules"`
	MatchRuleFormat  string              `json:"matchRuleFormat"`
	Deleted          bool                `json:"deleted"`
	RulesLocked      bool                `json:"rulesLocked"`
}

type Matches struct {
	MatchValue    string `json:"matchValue"`
	MatchOperator string `json:"matchOperator"`
	Negate        bool   `json:"negate"`
	CaseSensitive bool   `json:"caseSensitive"`
	MatchType     string `json:"matchType"`
}
type ForwardSettings struct {
	OriginID string `json:"originId"`
}
type MatchRules struct {
	Type            string          `json:"type"`
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Start           int             `json:"start"`
	End             int             `json:"end"`
	MatchURL        interface{}     `json:"matchURL"`
	Matches         []Matches       `json:"matches"`
	AkaRuleID       string          `json:"akaRuleId"`
	Location        string          `json:"location"`
	ForwardSettings ForwardSettings `json:"forwardSettings"`
}

func (alb *AlbService) PolicyVersion(policyId, version int) (pv *PolicyVersion, err error) {
	pv = &PolicyVersion{}
	url := fmt.Sprintf("/cloudlets/api/v2/policies/%d/versions/%d", policyId, version)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	alb.Response, err = alb.Session.Exec(req, &pv)
	return
}
