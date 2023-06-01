package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

type KeywordCreateRequest struct {
	NewsID  int64  `validate:"required,min=1"`
	Keyword string `validate:"required,min=1,max=50"`
}

func (r KeywordCreateRequest) RequestName() string {
	return "key-create-req"
}

type KeywordsGetByNewsIdRequest struct {
	NewsID []int32 `validate:"required,min=1"`
}

func (r KeywordsGetByNewsIdRequest) RequestName() string {
	return "key-get-by-news-id-req"
}

type KeywordDeleteRequest struct {
	Keyword string `validate:"required,min=1"`
}

func (r KeywordDeleteRequest) RequestName() string {
	return "key-delete-req"
}

func (srvc keywordService) Create(
	ctx context.Context, r *KeywordCreateRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.CreateKeyword(ctx, &model.CreateKeywordParams{
		NewsID:  r.NewsID,
		Keyword: r.Keyword,
	})
}

func (srvc keywordService) GetByNewsId(
	ctx context.Context, r *KeywordsGetByNewsIdRequest) ([]string, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetKeywordsByNewsId(ctx, r.NewsID)
}

func (srvc keywordService) Delete(
	ctx context.Context, r *KeywordDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteKeyword(ctx, r.Keyword)
}
