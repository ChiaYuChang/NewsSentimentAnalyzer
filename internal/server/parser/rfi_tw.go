package parser

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// RFIParser is a parser for Radio France Internationale news (https://www.rfi.fr/).
type RFIParser struct {
	JsonLDParser
	MetaParser
	multiSpaceRe *regexp.Regexp
	parseFunc    func(q *Query) *Query
}

func NewRFIParser() *RFIParser {
	p := &RFIParser{
		multiSpaceRe: regexp.MustCompile("\\s{2,}"),
	}
	p.parseFunc = BuildParseFunc(p, nil)
	p.JsonLDParser = JsonLDParser{
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.GUID = p.ToGUID(n.Link)
		},
	}

	p.MetaParser = MetaParser{}
	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["category"] = "meta[property='article:section']"
	p.MetaSelector["pubDate"] = "meta[property='article:published_time']"
	p.MetaSelector["modDate"] = "meta[property='article:modified_time']"
	p.MetaSelector["keywords"] = "meta[name='keywords']"

	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		n.Title = meta.Title
		n.Category = meta.Category
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

func (p RFIParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p RFIParser) Domain() []string {
	return []string{"www.rfi.fr"}
}

func (p RFIParser) ToGUID(u *url.URL) string {
	h := md5.New()
	h.Write([]byte(path.Base(u.Path)))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)[:15])
}

func (p RFIParser) ParseBody(item *News, body *goquery.Selection) error {
	article := body.Find("main article").First()
	if author, ok := article.Find("div.m-from-author a").Attr("title"); ok {
		item.Author = append(item.Author, author)
	}

	article.Find("div.t-content__body > p").Each(func(i int, s *goquery.Selection) {
		if c := strings.TrimSpace(s.Text()); c != "" {
			item.Content = append(item.Content, p.contentCleaner(c))
		}
	})

	var toGUIDErr error
	article.Find("div.t-content__tags ul.m-tags-list li.m-tags-list__tag a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, err := url.Parse(s.AttrOr("href", "nill")); err == nil {
				item.RelatedGUID = append(item.RelatedGUID, p.ToGUID(href))
			} else {
				toGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
				return false
			}
			return true
		})
	return toGUIDErr
}

func (p RFIParser) contentCleaner(c string) string {
	c = strings.ReplaceAll(c, "\n", "")
	c = p.multiSpaceRe.ReplaceAllString(c, " ")
	return strings.TrimSpace(c)
}
