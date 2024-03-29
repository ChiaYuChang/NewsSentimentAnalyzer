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
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/cookieMaker"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	errorcode "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	global.ReadConfig()

	logfile, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while OpenFile: %v", err)
		os.Exit(1)
	}
	global.NewGlobalLogger(logfile)
	fmt.Println(global.AppVar.String())

	pgSqlConn, err := global.ConnectToPostgres(context.TODO())
	if err != nil {
		global.Logger.Error().
			Str("db", "psql").
			Str("status", "ok").
			Err(err).
			Msg("failed to connect to postgresSQL server")
		os.Exit(1)
	}

	global.Logger.Info().
		Str("db", "psql").
		Str("status", "ok").
		Msg("connected to postgresSQL server")
	srvc := service.NewService(model.NewPGXStore(pgSqlConn), validator.Validate)

	rds := global.ConnectToRedis(context.TODO())

	rdsStatus := rds.Ping(context.Background())
	if err := rdsStatus.Err(); err != nil {
		global.Logger.Err(err).Send()
		os.Exit(1)
	}
	global.Logger.Err(err).Send()

	vw, err := view.NewViewWithDefaultTemplateFuncs(global.AppVar.App.Template...)
	if err != nil {
		global.Logger.
			Err(err).
			Msg("error while view.NewViewWithDefaultTemplateFuncs")
		os.Exit(1)
	}
	global.Logger.Info().Msg("Connected to redis")

	_ = global.SetMinifier(global.MinifierOpt{
		Mimetype: "application/javascript",
		Minifier: js.Minify,
	})

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
		true, true, http.SameSiteLaxMode,
	)
	cookieMaker.SetDefaultCookieMaker(cm)

	err = global.SetupMicroservice(global.AppVar.Microservice)
	if err != nil {
		global.Logger.
			Err(err).
			Msg("error while SetupMicroservice")
		os.Exit(1)
	}

	addr := fmt.Sprintf(
		"%s:%d",
		viper.GetString("APP_HOST"),
		viper.GetInt("APP_PORT"))

	global.Logger.Info().
		Str("addr", addr).
		Time("time", time.Now()).
		Msg("Server started")

	server := &http.Server{
		Addr:    addr,
		Handler: router.NewRouter(srvc, rds, vw, tm, cm),
	}

	serverCtx, serverCancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func(signalChan chan os.Signal) {
		sig := <-signalChan
		global.Logger.Info().
			Str("Signal", sig.String()).
			Msgf("Get %v\n", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				global.Logger.Error().
					Err(shutdownCtx.Err()).
					Msg("graceful shutdown timed out.. forcing exit.")
				os.Exit(0)
			}
		}()

		ec := make(chan error)
		go func() {
			ec <- server.Shutdown(shutdownCtx)
			global.Logger.Info().Msg("Server shutdown")
		}()

		go func() {
			ec <- srvc.Close(shutdownCtx)
			global.Logger.Info().Msg("PostgresSQL connection closed")
		}()

		go func() {
			ec <- rds.Close()
			global.Logger.Info().Msg("Redis connection closed")
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

	startAt := time.Now()
	global.Logger.Info().
		Str("addr", addr).
		Str("api verion", viper.GetString("APP_API_VERSION")).
		Str("cert", global.AppVar.App.Certificate.CertFilePath()).
		Str("key", global.AppVar.App.Certificate.KeyFilePath()).
		Msg("Start server")

	if err := server.ListenAndServeTLS(
		// "./secrets/server.crt",
		// "./secrets/server.key",
		global.AppVar.App.Certificate.CertFilePath(),
		global.AppVar.App.Certificate.KeyFilePath(),
	); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	<-serverCtx.Done()
	endAt := time.Now()
	global.Logger.Info().
		Dur("Uptime", endAt.Sub(startAt)).
		Msg("Shutdown gracefully")

	if err := logfile.Close(); err != nil {
		global.Logger.Err(err).Send()
	}
	os.Exit(0)
}
