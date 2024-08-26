package gotestrail

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type CasesAPI struct{ client *Client }

func NewCasesAPI(client *Client) CasesAPI { return CasesAPI{client: client} }

func (api CasesAPI) GetCasesAll(ctx context.Context, projectID uint, qry url.Values) (*CaseSet, error) {
	limit := LimitMax
	offset := uint(0)
	set := NewCaseSet()
	for {
		res, resp, err := api.GetCases(ctx, int(projectID), int(limit), int(offset), qry)
		if err != nil {
			return set, err
		} else if resp.StatusCode >= 300 {
			return set, errors.New("status code >= 300")
		} else {
			if len(res.Cases) == 0 {
				break
			} else {
				set.Add(res.Cases...)
				if res.Links.Next == nil {
					break
				} else {
					offset += 250
				}
			}
		}
	}
	return set, nil
}

func (api CasesAPI) GetCases(ctx context.Context, projectID, limit, offset int, qry url.Values) (*GetCasesResponse, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathCasesGet, projectID, limit, offset, qry),
	}
	res := GetCasesResponse{}
	if resp, err := api.client.sclient.Do(sreq); err != nil {
		return nil, nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else {
		err := os.WriteFile("cases_temp.json", b, 0600)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, &res)
		return &res, resp, err
	}
}

type GetCasesResponse struct {
	Offset uint   `json:"offset"`
	Limit  uint   `json:"limit"`
	Size   uint   `json:"size"`
	Links  Links  `json:"_links"`
	Cases  []Case `json:"cases"`
}
