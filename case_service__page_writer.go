package gotestrail

import (
	"errors"
	"net/http"

	"github.com/grokify/mogo/net/http/httputilmore"
)

// ProcessCaseAPIResponse processes a http.Response and returns a next page URL along with body
// bytes.
func ProcessCaseAPIResponse(r *http.Response, fn func(b []byte) error) (string, []byte, error) {
	if r == nil {
		return "", []byte{}, httputilmore.ErrHTTPResponseCannotBeNil
	} else if r.StatusCode >= 300 {
		return "", []byte{}, errors.New("http response status code >= 300")
	} else if rp, b, err := ParseResponsePaginationFromBody(r.Body); err != nil {
		return rp.LinkNext, b, err
	} else if fn == nil {
		return rp.LinkNext, b, nil
	} else {
		return rp.LinkNext, b, fn(b)
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
