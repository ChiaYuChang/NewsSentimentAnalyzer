package languagedetector

import (
	"context"
	"errors"
	"fmt"
	"io"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/language_detector"
	"google.golang.org/grpc"
)

func NewPbError(ecErr *ec.Error) *pb.Error {
	pbErr := &pb.Error{}
	pbErr.Code = int64(ecErr.ErrorCode)
	pbErr.Message = ecErr.Message
	pbErr.Details = ecErr.Details
	return pbErr
}

type LanguageDetectorClient struct {
	pb.LanguageDetectorClient
	DetectLanguageCallOpt []grpc.CallOption
	HealthCheckCallOpt    []grpc.CallOption
}

type LanguageDetectorResponse struct {
	pb.LanguageDetectResponse
	*ec.Error
}

var ErrNil = errors.New("receive/send a nil object")

func (cli LanguageDetectorClient) sendRequest(ctx context.Context, c <-chan *pb.LanguageDetectRequest) (
	pb.LanguageDetector_DetectLanguageClient, error, <-chan error) {
	stream, err := cli.LanguageDetectorClient.DetectLanguage(
		context.Background(), cli.DetectLanguageCallOpt...)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECgRPCClientError).
			WithDetails("error while calling DetectLanguage").
			WithDetails(err.Error())
		return nil, ecErr, nil
	}

	es := make(chan error)
	go func(ctx context.Context, c <-chan *pb.LanguageDetectRequest, e chan<- error) {
		defer stream.CloseSend()
		defer close(e)
		for {
			select {
			case <-ctx.Done():
				e <- ctx.Err()
				return
			case req, ok := <-c:
				if !ok {
					return
				}

				if req == nil {
					continue
				}

				err := stream.Send(req)
				if err != nil {
					ecErr := ec.MustGetEcErr(ec.ECgRPCClientError).
						WithDetails("stream.Send error").
						WithDetails(fmt.Sprintf("object id: %s", req.Id)).
						WithDetails(err.Error())
					e <- ecErr
				}
			}
		}
	}(ctx, c, es)
	return stream, nil, es
}

func (cli LanguageDetectorClient) recvResponse(ctx context.Context, stream pb.LanguageDetector_DetectLanguageClient) (
	<-chan *pb.LanguageDetectResponse, <-chan error) {
	c := make(chan *pb.LanguageDetectResponse)
	e := make(chan error)

	go func(stream pb.LanguageDetector_DetectLanguageClient, c chan<- *pb.LanguageDetectResponse, e chan<- error) {
		defer close(c)
		defer close(e)
		for {
			select {
			case <-ctx.Done():
				e <- ctx.Err()
				return
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					return
				}

				if err != nil {
					ecErr := ec.MustGetEcErr(ec.ECgRPCClientError).
						WithDetails("stream.Recv error")
					if resp != nil {
						ecErr.WithDetails(fmt.Sprintf("object id: %s", resp.Id))
					}
					ecErr.WithDetails(err.Error())
					e <- ecErr
					continue
				}

				if resp == nil {
					continue
				}
				c <- resp
			}
		}
	}(stream, c, e)
	return c, e
}

func (cli LanguageDetectorClient) mergeErrors(e1, e2 <-chan error) <-chan error {
	e3 := make(chan error)
	go func() {
		var e1Ok, e2Ok bool = true, true
		for e1Ok || e2Ok {
			var e1Err, e2Err error
			select {
			case e1Err, e1Ok = <-e1:
				if e1Ok && e1Err != nil {
					e3 <- e1Err
				}
			case e2Err, e2Ok = <-e2:
				if e2Ok && e2Err != nil {
					e3 <- e2Err
				}
			}
		}
		close(e3)
	}()
	return e3
}

func (cli LanguageDetectorClient) DetectLanguage(ctx context.Context,
	reqChan <-chan *pb.LanguageDetectRequest) (<-chan *pb.LanguageDetectResponse, error, <-chan error) {

	stream, err, sendErrChan := cli.sendRequest(ctx, reqChan)
	if err != nil {
		return nil, err, nil
	}
	respChan, recvErrChan := cli.recvResponse(ctx, stream)

	return respChan, nil, cli.mergeErrors(sendErrChan, recvErrChan)
}

func (cli LanguageDetectorClient) HealthCheck(ctx context.Context) error {
	ping := &pb.PingPong{Signal: false}
	pong, err := cli.LanguageDetectorClient.HealthCheck(ctx, ping, cli.HealthCheckCallOpt...)
	if err != nil {
		return err
	}

	if ping.GetSignal() == pong.GetSignal() {
		return errors.New("unknown error, expect true but get false")
	}
	return nil
}
