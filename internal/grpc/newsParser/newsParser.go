package newsparser

import (
	"context"
	"errors"
	"fmt"
	"sync"

	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/news_parser"
	"google.golang.org/grpc"
)

var cliS newsParserClientSingleton

type newsParserClientSingleton struct {
	cli   NewsParserClient
	Error error
	sync.Once
}

func SetupNewsParserClient(host string, port int,
	dailOpts []grpc.DialOption,
	parseUrlCallOpt []grpc.CallOption,
	getGUIDCallOpt []grpc.CallOption,
	healthCheckCallOpt []grpc.CallOption) (NewsParserClient, error) {

	cliS.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("%s:%d", host, port),
			dailOpts...,
		)
		if err != nil {
			cliS.Error = err
			return
		}

		cliS.cli = NewsParserClient{
			URLParserClient:    pb.NewURLParserClient(conn),
			ParseUrlCallOpt:    parseUrlCallOpt,
			HealthCheckCallOpt: healthCheckCallOpt,
			GetGUIDCallOpt:     getGUIDCallOpt,
		}
	})
	return cliS.cli, nil
}

func GetNewsParserClient() (NewsParserClient, error) {
	return cliS.cli, cliS.Error
}

type NewsParserClient struct {
	pb.URLParserClient
	ParseUrlCallOpt    []grpc.CallOption
	HealthCheckCallOpt []grpc.CallOption
	GetGUIDCallOpt     []grpc.CallOption
}

func (cli NewsParserClient) ParseURL(ctx context.Context, id int64, url string) (*pb.ParseURLResponse, error) {
	return cli.URLParserClient.ParseUrl(ctx, &pb.ParseURLRequest{Id: id, URL: url})
}

func (cli NewsParserClient) GetGUID(ctx context.Context, id int64, url string) (int64, string, error) {
	resp, err := cli.URLParserClient.GetGUID(ctx, &pb.GetGUIDRequest{Id: id, URL: url})
	if err != nil {
		return 0, "", err
	}
	return resp.Id, resp.GUID, err
}

func (cli NewsParserClient) HealthCheck(ctx context.Context) error {
	ping := &pb.PingPong{Signal: false}
	pong, err := cli.URLParserClient.HealthCheck(ctx, ping, cli.HealthCheckCallOpt...)
	if err != nil {
		return err
	}

	if ping.GetSignal() == pong.GetSignal() {
		return errors.New("unknown error, expect true but get false")
	}
	return nil
}
