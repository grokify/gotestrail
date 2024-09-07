package gotestrail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type CaseFieldAPI struct{ client *Client }

func NewCaseFieldAPI(client *Client) CaseFieldAPI { return CaseFieldAPI{client: client} }

func (api CaseFieldAPI) GetCaseFields(ctx context.Context) (CaseFields, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathCaseFieldsGet, -1, -1, -1, url.Values{}),
	}
	var res []CaseField
	if resp, err := api.client.sclient.Do(sreq); err != nil {
		return nil, nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else {
		return res, resp, json.Unmarshal(b, &res)
	}
}
