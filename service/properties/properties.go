package properties

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/util/osutil"
)

type Propsv struct {
	PapiClient      papi.PAPI
	CacheDir        string
	propertiesCache map[papi.GetPropertiesRequest]*papi.GetPropertiesResponse
}

func NewPropertyService(p papi.PAPI, cache string) (ps *Propsv) {
	ps = &Propsv{
		PapiClient:      p,
		CacheDir:        osutil.ExpandPath(cache),
		propertiesCache: make(map[papi.GetPropertiesRequest]*papi.GetPropertiesResponse),
	}
	return
}

func (p *Propsv) GetProperties(ctx context.Context, params papi.GetPropertiesRequest) (*papi.GetPropertiesResponse, error) {

	pr, found := p.propertiesCache[params]
	if found {
		return pr, nil
	}

	pr, err := p.PapiClient.GetProperties(ctx, params)
	if err == nil {
		p.propertiesCache[params] = pr
	}

	return pr, err
}

func (ps *Propsv) GetRuleTree(params papi.GetRuleTreeRequest) (tree *papi.GetRuleTreeResponse) {
	filename := fmt.Sprintf("%s/%s/property/%s_%d.json", ps.CacheDir, params.ContractID, params.PropertyID, params.PropertyVersion)

	if b := LoadFile(filename); b != nil {
		tree = &papi.GetRuleTreeResponse{}
		json.Unmarshal(*b, tree)
	} else {
		var err error
		tree, err = ps.PapiClient.GetRuleTree(context.Background(), params)

		if err != nil {
			log.Error(fmt.Sprint(err))
			return
		}

		SaveFile(filename, *tree)
	}
	return
}

type UsedBehavior struct {
	Behavior *papi.RuleBehavior
	Criteria [][]papi.RuleBehavior
}
type PropSum struct {
	Behaviors map[string][]UsedBehavior
}

func walkTree(psum *PropSum, r *papi.Rules, c [][]papi.RuleBehavior) {
	if c == nil {
		c = make([][]papi.RuleBehavior, 0, 10)
	}
	for i, b := range r.Behaviors {
		if psum.Behaviors == nil {
			psum.Behaviors = make(map[string][]UsedBehavior, 100)
		}
		pb, f := psum.Behaviors[b.Name]
		if !f {
			pb = make([]UsedBehavior, 0, 10)
		}
		u := UsedBehavior{
			Behavior: &r.Behaviors[i],
			Criteria: c,
		}
		pb = append(pb, u)
		psum.Behaviors[b.Name] = pb
	}
	for t := range r.Children {
		var c2 [][]papi.RuleBehavior
		copy(c2, c)
		c2 = append(c2, r.Children[t].Criteria)
		walkTree(psum, &r.Children[t], c2)
	}
}

func (ps *Propsv) FindBehaviors(r *papi.Rules) (propsum *PropSum) {
	propsum = &PropSum{}
	walkTree(propsum, r, nil)
	return
}

func LoadFile(filename string) (b *[]byte) {
	jsonFile, err := os.Open(filename)
	if err == nil {
		defer jsonFile.Close()
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Errorf("load error %w", err)
		}
		b = &byteValue
		return
	}
	return
}

func SaveFile(filename string, x interface{}) {
	os.MkdirAll(filepath.Dir(filename), 0750)

	byteblob, err2 := json.Marshal(x)
	if err2 != nil {
		log.Errorf("marshall error %w", err2)
	}
	err2 = os.WriteFile(filename, byteblob, 0644)
	if err2 != nil {
		log.Errorf("write error - %w", err2)
	}
}

func (ps *Propsv) GetPropertyVersionHostnames(params papi.GetPropertyVersionHostnamesRequest) (pr *papi.GetPropertyVersionHostnamesResponse, err error) {

	filename := fmt.Sprintf("%s/%s/property/%s_%d_host.json", ps.CacheDir, params.ContractID, params.PropertyID, params.PropertyVersion)

	if b := LoadFile(filename); b != nil {
		pr = &papi.GetPropertyVersionHostnamesResponse{}

		json.Unmarshal(*b, pr)
	} else {
		pr, err = ps.PapiClient.GetPropertyVersionHostnames(context.Background(), params)

		if err != nil {
			log.Error(fmt.Sprint(err))
			return
		}

		SaveFile(filename, *pr)
	}
	return
}
