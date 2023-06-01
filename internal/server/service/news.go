package service

import (
	"context"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

type NewsCreateRequest struct {
	Md5Hash     string    `validate:"required,md5"`
	Title       string    `validate:"required,min=1"`
	Url         string    `validate:"required,url"`
	Description string    `validate:"required,min=1"`
	Content     string    `validate:"required,min=1"`
	Source      string    `validate:"-"`
	PublishAt   time.Time `validate:"required,before_now"`
}

func (r NewsCreateRequest) RequestName() string {
	return "news-create-req"
}

type NewsDeleteRequest struct {
	ID int64 `validate:"required,min=1"`
}

type NewsDeletePublishBeforeRequest struct {
	Before time.Time `validate:"required"`
}

func (r NewsDeletePublishBeforeRequest) RequestName() string {
	return "news-delete-publish-before-req"
}

type NewsGetByKeywordsRequest struct {
	keywords []string `validate:"require,min=1"`
}

func (r NewsGetByKeywordsRequest) RequestName() string {
	return "news-get-by-keywords-req"
}

type NewsGetByPublishBetweenRequest struct {
	From time.Time `validate:"required"`
	To   time.Time `validate:"required"`
}

func (r NewsGetByPublishBetweenRequest) RequestName() string {
	return "news-get-by-publish-between-req"
}

type NewsGetByMD5HashRequest struct {
	MD5Hash string `validate:"required,md5"`
}

func (r NewsGetByMD5HashRequest) RequestName() string {
	return "news-get-by-md5-req"
}

type NewsListRequest struct {
	N int32 `validate:"required,min=1"`
}

func (r NewsListRequest) RequestName() string {
	return "news-list-req"
}

func (srvc newsService) Create(ctx context.Context, r *NewsCreateRequest) (id int64, err error) {
	if err := srvc.validate.Struct(srvc); err != nil {
		return 0, err
	}

	return srvc.store.CreateNews(ctx, &model.CreateNewsParams{
		Md5Hash:     r.Md5Hash,
		Title:       r.Title,
		Url:         r.Url,
		Description: r.Description,
		Content:     r.Content,
		Source:      StringToText(r.Source),
		PublishAt:   TimeToTimestamptz(r.PublishAt),
	})
}

func (srvc newsService) Delete(
	ctx context.Context, r *NewsDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteNews(ctx, r.ID)
}

func (srvc newsService) DeletePublishBefore(
	ctx context.Context, r *NewsDeletePublishBeforeRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteNewsPublishBefore(ctx, TimeToTimestamptz(r.Before))
}

func (srvc newsService) GetByKeywords(ctx context.Context, r *NewsGetByKeywordsRequest) ([]*model.GetNewsByKeywordsRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetNewsByKeywords(ctx, r.keywords)
}

func (srvc newsService) GetByPublishBetween(ctx context.Context, r *NewsGetByPublishBetweenRequest) ([]*model.GetNewsPublishBetweenRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetNewsPublishBetween(ctx, &model.GetNewsPublishBetweenParams{
		FromTime: TimeToTimestamptz(r.From),
		ToTime:   TimeToTimestamptz(r.To),
	})
}

func (srvc newsService) GetByMD5Hash(ctx context.Context, r *NewsGetByMD5HashRequest) (*model.GetNewsByMD5HashRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetNewsByMD5Hash(ctx, r.MD5Hash)
}

func (srvc newsService) ListRecentN(ctx context.Context, r *NewsListRequest) ([]*model.ListRecentNNewsRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.ListRecentNNews(ctx, r.N)
}
