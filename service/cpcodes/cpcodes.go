package cpcodes

import (
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
)

type CpcodeService struct {
	Session  session.Session
	Response *http.Response
}

type CpcodesRequest struct {
	Cpcodes []Cpcode `json:"cpcodes"`
}

func (cpr *CpcodesRequest) FindCpcode(cpcodeid int) (cpcodeinfo *Cpcode) {
	for i := range cpr.Cpcodes {
		if cpr.Cpcodes[i].CpcodeID == cpcodeid {
			cpcodeinfo = &cpr.Cpcodes[i]
			return
		}
	}
	return
}

type CpcodeTimezone struct {
	TimezoneID    string `json:"timezoneId"`
	TimezoneValue string `json:"timezoneValue"`
}
type CpcodeContract struct {
	ContractID string `json:"contractId"`
	Status     string `json:"status"`
}
type CpcodeProduct struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
}
type AccessGroupshort struct {
	GroupID    interface{} `json:"groupId"`
	ContractID string      `json:"contractId"`
}
type Cpcode struct {
	CpcodeID         int              `json:"cpcodeId"`
	CpcodeName       string           `json:"cpcodeName"`
	Purgeable        bool             `json:"purgeable"`
	AccountID        string           `json:"accountId"`
	DefaultTimezone  *CpcodeTimezone  `json:"defaultTimezone"`
	OverrideTimezone *CpcodeTimezone  `json:"overrideTimezone"`
	Type             string           `json:"type"`
	Contracts        []CpcodeContract `json:"contracts"`
	Products         []CpcodeProduct  `json:"products"`
	AccessGroup      AccessGroupshort `json:"accessGroup"`
}

type Repgroups struct {
	Groups []Repgroup `json:"groups"`
}
type Cpcodeshort struct {
	CpcodeID   int    `json:"cpcodeId"`
	CpcodeName string `json:"cpcodeName"`
}
type RepgroupContracts struct {
	ContractID string        `json:"contractId"`
	Cpcodes    []Cpcodeshort `json:"cpcodes"`
}

type Repgroup struct {
	ReportingGroupID   int                 `json:"reportingGroupId"`
	ReportingGroupName string              `json:"reportingGroupName"`
	Contracts          []RepgroupContracts `json:"contracts"`
	AccessGroup        AccessGroupshort    `json:"accessGroup"`
}

// FindByName takes the name and returns the id
func (rgs *Repgroups) FindByName(name string) (rg *Repgroup) {
	for i := range rgs.Groups {
		if rgs.Groups[i].ReportingGroupName == name {
			return &rgs.Groups[i]
		}
	}
	return
}

func (rg *Repgroup) FindCpcode(cpcode int) (found bool) {
	for i := range rg.Contracts {
		for j := range rg.Contracts[i].Cpcodes {
			if rg.Contracts[i].Cpcodes[j].CpcodeID == cpcode {
				found = true
				return
			}
		}
	}
	return
}

func NewCpcodeService(s session.Session) (cs *CpcodeService) {
	cs = &CpcodeService{
		Session: s,
	}
	return
}

func (cs *CpcodeService) GetCpcodes() (cpcodes *CpcodesRequest, err error) {

	cpcodes = &CpcodesRequest{}
	req, _ := http.NewRequest(http.MethodGet, "/cprg/v1/cpcodes", nil)
	cs.Response, err = cs.Session.Exec(req, &cpcodes)

	return
}

func (cs *CpcodeService) GetRepgroups(filter string) (repGroups *Repgroups, err error) {
	url := "/cprg/v1/reporting-groups"
	if filter != "" {
		url += "?" + filter
	}
	repGroups = &Repgroups{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	cs.Response, err = cs.Session.Exec(req, &repGroups)

	return
}

// MapCpcodeRepgroup creates a map indexed by the CPCode, returning a list of Reporting Group Ids the CPCode belongs too.
func (repGroups *Repgroups) MapCpcodeRepgroup() (m map[int][]int) {
	m = make(map[int][]int)
	for _, g := range repGroups.Groups {
		for _, c := range g.Contracts {
			for _, cp := range c.Cpcodes {
				_, f := m[cp.CpcodeID]
				if f {
					m[cp.CpcodeID] = append(m[cp.CpcodeID], g.ReportingGroupID)
				} else {
					m[cp.CpcodeID] = make([]int, 1, 3)
					m[cp.CpcodeID][0] = g.ReportingGroupID
				}
			}
		}
	}

	return
}

func (rg *Repgroups) FindByCpcode(cpcode int) (sg *Repgroups) {
	sg = &Repgroups{}
	for _, g := range rg.Groups {
		found := false
		for _, c := range g.Contracts {
			for _, cp := range c.Cpcodes {
				if cp.CpcodeID == cpcode {
					sg.Groups = append(sg.Groups, g)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	return
}
