package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/google/uuid"
)

func (srvc jobService) Service() Service {
	return Service(srvc)
}

type JobCreateRequest struct {
	Owner    uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Status   string    `validate:"required,job_status"`
	SrcApiID int16     `validate:"required,min=1"`
	SrcQuery string    `validate:"required,url,min=1"`
	LlmApiID int16     `validate:"required,min=1"`
	LlmQuery string    `validate:"required,json"`
}

func (r JobCreateRequest) RequestName() string {
	return "job-create-req"
}

type JobDeleteRequest struct {
	ID    int32     `validate:"required,min=1"`
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
}

func (r JobDeleteRequest) RequestName() string {
	return "job-delete-req"
}

type JobGetJobsByOwnerRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	N     int32     `validate:"required,min=1"`
}

func (r JobGetJobsByOwnerRequest) RequestName() string {
	return "job-get-by-owner-req"
}

type JobUpdateStatusRequest struct {
	Status string    `validate:"required,min=1,job_status"`
	ID     int32     `validate:"required,min=1"`
	Owner  uuid.UUID `validate:"not_uuid_nil,uuid4"`
}

func (r JobUpdateStatusRequest) RequestName() string {
	return "job-udpate-status-req"
}

func (srvc jobService) Create(ctx context.Context, r *JobCreateRequest) (id int32, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}

	return srvc.store.CreateJob(ctx, &model.CreateJobParams{
		Owner:    r.Owner,
		Status:   model.JobStatus(r.Status),
		SrcApiID: r.SrcApiID,
		SrcQuery: r.SrcQuery,
		LlmApiID: r.LlmApiID,
		LlmQuery: []byte(r.LlmQuery),
	})
}

func (srvc jobService) Delete(ctx context.Context, r *JobDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}
	return srvc.store.DeleteJob(ctx, &model.DeleteJobParams{
		ID: r.ID, Owner: r.Owner,
	})
}

func (srvc jobService) GetByOwner(ctx context.Context, r *JobGetJobsByOwnerRequest) ([]*model.GetJobsByOwnerRow, error) {
	if err := srvc.validate.Struct(r); err != nil {
		return nil, err
	}

	return srvc.store.GetJobsByOwner(ctx, &model.GetJobsByOwnerParams{
		Owner: r.Owner, N: r.N,
	})
}

func (srvc jobService) UpdateJobStatus(
	ctx context.Context, r *JobUpdateStatusRequest) (n int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}

	return srvc.store.UpdateJobStatus(ctx, &model.UpdateJobStatusParams{
		Status: model.JobStatus(r.Status),
		ID:     r.ID,
		Owner:  r.Owner,
	})
}

func (srvc jobService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpJobs(ctx)
}
