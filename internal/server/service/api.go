package service

import (
	"context"
	"errors"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (srvc apiService) Service() Service {
	return Service(srvc)
}

func (srvc apiService) Create(ctx context.Context, req *APICreateRequest) (id int16, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.CreateAPI(ctx, params)
}

func (srvc apiService) List(ctx context.Context, limit int32) ([]*model.ListAPIRow, error) {
	if err := srvc.validate.Var(limit, "required,min=1,max=100"); err != nil {
		return nil, err
	}
	return srvc.store.ListAPI(ctx, limit)
}

func (srvc apiService) Update(ctx context.Context, req *APIUpdateRequeset) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.UpdateAPI(ctx, params)
}

func (srvc apiService) Delete(ctx context.Context, id int16) (n int64, err error) {
	if err := srvc.validate.Var(id, "required,min=1"); err != nil {
		return 0, err
	}
	return srvc.store.DeleteAPI(ctx, id)
}

func (srvc apiService) Get(ctx context.Context, id int16) (*model.Api, error) {
	if err := srvc.validate.Var(id, "required,min=1"); err != nil {
		return nil, err
	}
	return srvc.store.GetAPI(ctx, id)
}

func (srvc apiService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpAPIs(ctx)
}

type APICreateRequest struct {
	Name string `validate:"required,min=3,max=20"`
	Type string `validate:"required,api_type"`
}

func (req APICreateRequest) RequestName() string {
	return "admin-api-create-req"
}

func (req APICreateRequest) ToParams() (*model.CreateAPIParams, error) {
	return &model.CreateAPIParams{
		Name: req.Name,
		Type: model.ApiType(req.Type),
	}, nil
}

type APIUpdateRequeset struct {
	Name string `validate:"required,min=3,max=20"`
	Type string `validate:"required,api_type"`
	ID   int16  `validate:"required,min=1"`
}

func (req APIUpdateRequeset) RequestName() string {
	return "admin-api-update-req"
}

func (req APIUpdateRequeset) ToParams() (*model.UpdateAPIParams, error) {
	return &model.UpdateAPIParams{
		Name: req.Name,
		Type: model.ApiType(req.Type),
		ID:   req.ID,
	}, nil
}

func ParsePgxError(err error, details ...string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ec.MustGetEcErr(ec.ECPgxErrNoRows)
	}

	if errors.Is(err, pgx.ErrTxClosed) {
		return ec.MustGetEcErr(ec.ECPgxErrTxClosed)
	}

	if errors.Is(err, pgx.ErrTxCommitRollback) {
		return ec.MustGetEcErr(ec.ECPgxErrTxCommitRollback)
	}

	var pgErr *pgconn.PgError
	var ecErr *ec.Error
	if errors.As(err, &pgErr) {
		ecErr = ec.NewErrorFromPgErr(pgErr)
	} else {
		ecErr = ec.MustGetEcErr(ec.ECServerError).
			WithDetails(err.Error())
	}
	return ecErr.WithDetails(details...)
}
