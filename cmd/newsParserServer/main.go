package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/parser"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/news_parser"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var Logger zerolog.Logger

type NewsParserServer struct {
	pb.URLParserServer
	Parser parser.Parser
}

func NewPbError(ecErr *ec.Error) *pb.Error {
	pbErr := &pb.Error{}
	pbErr.Code = int64(ecErr.ErrorCode)
	pbErr.Message = ecErr.Message
	pbErr.Details = ecErr.Details
	return pbErr
}

func NewServer(opts ...parser.Parser) NewsParserServer {
	var p parser.Parser
	if opts == nil {
		p = parser.GetDefaultParser()
	} else {
		p = parser.NewParserRepo(opts...)
	}
	return NewsParserServer{Parser: p}
}

func (srvr NewsParserServer) ParseUrl(ctx context.Context, req *pb.ParseURLRequest) (*pb.ParseURLResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		Logger.Info().
			Str("url", req.URL).
			Int64("id", req.GetId()).
			Time("time", time.Now()).
			Msg("get ParseURLRequest")
		q := parser.NewQueryWithId(int(req.GetId()), req.GetURL())
		q = srvr.Parser.Parse(q)

		return &pb.ParseURLResponse{
			Id: req.GetId(),
			NewsItem: &pb.NewsItem{
				Title:       q.News.Title,
				Link:        q.News.Link.String(),
				GUID:        q.News.GUID,
				Description: q.News.Description,
				Language:    q.News.Language,
				Author:      q.News.Author,
				Category:    q.News.Category,
				Content:     q.News.Content,
				PubDate:     timestamppb.New(q.News.PubDate),
				Tag:         q.News.Tag,
				RelatedGUID: q.News.RelatedGUID,
			},
		}, q.Error
	}
}

func (srvr NewsParserServer) GetGUID(ctx context.Context, req *pb.GetGUIDRequest) (*pb.GetGUIDResponse, error) {
	Logger.Info().
		Str("url", req.URL).
		Int64("id", req.GetId()).
		Time("time", time.Now()).
		Msg("get GetGUIDRequest")
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		u, err := url.Parse(req.URL)
		if err != nil {
			return nil, err
		}

		return &pb.GetGUIDResponse{
			Id:   req.Id,
			GUID: srvr.Parser.ToGUID(u),
		}, nil
	}
}

func (srvr NewsParserServer) HealthCheck(ctx context.Context, sig *pb.PingPong) (*pb.PingPong, error) {
	Logger.Info().Time("time", time.Now()).Msg("ping")
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		sig.Signal = !sig.GetSignal()
		return sig, nil
	}
}

func main() {
	pflag.StringP("host", "h", "localhost", "ip/host name for language detect server")
	pflag.StringP("certificate", "c", "", "ssl cretificate file")
	pflag.StringP("private-key", "k", "", "ssl private key file")
	pflag.StringP("log-file", "o", "log.json", "log file")
	pflag.IntP("port", "p", 50052, "port for new parser server")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logfile, err := os.OpenFile(
		viper.GetString("log-file"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while OpenFile: %v", err)
		os.Exit(1)
	}
	Logger = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, logfile)).
		With().
		Timestamp().
		Logger()

	lstnr, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%d", viper.GetString("host"), viper.GetInt("port")))
	if err != nil {
		Logger.Fatal().
			Str("host", viper.GetString("host")).
			Int("port", viper.GetInt("port")).
			Err(err).
			Msg("Failed to listen on given address")
		os.Exit(1)
	}

	grpcOpts := []grpc.ServerOption{}
	if crt, key := viper.GetString("certificate"), viper.GetString("private-key"); crt != "" && key != "" {
		creds, err := credentials.NewServerTLSFromFile(crt, key)
		if err != nil {
			Logger.Fatal().
				Str("certificate", crt).
				Str("private-key", key).
				Err(err).
				Msg("Failed to new TLS Server from file")
			os.Exit(1)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	s := grpc.NewServer(grpcOpts...)
	pb.RegisterURLParserServer(s, NewServer())
	startAt := time.Now()
	Logger.Info().
		Str("address", lstnr.Addr().String()).
		Time("start_at", startAt).
		Msg("News Parser server start")

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func(signalChan chan os.Signal) {
		sig := <-signalChan
		Logger.Info().
			Str("signal", sig.String()).
			Msg("receive syscall signal")
		s.GracefulStop()

		endAt := time.Now()
		Logger.Info().
			Time("end_at", endAt).
			TimeDiff("up_time", endAt, startAt).
			Msg("server gracefully stopped")
	}(signalChan)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(*sync.WaitGroup) {
		if err = s.Serve(lstnr); err != nil {
			Logger.Fatal().
				Err(err).
				Msg("Failed to serve")
			os.Exit(1)
		}
		wg.Done()
	}(wg)

	wg.Wait()
	os.Exit(0)
}
