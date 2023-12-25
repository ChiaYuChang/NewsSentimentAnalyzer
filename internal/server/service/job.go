package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

func (srvc jobService) Service() Service {
	return Service(srvc)
}

type JobCreateRequest struct {
	Owner    uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Ulid     string    `validate:"required,len=26"`
	Status   string    `validate:"required,job_status"`
	SrcApiID int16     `validate:"required,min=1"`
	SrcQuery string    `validate:"required,min=1"`
	LlmApiID int16     `validate:"required,min=1"`
	LlmQuery string    `validate:"required,json"`
}

func (r JobCreateRequest) RequestName() string {
	return "job-create-req"
}

type JobDeleteRequest struct {
	ID    int64     `validate:"required,min=1"`
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

type JobGetWithStatusFilterRequest struct {
	Owner   uuid.UUID       `validate:"not_uuid_nil,uuid4"`
	Next    int32           `validate:"min=0"`
	N       int32           `validate:"required,min=1"`
	JStatus model.JobStatus `validate:"required"`
}

func (r JobGetWithStatusFilterRequest) RequestName() string {
	return "job-get-with-status-filter-req"
}

func (r JobGetWithStatusFilterRequest) ToParams() (*model.GetJobsByOwnerFilterByStatusParams, error) {
	return &model.GetJobsByOwnerFilterByStatusParams{
		Owner: r.Owner, Next: r.Next, N: r.N, JStatus: r.JStatus,
	}, nil
}

type JobGetByJobIdRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Id    int64     `validate:"required,min=1"`
}

func (r JobGetByJobIdRequest) RequestName() string {
	return "job-get-by-jid-req"
}

func (r JobGetByJobIdRequest) ToParams() (*model.GetJobsByJobIdParams, error) {
	return &model.GetJobsByJobIdParams{
		Owner: r.Owner,
		ID:    r.Id,
	}, nil
}

type JobUpdateStatusRequest struct {
	Status string    `validate:"required,min=1,job_status"`
	ID     int64     `validate:"required,min=1"`
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

type JobGetWithJIdRangeFilterRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	FJid  int32     `validate:"min=1"`
	TJid  int32     `validate:"min=1"`
	N     int32     `validate:"min=1"`
}

func (r JobGetWithJIdRangeFilterRequest) RequestName() string {
	return "job-get-with-jid-range-filter-req"
}

func (r JobGetWithJIdRangeFilterRequest) ToParams() (*model.GetJobByOwnerFilterByJIdRangeParams, error) {
	return &model.GetJobByOwnerFilterByJIdRangeParams{
		Owner: r.Owner,
		FJid:  r.FJid,
		TJid:  r.TJid,
		N:     r.N,
	}, nil
}

type JobGetWithJIdAndStatusFilterRequest struct {
	Owner   uuid.UUID `validate:"not_uuid_nil,uuid4"`
	FJid    int32     `validate:"min=1"`
	TJid    int32     `validate:"min=1"`
	N       int32     `validate:"min=1"`
	JStatus model.JobStatus
}

func (r JobGetWithJIdAndStatusFilterRequest) RequestName() string {
	return "job-get-with-jid-and-status-filter-req"
}

func (r JobGetWithJIdAndStatusFilterRequest) ToParams() (*model.GetJobByOwnerFilterByJIdAndStatusParams, error) {
	return &model.GetJobByOwnerFilterByJIdAndStatusParams{
		Owner:   r.Owner,
		FJid:    r.FJid,
		TJid:    r.TJid,
		JStatus: r.JStatus,
		N:       r.N,
	}, nil
}

type JobGetByJIdsRequest struct {
	Owner uuid.UUID `validate:"not_uuid_nil,uuid4"`
	Ids   []int32   `validate:"min=1"`
	N     int32     `validate:"min=1"`
}

func (r JobGetByJIdsRequest) RequestName() string {
	return "job-get-with-jids-filter-req"
}

func (r JobGetByJIdsRequest) ToParams() (*model.GetJobByOwnerFilterByJIdsParams, error) {
	return &model.GetJobByOwnerFilterByJIdsParams{
		Owner: r.Owner,
		Ids:   r.Ids,
		N:     r.N,
	}, nil
}

// create job
func (srvc jobService) Create(ctx context.Context, r *JobCreateRequest) (id int64, err error) {
	if err := srvc.validate.Struct(r); err != nil {
		return 0, err
	}

	return srvc.store.CreateJob(ctx, &model.CreateJobParams{
		Owner:    r.Owner,
		Ulid:     r.Ulid,
		SrcApiID: r.SrcApiID,
		SrcQuery: r.SrcQuery,
		LlmApiID: r.LlmApiID,
		LlmQuery: []byte(r.LlmQuery),
	})
}

// delete job by job id
func (srvc jobService) Delete(ctx context.Context, req *JobDeleteRequest) (n int64, err error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}
	params, _ := req.ToParams()
	return srvc.store.DeleteJob(ctx, params)
}

type GetJobRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    model.JobStatus    `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type JobRow interface {
	*model.GetJobsByOwnerRow |
		*model.GetJobByOwnerFilterByJIdsRow |
		*model.GetJobByOwnerFilterByJIdRangeRow |
		*model.GetJobByOwnerFilterByJIdAndStatusRow |
		*model.GetJobsByOwnerFilterByStatusRow
}

func ToJobRow[T JobRow](r T) (*GetJobRow, error) {
	jsn, err := json.Marshal(r)
	if err != nil {
		global.Logger.Debug().
			Err(err).
			Msg("error while Marshal")
		return nil, err
	}
	var row GetJobRow
	err = json.Unmarshal(jsn, &row)

	if err != nil {
		global.Logger.Debug().
			Err(err).
			Msg("error while Unmarshal")
		return nil, err
	}
	return &row, nil
}

func ToJobRows[T JobRow](rs []T) ([]*GetJobRow, error) {
	jobs := make([]*GetJobRow, len(rs))
	for i, r := range rs {
		j, err := ToJobRow(r)
		if err != nil {
			global.Logger.Debug().
				Int("i", i).
				Err(err).
				Msg("error while ToJobRows")
			return nil, err
		}
		jobs[i] = j
	}
	return jobs, nil
}

func (srvc jobService) Get(ctx context.Context, owner uuid.UUID, jids []int32,
	fjid, tjid int32, status string, n int32, page int) ([]*GetJobRow, error) {

	global.Logger.Debug().Msg("Get Jobs")
	if len(jids) > 0 {
		// query jobs by given jids
		if rs, err := srvc.store.GetJobByOwnerFilterByJIds(ctx, &model.GetJobByOwnerFilterByJIdsParams{
			Owner: owner,
			Ids:   jids,
			N:     n,
		}); err != nil {
			return nil, ParsePgxError(err)
		} else {
			return ToJobRows(rs)
		}
	}

	var js model.JobStatus
	_ = js.Scan(status)
	var hasJStatusFilter bool
	switch js {
	case model.JobStatusCreated, model.JobStatusRunning, model.JobStatusDone, model.JobStatusCanceled, model.JobStatusFailed:
		hasJStatusFilter = true
	default:
		hasJStatusFilter = false
	}

	if fjid == tjid {
		global.Logger.
			Debug().
			Str("status", status).
			Int32("fjid", fjid).
			Bool("hasfilter", hasJStatusFilter).
			Msg("Next page")
		if hasJStatusFilter {
			if rs, err := srvc.getByOwnerWithStatusFilter(ctx,
				&JobGetWithStatusFilterRequest{
					Owner:   owner,
					Next:    fjid,
					N:       n,
					JStatus: js,
				},
			); err != nil {
				global.Logger.Debug().Err(err).Msg("error while calling ..getByOwnerWithStatusFilter method")
				return nil, ParsePgxError(err,
					fmt.Sprintf("error while calling .getByOwnerWithStatusFilter: %s", err.Error()))
			} else {
				return ToJobRows(rs)
			}
		} else {
			if rs, err := srvc.getByOwner(ctx,
				&JobGetByOwnerRequest{
					Owner: owner,
					Next:  fjid,
					N:     n,
				},
			); err != nil {
				return nil, ParsePgxError(err,
					fmt.Sprintf("error while calling .getByOwner: %s", err.Error()))
			} else {
				return ToJobRows(rs)
			}
		}
	}

	global.Logger.Debug().Msg("Page Cache")
	if hasJStatusFilter {
		if rs, err := srvc.getByOwnerWithIdAndStatusFilter(ctx,
			&JobGetWithJIdAndStatusFilterRequest{
				Owner:   owner,
				FJid:    fjid,
				TJid:    tjid,
				N:       n,
				JStatus: js,
			},
		); err != nil {
			return nil, ParsePgxError(err,
				fmt.Sprintf("error while calling .GetJobByOwnerFilterByJIdAndStatus: %s", err.Error()))
		} else {
			return ToJobRows(rs)
		}
	} else {
		if rs, err := srvc.getByOwnerWithJIdRangeFilter(ctx,
			&JobGetWithJIdRangeFilterRequest{
				Owner: owner,
				FJid:  fjid,
				TJid:  tjid,
				N:     n,
			},
		); err != nil {
			return nil, ParsePgxError(err,
				fmt.Sprintf("error while calling .getByOwnerWithJIdRangeFilter: %s", err.Error()))
		} else {
			return ToJobRows(rs)
		}
	}
}

// get job filter by owner
func (srvc jobService) getByOwner(ctx context.Context, req *JobGetByOwnerRequest) ([]*model.GetJobsByOwnerRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobsByOwner(ctx, params)
	return rows, ParsePgxError(err)
}

// get job within given job ids
func (srvc jobService) getByJobIds(ctx context.Context,
	req *JobGetByJIdsRequest) ([]*model.GetJobByOwnerFilterByJIdsRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobByOwnerFilterByJIds(ctx, params)
	return rows, ParsePgxError(err)
}

// get job id with job between [f_jid, t_jid]
func (srvc jobService) getByOwnerWithJIdRangeFilter(ctx context.Context,
	req *JobGetWithJIdRangeFilterRequest) ([]*model.GetJobByOwnerFilterByJIdRangeRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, ParsePgxError(err)
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobByOwnerFilterByJIdRange(ctx, params)
	return rows, ParsePgxError(err)
}

// get job id with job between [f_jid, t_jid] and job status = jstatus
func (srvc jobService) getByOwnerWithIdAndStatusFilter(ctx context.Context,
	req *JobGetWithJIdAndStatusFilterRequest) ([]*model.GetJobByOwnerFilterByJIdAndStatusRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobByOwnerFilterByJIdAndStatus(ctx, params)
	return rows, ParsePgxError(err)
}

// get job id with job status = jstatus
func (srvc jobService) getByOwnerWithStatusFilter(ctx context.Context,
	req *JobGetWithStatusFilterRequest) ([]*model.GetJobsByOwnerFilterByStatusRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobsByOwnerFilterByStatus(ctx, params)
	return rows, ParsePgxError(err)
}

// get job details by job id
func (srvc jobService) GetDetails(ctx context.Context, req *JobGetByJobIdRequest) (*model.GetJobsByJobIdRow, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return nil, err
	}

	params, _ := req.ToParams()
	rows, err := srvc.store.GetJobsByJobId(ctx, params)
	return rows, ParsePgxError(err)
}

// update job status
func (srvc jobService) UpdateStatus(ctx context.Context,
	req *JobUpdateStatusRequest) (int64, error) {
	if err := srvc.validate.Struct(req); err != nil {
		return 0, err
	}

	params, _ := req.ToParams()
	n, err := srvc.store.UpdateJobStatus(ctx, params)
	return n, ParsePgxError(err)
}

// clean job with deleted_at
func (srvc jobService) CleanUp(ctx context.Context) (n int64, err error) {
	return srvc.store.CleanUpJobs(ctx)
}

// count job belong to given owner
func (srvc jobService) Count(ctx context.Context, owner uuid.UUID) (*model.CountUserJobTxResult, error) {
	result, err := srvc.store.DoCountUserJobTx(ctx, owner)
	return result, ParsePgxError(err)
}

func (srvc jobService) GetOldestNCreatedJobsForEachUser(ctx context.Context, n int) (
	// []*model.GetOldestNCreatedJobsForEachUserRow, error) {
	[]CreatedJobsRow, error) {

	if n > math.MaxInt32 {
		n = math.MaxInt32
	}
	rows, err := srvc.store.GetOldestNCreatedJobsForEachUser(ctx, int32(n))
	// return rows, ParsePgxError(err)

	objs := make([]CreatedJobsRow, len(rows))
	for i := range rows {
		objs[i] = CreatedJobsRow{rows[i]}
	}
	return objs, ParsePgxError(err)
}

type CreatedJobsRow struct {
	*model.GetOldestNCreatedJobsForEachUserRow
}

func (r CreatedJobsRow) String() string {

	obj := struct {
		model.GetOldestNCreatedJobsForEachUserRow
		AnalyzerOption *AnalyzerOption `json:"llm_query"`
	}{
		GetOldestNCreatedJobsForEachUserRow: *r.GetOldestNCreatedJobsForEachUserRow,
		AnalyzerOption:                      r.AnalyzerOptions(),
	}

	b, _ := json.MarshalIndent(obj, "", "    ")
	return string(b)
}

func (r CreatedJobsRow) AnalyzerOptions() *AnalyzerOption {
	var aopt *AnalyzerOption
	json.Unmarshal(r.LlmQuery, &aopt)
	return aopt
}

type AnalyzerOption struct {
	IsTestData               bool   `                  json:"is_test_data,omitempty" redis:"is_test_data"`
	APIName                  string `form:"api"        json:"api"                    redis:"api"`
	APIId                    int    `form:"llm-api-id" json:"id,omitempty"           redis:"id"`
	EmbeddingOptions         `json:"embedding-options,omitempty"                     redis:"embedding-options"`
	SentimentAnalysisOptions `json:"sentiment-analysis-options,omitempty"            redis:"sentiment-analysis-options"`
}

func (opt AnalyzerOption) String() string {
	return opt.ToString("", "    ")
}

func (opt AnalyzerOption) ToString(prefix, indent string) string {
	b, _ := json.MarshalIndent(opt, prefix, indent)
	return string(b)
}

type EmbeddingOptions struct {
	Embedding      bool   `form:"do-embedding"    json:"embedding"             redis:"embedding"`
	InputType      string `form:"input-type"      json:"input_type,omitempty" redis:"input_type"`
	EmbeddingModel string `form:"embedding-model" json:"embedding_model"       redis:"embedding_model"`
}

type SentimentAnalysisOptions struct {
	Sentiment bool   `form:"do-sentiment" json:"sentiment"              redis:"sentiment"`
	MaxTokens int    `form:"max-tokens"   json:"max_tokens,omitempty"  redis:"max_tokens"`
	Truncate  string `form:"truncate"     json:"truncate,omitempty"    redis:"truncate"`
}
