package main

import (
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/google/uuid"
)

var (
	TIME_MIN, _ = time.Parse(time.DateOnly, "2020-01-01")
	TIME_MAX    = time.Now().UTC()
)

var (
	NUM_USER     = 50
	MAX_NUM_JOBS = 100
	MIN_NUM_JOBS = 20
)

var (
	API_KEY_LENGTH    = []int{12, 32, 52}
	API_KEY_CHAR_SET  = []rune{'0', 'x', '-'}
	API_KEY_CHAR_PROB = []float64{0.6, 0.35, 0.05}
)

const (
	APITypeSource = model.ApiTypeSource
	APITypeLLM    = model.ApiTypeLanguageModel
)

const (
	JobStatusCreated  = model.JobStatusCreated
	JobStatusRunning  = model.JobStatusRunning
	JobStatusDone     = model.JobStatusDone
	JobStatusFailed   = model.JobStatusFailed
	JobStatusCanceled = model.JobStatusCanceled
)

const (
	RoleUser  = model.RoleUser
	RoleAdmin = model.RoleAdmin
)

const (
	TEST_ADMIN_USER_EAMIL    = "admin@example.com"
	TEST_ADMIN_USER_PASSWORD = "password"
	TEST_USER_EAMIL          = "test@example.com"
	TEST_USER_PASSWORD       = "password"
)

var (
	TEST_ADMIN_USER_UID, _ = uuid.Parse("3e3eb7b4-f040-4656-8e65-adb85eff07b1")
	TEST_USER_UID, _       = uuid.Parse("bd7d1378-9eda-4031-96f4-419656cbcb16")
)
