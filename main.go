package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
)

func main() {
	if err := global.ReadAppVar(
		"./secret.json",
		"./config/option.json",
		"./config/endpoint.json",
	); err != nil {
		panic(fmt.Sprintf("error while reading secret: %s", err.Error()))
	}
	fmt.Println(global.AppVar)

	db := global.AppVar.Secret.Database["postgres"]
	options := make([]string, 0, len(db.Options))
	// for optName, optValue := range db.Options {
	// 	options = append(options, fmt.Sprintf("%s=%s", optName, optValue))
	// }
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		db.UserName, db.Password, db.Host, db.Port, db.DBName,
		url.QueryEscape(strings.Join(options, ",")),
	)

	conn, err := model.NewDBConnection(context.TODO(), connStr)
	if err != nil {
		panic(err)
	}

	service := service.NewService(model.NewPGXStore(conn), validator.Validate)
	vw, err := view.NewViewWithDefaultTemplateFuncs(global.AppVar.Server.TemplatePath...)
	if err != nil {
		panic(err)
	}

	fs := http.Dir(global.AppVar.Server.StaticFilePath)
	tm := tokenmaker.NewJWTMakerWithDefaultVal()
	cm := cookieMaker.NewTestCookieMaker()

	mux := router.NewRouter(service, vw, fs, tm, cm)

	errCh := make(chan error)
	go func(chan<- error) {
		errCh <- http.ListenAndServe(fmt.Sprintf(
			"%s:%d",
			global.AppVar.Server.Binding.Host,
			global.AppVar.Server.Binding.Port,
		), mux)
	}(errCh)

	err = <-errCh
	if err != nil {
		panic(err)
	}
}
