package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NYTimesParser struct {
	MetaParser
	multiSpaceRe  *regexp.Regexp
	getCategoryRe *regexp.Regexp
	parseFunc     func(q *Query) *Query
}

func NewNYTimesParser() *NYTimesParser {
	p := &NYTimesParser{
		multiSpaceRe:  regexp.MustCompile("\\s{2,}"),
		getCategoryRe: regexp.MustCompile("cn.nytimes.com/(\\w+)/.+"),
	}

	p.parseFunc = BuildParseFunc(p, nil)

	p.MetaParser = MetaParser{
		MetaSelector: map[string]string{
			"link":        "link[rel='canonical']",
			"description": "meta[name='description']",
			"pubDate":     "meta[property='article:published_time']",
			"author":      "meta#byline[name='byline']",
		},
		PostAssignFunc: func(n *News, meta *Meta) error {
			n.Title = strings.TrimSuffix(meta.Title, " - 紐約時報中文網")

			n.Description = meta.Description

			if u, err := url.Parse(meta.Link); err == nil {
				n.Link = u
			} else {
				return fmt.Errorf("error while url.Parse: %w", err)
			}

			n.Author = NewCSL(meta.Author)
			sm := p.getCategoryRe.FindStringSubmatch(meta.Link)
			if len(sm) == 2 {
				n.Category = sm[1]
			}
			n.GUID = p.ToGUID(n.Link)

			if pubDate, err := time.Parse(time.RFC3339, meta.PubDate); err == nil {
				n.PubDate = pubDate.UTC()
			} else {
				fmt.Printf("error while parsing time %s: %v\n", meta.PubDate, err)
				n.PubDate = time.Time{}
			}
			return nil
		},
	}
	return p
}

func (p NYTimesParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p NYTimesParser) ToGUID(u *url.URL) string {
	path := u.Path
	path = strings.TrimPrefix(path, "/cn.nytimes.com/")
	path = strings.TrimSuffix(path, "/zh-hant/")
	path = strings.Trim(path, "/")

	return strings.ReplaceAll(path, "/", "-")
}

func (p NYTimesParser) Domain() []string {
	return []string{"cn.nytimes.com"}
}

func (p NYTimesParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("main.main div.article-area article.article-content section.article-body div.article-body-item div.article-paragraph").
		Each(func(i int, s *goquery.Selection) {
			if s.Find("figcaption").Length() > 0 {
				return
			}
			content := p.cleanContent(s.Text())
			if content != "" {
				n.Content = append(n.Content, content)
			}
		})

	var rGUIDErr error
	body.Find("main.main div.article-area div.article-footer div.related-cont ul.refer-list li.article-refer div.refer-list-item a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if u, err := url.Parse(href); err == nil {
					n.RelatedGUID = append(n.RelatedGUID, p.ToGUID(u))
				} else {
					rGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
					return false
				}
			}
			return true
		})
	return rGUIDErr
}

func (p NYTimesParser) ParseJsonLD(n *News, s *goquery.Selection) error {
	// NY Times news doesn't have json linked data
	return nil
}

func (p NYTimesParser) cleanContent(c string) string {
	content := strings.TrimSpace(c)
	content = p.multiSpaceRe.ReplaceAllString(content, " ")
	content = strings.ReplaceAll(content, "\n", "")
	return content
}
