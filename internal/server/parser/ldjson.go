package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Comma Separated List
type CSL []string

// Check if string is unicode escaped
func IsUnicodeEscaped(s string) bool {
	return strings.Contains(s, "\\u")
}

// Create new comma separated list from string (eg. "a,b,c")
func NewCSL(s string) CSL {
	ss := strings.Split(s, ",")
	csl := make([]string, 0, len(ss))
	for _, s := range ss {
		if s = strings.TrimSpace(s); s != "" {
			if IsUnicodeEscaped(s) {
				if uqs, err := strconv.Unquote(`"` + s + `"`); err == nil {
					s = uqs
				}
			}
			csl = append(csl, s)
		}
	}
	return csl
}

func (csl CSL) MarshalJSON() ([]byte, error) {
	// s := strings.Join(csl, ",")
	// return []byte(fmt.Sprintf(`"%s"`, s)), nil
	return []byte(strconv.Quote(strings.Join(csl, ","))), nil
}

func (csl *CSL) UnmarshalJSON(b []byte) error {
	s := strings.ReplaceAll(string(b), "\n", "")
	if s[0] == '[' && s[len(s)-1] == ']' {
		var ss []string
		re := regexp.MustCompile("\"([^,]+?)\"")
		for _, t := range re.FindAllStringSubmatch(s, -1)[1:] {
			ss = append(ss, string(t[1]))
		}
		(*csl) = CSL(ss)
	} else {
		s = strings.Trim(s, "\"")
		(*csl) = NewCSL(s)
	}
	return nil
}

// json linked data
type JsonLD struct {
	Type           string           `json:"@type"`
	Headline       string           `json:"headline"`
	Description    string           `json:"description,omitempty"`
	URL            *url.URL         `json:"url"`
	ArticleSection string           `json:"articleSection"`
	Author         JsonLDObjectList `json:"author"`
	About          JsonLDObjectList `json:"about,omitempty"`
	Keywords       CSL              `json:"keywords"`
	PublishedAt    time.Time        `json:"datePublished"`
	UpdatedAt      time.Time        `json:"dateModified"`
}

// json linked data @type = Person
type JsonLDObject struct {
	Id      string   `json:"@id,omitempty"`
	Type    string   `json:"@type,omitempty"`
	Name    string   `json:"name,omitempty"`
	SameAs  []string `json:"sameAs,omitempty"`
	AltName string   `json:"alternateName,omitempty"`
	Width   int      `json:"width,omitempty"`
	Height  int      `json:"height,omitempty"`
}

type JsonLDObjectList []JsonLDObject

func (jsonld *JsonLDObjectList) UnmarshalJSON(data []byte) error {
	objs := []JsonLDObject{}
	data = bytes.TrimSpace(data)
	if data[0] == '{' && data[len(data)-1] == '}' {
		tmp := make([]byte, len(data)+2)
		tmp[0], tmp[len(tmp)-1] = '[', ']'
		copy(tmp[1:len(tmp)-1], data)
		data = tmp
	}

	err := json.Unmarshal(data, &objs)
	(*jsonld) = JsonLDObjectList(objs)

	return err
}

func (jsonld *JsonLD) UnmarshalJSON(data []byte) error {
	type InnerJsonLD JsonLD
	tmp := struct {
		*InnerJsonLD
		URL         string `json:"url"`
		PublishedAt string `json:"datePublished"`
		UpdatedAt   string `json:"dateModified"`
	}{
		InnerJsonLD: (*InnerJsonLD)(jsonld),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if u, err := url.Parse(tmp.URL); err != nil {
		return fmt.Errorf("error while url.Parse: %w", err)
	} else {
		tmp.InnerJsonLD.URL = u
	}

	if tmp.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, tmp.PublishedAt); err != nil {
			return fmt.Errorf("error while parsing time %s: %w", tmp.PublishedAt, err)
		} else {
			tmp.InnerJsonLD.PublishedAt = t.UTC()
		}
	}

	if tmp.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, tmp.UpdatedAt); err != nil {
			return fmt.Errorf("error while parsing time %s: %w", tmp.UpdatedAt, err)
		} else {
			tmp.InnerJsonLD.UpdatedAt = t.UTC()
		}
	}
	return nil
}

func (jsonld JsonLD) MarshalJSON() ([]byte, error) {
	type InnerJsonLD JsonLD
	tmp := struct {
		InnerJsonLD
		URL string `json:"url"`
	}{
		InnerJsonLD: (InnerJsonLD)(jsonld),
	}

	tmp.URL = tmp.InnerJsonLD.URL.String()
	return json.Marshal(tmp)
}

func (jsonld JsonLD) String() string {
	data, _ := json.MarshalIndent(jsonld, "", "\t")
	return string(data)
}

