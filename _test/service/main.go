package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/spf13/viper"
)

type Job struct {
	id  int
	err error
}

var ErrorAffectMoreThanExpectedRows = errors.New("affect more than expected rows")
var ErrorFieldNotMatch = errors.New("field in the object not match")

func main() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	pSqlConfig := viper.Get("postgres").(map[string]any)
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pSqlConfig["username"].(string),
		pSqlConfig["password"].(string),
		pSqlConfig["host"].(string),
		pSqlConfig["port"].(int),
		pSqlConfig["database"].(string),
	)

	// fmt.Println(connStr)
	conn, err := model.NewDBConnectionPools(context.Background(), connStr)
	if err != nil {
		if pgErr, ok := service.ToPgError(err); ok {
			fmt.Printf("- [%s] code: %s msg: %s\n", pgErr.Severity, pgErr.Code, pgErr.Message)
		} else {
			fmt.Println(fmt.Errorf("error while NewDBConnection: %w", err))
		}
		os.Exit(1)
	}
	defer conn.Close()

	if err := conn.Ping(context.Background()); err != nil {
		if pgErr, ok := service.ToPgError(err); ok {
			fmt.Printf("- [%s] code: %s msg: %s\n", pgErr.Severity, pgErr.Code, pgErr.Message)
		} else {
			fmt.Println(fmt.Errorf("error while Ping: %w", err))
		}
		os.Exit(0)
	}
	store := model.NewPGXPoolStore(conn)

	srvc := service.NewService(store, validator.Validate)
	testScript := []TestScript{
		CreateUpateThenDeleteUser{n: 10},
		CreateUpdateThenDeletAPI{n: 10},
		CreateGetThenDeletAPIKey{n: 10, nUser: 2, nAPI: 2},
		UserDeleteCascase{n: 10},
	}

	for i, ts := range testScript {
		fmt.Printf("Test %02d/%02d:\n", i, len(testScript))
		fmt.Println("Description:", ts.Description())

		logger := Logger{}
		logger.Verbose = true
		logger.Println("  > Sequentially")
		job := ts.Do(context.WithValue(context.Background(), "jid", 0), srvc, logger)
		if job.err != nil {
			panic(job.err)
		}

		logger.Println("  > Parallelly")
		logger.Verbose = false
		logger.RandomDelay = true
		ch := make(chan Job)

		nRT := ts.N()
		for i := 0; i < nRT; i++ {
			go func(ch chan<- Job) {
				ctx := context.WithValue(context.Background(), "jid", i)
				job := ts.Do(ctx, srvc, logger)
				ch <- job
			}(ch)
			time.Sleep(1 * time.Second)
		}

		es := make([]error, 0, nRT)
		for i := 0; i < nRT; i++ {
			job := <-ch
			if job.err != nil {
				es = append(es, err)
			}
			fmt.Printf("\t- Job %d Done (%02d/%02d)\n", job.id, i+1, nRT)
		}

		if len(es) > 0 {
			for _, e := range es {
				fmt.Println(e)
			}
			os.Exit(1)
		} else {
			fmt.Println("\t No error occured")
		}

		fmt.Println("  > Done")
	}
	fmt.Println("All Test PASSED")
}

type TestScript interface {
	Do(ctx context.Context, srvc service.Service, logger Logger) Job
	Description() string
	Steps() []string
	N() int
}
