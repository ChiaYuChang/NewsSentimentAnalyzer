package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type EtTodayParser struct {
	JsonLDParser
	MetaParser
	hiddenCharRe *regexp.Regexp
	paragraphRe  *regexp.Regexp
	getAuthorRe  *regexp.Regexp
	parseFunc    func(q *Query) *Query
}

func NewEtTodayParser() *EtTodayParser {
	p := &EtTodayParser{
		hiddenCharRe: regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`),
		paragraphRe:  regexp.MustCompile("^[^▸▲]"),
		getAuthorRe:  regexp.MustCompile(`／.+報導$`),
	}
	p.parseFunc = BuildParseFunc(p, nil)
	p.JsonLDParser = JsonLDParser{
		DataPreprocessFunc: func(data []byte) []byte {
			data = p.hiddenCharRe.ReplaceAll(data, []byte(""))
			return data
		},
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.Author = nil
			n.GUID = p.ToGUID(n.Link)
		},
	}

	p.MetaParser = MetaParser{}
	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical'][itemprop='mainEntityOfPage']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["category"] = "meta[name='section'][property='article:section']"
	p.MetaSelector["pubDate"] = "meta[name='pubdate'][property='article:published_time'][itemprop='datePublished']"
	p.MetaSelector["keywords"] = "meta[name='news_keywords'][itemprop='keywords']"

	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		title := strings.Split(meta.Title, " | ")
		if len(title) >= 3 {
			n.Title = title[0]
		}

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

func (p EtTodayParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p EtTodayParser) ToGUID(href *url.URL) string {
	ss := strings.Split(strings.TrimSuffix(href.Path, ".htm"), "/")
	return strings.Join(ss[2:], "-")
}

func (p EtTodayParser) Domain() []string {
	return []string{"www.ettoday.net"}
}

func (p EtTodayParser) parseJsonLD(n *News, data []byte) error {

	if !json.Valid(data) {
		return fmt.Errorf("invalid json")
	}
	return ParseJsonLD(n, data, p.ToGUID)
}

func (p EtTodayParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("article div.story[itemprop='articleBody'] > p").
		Not("img").Each(func(i int, s *goquery.Selection) {

		if s.Find("span strong").Length() > 0 {
			return
		}

		if strings.HasPrefix(s.Text(), "延伸閱讀") {
			return
		}

		if c := strings.TrimSpace(s.Text()); c != "" && p.paragraphRe.MatchString(c) {
			if len(n.Author) == 0 {
				n.Author = p.getAuthor(c)
			} else {
				n.Content = append(n.Content, c)
			}
		}
	})

	body.Find("div#hot_area div.tab_content div.piece div.part_list_3").
		First().Find("h3 a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			if u, err := url.Parse(href); err == nil {
				n.RelatedGUID = append(n.RelatedGUID, p.ToGUID(u))
			}
		}
	})
	return nil
}

func (p EtTodayParser) getAuthor(s string) []string {
	author := []string{}
	if p.getAuthorRe.MatchString(s) {
		s = p.getAuthorRe.ReplaceAllString(s, "")
		for _, a := range strings.Split(s, "、") {
			a = strings.TrimPrefix(a, "記者")
			author = append(author, strings.TrimSpace(a))
		}
	} else {
		ss := strings.Split(s, "／")
		if len(ss) > 1 {
			author = ss[1:]
		} else {
			author = ss
		}
	}
	return author
}
