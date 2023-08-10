package service

import (
	"context"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
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

func (r JobDeleteRequest) ToParams() (*model.DeleteJobParams, error) {
	return &model.DeleteJobParams{
		ID: r.ID, Owner: r.Owner,
	}, nil
}

type JobGetByOwnerRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Next  int32     `validate:"min=0"`
	N     int32     `validate:"required,min=1"`
}

func (r JobGetByOwnerRequest) RequestName() string {
	return "job-get-by-owner-req"
}

func (r JobGetByOwnerRequest) ToParams() (*model.GetJobsByOwnerParams, error) {
	return &model.GetJobsByOwnerParams{
		Owner: r.Owner, Next: r.Next, N: r.N,
	}, nil
}

type JobGetByJobIdRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Id    int32     `validate:"required,min=1"`
}

func (r JobGetByJobIdRequest) RequestName() string {
	return "job-get-by-owner-req"
}

func (r JobGetByJobIdRequest) ToParams() (*model.GetJobsByJobIdParams, error) {
	return &model.GetJobsByJobIdParams{
		Owner: r.Owner,
		ID:    r.Id,
	}, nil
}

type JobUpdateStatusRequest struct {
	Status string    `validate:"required,min=1,job_status"`
	ID     int32     `validate:"required,min=1"`
	Owner  uuid.UUID `validate:"not_uuid_nil,uuid4"`
}

func (r JobUpdateStatusRequest) RequestName() string {
	return "job-udpate-status-req"
}

func (r JobUpdateStatusRequest) ToParams() (*model.UpdateJobStatusParams, error) {
	return &model.UpdateJobStatusParams{
		Status: model.JobStatus(r.Status),
		ID:     r.ID,
		Owner:  r.Owner,
	}, nil
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

func (srvc jobService) Delete(ctx context.Context, req *JobDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.DeleteJob(ctx, params)
}

func (srvc jobService) GetByOwner(ctx context.Context, req *JobGetByOwnerRequest) ([]*model.GetJobsByOwnerRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	return srvc.store.GetJobsByOwner(ctx, params)
}

func (srvc jobService) GetByJobId(ctx context.Context, req *JobGetByJobIdRequest) (*model.GetJobsByJobIdRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	return srvc.store.GetJobsByJobId(ctx, params)
}

func (srvc jobService) UpdateJobStatus(
	ctx context.Context, req *JobUpdateStatusRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}

	params, _ := req.ToParams()
	return srvc.store.UpdateJobStatus(ctx, params)
}

func (srvc jobService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpJobs(ctx)
}
