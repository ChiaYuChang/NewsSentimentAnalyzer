package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
)

type NewsItemList []NewsItem

func (items NewsItemList) N() int {
	return len(items) + 1
}

func (items NewsItemList) Item() []NewsItem {
	return []NewsItem(items)
}

var NewsCategory = []string{
	"General",
	"World",
	"Nation",
	"Business",
	"Technology",
	"Entertainment",
	"Sports",
	"Science",
	"Health",
}

var NewsLanguage = []string{
	"en-US",
	"zh-TW",
	"zh-CN",
	"ja-JP",
	"ko-KR",
	"fr-FR",
	"de-DE",
	"es-ES",
	"ru-RU",
	"ar-SA",
	"pt-BR",
}

type NewsItem struct {
	Md5Hash     string    `json:"md5_hash"`
	Guid        string    `json:"guid"`
	Author      []string  `json:"author"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Content     []string  `json:"content"`
	Category    string    `json:"category"`
	Source      string    `json:"source"`
	RelatedGuid []string  `json:"related_guid"`
	PublishAt   time.Time `json:"publish_at"`
}

func ToArray(ss []string) string {
	for i := 0; i < len(ss); i++ {
		ss[i] = fmt.Sprintf("'%s'", ss[i])
	}
	return fmt.Sprintf("ARRAY [%s]", strings.Join(ss, ","))
}

func NewRandomNewsItem() NewsItem {
	u, _ := rg.GenRdnUrl()
	contents := make([]string, 3+rand.Intn(5))
	for j := 0; j < len(contents); j++ {
		contents[j] = rg.Must(rg.AlphaNum.GenRdmString(20 + rand.Intn(30)))
	}
	guid := strings.ReplaceAll(u.Path, "/", "-")
	guid = strings.TrimLeft(guid, "-")

	authors := make([]string, 1+rand.Intn(2))
	for i := 0; i < len(authors); i++ {
		authors[i] = rg.Must(rg.Alphabet.GenRdmString(5 + rand.Intn(10)))
	}

	item := NewsItem{
		Guid:        guid,
		Title:       rg.Must(rg.Alphabet.GenRdmString(10 + rand.Intn(30))),
		Link:        u.String(),
		Author:      authors,
		Description: rg.Must(rg.Alphabet.GenRdmString(50 + rand.Intn(150))),
		Language:    NewsLanguage[rand.Intn(len(NewsLanguage))],
		Content:     contents,
		Category:    NewsCategory[rand.Intn(len(NewsCategory))],
		Source:      u.Host,
		RelatedGuid: []string{},
		PublishAt:   rg.GenRdnTime(time.Now().AddDate(-1, 0, 0), time.Now()),
	}

	hasher := md5.New()
	hasher.Write([]byte(item.Title))
	hasher.Write([]byte(item.Link))
	hasher.Write([]byte(item.Description))
	for _, c := range item.Content {
		hasher.Write([]byte(c))
	}
	hasher.Write([]byte(item.PublishAt.UTC().Format(time.DateOnly)))
	item.Md5Hash = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return item
}

func NewRandomNewsItems(n int) NewsItemList {
	rgids := make([]string, 0, n)
	items := make([]NewsItem, n)
	for i := 0; i < n; i++ {
		items[i] = NewRandomNewsItem()
		rgids = append(rgids, items[i].Guid)
	}

	for i := 0; i < n; i++ {
		rgids := make([]string, 3+rand.Intn(n/100+1))
		for j := 0; j < len(rgids); j++ {
			rgids[j] = items[rand.Intn(n)].Guid
		}
		items[i].RelatedGuid = rgids
	}

	return items
}

func (item NewsItem) AuthorString() string {
	return ToArray(item.Author)
}

func (item NewsItem) ContentString() string {
	return ToArray(item.Content)
}

func (item NewsItem) RelatedGuidString() string {
	return ToArray(item.RelatedGuid)
}
