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

type BBCParser struct {
	JsonLDParser
	MetaParser
	multiSpaceRe *regexp.Regexp
	parseFunc    func(q *Query) *Query
}

func NewBBCParser() *BBCParser {
	p := &BBCParser{
		multiSpaceRe: regexp.MustCompile("\\s{2,}"),
	}

	p.parseFunc = BuildParseFunc(p, nil)
	p.JsonLDParser = JsonLDParser{}
	p.JsonLDParser.PostAssignFunc = func(n *News, jld *JsonLD) {
		tags := make([]string, len(jld.About))
		for i, about := range jld.About {
			tags[i] = about.Name
		}
		n.Tag = tags
	}
	p.JsonLDParser.DataPreprocessFunc = func(data []byte) []byte {
		var str, end int
		for str = 0; str < len(data); str++ {
			if data[str] == '[' {
				break
			}
		}

		for end = len(data) - 1; end >= 0; end-- {
			if data[end] == ']' {
				break
			}
		}
		data, err := p.JsonLDParser.FindTargetJsonLD(data[str:end+1], "NewsArticle")
		if err != nil {
			return []byte("{}")
		}
		return data
	}

	re := regexp.MustCompile("zhongwen/trad/(.+)-[0-9]{7,10}$")
	p.MetaParser = MetaParser{
		MetaSelector: map[string]string{
			"link":        "link[rel='canonical']",
			"description": "meta[name='description']",
			"pubDate":     "meta[name='article:published_time']",
			"modDate":     "meta[name='article:modified_time']",
			"keywords":    "meta[name='article:tag']",
		},
		PostAssignFunc: func(n *News, meta *Meta) error {
			n.Title = strings.TrimSuffix(meta.Title, " - BBC News 中文")

			n.Description = meta.Description
			n.Tag = NewCSL(meta.Tag)

			if u, err := url.Parse(meta.Link); err == nil {
				n.Link = u
			} else {
				return fmt.Errorf("error while url.Parse: %w", err)
			}

			sm := re.FindStringSubmatch(meta.Link)
			if len(sm) == 2 {
				n.Category = sm[1]
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

func (p BBCParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p BBCParser) ToGUID(u *url.URL) string {
	return path.Base(u.Path)
}

func (p BBCParser) Domain() []string {
	return []string{"www.bbc.com", "www.bbc.co.uk"}
}

func (p BBCParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("main[role='main'] div[dir='ltr'] div.bbc-1atl7vu.euvj3t14 ul.bbc-143s8qx[role='list'] li").Each(func(i int, s *goquery.Selection) {
		author := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(author, "BBC") {
			n.Author = append(n.Author, author)
		}
	})

	body.Find("main[role='main'] div.bbc-19j92fr[dir='ltr'] > p.bbc-w2hm1d,h2.bbc-z6r16b").Each(func(i int, s *goquery.Selection) {
		content := p.cleanContent(s.Text())
		if content != "" {
			n.Content = append(n.Content, content)
		}
	})

	RGUIDSet := collection.NewSet[string]()
	var tagSelector = []string{
		"main[role='main'] div.etpldq00 ul li a.focusIndicatorReducedWidth",
		"section[data-e2e='related-content-heading'] ul li a",
	}
	var urlParseErr error
	for _, tgs := range tagSelector {
		body.Find(tgs).EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if u, err := url.Parse(href); err == nil {
					RGUIDSet.Add(p.ToGUID(u))
				} else {
					urlParseErr = fmt.Errorf("error while parsing %s: %w", href, err)
					return false
				}
			}
			return true
		})
		if urlParseErr != nil {
			return urlParseErr
		}
	}

	n.RelatedGUID = RGUIDSet.Key()
	return nil
}

func (p BBCParser) cleanContent(content string) string {
	content = strings.TrimSpace(content)
	content = p.multiSpaceRe.ReplaceAllString(content, " ")
	content = strings.ReplaceAll(content, "\n", "")
	return content
}
