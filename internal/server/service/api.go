package service

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"golang.org/x/net/context"
)

func (srvc apiService) Service() Service {
	return Service(srvc)
}

type APICreateRequest struct {
	Name string `validate:"required,min=3,max=20"`
	Type string `validate:"required,api_type"`
}

func (r APICreateRequest) RequestName() string {
	return "admin-api-create-req"
}

type APIListRequest struct {
	N int `validate:"required,min=1,max=100"`
}

func (r APIListRequest) RequestName() string {
	return "admin-api-list-req"
}

type APIUpdateRequeset struct {
	Name string `validate:"required,min=3,max=20"`
	Type string `validate:"required,api_type"`
	ID   int16  `validate:"required,min=1"`
}

func (r APIUpdateRequeset) RequestName() string {
	return "admin-api-update-req"
}

type APIDeleteReqiest struct {
	ID int16 `validate:"required,min=1"`
}

func (r APIDeleteReqiest) RequestName() string {
	return "admin-api-delete-req"
}

type APICleanUpReqiest struct{}

func (r APICleanUpReqiest) RequestName() string {
	return "admin-api-clearn-up-req"
}

func (srvc apiService) Create(ctx context.Context, r *APICreateRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}

	return srvc.Store.CreateAPI(ctx, &model.CreateAPIParams{
		Name: r.Name, Type: model.ApiType(r.Type),
	})
}

func (srvc apiService) List(ctx context.Context, r *APIListRequest) ([]*model.ListAPIRow, error) {
	if err := srvc.Validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.Store.ListAPI(ctx, int32(r.N))
}

func (srvc apiService) Update(ctx context.Context, r *APIUpdateRequeset) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}
	return srvc.Store.UpdateAPI(ctx, &model.UpdateAPIParams{
		Name: r.Name,
		Type: model.ApiType(r.Type),
		ID:   r.ID,
	})
}

func (srvc apiService) Delete(ctx context.Context, r *APIDeleteReqiest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}
	return srvc.Store.DeleteAPI(ctx, r.ID)
}

func (srvc apiService) CleanUp(ctx context.Context, r *APICleanUpReqiest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}
	return srvc.CleanUpAPIs(ctx)
}
