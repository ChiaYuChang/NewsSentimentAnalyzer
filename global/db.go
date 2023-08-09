package global

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func ConnectToPostgres(ctx context.Context) (*pgx.Conn, error) {
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
	conn, err := model.NewDBConnection(ctx, connStr)
	if err != nil {
		Logger.Err(err).
			Msg("failed to connect to postgres")
		return nil, err
	}
	Logger.Info().
		Str("pgSqlConn", fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?%s",
			viper.GetString("POSTGRES_USERNAME"),
			strings.Repeat("x", len(viper.GetString("POSTGRES_PASSWORD"))),
			viper.GetString("POSTGRES_HOST"),
			viper.GetInt("POSTGRES_PORT"),
			viper.GetString("POSTGRES_DB_NAME"),
			options.Encode(),
		)).
		Msg("Connected to postgres")
	return conn, nil
}

func ConnectToRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network: viper.GetString("REDIS_NETWORK"),
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("REDIS_HOST"),
			viper.GetInt("REDIS_PORT"),
		),
		MaxRetries:   viper.GetInt("REDIS_MAX_RETRIES"),
		ReadTimeout:  viper.GetDuration("REDIS_READ_TIMEOUT"),
		WriteTimeout: viper.GetDuration("REDIS_WRITE_TIMEOUT"),
		PoolFIFO:     viper.GetBool("REDIS_FIFO"),
		PoolSize:     viper.GetInt("REDIS_POOLSIZE"),
	})
}
