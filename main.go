package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/spf13/viper"
)

func main() {
	global.ReadConfig()
	fmt.Println(global.AppVar)

	options := url.Values{}
	options.Add("sslmode", viper.GetString("POSTGRES_SSL_MODE"))
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		viper.GetString("POSTGRES_USERNAME"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetInt("POSTGRES_PORT"),
		viper.GetString("POSTGRES_DB_NAME"),
		options.Encode(),
	)

	fmt.Println("ConnStr:", connStr)
	conn, err := model.NewDBConnection(context.TODO(), connStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	service := service.NewService(model.NewPGXStore(conn), validator.Validate)

	vw, err := view.NewViewWithDefaultTemplateFuncs(global.AppVar.App.Template...)
	if err != nil {
		panic(err)
	}

	fs := http.Dir(global.AppVar.App.StaticFile.Path)
	tm := tokenmaker.NewJWTMakerWithDefaultVal()
	cm := cookieMaker.NewTestCookieMaker()

	mux := router.NewRouter(service, vw, fs, tm, cm)

	errCh := make(chan error)
	addr := fmt.Sprintf("%s:%d", viper.GetString("APP_HOST"), viper.GetInt("APP_PORT"))
	fmt.Println("Server start at:", addr)
	go func(chan<- error) {
		errCh <- http.ListenAndServe(addr, mux)
	}(errCh)

	err = <-errCh
	if err != nil {
		panic(err)
	}
}
