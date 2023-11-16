package googlecse

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	srv "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm/GoogleCSE"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
	"github.com/google/uuid"
)

const (
	qpKeyword                      api.Key = "q"
	qpEnableChineseSearch          api.Key = "c2coff"
	qpSearchEngineId               api.Key = "cx"
	qpGeoLocation                  api.Key = "gl"
	qpSaveLevel                    api.Key = "safe"
	qpLanguage                     api.Key = "lr"
	qpEnableDuplicateContentFilter api.Key = "filter"
	qpDateRestrict                 api.Key = "dateRestrict"
	qpSiteSearch                   api.Key = "siteSearch"
	qpSiteSearchFilter             api.Key = "siteSearchFilter"
	qpExactTerms                   api.Key = "exactTerms"
	qpExcludeTerms                 api.Key = "excludeTerms"
	qpPageSize                     api.Key = "num"
	qpStart                        api.Key = "start"
	qpAPIKey                       api.Key = "key"
)

var ErrRequiredFieldMissing = errors.New("required field is missing")

// cseChineseSearch is a type for Chinese search.
// See https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list
// for more information.
// 0: enabled (default)
// 1: disabled
type cseChineseSearch int

const (
	EnableChineseSearch  cseChineseSearch = 0 // default
	DisableChineseSearch cseChineseSearch = 1
)

func (c cseChineseSearch) String() string {
	if c == EnableChineseSearch {
		return "enabled"
	}
	return "disabled"
}

func (c cseChineseSearch) ToParam() string {
	if c == DisableChineseSearch {
		// default value
		return ""
	}
	return strconv.Itoa(int(c))
}

// cseSafeLevel is a type for safe search.
// See https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list
// for more information.
// "off": Disables SafeSearch filtering. (default)
// "active": Enables SafeSearch filtering.
type cseSafeLevel int

const (
	DisableSafeSearch cseSafeLevel = 0 // default
	EnableSafeSearch  cseSafeLevel = 1
)

func (s cseSafeLevel) String() string {
	if s == EnableSafeSearch {
		return "active"
	}
	return "off"
}

func (s cseSafeLevel) ToParam() string {
	if s == DisableSafeSearch {
		return ""
	}
	return "active"
}

// cseDuplicateContentFilter is a type for duplicate content filter.
// See https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list
// for more information.
// 0: disabled (default)
// 1: enabled
type cseDuplicateContentFilter int

const (
	DisableDuplicateContentFilter cseDuplicateContentFilter = 0
	EnableDuplicateContentFilter  cseDuplicateContentFilter = 1 // default
)

func (f cseDuplicateContentFilter) String() string {
	if f == EnableDuplicateContentFilter {
		return "enabled"
	}
	return "disabled"
}

func (f cseDuplicateContentFilter) ToParam() string {
	if f == EnableDuplicateContentFilter {
		return ""
	}
	return strconv.Itoa(int(f))
}

type cseSiteSearchFilter int

const (
	ExcludeSiteSearchFilter  cseSiteSearchFilter = 0
	IncludesSiteSearchFilter cseSiteSearchFilter = 1
)

func (f cseSiteSearchFilter) String() string {
	if f == ExcludeSiteSearchFilter {
		return "excluded"
	}
	return "included"
}

func (f cseSiteSearchFilter) ToParam() string {
	if f == ExcludeSiteSearchFilter {
		return "e"
	}
	return "i"
}

// Request is a wrapper of customsearch.Service.Cse.List().
// See https://developers.google.com/custom-search/v1/reference/rest/v1/cse/list
// for more information.
type Request struct {
	*api.RequestProto
	SearchEngineId string `json:"search_engine_id,omitempty"`
	Start          int    `json:"start,omit"`
	PageSize       int    `json:"page_size,omitempty"`
	// CallOpts       []googleapi.CallOption `json:"-"`
}

func NewRequest(apikey string, engId string) (*Request, error) {
	if apikey == "" || engId == "" {
		return nil, ErrRequiredFieldMissing
	}

	req := api.NewRequestProtoType(srv.API_NAME, "")
	req.SetApiKey(apikey)
	req.Add(qpEnableChineseSearch, EnableChineseSearch.ToParam())
	req.Add(qpSaveLevel, EnableSafeSearch.ToParam())
	req.Add(qpEnableDuplicateContentFilter, EnableDuplicateContentFilter.ToParam())
	return &Request{
		RequestProto:   req,
		SearchEngineId: engId,
		Start:          1,
		PageSize:       DEFAULT_PAGE_SIZE,
	}, nil
}

