package api

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取配置
func (server *Server) getConfig(ctx *gin.Context) {
	config := server.config.BusinessConfig
	ctx.JSON(http.StatusOK, config)
}

// updateConfigRequest 动态更新配置请求

// 即util.BusinessConfig成员使用指针类型
type updateConfigRequest struct {
	// 预约相关配置
	MaxReservationDuration          *time.Duration `json:"max_reservation_duration" binding:"omitempty"`
	MaxReservationAdvanceDuration   *time.Duration `json:"max_reservation_advance_duration" binding:"omitempty"`
	CancellableReservationDuration  *time.Duration `json:"cancellable_reservation_duration" binding:"omitempty"`
	ReservationRemindBeforeDuration *time.Duration `json:"reservation_remind_before_duration" binding:"omitempty"`
	ReservationRemindAfterDuration  *time.Duration `json:"reservation_remind_after_duration" binding:"omitempty"`
	ReservationViolationDuration    *time.Duration `json:"reservation_violation_duration" binding:"omitempty"`
}

func (server *Server) updateConfig(ctx *gin.Context) {
	var req updateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 更新配置
	reqValue := reflect.ValueOf(req)
	configValue := reflect.ValueOf(&server.config.BusinessConfig).Elem()

	for i := 0; i < reqValue.Type().NumField(); i++ {
		reqField := reqValue.Type().Field(i)
		configField := configValue.FieldByName(reqField.Name)

		// 如果字段存在且是指针类型，检查是否为nil
		reqFieldValue := reqValue.Field(i)
		if reqFieldValue.IsValid() && !reqFieldValue.IsNil() {
			// 更新配置字段的值
			configField.Set(reqFieldValue.Elem())
		}
	}

	ctx.JSON(http.StatusOK, server.config.BusinessConfig)
}
