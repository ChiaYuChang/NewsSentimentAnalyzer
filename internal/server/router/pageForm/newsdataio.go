package pageform

type NEWSDATAIOLatestNews struct {
	// IncludeContent // currently not yet support
	Keyword  string `form:"keyword"`
	Domains  string `form:"domains"`
	Language string `form:"language"`
	Country  string `form:"country"`
	Category string `form:"category"`
}

type NEWSDATAIONewsArchive struct {
	TimeRange
	Keyword  string `form:"keyword"`
	Domains  string `form:"domains"`
	Language string `form:"language"`
	Country  string `form:"country"`
	Category string `form:"category"`
}

type NEWSDATAIONewsSources struct {
	Language string `form:"language"`
	Country  string `form:"country"`
	Category string `form:"category"`
}

type IncludeContent struct {
	Image       bool `form:"image"`
	Video       bool `form:"video"`
	FullContent bool `form:"full-content"`
}
