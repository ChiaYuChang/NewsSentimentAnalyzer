package api

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NewsPreview struct {
	Id          int       `json:"id"               redis:"id"`
	Title       string    `json:"title"            redis:"title"`
	Link        string    `json:"link"             redis:"link"`
	Description string    `json:"description"      redis:"description"`
	Category    string    `json:"category"         redis:"category"`
	Content     string    `json:"content"          redis:"content"`
	PubDate     time.Time `json:"publication_date" redis:"publication_date"`
}

type PreviewCache struct {
	Query     CacheQuery    `json:"query"      redis:"query"`
	CreatedAt time.Time     `json:"created_at" redis:"created_at"`
	NewsItem  []NewsPreview `json:"news_item"  redis:"news_item"`
	end       int           `json:"-"          redis:"-"`
}

func (cache *PreviewCache) AppendNewsItem(items ...NewsPreview) {
	var i int
	cache.NewsItem = append(cache.NewsItem, items...)
	for i = cache.end; i < cache.Len(); i++ {
		cache.NewsItem[i].Id = i + 1
	}
	cache.end = i
}

func (cache PreviewCache) Len() int {
	return len(cache.NewsItem)
}

func (cache *PreviewCache) AddRandomSalt(l int) error {
	return cache.Query.AddRandomSalt(l)
}

func (cache PreviewCache) Key() string {
	salt, _ := hex.DecodeString(cache.Query.Salt)

	hasher := crypto.MD5.New()
	hasher.Write([]byte(cache.Query.UserId.String()))
	hasher.Write(salt)
	hasher.Write([]byte(cache.CreatedAt.UTC().Format(time.RFC3339)))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func (cache *PreviewCache) SetNextPage(token any) error {
	return cache.Query.SetNextPage(token)
}

type CacheQuery struct {
	UserId   uuid.UUID         `json:"user_id"              redis:"user_id"`
	Salt     string            `json:"salt"                 redis:"salt"`
	APIKey   string            `json:"api_key"              redis:"api_key"`
	APIEP    string            `json:"api_ep"               redis:"api_ep"`
	NextPage NextPageToken     `json:"next_page,omitempty"  redis:"next_page"`
	RawQuery string            `json:"raw_query"            redis:"raw_query"`
	Body     string            `json:"body"                 redis:"body"`
	Other    map[string]string `json:"other,omitempty"      redis:"other"`
}

func (cq *CacheQuery) SetNextPage(token any) error {
	switch val := token.(type) {
	default:
		return errors.New("Invalid token type")
	case int:
		cq.NextPage = IntNextPageToken(val)
	case string:
		cq.NextPage = StrNextPageToken(val)
	case IntNextPageToken:
		cq.NextPage = val
	case StrNextPageToken:
		cq.NextPage = val
	}
	return nil
}

func (cq *CacheQuery) AddRandomSalt(l int) error {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	cq.Salt = hex.EncodeToString(b)
	return nil
}

func (cq *CacheQuery) UnmarshalJSON(data []byte) error {
	type InnerCacheQuery CacheQuery
	tmp := struct {
		*InnerCacheQuery
		NextPage any `json:"next_page,omitempty"`
	}{
		InnerCacheQuery: (*InnerCacheQuery)(cq),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var token NextPageToken
	switch val := tmp.NextPage.(type) {
	default:
		return fmt.Errorf("unknown page token type")
	case string:
		token = StrNextPageToken(val)
	case int:
		token = IntNextPageToken(val)
	case float64:
		token = IntNextPageToken(int(val))
	}
	tmp.InnerCacheQuery.NextPage = token

	(*cq) = CacheQuery(*tmp.InnerCacheQuery)
	return nil
}

type Preview struct {
	CacheKey string        `json:"cache_key"`
	NewsItem []NewsPreview `json:"news_item"`
}

type NextPageToken interface {
	String() string
	Equal(token NextPageToken) bool
}

type IntNextPageToken int

func (t IntNextPageToken) String() string {
	return fmt.Sprintf("%d", t)
}

func (t1 IntNextPageToken) Equal(t2 NextPageToken) bool {
	t3, ok := t2.(IntNextPageToken)
	if !ok {
		return false
	}
	return t1 == t3
}

type StrNextPageToken string

func (t StrNextPageToken) String() string {
	return string(t)
}

func (t1 StrNextPageToken) Equal(t2 NextPageToken) bool {
	t3, ok := t2.(StrNextPageToken)
	if !ok {
		return false
	}
	return t1 == t3
}
