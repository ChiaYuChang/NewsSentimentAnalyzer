package parser

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/PuerkitoBio/goquery"
)

// PTSParser is a parser for Public Television Service news (https://news.pts.org.tw/).
// It implements the Parser interface.
type PTSParser struct {
	JsonLDParser
	MetaParser
	multiSpaceRe *regexp.Regexp
	parseFunc    func(q *Query) *Query
}

func NewPTSParser() *PTSParser {
	p := &PTSParser{
		multiSpaceRe: regexp.MustCompile("\\s{2,}"),
	}
	p.parseFunc = BuildParseFunc(p, nil)

	p.JsonLDParser = JsonLDParser{
		PostAssignFunc: func(n *News, jld *JsonLD) {
			if ss := strings.Split(n.Title, " ｜ "); len(ss) == 2 {
				n.Title = ss[0]
			}

			aSet := collection.NewSet[string]()
			for _, author := range jld.Author {
				if strings.ToLower(author.Type) == "person" {
					for _, a := range strings.Split(strings.Split(author.Name, "／")[0], " ") {
						aSet.Add(a)
					}
				}
			}
			n.Author = aSet.Key()
		},
	}

	p.MetaParser = MetaParser{}
	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["pubDate"] = "meta[property='pubdate']"
	p.MetaSelector["modDate"] = "meta[property='moddate']"
	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		title := strings.Split(meta.Title, " ｜ ")
		if len(title) == 2 {
			n.Title = title[0]
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
			n.PubDate = time.Time{}
		}
		return nil
	}
	return p
}

func (p PTSParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p PTSParser) ToGUID(href *url.URL) string {
	return path.Base(href.Path)
}

func (p PTSParser) Domain() []string {
	return []string{"news.pts.org.tw"}
}

func (p PTSParser) ParseBody(n *News, body *goquery.Selection) error {
	n.Category = body.Find("article nav[aria-label='breadcrumb'] ol.breadcrumb li.breadcrumb-item a").Get(1).FirstChild.Data

	paragraphs := body.Find("article.row div.post-article")
	n.Content = []string{}
	paragraphs.Find("p, h2, table tbody tr").Each(func(i int, s *goquery.Selection) {
		n.Content = append(n.Content, p.cleanContent(s.Text()))
	})

	rnSet := collection.NewSet[string]()
	var toGUIDErr error
	body.Find("div.issue-cords-1 div.d-flex div.issue-item-title h6 a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok && strings.Contains(href, "news.pts.org.tw/article") {
				if u, err := url.Parse(href); err == nil {
					rnSet.Add(p.ToGUID(u))
				} else {
					toGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
					return false
				}
			}
			return true
		})

	if toGUIDErr != nil {
		return toGUIDErr
	}

	body.Find("div.relative-news-list-content h4.m-0 a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok && strings.Contains(href, "news.pts.org.tw/article/") {
				if u, err := url.Parse(href); err == nil {
					rnSet.Add(p.ToGUID(u))
				} else {
					toGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
					return false
				}
			}
			return true
		})
	n.RelatedGUID = rnSet.Key()
	return toGUIDErr
}

func (p PTSParser) cleanContent(c string) string {
	c = p.multiSpaceRe.ReplaceAllString(c, " ")
	c = strings.TrimSpace(strings.ReplaceAll(c, "\n", " "))
	return c
}
