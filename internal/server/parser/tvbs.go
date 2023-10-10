package parser

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// TVBSParser is a parser for TVBS news (https://news.tvbs.com.tw/).
type TVBSParser struct {
	JsonLDParser
	MetaParser
	parseFunc func(q *Query) *Query
}

func NewTVBSParser() *TVBSParser {
	p := &TVBSParser{}
	p.parseFunc = BuildParseFunc(p, map[string]string{
		"head": "head",
		"body": "div#news_detail_div",
	})

	p.JsonLDParser = JsonLDParser{
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.Title = strings.Split(jld.Headline, "│")[0]
			for _, a := range jld.Author {
				if a.Type == "Person" {
					n.Author = append(n.Author, a.Name)
				}
			}
		},
	}

	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["language"] = "meta[http-equiv='content-language']"
	p.MetaSelector["category"] = "meta[name='section'][property='article:section'][itemprop='articleSection']"
	p.MetaSelector["pubDate"] = "meta[name='pubdate'][property='article:published_time']"
	p.MetaSelector["modDate"] = "meta[name='moddate'][property='article:modified_time']"
	p.MetaSelector["keywords"] = "meta[name='keywords'][itemprop='keywords']"
	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		n.Title = strings.Split(meta.Title, "│")[0]
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

func (p TVBSParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p TVBSParser) ToGUID(href *url.URL) string {
	return strings.Replace(href.Path[1:], "/", "-", 1)
}

func (p TVBSParser) Domain() []string {
	return []string{"news.tvbs.com.tw"}
}

func (p TVBSParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("div[align='center']").Remove()
	body.Find("div.img").Remove()
	body.Find("div.guangxuan").Remove()

	n.Content = append(n.Content, p.cleanContent(body.Text())...)
	return nil
}

func (p TVBSParser) cleanContent(c string) []string {
	c = strings.ReplaceAll(c, "\r\n", "\n")
	c = strings.ReplaceAll(c, "\u00a0", "\n")
	ss := []string{}
	for _, s := range strings.Split(c, "\n") {
		s = strings.TrimSpace(s)
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
