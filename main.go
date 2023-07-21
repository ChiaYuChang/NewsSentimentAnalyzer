package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	errorcode "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/spf13/viper"
)

func main() {
	global.ReadConfig()
	fmt.Println(global.AppVar)

	pgSqlConn, err := global.ConnectToPostgres(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srvc := service.NewService(model.NewPGXStore(pgSqlConn), validator.Validate)

	rds := global.ConnectToRedis()
	rdsStatus := rds.Ping(context.Background())
	if err := rdsStatus.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	vw, err := view.NewViewWithDefaultTemplateFuncs(global.AppVar.App.Template...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error while NewView: %s\n", err)
		os.Exit(1)
	}

	tm := tokenmaker.NewJWTMaker(
		global.AppVar.Token.Secret(),
		global.AppVar.Token.SignMethod.Algorthm,
		global.AppVar.Token.SignMethod.Size,
		global.AppVar.Token.ExpireAfter,
		global.AppVar.Token.ValidAfter,
	)

	cm := cookieMaker.NewCookieMaker(
		"/", viper.GetString("APP_HOST"),
		int(global.AppVar.Token.ExpireAfter.Seconds()),
		true, false, http.SameSiteLaxMode,
	)

	addr := fmt.Sprintf("%s:%d", viper.GetString("APP_HOST"), viper.GetInt("APP_PORT"))
	server := &http.Server{
		Addr:    addr,
		Handler: router.NewRouter(srvc, rds, vw, tm, cm),
	}

	serverCtx, serverCancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func(signalChan chan os.Signal) {
		sig := <-signalChan
		fmt.Printf("Get %v\n", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				fmt.Fprintln(os.Stderr, "graceful shutdown timed out.. forcing exit.")
				os.Exit(0)
			}
		}()

		ec := make(chan error)
		go func() {
			ec <- server.Shutdown(shutdownCtx)
			fmt.Println("Server shutdown")
		}()

		go func() {
			ec <- srvc.Close(shutdownCtx)
			fmt.Println("PostgresSQL connection closed")
		}()

		go func() {
			ec <- rds.Close()
			fmt.Println("Redis connection closed")
		}()

		var ecErr *errorcode.Error
		for i := 0; i < 3; i++ {
			err := <-ec
			if err != nil {
				if ecErr == nil {
					ecErr = errorcode.MustGetEcErr(errorcode.ECServerError)
				}
				ecErr.WithDetails(err.Error())
			}
		}

		if ecErr != nil {
			jsn, _ := ecErr.ToJson()
			fmt.Fprintln(os.Stderr, string(jsn))
		}
		serverCancel()
	}(signalChan)

	fmt.Println("Server start at:", addr)
	if err := server.ListenAndServeTLS(
		"./secrets/server.crt",
		"./secrets/server.key",
	); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	<-serverCtx.Done()
	fmt.Println("Shutdown gracefully")
	os.Exit(0)
}
