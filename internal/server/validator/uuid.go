package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func RegisterUUID(v *validator.Validate) error {
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(uuid.UUID); ok {
			return valuer.String()
		}
		return nil
	}, uuid.UUID{})

	return v.RegisterValidation("not_uuid_nil", func(fl validator.FieldLevel) bool {
		v := fl.Field().Interface().(string)
		if v == uuid.Nil.String() {
			return false
		}
		return true
	})
}
