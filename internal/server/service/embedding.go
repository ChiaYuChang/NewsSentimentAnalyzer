package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/pgvector/pgvector-go"
)

type CreateEmbeddingRequest struct {
	NewsId    int64           `validate:"required,min=1"`
	Model     string          `validate:"required,max=32"`
	Embedding []float32       `validate:"required,max=1536"`
	Sentiment model.Sentiment `validate:"required,oneof=positive negative neutral"`
}

func (req CreateEmbeddingRequest) RequestName() string {
	return "embedding-create-req"
}

func (req CreateEmbeddingRequest) ToParams() (*model.CreateEmbeddingParams, error) {
	return &model.CreateEmbeddingParams{
		NewsID:    req.NewsId,
		Model:     req.Model,
		Embedding: pgvector.NewVector(req.Embedding),
		Sentiment: req.Sentiment,
	}, nil
}

// Create creates a new embedding
func (srvc embeddingService) Create(ctx context.Context, req *CreateEmbeddingRequest) (int64, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}

	params, _ := req.ToParams()
	id, err := srvc.store.CreateEmbedding(ctx, params)
	return id, ParsePgxError(err)
}

// GetByJobId returns all embeddings by job id
func (srvc embeddingService) GetByJobId(ctx context.Context, jobId int64) ([]*model.GetEmbeddingByJobIdRow, error) {
	if err := srvc.validate.Var(jobId, "required,min=1"); err != nil {
		return nil, err
	}

	rows, err := srvc.store.GetEmbeddingByJobId(ctx, jobId)
	return rows, ParsePgxError(err)
}

type GetEmbeddingByNewsIdsAndModelRequest struct {
	Model   string  `validate:"required,max=32"`
	NewsIds []int32 `validate:"required,min=1,dive,min=1"`
}

func (req GetEmbeddingByNewsIdsAndModelRequest) RequestName() string {
	return "embedding-get-by-news-ids-and-model-req"
}

func (req GetEmbeddingByNewsIdsAndModelRequest) ToParams() (*model.GetEmbeddingByNewsIdsAndModelParams, error) {
	return &model.GetEmbeddingByNewsIdsAndModelParams{
		Model:   req.Model,
		NewsIds: req.NewsIds,
	}, nil
}

func (srvc embeddingService) GetEmbeddingByNewsIdsAndModel(ctx context.Context, req *GetEmbeddingByNewsIdsAndModelRequest) ([]*model.GetEmbeddingByNewsIdsAndModelRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetEmbeddingByNewsIdsAndModel(ctx, params)
	return rows, ParsePgxError(err)
}
