package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/PuerkitoBio/goquery"
)

// CNAParser is a parser for CNA news (https://www.cna.com.tw/).
// It implements the Parser interface.
type CNAParser struct {
	JsonLDParser
	MetaParser
	toGUIDRe    *regexp.Regexp
	getAuthorRe *regexp.Regexp
	jsonldRe    *regexp.Regexp
	parseFunc   func(q *Query) *Query
}

func NewCNAParser() *CNAParser {
	p := &CNAParser{
		toGUIDRe:    regexp.MustCompile("news/(\\w{3,4})/(\\d+).aspx"),
		getAuthorRe: regexp.MustCompile("。（(.+?)）\\d{6,7}"),
		jsonldRe:    regexp.MustCompile(`"headline":.*?,"(about)":.*?,"url":`),
	}
	p.parseFunc = BuildParseFunc(p, nil)

	p.JsonLDParser = JsonLDParser{
		DataPreprocessFunc: func(data []byte) []byte {
			repl := bytes.Replace(
				p.jsonldRe.Find(data),
				[]byte("about"),
				[]byte("description"), 1)
			return p.jsonldRe.ReplaceAll(data, repl)
		},
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.Author = append(n.Author, jld.Author[0].Name)
			n.GUID = p.ToGUID(n.Link)
			for i, t := range n.Tag {
				if t == "NewsArticle" {
					n.Tag[i] = ""
				}
			}
		},
	}

	p.MetaParser = MetaParser{}
	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["category"] = "meta[name='section'][property='article:section'][itemprop='articleSection']"
	p.MetaSelector["pubDate"] = "meta[property='article:published_time']"
	p.MetaSelector["modDate"] = "meta[property='article:modified_time']"
	p.MetaSelector["keywords"] = "meta[property='article:tag']"

	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		title := strings.Split(meta.Title, " | ")
		if len(title) == 3 {
			n.Title = title[0]
			n.Category = title[1]
		}

		n.Description = meta.Description
		n.Tag = NewCSL(meta.Tag)

		if u, err := url.Parse(meta.Link); err == nil {
			n.Link = u
		} else {
			return fmt.Errorf("error while url.Parse: %w", err)
		}

		n.GUID = p.ToGUID(n.Link)

		if pubDate, err := time.Parse(time.RFC3339, meta.PubDate); err == nil {
			n.PubDate = pubDate.UTC()
		} else {
			return fmt.Errorf("error while time.Parse: %w", err)
		}
		return nil
	}
	return p
}

func (p CNAParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p CNAParser) ToGUID(href *url.URL) string {
	return strings.Join(p.toGUIDRe.FindStringSubmatch(href.Path)[1:], "-")
}

func (p CNAParser) Domain() []string {
	return []string{"www.cna.com.tw"}
}

func (p CNAParser) ParseBody(n *News, body *goquery.Selection) error {
	tagSet := collection.NewSet[string]()
	toGuidErr := Errors{}

	body.Find("div.centralContent").Each(func(i int, s *goquery.Selection) {
		paragraph := s.Find("div.paragraph").First()
		paragraph.Find("p").Each(func(i int, s *goquery.Selection) {
			n.Content = append(n.Content, strings.TrimSpace(s.Text()))
		})

		paragraph.Find("div.keywordTag a").Each(func(i int, s *goquery.Selection) {
			tag := strings.TrimPrefix(s.Text(), "#")
			tagSet.Add(tag)
		})

		paragraph.Find("div.moreArticle a.moreArticle-link").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok && strings.Contains(href, "www.cna.com.tw/news") {
				if u, err := url.Parse(href); err == nil {
					n.RelatedGUID = append(n.RelatedGUID, p.ToGUID(u))
				} else {
					toGuidErr.Add(href, fmt.Errorf("error while parsing %s: %w", href, err))
				}
			}
		})
	})

	n.Tag = tagSet.Key()
	n.Author = p.getAuthor(n.Content[len(n.Content)-1])
	return toGuidErr
}

func (p CNAParser) getAuthor(s string) []string {
	ss := strings.Split(p.getAuthorRe.FindStringSubmatch(s)[1], "/")
	authors := make([]string, 0, len(ss))
	for _, a := range ss {
		if strings.Contains(a, "：") {
			authors = append(authors, strings.TrimSpace(strings.Split(a, "：")[1]))
		} else {
			authors = append(authors, strings.TrimSpace(a))
		}
	}
	return authors
}
