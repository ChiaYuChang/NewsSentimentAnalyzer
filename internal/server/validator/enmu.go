package validator

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	val "github.com/go-playground/validator/v10"
)

type Enmus[T ~string] struct {
	tag  string
	eMap map[T]struct{}
}

func NewEnmus[T ~string](tag string, es ...T) Enmus[T] {
	eMap := make(map[T]struct{}, len(es))
	for _, e := range es {
		eMap[e] = struct{}{}
	}
	return Enmus[T]{tag: tag, eMap: eMap}
}

func (e Enmus[T]) ValFun() val.Func {
	return func(fl val.FieldLevel) bool {
		v, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		_, ok = e.eMap[T(v)]
		if !ok {
			return false
		}
		return true
	}
}

func (e Enmus[T]) Tag() string {
	return e.tag
}

var EnmusRole = NewEnmus(
	"role",
	model.RoleUser,
	model.RoleAdmin,
)

var EnmusJobStatus = NewEnmus(
	"job_status",
	model.JobStatusCreated,
	model.JobStatusRunning,
	model.JobStatusDone,
	model.JobStatusFailure,
	model.JobStatusCanceled,
)

var EnmusApiType = NewEnmus(
	"api_type",
	model.ApiTypeLanguageModel,
	model.ApiTypeSource,
)

var EnmusEventType = NewEnmus(
	"event_type",
	model.EventTypeSignIn,
	model.EventTypeSignOut,
	model.EventTypeAuthorization,
	model.EventTypeApiKey,
	model.EventTypeQuery,
)
