package parser

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/PuerkitoBio/goquery"
)

// CTParser is a parser for China Times news (https://www.chinatimes.com/).
// It implements the Parser interface.
type CTParser struct {
	JsonLDParser
	MetaParser
	parseFunc func(q *Query) *Query
}

func NewCTParser() *CTParser {
	p := &CTParser{}
	p.parseFunc = BuildParseFunc(p, nil)
	p.JsonLDParser = JsonLDParser{
		PostAssignFunc: func(n *News, jld *JsonLD) {
			n.GUID = p.ToGUID(n.Link)
		},
	}

	p.MetaParser = MetaParser{}
	p.MetaSelector, _ = GetDefaultMetaSelectorCopy()
	p.MetaSelector["description"] = "meta[name='description']"
	p.MetaParser.PostAssignFunc = func(n *News, meta *Meta) error {
		title := strings.Split(meta.Title, " - ")
		n.Title = title[0]
		n.Category = title[1]

		n.Language = meta.Language
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

// parse news
func (p CTParser) Parse(q *Query) *Query {
	return p.parseFunc(q)
}

// return guid from given url
func (p CTParser) ToGUID(herf *url.URL) string {
	return path.Base(herf.Path)
}

// return domains that the parser can parse
func (p CTParser) Domain() []string {
	return []string{"www.chinatimes.com"}
}

// extract content, related articles from body
func (p CTParser) ParseBody(item *News, body *goquery.Selection) error {
	article := body.Find("article.article-box")
	if len(item.Author) == 0 {
		article.Find("div.meta-info-wrapper div.meta-info div.author").Each(func(i int, s *goquery.Selection) {
			item.Author = append(item.Author, p.getAuthor(s.Text()))
		})
	}

	article.Find("div.article-body[itemprop='articleBody'] p").Each(func(i int, s *goquery.Selection) {
		if c := strings.TrimSpace(s.Text()); c != "" {
			item.Content = append(item.Content, c)
		}
	})

	rGUID := collection.NewSet[string]()
	var toGUIDErr error
	article.
		Find("div.article-body[itemprop='articleBody'] div.promote-word a").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if u, err := url.Parse(href); err == nil {
					rGUID.Add(p.ToGUID(u))
				} else {
					toGUIDErr = fmt.Errorf("error while parse %s: %w", href, err)
					return false
				}
			}
			return true
		})

	item.RelatedGUID = rGUID.Key()
	return toGUIDErr
}

// extract author
func (p CTParser) getAuthor(s string) string {
	var author string = s
	if strings.Contains(s, "、") {
		author = strings.Split(s, "、")[0]
	}
	if strings.Contains(s, "_") {
		author = strings.Split(author, "_")[0]
	}
	return strings.TrimSpace(author)
}