func (req *Request) SetKeyword(keyword string) *Request {
	req.Set(qpKeyword, keyword)
	return req
}

func (req *Request) SetChineseSearch(c2off cseChineseSearch) *Request {
	req.Set(qpEnableChineseSearch, c2off.ToParam())
	return req
}

func (req *Request) SetGeoLocation(gl string) *Request {
	req.Set(qpGeoLocation, gl)
	return req
}

func (req *Request) SetSafeLevel(safe cseSafeLevel) *Request {
	req.Set(qpSaveLevel, safe.ToParam())
	return req
}

func (req *Request) SetLanguage(lr string) *Request {
	req.Set(qpLanguage, lr)
	return req
}

func (req *Request) SetEngineId(engId string) *Request {
	req.SearchEngineId = engId
	return req
}

func (req *Request) SetDuplicateContentFilter(filter cseDuplicateContentFilter) *Request {
	req.Set(qpEnableDuplicateContentFilter, filter.ToParam())
	return req
}

func (req *Request) SetDateRestict(dateRestict string) *Request {
	req.Set(qpDateRestrict, dateRestict)
	return req
}

func (req *Request) SetSiteSearch(site string, filter cseSiteSearchFilter) *Request {
	req.Set(qpSiteSearch, site)
	req.Set(qpSiteSearchFilter, filter.ToParam())
	return req
}

func (req *Request) WithExactTerm(term string) *Request {
	req.Add(qpExactTerms, term)
	return req
}

func (req *Request) WithExcludeTerm(term string) *Request {
	req.Add(qpExcludeTerms, term)
	return req
}

func (req *Request) SetPageSize(pagesize int) *Request {
	req.Set(qpPageSize, strconv.Itoa(pagesize))
	return req
}

func (req *Request) SetStart(start int) *Request {
	req.Set(qpStart, strconv.Itoa(start))
	return req
}

func (req Request) String() string {
	return req.Encode()
}

func (req *Request) SetEndpoint(ep string) (*Request, error) {
	switch ep {
	case srv.EPCustomSearch, EPCustomSearch:
		req.RequestProto.SetEndpoint(EPCustomSearch)
	case srv.EPSiteRestricted, EPSiteRestricted:
		req.RequestProto.SetEndpoint(EPSiteRestricted)
	default:
		return nil, client.ErrUnknownEndpoint
	}
	return req, nil
}

func (req *Request) ToHttpRequest() (*http.Request, error) {
	httpReq, err := http.NewRequest(API_METHOD, API_URL, nil)
	if err != nil {
		return nil, err
	}

	p, err := req.Params.Clone()
	if err != nil {
		return nil, err
	}
	p.Add(qpSearchEngineId, req.SearchEngineId)
	p.Add(qpAPIKey, req.APIKey())

	if req.PageSize != 10 {
		p.Add(qpPageSize, strconv.Itoa(req.PageSize))
	}

	if req.Start != 1 {
		p.Add(qpStart, strconv.Itoa(req.Start))
	}

	httpReq.URL.RawQuery = p.Encode()
	return httpReq, nil
}

func (req Request) ToPreviewCache(uid uuid.UUID) (cKey string, c *api.PreviewCache) {
	other := map[string]string{}
	other[qpSearchEngineId.String()] = req.SearchEngineId
	other[qpPageSize.String()] = strconv.Itoa(req.PageSize)

	return req.RequestProto.ToPreviewCache(
		uid, api.IntNextPageToken(req.Start+req.PageSize), other)
}

func RequestFromPreviewCache(cq api.CacheQuery) (api.Request, error) {
	if cq.NextPage.Equal(api.IntLastPageToken) {
		// last page
		return nil, api.ErrNotNextPage
	}

	engId, ok := cq.Other[qpSearchEngineId.String()]
	if !ok {
		return nil, ErrRequiredFieldMissing
	}

	req, err := NewRequest(cq.API.Key, engId)
	if err != nil {
		return nil, err
	}

	req.Values, err = url.ParseQuery(cq.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("error while parsing raw query: %w", err)

	}

	req.PageSize = DEFAULT_PAGE_SIZE
	ps, ok := cq.Other[qpPageSize.String()]
	if ok && ps != "10" {
		// default value is 10, so we ignore it if it's not 10.
		if i, err := convert.StrTo(ps).Int(); err == nil {
			req.PageSize = i
		}
	}

	req.Start = int(cq.NextPage.(api.IntNextPageToken))
	return req, nil
}
