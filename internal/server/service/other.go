package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/google/uuid"
)

func (srvc otherService) Service() Service {
	return Service(srvc)
}

type NewsJobCreateRequest struct {
	JobID  int `validate:"required,min=1"`
	NewsID int `validate:"required,min=1"`
}

func (r NewsJobCreateRequest) RequestName() string {
	return "newsjob-create-req"
}

type LogCreateRequest struct {
	UserID  uuid.UUID `validate:"required,uuid4"`
	Type    string    `validate:"required,event_type"`
	Message string    `validate:"required"`
}

func (r LogCreateRequest) RequestName() string {
	return "log-create-req"
}

func (srvc otherService) CreateNewsJob(
	ctx context.Context, r *NewsJobCreateRequest) (id int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}

	return srvc.store.CreateNewsJob(ctx, &model.CreateNewsJobParams{
		JobID: int64(r.JobID), NewsID: int64(r.NewsID),
	})
}

func (srvc otherService) CreateEvent(
	ctx context.Context, r *LogCreateRequest) (int64, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}

	return srvc.store.CreateLog(ctx, &model.CreateLogParams{
		UserID:  r.UserID,
		Type:    model.EventType(r.Type),
		Message: r.Message,
	})
}
