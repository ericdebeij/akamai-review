package aksv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/yearmonth"
)

type ProductInfo struct {
	ProductName string
	StartMonth  string
	EndMonth    string
}
type ContractInfo struct {
	ContractID string
	ProductMap map[string]ProductInfo
}

type BillingService struct {
	Session  session.Session
	Response *http.Response
}

func NewBillingService(s session.Session) (bs *BillingService) {
	bs = &BillingService{
		Session: s,
	}
	return
}

type GetMonthlySummaryResponse struct {
	Start              string         `json:"start"`
	End                string         `json:"end"`
	RequestDate        time.Time      `json:"requestDate"`
	AccountID          string         `json:"accountId"`
	ContractID         string         `json:"contractId"`
	ProductID          string         `json:"productId"`
	ProductName        string         `json:"productName"`
	ReportingGroupID   int            `json:"reportingGroupId,omitempty"`
	ReportingGroupName string         `json:"reportingGroupName,omitempty"`
	UsagePeriods       []UsagePeriods `json:"usagePeriods"`
}
type Stats struct {
	StatType   string  `json:"statType"`
	Unit       string  `json:"unit"`
	IsBillable bool    `json:"isBillable"`
	Value      float64 `json:"value"`
}
type UsagePeriods struct {
	Month       string        `json:"month"`
	Start       string        `json:"start"`
	End         string        `json:"end"`
	Region      string        `json:"region"`
	DataStatus  string        `json:"dataStatus"`
	CpCodes     []int         `json:"cpCodes,omitempty"`
	Stats       []Stats       `json:"stats,omitempty"`
	CpCodeStats []CpCodeStats `json:"cpCodeStats,omitempty"`
}

type CpCodeStats struct {
	Stats  []Stats `json:"stats"`
	CpCode int     `json:"cpCode"`
}

func (bs *BillingService) GetMonthlySummary(contractId string, rgroupId int, productId string, startMonth string, endMonth string) (msum *GetMonthlySummaryResponse, err error) {
	var url string
	s := startMonth
	c := 204
	for c == 204 && s <= endMonth {
		if rgroupId == 0 {
			url = fmt.Sprintf("/billing/v1/contracts/%s/products/%s/usage/monthly-summary?start=%s&end=%s",
				contractId, productId, s, endMonth)
		} else {
			url = fmt.Sprintf("/billing/v1/reporting-groups/%d/products/%s/usage/monthly-summary?start=%s&end=%s",
				rgroupId, productId, s, endMonth)
		}
		log.Debugf("monthly summary: %s", url)
		msum = &GetMonthlySummaryResponse{}
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		bs.Response, err = bs.Session.Exec(req, &msum)

		c = bs.Response.StatusCode

		// Workaround for 503 instead of 204
		if c == 503 {
			log.Errorf("unexpected 503 in %s, assume there is no data", url)
			c = 204
		}

		if c == 204 {
			s = yearmonth.Add(s, 1)
		}
	}
	if c >= 202 {
		err = fmt.Errorf("billing service no data(%d): %s %d %s %s %s", c, contractId, rgroupId, productId, startMonth, endMonth)
		msum = nil
	}
	return
}

type GetUsageProductsResponse struct {
	Start        string    `json:"start"`
	End          string    `json:"end"`
	RequestDate  time.Time `json:"requestDate"`
	AccountID    string    `json:"accountId"`
	ContractID   string    `json:"contractId"`
	UsagePeriods []struct {
		Month         string    `json:"month"`
		UsageProducts []Product `json:"usageProducts"`
	} `json:"usagePeriods"`
}

type Product struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
}

func (bs *BillingService) GetUsageProducts(contractId string, startMonth, endMonth string) (products *GetUsageProductsResponse, err error) {

	m := startMonth
	for endMonth > string(m) {
		url := fmt.Sprintf("/billing/v1/contracts/%s/products?start=%s&end=%s", contractId, string(m), endMonth)
		log.Debug(url)
		req, err2 := http.NewRequest(http.MethodGet, url, nil)
		if err2 != nil {
			err = err2
			return
		}
		products = &GetUsageProductsResponse{}
		bs.Response, err = bs.Session.Exec(req, products)

		if bs.Response.StatusCode != 204 && (bs.Response.StatusCode != 400 || string(m) == endMonth) {
			if bs.Response.StatusCode != 200 {
				err = fmt.Errorf("usage error (%d), %s %s %s", bs.Response.StatusCode, contractId, string(m), endMonth)
				products = nil
			}
			return
		}

		m = yearmonth.Add(m, 1)
	}
	products = nil
	err = fmt.Errorf("no usage products: %s %s %s", contractId, startMonth, endMonth)
	return
}
func (bs *BillingService) GetUsageCpcode(contractId, productId, startMonth, endMonth string) (msum *GetMonthlySummaryResponse, err error) {

	url := fmt.Sprintf("/billing/v1/contracts/%s/products/%s/usage/by-cp-code/monthly-summary?start=%s&end=%s", contractId, productId, startMonth, endMonth)
	log.Debug(url)
	req, err2 := http.NewRequest(http.MethodGet, url, nil)
	if err2 != nil {
		log.Fatalf("request error for : %w", err)
		err = err2
		return
	}
	msum = &GetMonthlySummaryResponse{}

	bs.Response, err = bs.Session.Exec(req, msum)
	if err != nil {
		log.Fatalf("response error: %w", err)
		return
	}

	return
}

func (bs *BillingService) GetContractInfo(contractId string, startMonth, toMonth string) (cinfo *ContractInfo) {
	cinfo = &ContractInfo{
		ContractID: contractId,
		ProductMap: map[string]ProductInfo{},
	}

	prods, err := bs.GetUsageProducts(contractId, startMonth, toMonth)

	if err != nil {
		log.Errorf("usageProducts: %s", err)
		return
	}
	for _, uperiod := range prods.UsagePeriods {
		for _, uprod := range uperiod.UsageProducts {
			t, ok := cinfo.ProductMap[uprod.ProductID]
			if !ok {
				t = ProductInfo{
					ProductName: uprod.ProductName,
					StartMonth:  uperiod.Month,
					EndMonth:    uperiod.Month,
				}
			}
			if uperiod.Month > t.EndMonth {
				t.EndMonth = uperiod.Month
			}
			if uperiod.Month < t.StartMonth {
				t.StartMonth = uperiod.Month
			}
			cinfo.ProductMap[uprod.ProductID] = t
		}
	}
	return
}
