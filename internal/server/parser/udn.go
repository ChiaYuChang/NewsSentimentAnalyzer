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

// UDNParser is a parser for UDN news (https://udn.com/).
// It implements the Parser interface.
type UDNParser struct {
	JsonLDParser
	toGUIDRe      *regexp.Regexp
	parseFunc     func(q *Query) *Query
	contentParser map[string]ContentParsingFunc
}

// NewUDNParser returns a UDNParser.
// It registers sub parsers for udn.com and global.udn.com.
func NewUDNParser() *UDNParser {
	p := &UDNParser{
		contentParser: map[string]ContentParsingFunc{},
		toGUIDRe:      regexp.MustCompile("[a-z]+/[a-z]+/([0-9]{4,7})/([0-9]{4,8})"),
	}
	p.RegisterSubParser("udn.com", p.parseNewsBody)
	p.RegisterSubParser("global.udn.com", p.parseGlobalStory)

	p.parseFunc = BuildParseFunc(p, nil)
	p.JsonLDParser = JsonLDParser{}
	p.PostAssignFunc = func(n *News, jld *JsonLD) {
		n.GUID = p.ToGUID(n.Link)
	}

	return p
}

// RegisterSubParser registers a sub parser for a domain.
func (p UDNParser) RegisterSubParser(domain string, spFunc ContentParsingFunc) {
	p.contentParser[domain] = spFunc
}

func (p UDNParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p UDNParser) ToGUID(herf *url.URL) string {
	if herf == nil {
		return ""
	}

	// re := regexp.MustCompile("[a-z]+/[a-z]+/([0-9]{4,7})/([0-9]{4,8})")
	sm := p.toGUIDRe.FindStringSubmatch(herf.Path)
	if len(sm) != 3 {
		return strings.ReplaceAll(herf.Path, "/", "-")
	}
	return path.Base(strings.Join([]string{sm[1], sm[2]}, "-"))
}

func (p UDNParser) Domain() []string {
	keys := make([]string, 0, len(p.contentParser))
	for k := range p.contentParser {
		keys = append(keys, k)
	}
	return keys
}

func (p UDNParser) ParseMeta(n *News, head *goquery.Selection) error {
	ms := MetaSelector{
		"link": NewBuilder().
			Append("link", nil, map[string]string{
				"rel": "canonical",
			}).Build(),
		"description": NewBuilder().
			Append("meta", nil, map[string]string{
				"name":     "description",
				"itemprop": "description",
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
		"keywords": NewBuilder().
			Append("meta", nil, map[string]string{
				"itemprop": "keywords",
			}).Build(),
	}
	fmt.Println(ms["author"])
	meta := ParseMeta(ms, head)

	title := strings.Split(meta.Title, " | ")
	if len(title) >= 3 {
		n.Title = title[0]
	}

	n.Description = meta.Description
	n.Category = meta.Category
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

func (p UDNParser) ParseBody(n *News, body *goquery.Selection) error {
	fmt.Println(n.Link.Host)
	if sp, ok := p.contentParser[n.Link.Host]; ok {
		sp(n, body)
	} else {
		fmt.Println("default handler")
		p.parseNewsBody(n, body)
	}
	return nil
}

func (p UDNParser) parseNewsBody(item *News, body *goquery.Selection) {
	item.Author = []string{}
	body.Find("section.authors span.article-content__author > a").Each(func(i int, s *goquery.Selection) {
		item.Author = append(item.Author, s.Text())
	})
	body.Find("article.article-content div.article-content__paragraph p").Each(func(i int, s *goquery.Selection) {
		if c := s.Text(); c != "" {
			item.Content = append(item.Content, strings.TrimSpace(c))
		}
	})

	rGUID := collection.NewSet[string]()
	body.Find("section.more-news div.context-box__content div.story-list__news div.story-list__text a").Each(func(i int, s *goquery.Selection) {
		if u, err := url.Parse(s.AttrOr("href", "null")); err == nil {
			rGUID.Add(p.ToGUID(u))
		}
	})
	item.RelatedGUID = rGUID.Key()
	return
}

func (p UDNParser) parseGlobalStory(n *News, body *goquery.Selection) {
	fmt.Println("use global.udn.com parser")
	body.Find("div#story_body div.story_body_content > p").Each(func(i int, s *goquery.Selection) {
		if c := s.Text(); c != "" {
			n.Content = append(n.Content, strings.TrimSpace(c))
		}
	})

	rGUID := collection.NewSet[string]()
	body.Find("div#story_also dt a[data-slotname='list_推薦閱讀']").Each(func(i int, s *goquery.Selection) {
		if u, err := url.Parse(s.AttrOr("href", "null")); err == nil {
			rGUID.Add(p.ToGUID(u))
		}
	})

	body.Find("div#story_bady_info > span").
		Each(func(i int, s *goquery.Selection) {
			if author := s.Text(); author != "" {
				n.Author = append(n.Author, strings.TrimSpace(author))
			}
		})

	n.RelatedGUID = rGUID.Key()
	return
}