type JsonLDList struct {
	JsonLD      []*JsonLD
	TypeToIndex map[string]int
}

func (jsonldLst *JsonLDList) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if data[0] != '[' && data[len(data)-1] != ']' {
		data = bytes.Join([][]byte{{'['}, data, {']'}}, []byte{})
	}

	jld := []*JsonLD{}
	if err := json.Unmarshal(data, &jld); err != nil {
		return err
	}

	jsonldLst.JsonLD = jld
	jsonldLst.TypeToIndex = make(map[string]int, len(jld))
	for i := range jld {
		jsonldLst.TypeToIndex[jld[i].Type] = i
	}
	return nil
}

func (jsonldLst JsonLDList) MarshalJSON() ([]byte, error) {
	if len(jsonldLst.JsonLD) == 1 {
		return json.Marshal(jsonldLst.JsonLD[0])
	}
	return json.Marshal(jsonldLst.JsonLD)
}

func (jsonldLst JsonLDList) GetByIndex(i int) *JsonLD {
	return jsonldLst.JsonLD[i]
}

func (jsonldLst JsonLDList) GetByType(t string) *JsonLD {
	return jsonldLst.JsonLD[jsonldLst.TypeToIndex[t]]
}

func (jsonldLst JsonLDList) String() string {
	data, _ := json.MarshalIndent(jsonldLst, "", "    ")
	return string(data)
}

func (jsonldLst JsonLDList) Len() int {
	return len(jsonldLst.JsonLD)
}

func (jsonldLst *JsonLDList) Append(jld ...*JsonLD) {
	l := jsonldLst.Len()
	for i := range jld {
		jsonldLst.JsonLD = append(jsonldLst.JsonLD, jld[i])
		jsonldLst.TypeToIndex[jld[i].Type] = l + i
	}
}

func (jsonldLst1 *JsonLDList) Merge(jsonldLst2 JsonLDList) {
	jsonldLst1.Append(jsonldLst2.JsonLD...)
}

var ErrTargetJsonLDNotFound = fmt.Errorf("target jsonld not found")

type JsonLDParser struct {
	DataPreprocessFunc func(data []byte) []byte
	PostAssignFunc     func(n *News, jld *JsonLD)
}

// FindTargetJsonLD find target jsonld by @type
func (jldp JsonLDParser) FindTargetJsonLD(data []byte, target string) ([]byte, error) {
	// find target jsonld by @type
	var str, end int
	var isTarget bool

	for str = 0; str < len(data); str++ {
		if data[str] == '{' {
			lev := 1

			// search for the end of json object
			for end = str + 1; end < len(data); end++ {
				if data[end] == '{' {
					lev++
				}

				if data[end] == '}' {
					lev--
					if lev == 0 {
						break
					}
				}
			}

			isTarget = jldp.IsTargetType(data[str:end], target)
			if isTarget {
				return data[str : end+1], nil
			}
			str = end
		}
	}
	return nil, ErrTargetJsonLDNotFound
}

func (jldp JsonLDParser) IsTargetType(data []byte, target string) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == '@' && string(data[i:i+5]) == "@type" {
			tagStr, tagEnd := i, i+1
			for ; tagEnd < len(data); tagEnd++ {
				if data[tagEnd] == ',' {
					break
				}
			}

			if strings.Contains(string(data[tagStr:tagEnd]), target) {
				return true
			}
		}
	}
	return false
}

func (jldp JsonLDParser) ParseJsonLD(n *News, s *goquery.Selection) error {
	var parsingErr error
	var jld JsonLD
	s.EachWithBreak(func(i int, s *goquery.Selection) bool {
		jlddata, findTargetErr := jldp.FindTargetJsonLD([]byte(s.Text()), "NewsArticle")

		if findTargetErr != nil {
			return true
		}
		if jldp.DataPreprocessFunc != nil {
			jlddata = jldp.DataPreprocessFunc(jlddata)
		}

		if !json.Valid(jlddata) {
			parsingErr = fmt.Errorf("invalid json object")
			return false
		}

		parsingErr = json.Unmarshal(jlddata, &jld)
		return false
	})

	if parsingErr != nil {
		return parsingErr
	}

	jldp.defaultFieldAssignFunc(n, &jld)
	if jldp.PostAssignFunc != nil {
		jldp.PostAssignFunc(n, &jld)
	}
	return nil
}

func (jldp JsonLDParser) defaultFieldAssignFunc(n *News, jld *JsonLD) {
	n.Title = jld.Headline
	n.Category = jld.ArticleSection
	n.Description = jld.Description
	n.Link = jld.URL
	n.PubDate = jld.PublishedAt
	n.Tag = jld.Keywords
}

func ParseJsonLD(*News, []byte, func(href *url.URL) string) error {
	return nil
}
