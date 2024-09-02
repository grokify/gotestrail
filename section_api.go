package gotestrail

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type SectionAPI struct{ client *Client }

func NewSectionAPI(client *Client) SectionAPI { return SectionAPI{client: client} }

func (api SectionAPI) GetSectionSetAll(ctx context.Context, projectID uint, qry url.Values) (*SectionSet, error) {
	limit := LimitMax
	offset := uint(0)
	ss := NewSectionSet()
	for {
		res, resp, err := api.GetSections(ctx, projectID, limit, offset, qry)
		if err != nil {
			return ss, err
		} else if resp.StatusCode >= 300 {
			return ss, errors.New("status code >= 300")
		} else {
			if len(res.Sections) == 0 {
				break
			} else {
				ss.Add(res.Sections...)
				if res.Links.Next == nil {
					break
				} else {
					offset += 250
				}
			}
		}
	}
	return ss, nil
}

func (api SectionAPI) GetSections(ctx context.Context, projectID, limit, offset uint, qry url.Values) (*GetSectionsResponse, *http.Response, error) {
	sreq := httpsimple.Request{
		Method: http.MethodGet,
		URL:    BuildAPIURL(api.client.sclient.BaseURL, APIPathSectionsGet, int(projectID), int(limit), int(offset), qry),
	}
	res := GetSectionsResponse{}
	if resp, err := api.client.sclient.Do(sreq); err != nil {
		return nil, nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else {
		err = json.Unmarshal(b, &res)
		return &res, resp, err
	}
}

// "\/api\/v2\/get_sections\/13&limit=250&offset=250",

type GetSectionsResponse struct {
	Offset   uint      `json:"offset"`
	Limit    uint      `json:"limit"`
	Size     uint      `json:"size"`
	Links    Links     `json:"_links"`
	Sections []Section `json:"sections"`
}
