package googlecse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"google.golang.org/api/googleapi"
)

const (
	qpKeyword                      = "q"
	qpEnableChineseSearch          = "c2coff"
	qpSearchEngineId               = "cx"
	qpGeoLocation                  = "gl"
	qpSaveLevel                    = "safe"
	qpLanguage                     = "lr"
	qpEnableDuplicateContentFilter = "filter"
	qpDateRestrict                 = "dateRestrict"
	qpSiteSearch                   = "siteSearch"
	qpSiteSearchFilter             = "siteSearchFilter"
	qpExactTerms                   = "exactTerms"
	qpExcludeTerms                 = "excludeTerms"
	qpPageSize                     = "num"
	qpStart                        = "start"
)

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
	// APIKey                 string                    `json:"apikey,omitempty"`
	// Query                  string                    `json:"query"`
	// ChineseSearch          cseChineseSearch          `json:"chinese_search,omitempty"`
	// GeoLocation            string                    `json:"geo_location,omitempty"`
	// SafeLevel              cseSafeLevel              `json:"safe_level,omitempty"`
	// Language               string                    `json:"language,omitempty"`
	// DuplicateContentFilter cseDuplicateContentFilter `json:"duplicate_content_filter,omitempty"`
	// DateRestrict           string                    `json:"date_restrict,omitempty"`
	// SiteSearch             string                    `json:"site_search,omitempty"`
	// SiteSearchFilter       cseSiteSearchFilter       `json:"site_search_filter,omitempty"`
	// ExactTerms             string                    `json:"exact_terms,omitempty"`
	// ExcludeTerms           string                    `json:"exclude_terms,omitempty"`
	// PageSize               int64                     `json:"page_size,omitempty"`
	// More                   [][2]string               `json:"more,omitempty"`
	// Start                  uint32                    `json:"start,omit"`
	CallOpts []googleapi.CallOption `json:"-"`
	// params                 api.ParamsMap             `json:"-"`
}

func NewSearchRequest(ctx context.Context, apikey string, engId string,
	opt ...googleapi.CallOption) (*Request, error) {

	if apikey == "" {
		return nil, ErrRequiredFieldMissing
	}

	if engId == "" {
		return nil, ErrRequiredFieldMissing
	}

	req := api.NewRequestProtoType("")
	req.SetApiKey(apikey)
	req.Add(qpEnableChineseSearch, EnableChineseSearch.ToParam())
	req.Add(qpSaveLevel, EnableSafeSearch.ToParam())
	req.Add(qpEnableDuplicateContentFilter, EnableDuplicateContentFilter.ToParam())
	req.Set(qpPageSize, "10")
	req.Set(qpStart, "1")
	return &Request{
		RequestProto:   req,
		SearchEngineId: engId,
		CallOpts:       opt,
	}, nil
}

func (r *Request) SetKeyword(keyword string) *Request {
	r.Set(qpKeyword, keyword)
	return r
}

func (r *Request) SetChineseSearch(c2off cseChineseSearch) *Request {
	r.Set(qpEnableChineseSearch, c2off.ToParam())
	return r
}

func (r *Request) SetGeoLocation(gl string) *Request {
	r.Set(qpGeoLocation, gl)
	return r
}

func (r *Request) SetSafeLevel(safe cseSafeLevel) *Request {
	r.Set(qpSaveLevel, safe.ToParam())
	return r
}

func (r *Request) SetLanguage(lr string) *Request {
	r.Set(qpLanguage, lr)
	return r
}

func (r *Request) SetEngineId(engId string) *Request {
	r.SearchEngineId = engId
	return r
}

func (r *Request) SetDuplicateContentFilter(filter cseDuplicateContentFilter) *Request {
	r.Set(qpEnableDuplicateContentFilter, filter.ToParam())
	return r
}

func (r *Request) SetDateRestict(dateRestict string) *Request {
	r.Set(qpDateRestrict, dateRestict)
	return r
}

func (r *Request) SetSiteSearch(site string, filter cseSiteSearchFilter) *Request {
	r.Set(qpSiteSearch, site)
	r.Set(qpSiteSearchFilter, filter.ToParam())
	return r
}

func (r *Request) WithExactTerm(term string) *Request {
	r.Add(qpExactTerms, term)
	return r
}

func (r *Request) WithExcludeTerm(term string) *Request {
	r.Add(qpExcludeTerms, term)
	return r
}

func (r *Request) WithPageSize(pagesize int) *Request {
	r.Set(qpPageSize, strconv.Itoa(pagesize))
	return r
}

func (r *Request) SetStart(start int) *Request {
	r.Set(qpStart, strconv.Itoa(start))
	return r
}

func (r *Request) String() string {
	return r.Encode()
}

func (q *Request) Endpoint() string {
	return EPCustomSearch
}

func (q *Request) ToHttpRequest() (*http.Request, error) {
	b, err := json.Marshal(q.Values)
	if err != nil {
		return nil, err
	}

	return q.RequestProto.ToHTTPRequest(API_URL, API_METHOD, bytes.NewBuffer(b))
}

var ErrRequiredFieldMissing = errors.New("required field is missing")
