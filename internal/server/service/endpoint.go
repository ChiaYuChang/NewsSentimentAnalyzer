package service

import (
	"context"
	"errors"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/jackc/pgx/v5"
)

func (srvc endpointService) Service() Service {
	return Service(srvc)
}

func (srvc endpointService) ListEndpointByOwner(ctx context.Context, owner int32) ([]*model.ListEndpointByOwnerRow, error) {
	if err := srvc.validate.Var(owner, "required,min=1"); err != nil {
		return nil, err
	}
	return srvc.store.ListEndpointByOwner(ctx, owner)
}

type EndpointCreateRequest struct {
	Name         string `validate:"required,max=32"`
	ApiID        int16  `validate:"required,min=1"`
	TemplateName string `validate:"required,max=32"`
}

func (req EndpointCreateRequest) RequestName() string {
	return "apikey-create-req"
}

func (req EndpointCreateRequest) ToParams() (*model.CreateEndpointParams, error) {
	return &model.CreateEndpointParams{
		Name:         req.Name,
		ApiID:        req.ApiID,
		TemplateName: req.TemplateName,
	}, nil
}

func (srvc endpointService) Create(ctx context.Context, req *EndpointCreateRequest) (int32, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.CreateEndpoint(ctx, params)
}

func (srvc endpointService) Delete(ctx context.Context, endpointId int32) (int64, error) {
	if err := srvc.validate.Var(endpointId, "required,min=1"); err != nil {
		return 0, err
	}
	return srvc.store.DeleteEndpoint(ctx, endpointId)
}

func (srvc endpointService) ListAll(ctx context.Context, limit int32, rowChan chan<- *model.ListAllEndpointRow) error {
	defer close(rowChan)
	params := model.ListAllEndpointParams{
		Limit: limit,
		Next:  0,
	}

	for {
		rows, err := srvc.store.ListAllEndpoint(ctx, &params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				break
			}
			return err
		}
		if len(rows) > 0 {
			params.Next = rows[len(rows)-1].EndpointID
		} else {
			break
		}

		for _, row := range rows {
			rowChan <- row
		}

	}
	return nil
}
