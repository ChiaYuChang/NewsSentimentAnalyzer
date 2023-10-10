package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// LTNParser is a parser for Liberty Times Net news (https://news.ltn.com.tw/).
// It implements the Parser interface.
type LTNParser struct {
	JsonLDParser
	MetaParser
	mutliSpaceRe  *regexp.Regexp
	getAuthorRe   *regexp.Regexp
	parseFunc     func(q *Query) *Query
	contentParser map[string]ContentParsingFunc
}

func NewLTNParser() *LTNParser {
	p := &LTNParser{
		contentParser: map[string]ContentParsingFunc{},
		mutliSpaceRe:  regexp.MustCompile("\\s{2,}"),
		getAuthorRe:   regexp.MustCompile("^〔(.+)／.{1,10}〕"),
	}
	p.RegisterSubParser("news.ltn.com.tw", p.parseNewsBody)
	p.RegisterSubParser("sports.ltn.com.tw", p.parseSportsNewsBody)
	p.parseFunc = BuildParseFunc(p, map[string]string{
		"head": "html", // bug while parsing head "head": "head meta" get nothing
	})

	p.JsonLDParser = JsonLDParser{
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.GUID = p.ToGUID(n.Link)
		},
	}
	ms := map[string]string{}
	ms["link"] = "link[rel='canonical']"
	ms["description"] = "meta[name='description'][itemprop='description']"
	ms["category"] = "meta[name='section'][property='article:section'][itemprop='articleSection']"
	ms["pubDate"] = "meta[name='pubdate'][property='article:published_time'][itemprop='datePublished']"
	ms["keywords"] = "meta[name='news_keywords'][itemprop='keywords']"

	p.MetaParser = MetaParser{
		MetaSelector: ms,
		PostAssignFunc: func(n *News, meta *Meta) error {
			title := strings.Split(meta.Title, " - ")
			if len(title) > 1 {
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
				n.PubDate = time.Time{}
			}
			return nil
		},
	}
	return p
}

func (p LTNParser) RegisterSubParser(domain string, spFunc ContentParsingFunc) {
	p.contentParser[domain] = spFunc
}

func (p LTNParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p LTNParser) ToGUID(href *url.URL) string {
	return strings.ReplaceAll(strings.TrimLeft(href.Path, "/"), "/", "-")
}

func (p LTNParser) Domain() []string {
	keys := make([]string, 0, len(p.contentParser))
	for k := range p.contentParser {
		keys = append(keys, k)
	}
	return keys
}

func (p LTNParser) ParseBody(n *News, body *goquery.Selection) error {
	if sp, ok := p.contentParser[n.Link.Host]; ok {
		sp(n, body)
	} else {
		p.parseNewsBody(n, body)
	}

	var toGUIDErr error
	body.Find("div.content div.related[data-desc='相關新聞'] a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if _, ok := s.Attr("data-desc"); ok {
				if href, ok := s.Attr("href"); ok {
					if u, err := url.Parse(href); err == nil {
						n.RelatedGUID = append(n.RelatedGUID, p.ToGUID(u))
					} else {
						toGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
						return false
					}
				}
			}
			return true
		})

	n.Author = p.getAuthor(n.Content[0])
	return toGUIDErr
}

func (p LTNParser) parseNewsBody(item *News, body *goquery.Selection) {
	body.Find("div.content div.whitecon[itemprop='articleBody'] div.text > p").Each(func(i int, s *goquery.Selection) {
		if _, ok := s.Attr("class"); !ok {
			item.Content = append(item.Content, p.cleanContent(s.Text()))
		}
	})
}

func (p LTNParser) parseSportsNewsBody(item *News, body *goquery.Selection) {
	body.Find("div.content div.whitecon[data-desc='內文'] div.text > p").Each(func(i int, s *goquery.Selection) {
		if len(s.Children().Nodes) == 0 {
			if _, ok := s.Attr("class"); !ok && s.Text() != "" {
				item.Content = append(item.Content, p.cleanContent(s.Text()))
			}
		}
	})
}

func (p LTNParser) getAuthor(s string) []string {
	authors := strings.TrimLeft(p.getAuthorRe.FindStringSubmatch(s)[1], "記者")
	return strings.Split(authors, "、")
}

func (p LTNParser) cleanContent(c string) string {
	c = strings.ReplaceAll(c, "\n", "")
	c = p.mutliSpaceRe.ReplaceAllString(c, " ")
	return strings.TrimSpace(c)
}
