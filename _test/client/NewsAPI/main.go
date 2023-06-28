package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	newsapi "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/newsAPI"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/code"
)

func main() {

	cli := newsapi.Client{
		ApiKey: global.Secrets.API.NewsSource["NewsAPI"],
		Client: http.DefaultClient,
	}
	fmt.Println("apikey:", global.Secrets)

	nw := time.Now()
	req, err := cli.NewQuery("Microsoft", newsapi.EPEverything).
		AppendParams(newsapi.EverythingParams{
			SearchIn:       nil,
			Domains:        nil,
			ExcludeDomains: nil,
			From:           nw.Add(10 * 24 * time.Hour),
			To:             nw,
			Language:       code.Language("en"),
			SortedBy:       newsapi.ByPopularity,
		}).ToHTTPRequest(context.Background(), cli.ApiKey)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(req.URL)

	// resp, err := cli.Client.Do(req)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// obj, err := newsapi.ParseHTTPResponse(resp)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(obj)
}
