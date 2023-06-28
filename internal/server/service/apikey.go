package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

func (srvc apikeyService) Service() Service {
	return Service(srvc)
}

func (srvc apikeyService) Create(ctx context.Context, req *APIKeyCreateRequest) (id int32, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.CreateAPIKey(ctx, params)
}

func (srvc apikeyService) Delete(ctx context.Context, req *APIKeyDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.DeleteAPIKey(ctx, params)
}

func (srvc apikeyService) Get(ctx context.Context, r *APIKeyGetRequest) (*model.GetAPIKeyRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}

	return srvc.store.GetAPIKey(
		ctx,
		&model.GetAPIKeyParams{
			Owner: r.Owner,
			ApiID: r.ApiID,
		})
}

func (srvc apikeyService) List(ctx context.Context, owner int32) ([]*model.ListAPIKeyRow, error) {
	if err := srvc.validate.Var(owner, "required,min=1"); err != nil {
		return nil, err
	}
	return srvc.store.ListAPIKey(ctx, owner)
}

func (srvc apikeyService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpAPIKey(ctx)
}

func (srvc apikeyService) CreateOrUpdate(ctx context.Context, r *APIKeyCreateOrUpdateRequest) (*model.CreateOrUpdateAPIKeyTxResults, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	params, _ := r.ToParams()
	return srvc.store.DoCreateOrUpdateAPIKeyTx(ctx, params)
}

type APIKeyCreateRequest struct {
	Owner int32  `validate:"required"`
	ApiID int16  `validate:"required"`
	Key   string `validate:"required"`
}

func (req APIKeyCreateRequest) RequestName() string {
	return "apikey-create-req"
}

func (req APIKeyCreateRequest) ToParams() (*model.CreateAPIKeyParams, error) {
	return &model.CreateAPIKeyParams{
		Owner: req.Owner,
		ApiID: req.ApiID,
		Key:   req.Key,
	}, nil
}

type APIKeyDeleteRequest struct {
	Owner int32 `validate:"required,min=1"`
	ApiID int16 `validate:"required,min=1"`
}

func (req APIKeyDeleteRequest) RequestName() string {
	return "apikey-delete-req"
}

func (req APIKeyDeleteRequest) ToParams() (*model.DeleteAPIKeyParams, error) {
	return &model.DeleteAPIKeyParams{
		Owner: req.Owner,
		ApiID: req.ApiID,
	}, nil
}

type APIKeyGetRequest struct {
	Owner int32 `validate:"required,min=1"`
	ApiID int16 `validate:"required,min=1"`
}

func (req APIKeyGetRequest) RequestName() string {
	return "apikey-get-req"
}

func (req APIKeyGetRequest) ToParams() *APIKeyGetRequest {
	return &APIKeyGetRequest{
		Owner: req.Owner,
		ApiID: req.ApiID,
	}
}

type APIKeyUpdateRequest struct {
	Key      string `validate:"required,min=1"`
	Owner    int32  `validate:"required,min=1"`
	OldApiID int16  `validate:"required,min=1"`
	NewApiID int16  `validate:"required,min=1"`
}

func (req APIKeyUpdateRequest) RequestName() string {
	return "apikey-update-req"
}

type APIKeyCreateOrUpdateRequest struct {
	Owner int32  `validate:"required,min=1"`
	ApiID int16  `validate:"required,min=1"`
	Key   string `validate:"required,min=32,max=64"`
}

func (req APIKeyCreateOrUpdateRequest) RequestName() string {
	return "apikey-update/create-req"
}

func (req APIKeyCreateOrUpdateRequest) ToParams() (*model.CreateOrUpdateAPIKeyTxParams, error) {
	return &model.CreateOrUpdateAPIKeyTxParams{
		Owner: req.Owner,
		ApiID: req.ApiID,
		Key:   req.Key,
	}, nil
}
