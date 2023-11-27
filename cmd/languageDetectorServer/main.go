package main

import (
	"context"
	"fmt"
	"io"

	// "log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto"
	"github.com/pemistahl/lingua-go"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Logger zerolog.Logger

type LanguageDetectServer struct {
	pb.LanguageDetectorServer
	Detector lingua.LanguageDetector
}

func NewPbError(ecErr *ec.Error) *pb.Error {
	pbErr := &pb.Error{}
	pbErr.Code = int64(ecErr.ErrorCode)
	pbErr.Message = ecErr.Message
	pbErr.Details = ecErr.Details
	return pbErr
}

func NewServer(lang []lingua.Language) LanguageDetectServer {
	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(lang...).
		Build()

	return LanguageDetectServer{Detector: detector}
}

func (srvr LanguageDetectServer) DetectLanguage(stream pb.LanguageDetector_DetectLanguageServer) error {
	detector := srvr.Detector
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		Logger.Debug().
			Str("id", req.GetId()).
			Int("len", len(req.GetText())).
			Msgf("get request")

		resp := &pb.LanguageDetectResponse{
			Id:          req.GetId(),
			Probability: 0.00,
			Language:    int64(lingua.Unknown),
		}

		if err != nil {
			ecErr := ec.MustGetEcErr(ec.ECgRPCServerError).
				WithDetails(fmt.Sprintf("error while reading client stream %v", err)).
				WithDetails(err.Error())
			return ecErr
		}

		if resp.Error == nil {
			select {
			case <-stream.Context().Done():
				resp.Language = int64(lingua.Unknown)
				resp.Probability = 0.00
				ecErr := ec.MustGetEcErr(ec.ECgRPCServerError).
					WithDetails(stream.Context().Err().Error())
				resp.Error = NewPbError(ecErr)
			default:
				if req.GetLanguageOption() != nil {
					Logger.Debug().
						Ints64("language_id", req.GetLanguageOption().GetLanguageOpt()).
						Msg("new language server")
					langOpt := []lingua.Language{}
					for _, lang := range req.GetLanguageOption().GetLanguageOpt() {
						langOpt = append(langOpt, lingua.Language(lang))
					}
					detector = lingua.NewLanguageDetectorBuilder().
						FromLanguages(langOpt...).
						Build()
				}

				lcvs := detector.ComputeLanguageConfidenceValues(req.GetText())
				for _, lcv := range lcvs {
					prob := lcv.Value()
					lang := lcv.Language()

					if prob > resp.Probability {
						resp.Language = int64(lang)
						resp.Probability = prob
					}
					resp.ConfidenceValue = append(
						resp.ConfidenceValue,
						&pb.ConfidenceValue{
							Language: lang.IsoCode639_1().String(),
							Value:    prob,
						})
				}

				if resp.Probability < req.GetThreshold() {
					resp.Language = int64(lingua.Unknown)
				}

			}
		}

		Logger.Debug().
			Str("language", lingua.Language(resp.Language).String()).
			Int64("language_id", resp.GetLanguage()).
			Msg("send response")
		err = stream.Send(resp)
		if err != nil {
			ecErr := ec.MustGetEcErr(ec.ECgRPCServerError).
				WithDetails(fmt.Sprintf("error while sending data to client: %v", err)).
				WithDetails(err.Error())
			return ecErr
		}
	}
}

func (srvr LanguageDetectServer) HealthCheck(ctx context.Context, sig *pb.PingPong) (*pb.PingPong, error) {
	Logger.Info().Msg("ping")
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		sig.Signal = !sig.GetSignal()
		return sig, nil
	}
}

func main() {
	defaultLangs := []string{}
	for _, lang := range lingua.AllLanguages() {
		defaultLangs = append(defaultLangs, lang.IsoCode639_1().String())
	}

	pflag.StringP("host", "h", "localhost", "ip/host name for language detect server")
	pflag.StringP("certificate", "c", "", "ssl cretificate file")
	pflag.StringP("private-key", "k", "", "ssl private key file")
	pflag.StringP("log-file", "o", "log.json", "log file")

	pflag.IntP("port", "p", 50051, "port for language detect server")
	pflag.StringSliceP("supprot-language", "l", defaultLangs, "support language for the server")

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

	langs := []lingua.Language{}
	for _, langCode := range viper.GetStringSlice("supprot-language") {
		langs = append(langs, lingua.GetLanguageFromIsoCode639_1(
			lingua.GetIsoCode639_1FromValue(langCode)))
	}
	Logger.Info().
		Strs("support language", viper.GetStringSlice("supprot-language")).
		Msg("support language")

	s := grpc.NewServer(grpcOpts...)
	pb.RegisterLanguageDetectorServer(s, NewServer(langs))
	startAt := time.Now()
	Logger.Info().
		Str("address", lstnr.Addr().String()).
		Time("start_at", startAt).
		Msg("Language detect server start")

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
