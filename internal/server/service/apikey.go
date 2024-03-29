package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/google/uuid"
)

func (srvc apikeyService) Service() Service {
	return Service(srvc)
}

func (srvc apikeyService) Create(ctx context.Context, req *APIKeyCreateRequest) (id int32, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	id, err = srvc.store.CreateAPIKey(ctx, params)
	return id, ParsePgxError(err)
}

func (srvc apikeyService) Delete(ctx context.Context, req *APIKeyDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	n, err = srvc.store.DeleteAPIKey(ctx, params)
	return n, ParsePgxError(err)
}

func (srvc apikeyService) Get(ctx context.Context, r *APIKeyGetRequest) (*model.GetAPIKeyRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}

	row, err := srvc.store.GetAPIKey(
		ctx,
		&model.GetAPIKeyParams{
			Owner: r.Owner,
			ApiID: r.ApiID,
		})
	return row, ParsePgxError(err)
}

func (srvc apikeyService) List(ctx context.Context, owner uuid.UUID) ([]*model.ListAPIKeyRow, error) {
	if err := srvc.validate.Var(owner, "not_uuid_nil,uuid4"); err != nil {
		return nil, err
	}
	rows, err := srvc.store.ListAPIKey(ctx, owner)
	return rows, ParsePgxError(err)
}

func (srvc apikeyService) CleanUp(ctx context.Context) (n int64, err error) {
	n, err = srvc.store.CleanUpAPIKey(ctx)
	return n, ParsePgxError(err)
}

func (srvc apikeyService) CreateOrUpdate(ctx context.Context, r *APIKeyCreateOrUpdateRequest) (*model.CreateOrUpdateAPIKeyTxResults, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}
	params, _ := r.ToParams()
	return srvc.store.DoCreateOrUpdateAPIKeyTx(ctx, params)
}

type APIKeyCreateRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	ApiID int16     `validate:"required"`
	Key   string    `validate:"required"`
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
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	ApiID int16     `validate:"required,min=1"`
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
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	ApiID int16     `validate:"required,min=1"`
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
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	ApiID int16     `validate:"required,min=1"`
	Key   string    `validate:"required,min=32,max=64"`
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
