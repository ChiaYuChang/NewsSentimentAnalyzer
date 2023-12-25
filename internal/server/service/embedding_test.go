package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	dburl := fmt.Sprintf("postgres://admin:%s@localhost:5434/nsa_development?sslmode=disable", "admin")
	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	val, _ := validator.GetDefaultValidate()
	srvc := service.NewService(model.NewPGXStore(conn), val)

	embd := make([]float32, 1536)
	for i := 0; i < 384; i++ {
		embd[i] = rand.Float32()
	}

	req := &service.CreateEmbeddingRequest{
		NewsId:    1,
		Model:     "embed-multilingual-light-v3.0",
		Embedding: embd,
		Sentiment: model.SentimentPositive,
	}

	id, err := srvc.Embedding().Create(context.Background(), req)
	require.NoError(t, err)
	require.NotEqual(t, 0, id)
}
