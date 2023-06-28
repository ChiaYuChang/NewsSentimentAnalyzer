package pageform

const NEWSDATAIOName = "NewsAPI"

const (
	NEWSDATAIOEPLatestNews  string = "latest-news"
	NEWSDATAIOEPNewsArchive        = "archive"
	NEWSDATAIOEPNewsSources        = "sources"
	NEWSDATAEPCrypto               = "crypto"
)

type NEWSDATAIOLatestNews struct {
	// IncludeContent // currently not yet support
	Keyword  string   `form:"keyword"  validate:"max=512"`
	Domains  string   `form:"domains"`
	Language []string `form:"language"`
	Country  []string `form:"country"`
	Category []string `form:"category"`
}

func (f NEWSDATAIOLatestNews) Endpoint() string {
	return NEWSDATAIOEPLatestNews
}

func (f NEWSDATAIOLatestNews) API() string {
	return NEWSDATAIOName
}

type NEWSDATAIONewsArchive struct {
	TimeRange
	Keyword  string   `form:"keyword"`
	Domains  string   `form:"domains"`
	Language []string `form:"language"`
	Country  []string `form:"country"`
	Category []string `form:"category"`
}

func (f NEWSDATAIONewsArchive) Endpoint() string {
	return NEWSDATAIOEPNewsArchive
}

func (f NEWSDATAIONewsArchive) API() string {
	return NEWSDATAIOName
}

type NEWSDATAIONewsSources struct {
	Language []string `form:"language"`
	Country  []string `form:"country"`
	Category []string `form:"category"`
}

func (f NEWSDATAIONewsSources) Endpoint() string {
	return NEWSDATAIOEPNewsSources
}

func (f NEWSDATAIONewsSources) API() string {
	return NEWSDATAIOName
}

type IncludeContent struct {
	Image       bool `form:"image"`
	Video       bool `form:"video"`
	FullContent bool `form:"full-content"`
}
