package aksv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v2/internal/akutil"
)

type Propsv struct {
	PapiClient      papi.PAPI
	CacheDir        string
	propertiesCache map[papi.GetPropertiesRequest]*papi.GetPropertiesResponse
}

func NewPropertyService(p papi.PAPI, cache string) (ps *Propsv) {
	ps = &Propsv{
		PapiClient:      p,
		CacheDir:        akutil.ExpandPath(cache),
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

func (p *Propsv) GetPropertyVersionHostnames(ctx context.Context, params papi.GetPropertyVersionHostnamesRequest) (*papi.GetPropertyVersionHostnamesResponse, error) {
	filename := fmt.Sprintf("%s/prophosts/%s_%d.json", p.CacheDir, params.PropertyID, params.PropertyVersion)
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err == nil {
		defer jsonFile.Close()

		pr := &papi.GetPropertyVersionHostnamesResponse{}
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal(byteValue, pr)
		return pr, nil

	} else {
		if !os.IsNotExist(err) {
			return nil, err
		}

		pr, err := p.PapiClient.GetPropertyVersionHostnames(ctx, params)

		// Now we have this, lets check if this is a locked configuration, in which case we can store it in cache.
		gpp := papi.GetPropertiesRequest{
			ContractID: params.ContractID,
			GroupID:    params.GroupID,
		}

		gpr, found := p.propertiesCache[gpp]

		if found {
			for _, prop := range gpr.Properties.Items {
				if params.PropertyID == prop.PropertyID {

					if params.PropertyVersion == *prop.ProductionVersion || params.PropertyVersion == *prop.StagingVersion {
						byteblob, err2 := json.Marshal(pr)
						if err2 != nil {
							log.Errorf("marshall ", err2)
						}
						os.MkdirAll(filepath.Dir(filename), 0750)
						err2 = os.WriteFile(filename, byteblob, 0644)
						if err2 != nil {
							log.Errorf("write %w", err2)
						}
						break
					}
				}
			}
		}
		return pr, err
	}
}