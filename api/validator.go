package api

import (
	"man/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

var (
	validPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
		if password, ok := fieldLevel.Field().Interface().(string); ok {
			if err := util.ValidatePassword(password); err == nil {
				return true
			}
		}
		return false
	}

	validEmail validator.Func = func(fieldLevel validator.FieldLevel) bool {
		if email, ok := fieldLevel.Field().Interface().(string); ok {
			if err := util.ValidateEmail(email); err == nil {
				return true
			}
		}
		return false
	}
)

func registerValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validPassword)
		v.RegisterValidation("email", validEmail)
	}
}
