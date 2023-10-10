package parser

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/PuerkitoBio/goquery"
)

// UPParser is a parser for UP media news (https://www.upmedia.mg/).
// It implements the Parser interface.
type UPParser struct {
	MetaParser
	parseFunc func(q *Query) *Query
}

func NewUPParser() *UPParser {
	p := &UPParser{}
	p.parseFunc = BuildParseFunc(p, nil)

	p.MetaParser = MetaParser{}
	p.MetaSelector = map[string]string{}
	p.MetaSelector["link"] = "link[rel='canonical']"
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaSelector["category"] = "meta[itemprop='articleSection']"
	p.MetaSelector["pubDate"] = "meta[itemprop='datePublished']"
	p.MetaSelector["modDate"] = "meta[itemprop='dateModified']"
	p.MetaSelector["keywords"] = "meta[itemprop='keywords']"
	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		if ss := strings.Split(meta.Title, "--"); len(ss) == 2 {
			n.Title = strings.TrimSpace(ss[0])
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

func (p UPParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

func (p UPParser) ToGUID(herf *url.URL) string {
	q := herf.Query()
	return q.Get("Type") + "-" + q.Get("SerialNo")
}

func (p UPParser) Domain() []string {
	return []string{"www.upmedia.mg"}
}

func (p UPParser) ParseJsonLD(n *News, s *goquery.Selection) error {
	// UP news doesn't have json linked data
	return nil
}

func (p UPParser) ParseBody(n *News, body *goquery.Selection) error {
	body.Find("div#news-info div.author > a").Each(func(i int, s *goquery.Selection) {
		a := strings.TrimPrefix(s.Text(), "上報快訊／")
		n.Author = append(n.Author, a)
	})

	body.Find("div#news-info div.editor > p").Each(func(i int, s *goquery.Selection) {
		if c := strings.TrimSpace(s.Text()); c != "" && !strings.Contains(c, "作者為《上報》總主筆") {
			n.Content = append(n.Content, c)
		}
	})

	rGUID := collection.NewSet[string]()
	var toGUIDErr error
	body.Find("div#news-info div.related ul li > a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if u, err := url.Parse(href); err == nil {
					rGUID.Add(p.ToGUID(u))
				} else {
					toGUIDErr = err
					return false
				}
			}
			return true
		})
	n.RelatedGUID = rGUID.Key()
	return toGUIDErr
}
