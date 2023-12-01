package parser

import (
	"fmt"
	"net/url"
	"sort"
	"sync"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/PuerkitoBio/goquery"
)

var defaultParserRepo struct {
	Parser
	sync.Once
}

func ToGUID(href *url.URL) string {
	return GetDefaultParser().ToGUID(href)
}

func Has(domain string) bool {
	switch p := GetDefaultParser().(type) {
	case ParserRepo:
		return p.Has(domain)
	default:
		for _, d := range p.Domain() {
			if domain == d {
				return true
			}
		}
		return false
	}
}

func ParseRawURL(rawurl string) *Query {
	return Parse(&Query{RawURL: rawurl})
}

func ParseURL(u *url.URL) *Query {
	return Parse(&Query{RawURL: u.String(), URL: u})
}

func Parse(q *Query) *Query {
	return GetDefaultParser().Parse(q)
}

func GetDefaultParser() Parser {
	defaultParserRepo.Do(func() {
		defaultParserRepo.Parser = NewParserRepo(
			// 國內媒體
			NewEtTodayParser(), // EtToday 新聞雲
			NewTVBSParser(),    // TVBS 新聞網
			NewSETNParser(),    // 三立新聞網
			NewUPParser(),      // 上報
			NewCNAParser(),     // 中央社
			NewCTParser(),      // 中國時報
			NewPTSParser(),     // 公視新聞網
			NewLTNParser(),     // 自由時報
			NewUDNParser(),     // 聯合新聞網
			// 外媒
			NewBBCParser(),     // BBC News 中文
			NewNYTimesParser(), // 紐約時報中文網
			NewRFIParser(),     // 法廣台灣
		)
		fmt.Println("set up parser repo singleton")
	})
	return defaultParserRepo.Parser
}

// Parser is an interface for parsing news.
type Parser interface {
	// Parse parses a news from a io.ReadCloser.
	Parse(q *Query) *Query
	// ToGUID returns a GUID from a URL.
	ToGUID(href *url.URL) string
	// Domain returns a list of domains that the parser can parse.
	Domain() []string
}

type StdParseProcess interface {
	ParseJsonLD(item *News, jsonld *goquery.Selection) error
	ParseMeta(item *News, s *goquery.Selection) error
	ParseBody(item *News, s *goquery.Selection) error
	ToGUID(href *url.URL) string
}

func BuildParseFunc(p StdParseProcess, selector map[string]string) func(q *Query) *Query {
	headselector := "head"
	bodyselector := "body"
	jsonldselector := "script[type='application/ld+json']"

	if selector != nil {
		if v, ok := selector["head"]; ok {
			headselector = v
		}
		if v, ok := selector["body"]; ok {
			bodyselector = v
		}
		if v, ok := selector["jsonld"]; ok {
			jsonldselector = v
		}
	}

	return func(q *Query) *Query {
		content, err := q.Content()
		if err != nil {
			q.Error = err
			return q
		}

		doc, err := ToDoc(content)
		if err != nil {
			q.Error = err
			return q
		}

		n1 := &News{}
		n2 := &News{}
		if lang, ok := doc.Find("html").First().Attr("lang"); ok {
			n1.Language = lang
			n2.Language = lang
		}

		if err := p.ParseJsonLD(n1, doc.Find(jsonldselector)); err != nil {
			q.Error = err
			return q
		}

		if err := p.ParseMeta(n2, doc.Find(headselector).First()); err != nil {
			q.Error = err
			return q
		}

		q.News = MergeNewsItem(n1, n2)
		sort.Sort(sort.StringSlice(q.News.Tag))

		p.ParseBody(q.News, doc.Find(bodyselector).First())
		return q
	}
}

// ContentParsingFunc is a function for parsing news content.
// It could be used when a news site has multiple domains.
type ContentParsingFunc func(i *News, s *goquery.Selection)

// ParserRepo is a repository for Parser.
// It implements the Parser interface.
type ParserRepo map[string]Parser

// NewNewsParser returns a ParserRepo.
// opts are optional parsers, which can be registered to the ParserRepo.
func NewParserRepo(opts ...Parser) ParserRepo {
	pr := ParserRepo{}
	for _, opt := range opts {
		pr.RegisterDomainParser(opt)
	}
	return pr
}

// RegisterDomainParser registers a parser for a domain.
func (repo ParserRepo) RegisterDomainParser(parser ...Parser) ParserRepo {
	for _, p := range parser {
		for _, domain := range p.Domain() {
			global.Logger.Info().
				Str("domain", domain).
				Msg("register parser")
			repo[domain] = p
		}
	}
	return repo
}

// Domain returns a list of domains that the parser can parse.
func (repo ParserRepo) Domain() []string {
	domains := make([]string, 0, len(repo))
	for domain := range repo {
		domains = append(domains, domain)
	}
	return domains
}

// Parse parses a news from a query.
func (repo ParserRepo) Parse(q *Query) *Query {
	if q.URL == nil {
		q.URL, q.Error = url.Parse(q.RawURL)
		if q.Error != nil {
			return q
		}
	}

	if q.RawURL == "" {
		q.RawURL = q.URL.String()
	}

	if parser, ok := repo[q.URL.Host]; ok {
		q = parser.Parse(q)
	} else {
		q.Error = fmt.Errorf("%w: %s", ErrParserNotFound, q.URL.Host)
	}
	return q
}

// ToGUID returns a GUID from a URL.
func (repo ParserRepo) ToGUID(href *url.URL) string {
	if parser, ok := repo[href.Host]; ok {
		return parser.ToGUID(href)
	}
	return href.Path
}

func (repo ParserRepo) Has(domain string) bool {
	_, ok := repo[domain]
	return ok
}
