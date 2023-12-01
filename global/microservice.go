package global

import (
	"context"
	"fmt"
	"strconv"
	"time"

	ld "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/languageDetector"
	np "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/newsParser"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APIType string

const (
	APIUnknown  APIType = "Unknown"
	APITypegRPC APIType = "gRPC"
	APITypeREST APIType = "REST"
)

func (t APIType) String() string {
	return string(t)
}

func (t *APIType) UnmarshalJSON(data []byte) error {
	s, _ := strconv.Unquote(string(data))
	switch s {
	case APITypegRPC.String():
		(*t) = APITypegRPC
	case APITypeREST.String():
		(*t) = APITypeREST
	default:
		(*t) = APIUnknown
	}
	return nil
}

type Microservice struct {
	Name   string            `mapstructure:"name"   json:"name"`
	Type   APIType           `mapstructure:"type"   json:"type"`
	Host   string            `mapstructure:"host"   json:"host"`
	Port   int               `mapstructure:"port"   json:"port"`
	Params map[string]string `mapstructure:"params" json:"params,omitempty"`
}

func SetupMicroservice(ms map[string]Microservice) error {
	if ms, ok := AppVar.Microservice["language-detector"]; ok {
		_, err := ld.SetupServiceLanguageDetectorClient(
			ms.Host, ms.Port,
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}, nil, nil,
		)

		if err != nil {
			Logger.Error().Err(err).Msg("error while setting up language-detector client")
			return err
		}
		Logger.Info().
			Str("name", ms.Name).
			Str("host", ms.Host).
			Int("port", ms.Port).
			Str("type", ms.Type.String()).
			Msg("add microservice")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cli, _ := ld.GetLanguageDetectorClient()
		if err := cli.HealthCheck(ctx); err != nil {
			return fmt.Errorf("error while checking language-detector health: %w", err)
		}
	}

	if ms, ok := AppVar.Microservice["news-parser"]; ok {
		_, err := np.SetupNewsParserClient(
			ms.Host, ms.Port,
			[]grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}, nil, nil, nil)

		if err != nil {
			Logger.Error().Err(err).Msg("error while setting up language-detector client")
			return err
		}
		Logger.Info().
			Str("name", ms.Name).
			Str("host", ms.Host).
			Int("port", ms.Port).
			Str("type", ms.Type.String()).
			Msg("add microservice")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cli, _ := np.GetNewsParserClient()
		if err := cli.HealthCheck(ctx); err != nil {
			return fmt.Errorf("error while checking news-paresr health: %w", err)
		}
	}
	return nil
}
