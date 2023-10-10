package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/PuerkitoBio/goquery"
)

type selectorEl struct {
	Tag   string
	Id    string
	Class []string
	Attr  map[string]string
}

func (sn selectorEl) String() string {
	s := sn.Tag
	attrs := make([]string, 0, len(sn.Attr))
	for k, v := range sn.Attr {
		attrs = append(attrs, fmt.Sprintf("[%s='%s']", k, v))
	}

	if sn.Id != "" {
		s += "#" + sn.Id
	}

	if len(sn.Class) > 0 {
		s += "." + strings.Join(sn.Class, ".")
	}

	if len(attrs) > 0 {
		s += strings.Join(attrs, "")
	}

	return s
}

// builder for goquery selector
type SelectorBuilder []selectorEl

// new builder for goquery selector
func NewBuilder() *SelectorBuilder {
	return &SelectorBuilder{}
}

// return goquery selector string
func (b SelectorBuilder) Build() string {
	ss := make([]string, len(b))
	for i, s := range b {
		ss[i] = s.String()
	}
	return strings.Join(ss, " ")
}

// append new selector string
func (b *SelectorBuilder) Append(tag string, class []string, attr map[string]string) *SelectorBuilder {
	(*b) = append((*b), selectorEl{Tag: tag, Class: class, Attr: attr})
	return b
}

type Meta struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Author      string `json:"author,omitempty"`
	Category    string `json:"category,omitempty"`
	PubDate     string `json:"pubDate,omitempty"`
	ModDate     string `json:"modDate,omitempty"`
	Tag         string `json:"tag"`
}

func (m Meta) String() string {
	return strings.Join([]string{
		fmt.Sprintf("Title       : %s", m.Title),
		fmt.Sprintf("Link        : %s", m.Link),
		fmt.Sprintf("Description : %s", m.Description),
		fmt.Sprintf("Language    : %s", m.Language),
		fmt.Sprintf("Author      : %s", m.Author),
		fmt.Sprintf("Category    : %s", m.Category),
		fmt.Sprintf("PubDate     : %s", m.PubDate),
		fmt.Sprintf("ModDate     : %s", m.ModDate),
		fmt.Sprintf("Tag         : %s", m.Tag),
	}, "\n")
}

// meta selector is used to parse meta data from html
// key is the name of meta data, value is the selector string
// if value is empty, it will not be parsed
type MetaSelector map[string]string

// default meta selector
var defaultMetaSelector = MetaSelector{
	"link": NewBuilder().
		Append("link", nil, map[string]string{
			"rel": "canonical",
		}).Build(),
	"description": NewBuilder().
		Append("meta", nil, map[string]string{
			"name":     "description",
			"itemprop": "description",
		}).Build(),
	"language": NewBuilder().
		Append("meta", nil, map[string]string{
			"http-equiv": "content-language",
		}).Build(),
	"author": NewBuilder().
		Append("meta", nil, map[string]string{
			"name":     "author",
			"itemprop": "author",
		}).Build(),
	"category": NewBuilder().
		Append("meta", nil, map[string]string{
			"name":     "section",
			"property": "article:section",
			"itemprop": "articleSection",
		}).Build(),
	"pubDate": NewBuilder().
		Append("meta", nil, map[string]string{
			"name":     "pubdate",
			"property": "article:published_time",
			"itemprop": "datePublished",
		}).Build(),
	"modDate": NewBuilder().
		Append("meta", nil, map[string]string{
			"name":     "lastmod",
			"property": "article:modified_time",
			"itemprop": "dateModified",
		}).Build(),
	"keywords": NewBuilder().
		Append("meta", nil, map[string]string{
			"itemprop": "keywords",
		}).Build(),
}

// GetDefaultMetaSelectorCopy return a copy of defaultMetaSelector
// if you want to modify defaultMetaSelector, use this function to get a copy
func GetDefaultMetaSelectorCopy() (MetaSelector, error) {
	m := MetaSelector{}
	data, err := json.Marshal(defaultMetaSelector)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &m)
	return m, err
}

// parse meta, if MetaSelector is nil, use defaultMetaSelector
// if you want to modify defaultMetaSelector, use GetDefaultMetaSelectorCopy to get a copy and modify it
func ParseMeta(selector MetaSelector, head *goquery.Selection) *Meta {
	if selector == nil {
		selector = defaultMetaSelector
	}
	meta := Meta{}
	meta.Title = head.Find("title").First().Text()

	items := []struct {
		pointer  *string
		selector string
		attr     string
	}{
		{&meta.Link, "link", "href"},
		{&meta.Description, "description", "content"},
		{&meta.Language, "language", "content"},
		{&meta.Author, "author", "content"},
		{&meta.Category, "category", "content"},
		{&meta.PubDate, "pubDate", "content"},
		{&meta.ModDate, "modDate", "content"},
		{&meta.Tag, "keywords", "content"},
	}

	for _, item := range items {
		if slctr, ok := selector[item.selector]; ok && slctr != "" {
			ss := collection.NewSet[string]()
			head.Find(slctr).Each(func(i int, s *goquery.Selection) {
				if val, ok := s.Attr(item.attr); ok {
					ss.Add(val)
				}
			})
			(*item.pointer) = strings.Join(ss.Key(), ",")
		}
	}
	return &meta
}

type MetaParser struct {
	MetaSelector
	PostAssignFunc func(n *News, meta *Meta) error
}

func (mp MetaParser) ParseMeta(n *News, head *goquery.Selection) error {
	meta := ParseMeta(mp.MetaSelector, head)

	if mp.PostAssignFunc != nil {
		mp.PostAssignFunc(n, meta)
	}
	return nil
}
