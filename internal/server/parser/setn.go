package parser

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SETNParser struct {
	JsonLDParser
	MetaParser
	hiddenCharRe *regexp.Regexp
	paragraphRe  *regexp.Regexp
	multiSpaceRe *regexp.Regexp
	parseFunc    func(q *Query) *Query
}

func NewSETNParser() *SETNParser {
	p := &SETNParser{
		hiddenCharRe: regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`),
		paragraphRe:  regexp.MustCompile("^[^▸▲]"),
		multiSpaceRe: regexp.MustCompile(`\s{2,}`),
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
	p.MetaSelector["description"] = "meta[name='Description']"
	p.MetaSelector["language"] = "meta[http-equiv='content-language']"
	p.MetaSelector["category"] = "meta[property='article:section']"
	p.MetaSelector["pubDate"] = "meta[name='pubdate']"
	p.MetaSelector["modDate"] = "meta[name='moddate']"
	p.MetaSelector["keywords"] = "meta[name='news_keywords'][itemprop='keywords']"
	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		if title := strings.Split(meta.Title, " | "); len(title) >= 3 {
			n.Title = strings.TrimSpace(p.hiddenCharRe.ReplaceAllString(title[0], ""))
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

		if pubDate, err := time.Parse(time.RFC3339, meta.PubDate+"+08:00"); err == nil {
			n.PubDate = pubDate.UTC()
		} else {
			return fmt.Errorf("error while time.Parse: %w", err)
		}
		return nil
	}
	return p
}

func (p SETNParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p SETNParser) ToGUID(href *url.URL) string {
	// url may look like this: https://www.setn.com/News.aspx?NewsID=1348480
	return href.Query().Get("NewsID")
}

func (p SETNParser) Domain() []string {
	return []string{"www.setn.com"}
}

func (p SETNParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("div#ckuse article div > p").
		Not("p[style='text-align: center;']").
		Not("p[style='text-align:center;']").Each(func(i int, s *goquery.Selection) {
		if p.paragraphRe.MatchString(s.Text()) {
			if len(n.Author) == 0 {
				n.Author = p.getAuthor(s.Text())
			} else {
				c := p.multiSpaceRe.ReplaceAllString(strings.TrimSpace(p.hiddenCharRe.ReplaceAllString(s.Text(), "")), " ")
				if !strings.HasPrefix(c, "延伸閱讀") {
					n.Content = append(n.Content, c)
				}
			}
		}
	})

	var toGUIDErr error
	body.Find("div.tagNewsArea div.tagNewsBox div.tagNews a.gt").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if u, err := url.Parse(href); err == nil {
					n.RelatedGUID = append(n.RelatedGUID, p.ToGUID(u))
				} else {
					toGUIDErr = fmt.Errorf("error while parsing %s: %w", href, err)
					return false
				}
			}
			return true
		})
	return toGUIDErr
}

func (p SETNParser) getAuthor(s string) []string {
	reFmt1 := regexp.MustCompile(`.{2,4}中心／(.+)報導$`)
	reFmt2 := regexp.MustCompile(`文／(.+)`)
	author := []string{}

	if reFmt1.MatchString(s) {
		ss := reFmt1.FindStringSubmatch(s)
		for _, a := range strings.Split(ss[1], "、") {
			author = append(author, strings.TrimSpace(a))
		}
	} else if reFmt2.MatchString(s) {
		author = reFmt2.FindStringSubmatch(s)[1:]
	} else {
		ss := strings.Split(s, "／")[0]
		for _, a := range strings.Split(ss, "、") {
			author = append(author, strings.TrimSpace(strings.TrimPrefix(a, "記者")))
		}
	}
	return author
}
