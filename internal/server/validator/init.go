package validator

import (
	"fmt"
	"os"
	"sync"

	val "github.com/go-playground/validator/v10"
)

var Validate *val.Validate

var sv validateSingleton

type validateSingleton struct {
	Validate *val.Validate
	Error    error
	sync.Once
}

func init() {
	var err error
	Validate, err = NewValiateWithDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while NewValiateWithDefault: %v", err)
		os.Exit(1)
	}
}

func GetDefaultValidate() (*val.Validate, error) {
	sv.Do(func() {
		sv.Validate = val.New()

		if err := RegisterUUID(sv.Validate); err != nil {
			sv.Error = fmt.Errorf("error while RegisterUUID: %w", err)
			return
		}

		if err := RegisterValidator(
			sv.Validate,
			NewDefaultPasswordValidator(),
			NotBeforeNow,
			EnmusRole,
			EnmusJobStatus,
			EnmusApiType,
			EnmusEventType,
		); err != nil {
			sv.Error = fmt.Errorf("error while register validators: %v", err)
			return
		}
	})
	return sv.Validate, sv.Error
}

func NewValiateWithDefault() (*val.Validate, error) {
	v := val.New()

	if err := RegisterUUID(v); err != nil {
		return nil, fmt.Errorf("error while RegisterUUID: %w", err)
	}

	if err := RegisterValidator(
		v,
		NewDefaultPasswordValidator(),
		NotBeforeNow,
		EnmusRole,
		EnmusJobStatus,
		EnmusApiType,
		EnmusEventType,
	); err != nil {
		return nil, fmt.Errorf("error while register validators: %v", err)
	}
	return v, nil
}
