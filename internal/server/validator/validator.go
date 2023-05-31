package validator

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	val "github.com/go-playground/validator/v10"
)

var Validate *val.Validate

func init() {
	Validate = val.New()
}

func RegisterEnmusValidator(val *val.Validate) {
	enmusRole := NewEnmus(
		"role",
		model.RoleUser,
		model.RoleAdmin,
	)
	val.RegisterValidation(
		enmusRole.Tag(),
		enmusRole.ValFun(),
	)

	enmusJobStatus := NewEnmus(
		"job_status",
		model.JobStatusCreated,
		model.JobStatusRunning,
		model.JobStatusDone,
		model.JobStatusFailure,
		model.JobStatusCanceled,
	)
	val.RegisterValidation(
		enmusJobStatus.Tag(),
		enmusJobStatus.ValFun(),
	)

	enmusApiType := NewEnmus(
		"api_type",
		model.ApiTypeLanguageModel,
		model.ApiTypeSource,
	)
	val.RegisterValidation(
		enmusApiType.Tag(),
		enmusApiType.ValFun(),
	)

	enmusEventType := NewEnmus(
		"event_type",
		model.EventTypeSignIn,
		model.EventTypeSignOut,
		model.EventTypeAuthorization,
		model.EventTypeApiKey,
		model.EventTypeQuery,
	)
	val.RegisterValidation(
		enmusEventType.Tag(),
		enmusEventType.ValFun(),
	)
}

func RegisterPasswordValidator(pv PasswordValidator, val *val.Validate) {
	val.RegisterValidation(
		pv.Tag(),
		pv.PasswordValFun(),
	)
}
