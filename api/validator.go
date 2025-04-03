package api

import (
	"man/util"
	"time"

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

	validateStartTime validator.Func = func(fieldLevel validator.FieldLevel) bool {
		if startTime, ok := fieldLevel.Field().Interface().(time.Time); ok {
			if err := util.ValidateStartTime(startTime); err == nil {
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
		v.RegisterValidation("start_time", validateStartTime)
	}
}
