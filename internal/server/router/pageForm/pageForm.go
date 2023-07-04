package pageform

import (
	"errors"
	"net/url"
	"reflect"

	"github.com/go-playground/form"
	val "github.com/go-playground/validator/v10"
)

var ErrUnregisteredPageForm = errors.New("unregistered pageform")

type pageFormRepoKey struct {
	API      string
	Endpoint string
}

type pageFormRepo map[pageFormRepoKey]PageForm

func (pfr pageFormRepo) Add(pf PageForm) {
	pfr[pageFormRepoKey{pf.API(), pf.Endpoint()}] = pf
}

func (pfr pageFormRepo) Set(api, endpoint string, pf PageForm) {
	pfr[pageFormRepoKey{api, endpoint}] = pf
}

func (pfr pageFormRepo) Get(api, endpoint string) (PageForm, error) {
	pf, ok := pfr[pageFormRepoKey{api, endpoint}]
	if !ok {
		return nil, ErrUnregisteredPageForm
	}

	return reflect.New(reflect.TypeOf(pf)).Interface().(PageForm), nil
}

var PageFormRepo = make(pageFormRepo)

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
