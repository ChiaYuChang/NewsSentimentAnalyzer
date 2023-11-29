package service

import (
	"context"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/convert"
)

type NewsCreateRequest struct {
	Md5Hash     string    `validate:"required,md5"`
	Guid        string    `validate:"required"`
	Author      []string  `validate:"-"`
	Title       string    `validate:"required"`
	Link        string    `validate:"required,url"`
	Description string    `validate:"required"`
	Language    string    `validate:"-"`
	Content     []string  `validate:"required"`
	Category    string    `validate:"-"`
	Source      string    `validate:"required"`
	RelatedGuid []string  `validate:"-"`
	PublishedAt time.Time `validate:"required,before_now"`
}

func (r NewsCreateRequest) RequestName() string {
	return "news-create-req"
}

func (r NewsCreateRequest) ToParams() (*model.CreateNewsParams, error) {

	param := &model.CreateNewsParams{
		Md5Hash:     r.Md5Hash,
		Guid:        r.Guid,
		Author:      r.Author,
		Title:       r.Title,
		Link:        r.Link,
		Description: r.Description,
		Content:     r.Content,
		Category:    r.Category,
		Source:      r.Source,
		RelatedGuid: r.RelatedGuid,
	}

	if len(r.Author) > 0 {
		param.Author = make([]string, len(r.Author))
		copy(param.Author, r.Author)
	}

	if len(r.Content) > 0 {
		param.Content = make([]string, len(r.Content))
		copy(param.Content, r.Content)
	}

	if len(r.RelatedGuid) > 0 {
		param.RelatedGuid = make([]string, len(r.RelatedGuid))
		copy(param.RelatedGuid, r.RelatedGuid)
	}

	if r.Language != "" {
		param.Language = convert.StrTo(r.Language).PgText()
	}

	param.PublishAt = convert.TimeTo(r.PublishedAt).ToPgTimeStampZ()
	return param, nil
}

type NewsDeleteRequest struct {
	ID int64 `validate:"required,min=1"`
}

func (srvc newsService) Delete(ctx context.Context, r *NewsDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteNews(ctx, r.ID)
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

type NewsGetByMD5HashsRequest struct {
	MD5Hash []string `validate:"required"`
}

func (r NewsGetByMD5HashsRequest) RequestName() string {
	return "news-get-by-md5s-req"
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
	params, err := r.ToParams()
	if err != nil {
		return 0, err
	}
	return srvc.store.CreateNews(ctx, params)
}

func (srvc newsService) DeletePublishBefore(
	ctx context.Context, r *NewsDeletePublishBeforeRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteNewsPublishBefore(ctx, convert.TimeTo(r.Before).ToPgTimeStampZ())
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
		FromTime: convert.TimeTo(r.From).ToPgTimeStampZ(),
		ToTime:   convert.TimeTo(r.To).ToPgTimeStampZ(),
	})
}

func (srvc newsService) GetByMD5Hash(ctx context.Context, r *NewsGetByMD5HashRequest) (*model.GetNewsByMD5HashRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetNewsByMD5Hash(ctx, r.MD5Hash)
}

func (srvc newsService) GetByMD5Hashs(ctx context.Context, r *NewsGetByMD5HashsRequest) ([]int64, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.GetNewsByMD5Hashs(ctx, r.MD5Hash)
}

func (srvc newsService) ListRecentN(ctx context.Context, r *NewsListRequest) ([]*model.ListRecentNNewsRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.store.ListRecentNNews(ctx, r.N)
}
