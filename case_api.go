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

type CaseAPI struct{ client *Client }

func NewCaseAPI(client *Client) CaseAPI { return CaseAPI{client: client} }

func (api CaseAPI) GetCaseSetAll(ctx context.Context, projectID uint, qry url.Values) (*CaseSet, error) {
	limit := LimitMax
	offset := uint(0)
	set := NewCaseSet()
	for {
		res, resp, err := api.GetCases(ctx, int(projectID), int(limit), int(offset), qry, true)
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

func (api CaseAPI) GetWriteFileCaseSetAll(ctx context.Context, filename string, perm os.FileMode, prefix, indent string, projectID uint, qry url.Values) (*CaseSet, error) {
	if set, err := api.GetCaseSetAll(ctx, projectID, qry); err != nil {
		return nil, err
	} else {
		return set, set.WriteFileJSON(filename, perm, prefix, indent)
	}
}

func (api CaseAPI) GetCases(ctx context.Context, projectID, limit, offset int, qry url.Values, parseResponse bool) (*GetCasesResponse, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathCasesGet, projectID, limit, offset, qry),
	}
	resp, err := api.client.sclient.Do(sreq)
	if err != nil {
		return nil, nil, err
	} else if !parseResponse {
		return nil, resp, err
	}

	res := GetCasesResponse{}
	if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else {
		return &res, resp, json.Unmarshal(b, &res)
	}
}

type GetCasesResponse struct {
	Offset uint   `json:"offset"`
	Limit  uint   `json:"limit"`
	Size   uint   `json:"size"`
	Links  Links  `json:"_links"`
	Cases  []Case `json:"cases"`
}
