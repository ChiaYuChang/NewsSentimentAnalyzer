package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGetOldestNCreatedJobsForEachUser(t *testing.T) {
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)

	conn, err := model.NewDBConnection(context.Background(), os.Getenv("POSTGRESQL_URL"))
	require.NoError(t, err)
	require.NotNil(t, conn)

	val, err := validator.GetDefaultValidate()
	require.NoError(t, err)
	require.NotNil(t, val)

	s := service.NewService(model.NewPGXStore(conn), val)
	require.NotNil(t, s)

	jobs, err := s.Job().GetOldestNCreatedJobsForEachUser(context.Background(), 2)
	require.NoError(t, err)
	require.NotNil(t, jobs)

	t.Logf("%d jobs", len(jobs))

	for i, j := range jobs {
		t.Logf("%d: %v", i, j)
		if i > 10 {
			break
		}
	}
}
