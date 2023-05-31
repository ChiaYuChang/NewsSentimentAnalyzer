package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
)

func (srvc jobService) Service() Service {
	return Service(srvc)
}

type JobCreateRequest struct {
	Owner    int32  `validate:"required,min=1"`
	Status   string `validate:"required,job_status"`
	SrcApiID int16  `validate:"required,min=1"`
	SrcQuery string `validate:"required,min=1,max=2048"`
	LlmApiID int16  `validate:"required,min=1"`
	LlmQuery string `validate:"required,min=1,max=2048"`
}

func (r JobCreateRequest) RequestName() string {
	return "job-create-req"
}

type JobDeleteRequest struct {
	ID    int32 `validate:"required,min=1"`
	Owner int32 `validate:"required,min=1"`
}

func (r JobDeleteRequest) RequestName() string {
	return "job-delete-req"
}

type JobGetJobsByOwnerRequest struct {
	Owner int32 `validate:"required,min=1"`
	N     int32 `validate:"required,min=1"`
}

func (r JobGetJobsByOwnerRequest) RequestName() string {
	return "job-get-by-owner-req"
}

type JobUpdateStatusRequest struct {
	Status string `validate:"required,min=1,job_status"`
	ID     int32  `validate:"required,min=1"`
	Owner  int32  `validate:"required,min=1"`
}

func (r JobUpdateStatusRequest) RequestName() string {
	return "job-udpate-status-req"
}

func (srvc jobService) Create(ctx context.Context, r *JobCreateRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}

	return srvc.Store.CreateJob(ctx, &model.CreateJobParams{
		Owner:    r.Owner,
		Status:   model.JobStatus(r.Status),
		SrcApiID: r.SrcApiID,
		SrcQuery: r.SrcQuery,
		LlmApiID: r.LlmApiID,
		LlmQuery: r.LlmQuery,
	})
}

func (srvc jobService) Delete(ctx context.Context, r *JobDeleteRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}
	return srvc.DeleteJob(ctx, &model.DeleteJobParams{
		ID: r.ID, Owner: r.Owner,
	})
}

func (srvc jobService) GetByOwner(ctx context.Context, r *JobGetJobsByOwnerRequest) ([]*model.GetJobsByOwnerRow, error) {
	if err := srvc.Validate.Struct(r); err != nil {
		return nil, err
	}

	return srvc.Store.GetJobsByOwner(ctx, &model.GetJobsByOwnerParams{
		Owner: r.Owner, N: r.N,
	})
}

func (srvc jobService) UpdateJobStatus(ctx context.Context, r *JobUpdateStatusRequest) error {
	if err := srvc.Validate.Struct(r); err != nil {
		return err
	}

	return srvc.Store.UpdateJobStatus(ctx, &model.UpdateJobStatusParams{
		Status: model.JobStatus(r.Status),
		ID:     r.ID,
		Owner:  r.Owner,
	})
}
