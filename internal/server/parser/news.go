package parser

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
)

var ErrParserNotFound = errors.New("could not found parser for given domain")

type News struct {
	Title       string    `json:"title"`
	Link        *url.URL  `json:"link"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Author      []string  `json:"author,omitempty"`
	Category    string    `json:"category,omitempty"`
	GUID        string    `json:"guid,omitempty"`
	PubDate     time.Time `json:"pubDate,omitempty"`
	Content     []string  `json:"content"`
	Tag         []string  `json:"tag"`
	RelatedGUID []string  `json:"related_guid"`
}

func (n News) String() string {
	sb := bytes.NewBufferString("Item:\n")
	fmt.Fprintf(sb, " - Title       : %s\n", n.Title)
	fmt.Fprintf(sb, " - Author      : %s\n", strings.Join(n.Author, ", "))

	if n.Link != nil {
		fmt.Fprintf(sb, " - Link        : %s\n", n.Link.String())
	} else {
		fmt.Fprintf(sb, " - Link        : %s\n", "")
	}

	if d := []rune(n.Description); len(d) > 50 {
		fmt.Fprintf(sb, " - Description : %s (...%d words)\n", string(d[:50]), len(d)-50)
	} else {
		fmt.Fprintf(sb, " - Description : %s\n", n.Description)
	}

	fmt.Fprintf(sb, " - Category    : %s\n", n.Category)
	fmt.Fprintf(sb, " - GUID        : %s\n", n.GUID)
	fmt.Fprintf(sb, " - Language    : %s\n", n.Language)
	fmt.Fprintf(sb, " - PubDate     : %s\n", n.PubDate.Format(time.DateTime))
	fmt.Fprintf(sb, " - Tags        : %s\n", strings.Join(n.Tag, ", "))

	if c := []rune(strings.Join(n.Content, "")); len(c) > 50 {
		fmt.Fprintf(sb, " - Content     : %s (...%d words)\n", string(c[:50]), len(c)-50)
	} else {
		fmt.Fprintf(sb, " - Content     : %s\n", string(c))
	}

	fmt.Fprintf(sb, " - Related GUID: %s\n", strings.Join(n.RelatedGUID, ", "))
	return sb.String()
}

// MergeNewsItem merge two news item into a new one.
func MergeNewsItem(n1, n2 *News) *News {

	n3 := &News{}
	n3.Title = IfElse(n1.Title != "", n1.Title, n2.Title)
	n3.Link = IfElse(n1.Link != nil, n1.Link, n2.Link)
	n3.Description = IfElse(n1.Description != "", n1.Description, n2.Description)
	n3.Language = IfElse(n1.Language != "", n1.Language, n2.Language)
	n3.Category = IfElse(n1.Category != "", n1.Category, n2.Category)
	n3.GUID = IfElse(n1.GUID != "", n1.GUID, n2.GUID)
	n3.PubDate = IfElse(!n1.PubDate.IsZero(), n1.PubDate, n2.PubDate)

	author := IfElse(len(n1.Author) > 0, n1.Author, n2.Author)
	n3.Author = make([]string, len(author))
	copy(n3.Author, author)

	content := IfElse(len(n1.Content) > 0, n1.Content, n2.Content)
	n3.Content = make([]string, len(content))
	copy(n3.Content, content)

	tagSet := collection.NewSet(n1.Tag...)
	for i := range n2.Tag {
		tagSet.Add(n2.Tag[i])
	}
	n3.Tag = tagSet.Key()

	rGIDSet := collection.NewSet(n1.RelatedGUID...)
	for i := range n2.RelatedGUID {
		rGIDSet.Add(n2.RelatedGUID[i])
	}
	n3.RelatedGUID = rGIDSet.Key()
	return n3
}

// IfElse is a ternary operator. If test is true, return yes, otherwise return no.
func IfElse[T any](test bool, yes, no T) T {
	if test {
		return yes
	}
	return no
}
