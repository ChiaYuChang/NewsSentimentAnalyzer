package pageform

import (
	"errors"
	"net/url"
	"reflect"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

func init() {
	validator.Validate.RegisterValidation(
		LocationValidator.Tag(),
		LocationValidator.ValFunc(),
	)
}

var ErrUnregisteredPageForm = errors.New("unregistered pageform")

type pageFormRepoKey [2]string

func NewPageFormRepoKey(api, endpoint string) pageFormRepoKey {
	return pageFormRepoKey{api, endpoint}
}

func (k pageFormRepoKey) API() string {
	return k[0]
}

func (k pageFormRepoKey) Endpoint() string {
	return k[1]
}

var PageFormRepo = make(pageFormRepo)

func Add(pf PageForm) {
	PageFormRepo.Add(pf)
}

func Set(api, endpoint string, pf PageForm) {
	PageFormRepo.Set(api, endpoint, pf)
}

func Get(api, endpoint string) (PageForm, error) {
	return PageFormRepo.Get(api, endpoint)
}

type pageFormRepo map[pageFormRepoKey]PageForm

func (pfr pageFormRepo) Add(pf PageForm) {
	pfr[NewPageFormRepoKey(pf.API(), pf.Endpoint())] = pf
}

func (pfr pageFormRepo) Set(api, endpoint string, pf PageForm) {
	pfr[NewPageFormRepoKey(api, endpoint)] = pf
}

func (pfr pageFormRepo) Get(api, endpoint string) (PageForm, error) {
	pf, ok := pfr[NewPageFormRepoKey(api, endpoint)]
	if !ok {
		return nil, ErrUnregisteredPageForm
	}
	return reflect.New(reflect.TypeOf(pf)).Interface().(PageForm), nil
}

type PageForm interface {
	Endpoint() string
	FormDecodeAndValidate(decoder *form.Decoder, val *val.Validate, postForm url.Values) (PageForm, error)
	API() string
	String() string
}

func FormDecode[T PageForm](decoder *form.Decoder, postForm url.Values) (T, error) {
	var pageForm T
	err := decoder.Decode(&pageForm, postForm)
	return pageForm, err
}

func FormValidate[T PageForm](val *val.Validate, pageForm T) error {
	err := val.Struct(pageForm)
	return err
}

func FormDecodeAndValidate[T PageForm](
	decoder *form.Decoder, val *val.Validate, postForm url.Values) (PageForm, error) {
	pageform, err := FormDecode[T](decoder, postForm)
	if err != nil {
		return nil, err
	}
	if err = FormValidate(val, pageform); err != nil {
		return nil, err
	}
	return pageform, nil
}
