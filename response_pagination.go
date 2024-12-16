package gotestrail

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/grokify/mogo/pointer"
)

// ResponsePagination is a container for pagination when only pagination
// parsing is needed.
type ResponsePagination struct {
	Offset   int
	Limit    int
	Size     int
	LinkNext string
	LinkPrev string
}

func ParseResponsePaginationFromBody(r io.Reader) (ResponsePagination, []byte, error) {
	if b, err := io.ReadAll(r); err != nil {
		return ResponsePagination{}, []byte{}, err
	} else {
		rp, err := ParseResponsePaginationFromRaw(b)
		return rp, b, err
	}
}

func ParseResponsePaginationFromRaw(b []byte) (ResponsePagination, error) {
	raw := responsePaginationRaw{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return ResponsePagination{}, err
	} else {
		return raw.ResponsePagination(), nil
	}
}

type responsePaginationRaw struct {
	Offset int                        `json:"offset"`
	Limit  int                        `json:"limit"`
	Size   int                        `json:"size"`
	Links  responsePaginationRawLinks `json:"_links"`
}

type responsePaginationRawLinks struct {
	Next *string `json:"next"`
	Prev *string `json:"prev"`
}

func (raw responsePaginationRaw) ResponsePagination() ResponsePagination {
	return ResponsePagination{
		Offset:   raw.Offset,
		Limit:    raw.Limit,
		Size:     raw.Size,
		LinkNext: strings.TrimSpace(pointer.Dereference(raw.Links.Next)),
		LinkPrev: strings.TrimSpace(pointer.Dereference(raw.Links.Prev)),
	}
}

/*
  "offset": 0,
  "limit": 250,
  "size": 250,
  "_links": {
    "next": "\/api\/v2\/get_cases\/13&limit=250&offset=250",
    "prev": null
  },
*/
