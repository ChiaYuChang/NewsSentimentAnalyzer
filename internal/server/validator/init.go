package validator

import (
	"fmt"

	val "github.com/go-playground/validator/v10"
)

var Validate *val.Validate

func init() {
	Validate = val.New()
	if err := RegisterValidator(
		Validate,
		NewDefaultPasswordValidator(),
		NotBeforeNow,
		EnmusRole,
		EnmusJobStatus,
		EnmusApiType,
		EnmusEventType,
	); err != nil {
		panic(fmt.Errorf("error while register validators: %w", err))
	}
}
