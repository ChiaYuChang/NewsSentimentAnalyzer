package object_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIntNextPageToken(t *testing.T) {
	uid := uuid.New()
	token := 1
	ctime := "2023-11-06T13:28:18Z"
	salt := "8b048a4c42da1b7fc4300978b85eb5bd72192573d99d4eced28eebceab3c02bc"
	jsn := fmt.Sprintf(`{
		"query": {
			"user_id": "%s",
			"next_page": %d,
			"salt": "%s"
		},
		"created_at": "%s"
	}`, uid.String(), token, salt, ctime)

	var s1 api.PreviewCache
	var s2 api.PreviewCache
	err := json.Unmarshal([]byte(jsn), &s1)
	require.NoError(t, err)

	b, err := json.MarshalIndent(s1, "", "    ")
	require.NoError(t, err)

	require.Equal(t, api.IntNextPageToken(1), s1.Query.NextPage)
	require.Equal(t, uid, s1.Query.UserId)

	err = json.Unmarshal(b, &s2)
	require.NoError(t, err)

	require.Equal(t, s1.Query.UserId, s1.Query.UserId)
	require.Equal(t, s1.Query.NextPage, s2.Query.NextPage)
}

func TestStrNextPageToken(t *testing.T) {
	uid := uuid.New()
	token := "xxxx-xx-xxxx-xx-xx"
	ctime := "2023-11-06T13:28:18Z"
	salt := "8b048a4c42da1b7fc4300978b85eb5bd72192573d99d4eced28eebceab3c02bc"
	jsn := fmt.Sprintf(`{
		"query": {
			"user_id": "%s",
			"next_page": "%s",
			"salt": "%s"
		},
		"created_at": "%s"
	}`, uid.String(), token, salt, ctime)

	var s1 api.PreviewCache
	var s2 api.PreviewCache
	err := json.Unmarshal([]byte(jsn), &s1)
	require.NoError(t, err)

	b, err := json.Marshal(s1)
	require.NoError(t, err)
	require.Equal(t, uid, s1.Query.UserId)
	require.Equal(t, api.StrNextPageToken(token), s1.Query.NextPage)
	require.Equal(t, salt, s1.Query.Salt)
	require.Equal(t, ctime, s1.CreatedAt.Format(time.RFC3339))

	err = json.Unmarshal(b, &s2)
	require.NoError(t, err)

	require.Equal(t, s1.Query.UserId, s2.Query.UserId)
	require.Equal(t, s1.Query.NextPage, s2.Query.NextPage)

	nps := []api.NewsPreview{
		{},
		{},
		{},
		{},
		{},
		{},
	}

	s1.AppendNewsItem(nps...)
	require.Equal(t, len(nps), s1.Len())
	s1.AppendNewsItem(nps...)
	require.Equal(t, 2*len(nps), s1.Len())

	for i, item := range s1.NewsItem {
		require.Equal(t, i+1, item.Id)
	}

}

func TestHTMLElementToJson(t *testing.T) {
	o0 := object.NewHTMLElement("div")
	o0.AddVal("checked")
	o0.AddPair("class", "btn")
	o0.Content = "text"

	t.Log("o0:", o0.ToHTML())

	b, err := json.Marshal(o0)
	require.NoError(t, err)

	t.Log(string(b))

	var o1 object.HTMLElement
	err = json.Unmarshal(b, &o1)
	require.NoError(t, err)

	t.Log("o1:", o1.ToHTML())

	require.Equal(t, o0.IsSelfClosing, o1.IsSelfClosing)
	require.Equal(t, o0.Content, o1.Content)
	require.Equal(t, o0.Tag, o1.Tag)
	require.Equal(t, o0.Content, o1.Content)
	require.ElementsMatch(t, o0.Attr, o1.Attr)
}

func TestHTMLElementListToJson(t *testing.T) {
	meta0 := object.NewHTMLElementList("meta")

	meta0.NewHTMLElement().
		AddPair("charset", "UTF-8").
		AddVal("checked")

	meta0.NewHTMLElement().
		AddPair("http-equiv", "X-UA-Compatible").
		AddPair("content", "IE=edge")

	meta0.NewHTMLElement().
		AddPair("name", "viewport").
		AddPair("content", "width=device-width, initial-scale=1.0")

	b, err := json.MarshalIndent(meta0, "", "    ")
	require.NoError(t, err)
	t.Log(string(b))

	var meta1 object.HTMLElementList
	err = json.Unmarshal(b, &meta1)
	require.NoError(t, err)

	require.Equal(t, meta0.Tag, meta1.Tag)
	require.Equal(t, meta0.Ele, meta1.Ele)
}
