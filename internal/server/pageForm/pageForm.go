package pageform

import (
	"errors"
	"net/url"
	"reflect"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-playground/form"
	"github.com/go-playground/mold/v4"
	val "github.com/go-playground/validator/v10"
)

func init() {
	validator.Validate.RegisterValidation(
		LocationValidator.Tag(),
		LocationValidator.ValFunc(),
	)

	Decoder = form.NewDecoder()
	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		if len(vals[0]) == 0 {
			return time.Time{}, nil
		}
		return time.Parse(time.DateOnly, vals[0])
	}, time.Time{})

	Modifier = mold.New()
}

var Decoder *form.Decoder
var Modifier *mold.Transformer

var ErrUnregisteredPageForm = errors.New("unregistered pageform")

type PageFormRepoKey [2]string

func NewPageFormRepoKey(api, endpoint string) PageFormRepoKey {
	return PageFormRepoKey{api, endpoint}
}

func (k PageFormRepoKey) API() string {
	return k[0]
}

func (k PageFormRepoKey) Endpoint() string {
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

type pageFormRepo map[PageFormRepoKey]PageForm

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
	SelectionOpts() []object.SelectOpts
	Key() PageFormRepoKey
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
