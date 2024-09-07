package gotestrail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

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
