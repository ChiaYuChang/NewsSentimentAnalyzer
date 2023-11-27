package languagedetector_test

import (
	"context"
	"testing"
	"time"

	ld "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/languageDetector"
	pb "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto"
	"github.com/google/uuid"
	"github.com/pemistahl/lingua-go"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPingPong(t *testing.T) {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)

	cli := ld.LanguageDetectorClient{
		pb.NewLanguageDetectorClient(conn), nil, nil,
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = cli.HealthCheck(ctx)
	require.NoError(t, err)
}

func TestLanguageDetect(t *testing.T) {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	require.NotNil(t, conn)

	cli := ld.LanguageDetectorClient{
		pb.NewLanguageDetectorClient(conn), nil, nil,
	}

	type testCase struct {
		Id   uuid.UUID
		Text string
		Lang lingua.Language
	}

	tcs := []testCase{
		{
			Id:   uuid.New(),
			Text: "Hello, world! This sentence is written in English.",
			Lang: lingua.English,
		},
		{
			Id:   uuid.New(),
			Text: "你好世界！這句話是用中文寫的。",
			Lang: lingua.Chinese,
		},
		{
			Id:   uuid.New(),
			Text: "こんにちは世界！この文は日本語で書かれています。",
			Lang: lingua.Japanese,
		},
		{
			Id:   uuid.New(),
			Text: "¡Hola Mundo! Esta frase está escrita en español.",
			Lang: lingua.Spanish,
		},
		{
			Id:   uuid.New(),
			Text: "Bonjour le monde! Cette phrase est écrite en français.",
			Lang: lingua.French,
		},
		{
			Id:   uuid.New(),
			Text: "Hallo Welt! Dieser Satz ist auf Deutsch verfasst.",
			Lang: lingua.German,
		},
	}

	m := map[string]string{}
	reqChan := make(chan *pb.LanguageDetectRequest)
	go func(tcs []testCase, m map[string]string) {
		defer close(reqChan)
		thr := float64(0.9)
		langOpt := &pb.LanguageOption{
			LanguageOpt: []int64{
				int64(lingua.English),
				int64(lingua.Chinese),
				int64(lingua.Japanese),
				int64(lingua.Spanish),
				int64(lingua.French),
				int64(lingua.German),
			},
		}
		for i := range tcs {
			tc := tcs[i]
			m[tc.Id.String()] = tc.Lang.String()
			t.Logf("send %s\n", tc.Id.String())

			reqChan <- &pb.LanguageDetectRequest{
				Id:             tc.Id.String(),
				Text:           tc.Text,
				LanguageOption: langOpt,
				Threshold:      &thr,
			}
			langOpt = nil
		}
	}(tcs, m)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	respChan, err, errChan := cli.DetectLanguage(ctx, reqChan)
	require.NoError(t, err)
	defer cancel()

	for rOk, eOk, i := true, true, 0; (rOk || eOk) && i < 20; i++ {
		var resp *pb.LanguageDetectResponse
		var err error
		select {
		case resp, rOk = <-respChan:
			if rOk {
				t.Logf("resp: %s expected: %s actual: %s\n",
					resp.Id, m[resp.Id],
					lingua.Language(resp.Language).String())
				if resp.Language == int64(lingua.Unknown) {
					for _, c := range resp.GetConfidenceValue() {
						t.Logf("%s: %5.4f\n", c.Language, c.Value)
					}
				}
			}
		case err, eOk = <-errChan:
			if eOk {
				require.NoError(t, err)
			}
		}
	}
}
