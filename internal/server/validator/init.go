package validator

import (
	"fmt"
	"os"

	val "github.com/go-playground/validator/v10"
)

var Validate *val.Validate

func init() {
	var err error
	Validate, err = NewValiateWithDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while NewValiateWithDefault: %v", err)
		os.Exit(1)
	}
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
