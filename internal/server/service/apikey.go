package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

func (srvc apikeyService) Service() Service {
	return Service(srvc)
}

type APIKeyCleanUpRequest struct{}

func (r APIKeyCleanUpRequest) RequestName() string {
	return "apikey-clean-up-req"
}

type APIKeyCreateRequest struct {
	Owner int32  `validate:"required,min=1"`
	ApiID int16  `validate:"required,min=1"`
	Key   string `validate:"required,min=1"`
}

func (r APIKeyCreateRequest) RequestName() string {
	return "apikey-create-req"
}

type APIKeyDeleteRequest struct {
	Owner int32 `validate:"required,min=1"`
	ApiID int16 `validate:"required,min=1"`
}

func (r APIKeyDeleteRequest) RequestName() string {
	return "apikey-delete-req"
}

type APIKeyGetRequest struct {
	ID    int32  `validate:"required,min=1"`
	Owner int32  `validate:"required,min=1"`
	ApiID int16  `validate:"required,min=1"`
	Key   string `validate:"required,min=1"`
}

func (r APIKeyGetRequest) RequestName() string {
	return "apikey-get-req"
}

type APIKeyListRequest struct {
	ID    int32  `validate:"required,min=1"`
	Owner int32  `validate:"required,min=1"`
	ApiID int16  `validate:"required,min=1"`
	Name  string `validate:"required,min=1"`
	Type  string `validate:"required,api_type"`
	Key   string `validate:"required"`
}

func (r APIKeyListRequest) RequestName() string {
	return "apikey-list-req"
}

type APIKeyUpdateRequest struct {
	Key      string `validate:"required,min=1"`
	Owner    int32  `validate:"required,min=1"`
	OldApiID int16  `validate:"required,min=1"`
	NewApiID int16  `validate:"required,min=1"`
}

func (r APIKeyUpdateRequest) RequestName() string {
	return "apikey-update-req"
}

func (srvc apikeyService) CleanUp(ctx context.Context) error {
	return srvc.Store.CleanUpAPIKey(ctx)
}

func (srvc apikeyService) Create(ctx context.Context, r *APIKeyCreateRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}

	return srvc.Store.CreateAPIKey(ctx, &model.CreateAPIKeyParams{
		Owner: r.Owner, ApiID: r.ApiID, Key: r.Key,
	})
}

func (srvc apikeyService) Delete(ctx context.Context, r *APIKeyDeleteRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}

	return srvc.Store.DeleteAPIKey(ctx, &model.DeleteAPIKeyParams{
		Owner: r.Owner, ApiID: r.ApiID,
	})
}

func (srvc apikeyService) Get(ctx context.Context, r *APIKeyGetRequest) ([]*model.GetAPIKeyRow, error) {
	if err := srvc.Validate.Struct(r); err != nil {
		return nil, err
	}

	return srvc.GetAPIKey(ctx, &model.GetAPIKeyParams{
		Owner: r.Owner, ApiID: r.ApiID,
	})
}

func (srvc apikeyService) List(ctx context.Context, r *APIKeyListRequest) ([]*model.ListAPIKeyRow, error) {
	if err := srvc.Validate.Struct(r); err != nil {
		return nil, err
	}
	return srvc.ListAPIKey(ctx, r.Owner)
}
