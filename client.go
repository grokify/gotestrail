package gotestrail

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

type Client struct {
	sclient     httpsimple.Client
	CaseAPI     CaseAPI
	CaseTypeAPI CaseTypeAPI
	SectionAPI  SectionAPI
}

func NewClient(serverURL, username, password string) (*Client, error) {
	c := &Client{}
	hclient, err := authutil.NewClientBasicAuth(username, password, false)
	if err != nil {
		return nil, err
	}
	sclient := httpsimple.NewClient(hclient, serverURL)
	c.sclient = sclient
	c.CaseAPI = NewCaseAPI(c)
	c.CaseTypeAPI = NewCaseTypeAPI(c)
	c.SectionAPI = NewSectionAPI(c)
	return c, nil
}

func BuildAPIURL(baseURL, apiPath string, id, limit, offset int, qry url.Values) string {
	if limit > 0 {
		qry[ParamLimit] = []string{strconv.Itoa(limit)}
	} else if limit == 0 {
		qry[ParamLimit] = []string{strconv.Itoa(int(LimitMax))}
	}
	if offset >= 0 {
		qry[ParamOffset] = []string{strconv.Itoa(offset)}
	}
	var qryenc string
	if len(qry) > 0 {
		qryenc = "&" + qry.Encode()
	}
	apiPath = IndexPath + apiPath
	path := []string{baseURL, apiPath}
	if id >= 0 {
		path = append(path, strconv.Itoa(id))
	}
	fmt.Printf("URL: %s\n", urlutil.JoinAbsolute(path...)+qryenc)
	return urlutil.JoinAbsolute(path...) + qryenc
}

type Links struct {
	Next *string `json:"next"`
	Prev *string `json:"prev"`
}
