package global

import (
	"fmt"
	"os"
	"sync"

	ld "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/languageDetector"
	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/language_detector"
	"google.golang.org/grpc"
)

var lds languageDetectorClientSingleton

type languageDetectorClientSingleton struct {
	cli ld.LanguageDetectorClient
	sync.Once
}

func NewLanguageDetectorClient(host string, port int,
	dailOpts []grpc.DialOption, detectLanguageCallOpt []grpc.CallOption,
	healthCheckCallOpt []grpc.CallOption) (ld.LanguageDetectorClient, error) {

	lds.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", host, port),
			dailOpts...,
		)
		if err != nil {
			Logger.Fatal().
				Str("address", host).
				Int("port", port).
				Err(err).
				Msg("failed to create language detector client")
			os.Exit(1)
		}

		lds.cli = ld.LanguageDetectorClient{
			LanguageDetectorClient: pb.NewLanguageDetectorClient(conn),
			DetectLanguageCallOpt:  detectLanguageCallOpt,
			HealthCheckCallOpt:     healthCheckCallOpt,
		}
	})
	return lds.cli, nil
}

func LanguageDetectorClient() ld.LanguageDetectorClient {
	return lds.cli
}
