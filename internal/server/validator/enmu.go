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

func NewEnmusFromMap[T ~string](tag string, m map[T]T, which string) Enmus[T] {
	es := make([]T, 0, len(m))
	for key, val := range m {
		if which == "val" {
			es = append(es, T(val))
		} else {
			es = append(es, T(key))
		}
	}
	return NewEnmus[T](tag, es...)
}

func (e Enmus[T]) Map() map[T]struct{} {
	return e.eMap
}

func (e Enmus[T]) ValFun() val.Func {
	return func(fl val.FieldLevel) bool {
		switch val := fl.Field().Interface().(type) {
		case string:
			if val != "" {
				_, ok := e.eMap[T(val)]
				if !ok {
					return false
				}
			}
		case []string:
			for _, v := range val {
				if v == "" {
					continue
				}
				_, ok := e.eMap[T(v)]
				if !ok {
					return false
				}
			}
		default:
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
