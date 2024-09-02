package gotestrail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type CaseTypeAPI struct{ client *Client }

func NewCaseTypeAPI(client *Client) CaseTypeAPI { return CaseTypeAPI{client: client} }

func (api CaseTypeAPI) GetCaseTypes(ctx context.Context) ([]CaseType, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathCaseTypesGet, -1, -1, -1, url.Values{}),
	}
	var res []CaseType
	if resp, err := api.client.sclient.Do(sreq); err != nil {
		return nil, nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else {
		return res, resp, json.Unmarshal(b, &res)
	}
}

type CaseType struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

type CaseTypeSet struct {
	CaseTypes map[uint]CaseType `json:"caseTypes"`
}

func NewCaseTypeSet() *CaseTypeSet {
	return &CaseTypeSet{CaseTypes: map[uint]CaseType{}}
}

func (set *CaseTypeSet) Add(items ...CaseType) {
	if set.CaseTypes == nil {
		set.CaseTypes = map[uint]CaseType{}
	}
	for _, item := range items {
		set.CaseTypes[item.ID] = item
	}
}

func ReadFileCaseTypeSet(filename string) (*CaseTypeSet, error) {
	set := NewCaseTypeSet()
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return set, json.Unmarshal(b, set)
	}
}

func (set *CaseTypeSet) ReadFileJSON(filename string) error {
	if new, err := ReadFileCaseTypeSet(filename); err != nil {
		return err
	} else {
		for k, v := range new.CaseTypes {
			set.CaseTypes[k] = v
		}
	}
	return nil
}
