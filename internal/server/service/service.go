package service

import (
	"errors"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
)

var DefaultService *Service
var ErrEndOfData = errors.New("end of data")
var ErrInvalidParams = errors.New("invalid parameters")

type Service struct {
	store    model.Store
	validate *val.Validate
}

type Request interface {
	RequestName() string
}

func NewService(store model.Store, val *val.Validate) Service {
	return Service{store: store, validate: val}
}

func NewServiceWithDefautlVal(store model.Store) Service {
	return NewService(store, validator.Validate)
}

type authService Service

func (srvc Service) Auth() authService {
	return authService(srvc)
}

type userService Service

func (srvc Service) User() userService {
	return userService(srvc)
}

type apiService Service

func (srvc Service) API() apiService {
	return apiService(srvc)
}

type apikeyService Service

func (srvc Service) APIKey() apikeyService {
	return apikeyService(srvc)
}

type jobService Service

func (srvc Service) Job() jobService {
	return jobService(srvc)
}

type keywordService Service

func (srvc Service) Keyword() keywordService {
	return keywordService(srvc)
}

type newsService Service

func (srvc Service) News() newsService {
	return newsService(srvc)
}

type adminService Service

func (srvc Service) Admin() adminService {
	return adminService(srvc)
}

type otherService Service

func (srvc Service) NewsJob() otherService {
	return otherService(srvc)
}
