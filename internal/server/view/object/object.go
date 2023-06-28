package object

type HeadConent struct {
	Meta   *HTMLElementList
	Link   *HTMLElementList
	Script *HTMLElementList
}

type Page struct {
	HeadConent
	Title string
}

type ErrorPage struct {
	Page
	ErrorCode          int
	ErrorMessage       string
	ErrorDetail        string
	ShouldAutoRedirect bool
	RedirectPageUrl    string
	RedirectPageName   string
	CountDownFrom      int // second
}

type EndPoint struct {
	HeadConent
	API      string
	EndPoint string
}

type SelectOpts struct {
	OptMap         [][2]string
	MaxDiv         int
	DefaultValue   string
	DefaultText    string
	InsertButtonId string
	DeleteButtonId string
	PositionId     string
	AlertMessage   string
}
