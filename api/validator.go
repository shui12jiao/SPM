package api

import (
	"man/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// binding 用于注册自定义验证器，只用来检查数据的合法性
// 不包含复杂的业务逻辑验证

var (
	validPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
		if password, ok := fieldLevel.Field().Interface().(string); ok {
			if err := util.ValidatePassword(password); err == nil {
				return true
			}
		}
		return false
	}
)

func registerValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validPassword)
		// v.RegisterValidation("email", validEmail) // binding现在已经支持了email的验证
	}
}
