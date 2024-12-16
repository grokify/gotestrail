package gotestrail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type CaseAPI struct{ client *Client }

func NewCaseAPI(client *Client) CaseAPI { return CaseAPI{client: client} }

func BuildFuncRespSavePages(dir string, projectID uint) func(qry url.Values, respBody []byte) error {
	return func(qry url.Values, respBody []byte) error {
		secID := qry.Get(QueryParamSectionID)
		offset := qry.Get(QueryParamOffset)
		filename := fmt.Sprintf("cases_proj-%d_sec-%s_offset-%s.json", projectID, secID, offset)
		absFilename := filepath.Join(dir, filename)
		return os.WriteFile(absFilename, respBody, 0600)
	}
}

func (api CaseAPI) GetCaseSetSectionsFunc(ctx context.Context, projectID uint, sectionIDs []uint, fnResp func(qry url.Values, respBody []byte) error) error {
	for _, secID := range sectionIDs {
		v := url.Values{}
		v.Add(QueryParamSectionID, fmt.Sprintf("%d", secID))
		_, err := api.GetCaseSetAllFunc(ctx, projectID, v, fnResp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (api CaseAPI) GetCaseSetAllFunc(ctx context.Context, projectID uint, qry url.Values, fnResp func(qry url.Values, respBody []byte) error) (*CaseSet, error) {
	limit := LimitMax
	offset := uint(0)
	set := NewCaseSet()
	for {
		casesRes, b, resp, err := api.GetCases(ctx, int(projectID), int(limit), int(offset), qry, true)
		if err != nil {
			return set, err
		} else if resp.StatusCode >= 300 {
			return set, errors.New("status code >= 300")
		} else {
			if fnResp != nil {
				fmt.Println(string(b))
				fmt.Printf("GOT_CASES (%d)\n", len(casesRes.Cases))
				qry.Add(QueryParamOffset, fmt.Sprintf("%d", offset))
				if err := fnResp(qry, b); err != nil {
					return nil, err
				}
			}
			if len(casesRes.Cases) == 0 {
				break
			} else {
				set.Add(casesRes.Cases...)
				if casesRes.Links.Next == nil {
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
	if set, err := api.GetCaseSetAllFunc(ctx, projectID, qry, nil); err != nil {
		return nil, err
	} else {
		return set, set.WriteFileJSON(filename, perm, prefix, indent)
	}
}

func (api CaseAPI) GetCases(ctx context.Context, projectID, limit, offset int, qry url.Values, parseResponse bool) (*GetCasesResponse, []byte, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathCasesGet, projectID, limit, offset, qry),
	}
	resp, err := api.client.sclient.Do(sreq)
	if err != nil {
		return nil, []byte{}, nil, err
	} else if !parseResponse {
		return nil, []byte{}, resp, err
	}
	res := GetCasesResponse{}
	if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, b, nil, err
	} else {
		return &res, b, resp, json.Unmarshal(b, &res)
	}
}

type GetCasesResponse struct {
	Offset uint   `json:"offset"`
	Limit  uint   `json:"limit"`
	Size   uint   `json:"size"`
	Links  Links  `json:"_links"`
	Cases  []Case `json:"cases"`
}

func ReadFileAPIResponseGetCases(filename string) (*GetCasesResponse, error) {
	r := GetCasesResponse{}
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}
