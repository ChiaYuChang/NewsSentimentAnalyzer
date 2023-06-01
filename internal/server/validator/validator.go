package validator

import (
	"time"

	val "github.com/go-playground/validator/v10"
)

type Validator interface {
	Tag() string
	ValFun() val.Func
}

func RegisterValidator(val *val.Validate, validators ...Validator) error {
	var err error
	for _, v := range validators {
		err = val.RegisterValidation(v.Tag(), v.ValFun())
		if err != nil {
			return err
		}
	}
	return nil
}

type TimeValidator struct {
	tag     string
	valFunc val.Func
}

func (tv TimeValidator) Tag() string {
	return tv.tag
}

func (tv TimeValidator) ValFun() val.Func {
	return tv.valFunc
}

func NotBeforeTime(tag string, when func() time.Time) TimeValidator {
	return TimeValidator{
		tag: tag,
		valFunc: func(fl val.FieldLevel) bool {
			t, ok := fl.Field().Interface().(time.Time)
			if !ok {
				return false
			}

			if t.After(when()) {
				return false
			}
			return true
		},
	}
}

var NotBeforeNow = NotBeforeTime("before_now", time.Now)
