package main

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

var JobStatus = []model.JobStatus{
	JobStatusCreated,
	JobStatusRunning,
	JobStatusDone,
	JobStatusFailed,
	JobStatusCanceled,
}

type Job struct {
	Item []JobItem
	N    int
}

type JobItem struct {
	Id         int             `json:"id"`
	ULID       string          `json:"ulid"`
	Owner      uuid.UUID       `json:"owner"`
	Status     model.JobStatus `json:"status"               mod:"trim"`
	SrcAPIName string          `json:"src_api_name"         mod:"trim"`
	SrcApiId   int             `json:"src_api_id"`
	SrcQuery   string          `json:"src_query,omitempty"  mod:"trim"`
	LlmAPIName string          `json:"llm_api_name"         mod:"trim"`
	LlmApiId   int             `json:"llm_api_id"`
	LlmQuery   string          `json:"llm_query,omitempty"  mod:"trim"`
	CreatedAt  time.Time       `json:"created_at,omitempty" mod:"empty"`
	UpdatedAt  time.Time       `json:"updated_at,omitempty" mod:"empty"`
	DeletedAt  NullableTime    `json:"-"`
}

type NullableTime struct {
	time.Time
	Valid bool
}

func (t NullableTime) Format(layout string) string {
	if !t.Valid {
		return "null"
	}
	return "'" + t.Time.Format(layout) + "'"
}

func NewJobs(maxJobN int, apis []APIItem, users []UserItem) Job {
	minJobN := MIN_NUM_JOBS

	srcApis := []int{}
	llmApis := []int{}
	for _, api := range apis {
		if api.Type == APITypeSource {
			srcApis = append(srcApis, api.Id)
		}
		if api.Type == APITypeLLM {
			llmApis = append(llmApis, api.Id)
		}
	}

	jobs := Job{}
	jobs.Item = make([]JobItem, 0, maxJobN*len(users))

	jss := NewSampler(JobStatus, []float64{0.3, 0.3, 0.3, 0.07, 0.03})
	for i := 0; i < len(users); i++ {
		owner := users[i]
		n := mrand.Intn(maxJobN-minJobN) + minJobN
		if owner.Id == TEST_ADMIN_USER_UID || owner.Id == TEST_USER_UID {
			n = maxJobN
		}

		for j := 0; j < n; j++ {
			srcApi := apis[srcApis[mrand.Intn(len(srcApis))]-1]
			llmApi := apis[llmApis[mrand.Intn(len(llmApis))]-1]

			jobs.Item = append(jobs.Item, JobItem{
				Owner:    owner.Id,
				ULID:     ulid.MustNew(ulid.Timestamp(time.Now()), crand.Reader).String(),
				Status:   jss.Get(),
				SrcApiId: srcApi.Id,
				SrcQuery: rg.Must(rg.AlphaNum.GenRdmString(mrand.Intn(10) + 20)),
				LlmApiId: llmApi.Id,
				LlmQuery: `{"is_test_data": true}`,
			})
		}
	}

	mrand.Shuffle(len(jobs.Item), func(i, j int) {
		jobs.Item[i], jobs.Item[j] = jobs.Item[j], jobs.Item[i]
	})
	cts := rg.GenRdnTimes(len(jobs.Item), TIME_MIN, TIME_MAX)
	for i := 0; i < len(jobs.Item); i++ {
		var ct, ut time.Time
		// var dt NullableTime
		ct = cts[i]
		switch jobs.Item[i].Status {
		default:
			ut = ct
		case JobStatusRunning, JobStatusDone, JobStatusFailed:
			ut = rg.GenRdnTime(cts[i], TIME_MAX)
		case JobStatusCanceled:
			ut = rg.GenRdnTime(cts[i], TIME_MAX)
			// dt = NullableTime{Time: ut, Valid: true}
		}
		jobs.Item[i].Id = i + 1
		jobs.Item[i].CreatedAt = cts[i]
		jobs.Item[i].UpdatedAt = ut
		// jobs[i].DeletedAt = dt
	}

	jobs.N = len(jobs.Item) + 1
	return jobs
}
