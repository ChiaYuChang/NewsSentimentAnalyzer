package api

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type NewsPreview struct {
	Id          ulid.ULID `json:"id"                  redis:"id"`
	Title       string    `json:"title"               redis:"title"`
	Link        string    `json:"link"                redis:"link"`
	Description string    `json:"description"         redis:"description"`
	Category    string    `json:"category,omitempty"  redis:"category"`
	Content     string    `json:"content,omitempty"   redis:"content"`
	PubDate     time.Time `json:"publication_date"    redis:"publication_date"`
}

func (np NewsPreview) ToNewsCreateRequest(guid, language, source string, relatedGuid ...string) *service.NewsCreateRequest {
	md5Hash, _ := MD5Hash(np.Title, np.PubDate, np.Content)
	req := &service.NewsCreateRequest{
		Md5Hash:     md5Hash,
		Guid:        guid,
		Author:      nil,
		Title:       np.Title,
		Link:        np.Link,
		Description: np.Description,
		Language:    language,
		Content:     []string{np.Content},
		Category:    np.Category,
		Source:      source,
		RelatedGuid: relatedGuid,
		PublishedAt: np.PubDate,
	}
	return req
}

type PreviewCache struct {
	Query     CacheQuery    `json:"query"      redis:"query"`
	CreatedAt time.Time     `json:"created_at" redis:"created_at"`
	NewsItem  []NewsPreview `json:"news_item"  redis:"news_item"`
}

func (cache *PreviewCache) AppendNewsItem(items ...NewsPreview) {
	cache.NewsItem = append(cache.NewsItem, items...)
}

func (cache PreviewCache) Len() int {
	return len(cache.NewsItem)
}

func (cache *PreviewCache) AddRandomSalt(l int) error {
	return cache.Query.AddRandomSalt(l)
}

func (cache PreviewCache) Key(prefix, suffix string) string {
	salt, _ := hex.DecodeString(cache.Query.Salt)

	hasher := crypto.MD5.New()
	hasher.Write([]byte(cache.Query.UserId.String()))
	hasher.Write(salt)
	hasher.Write([]byte(cache.CreatedAt.UTC().Format(time.RFC3339)))
	return fmt.Sprintf("%s%s%s", prefix, hex.EncodeToString(hasher.Sum(nil)), suffix)
}

func (cache *PreviewCache) SetNextPage(token any) error {
	return cache.Query.SetNextPage(token)
}

type CacheQuery struct {
	UserId   uuid.UUID         `json:"user_id"              redis:"user_id"`
	Salt     string            `json:"salt"                 redis:"salt"`
	API      API               `json:"api"                  redis:"api"`
	NextPage NextPageToken     `json:"next_page,omitempty"  redis:"next_page"`
	RawQuery string            `json:"raw_query,omitempty"  redis:"raw_query"`
	Body     string            `json:"body,omitempty"       redis:"body"`
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

type API struct {
	Key      string `json:"key"       redis:"key"`
	Name     string `json:"name"      redis:"name"`
	Endpoint string `json:"endpoint"  redis:"endpoint"`
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
