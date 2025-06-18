package api

import (
	"database/sql"
	"errors"
	"fmt"
	"man/db"
	"man/task"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 查询预约详情

type getReservationRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

// getReservation 获取预约详情
// @Summary 获取预约详情
// @Description 根据预约 ID 获取预约信息
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path string true "预约ID"
// @Success 200 {object} db.Reservation
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/reservation/{id} [get]
// @Router /admin/reservation/{id} [get]
func (server *Server) getReservation(ctx *gin.Context) {
	var req getReservationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	seat, err := server.store.GetReservation(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seat)
}

type listReservationRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
	// 可为空参数
	StartTime *time.Time `form:"start_time" binding:"omitempty"`
	EndTime   *time.Time `form:"end_time" binding:"omitempty"`
	UserID    *int32     `form:"user_id" binding:"omitempty,min=1"`
	SeatID    *int32     `form:"seat_id" binding:"omitempty,min=1"`
	Status    *string    `form:"status" binding:"omitempty"`
}

// listReservation 获取预约列表
// @Summary 条件获取预约列表
// @Description 支持按时间、用户、座位、状态过滤，并分页返回预约信息
// @Tags Reservation
// @Accept json
// @Produce json
// @Param start_time query string false "起始时间 (ISO8601)"
// @Param end_time query string false "结束时间 (ISO8601)"
// @Param user_id query int false "用户ID"
// @Param seat_id query int false "座位ID"
// @Param status query string false "预约状态"
// @Param page query int true "页码，从1开始"
// @Param page_size query int true "每页数量 (5–50)"
// @Success 200 {array} db.Reservation
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/reservation [get]
func (server *Server) listReservation(ctx *gin.Context) {
	var req listReservationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListReservationParams{
		Limit:     req.PageSize,
		Offset:    (req.Page - 1) * req.PageSize,
		StartTime: db.ToNull[sql.NullTime](req.StartTime),
		EndTime:   db.ToNull[sql.NullTime](req.EndTime),
		UserID:    db.ToNull[sql.NullInt32](req.UserID),
		SeatID:    db.ToNull[sql.NullInt32](req.SeatID),
		Status:    db.ToNull[sql.NullString](req.Status),
	}

	seats, err := server.store.ListReservation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seats)
}

type listMyReservationRequest = Pagination

// listMyReservation 获取当前用户的预约列表
// @Summary 获取当前用户的预约记录（分页）
// @Tags Reservation
// @Accept json
// @Produce json
// @Param page query int true "页码，从1开始"
// @Param page_size query int true "每页数量 (5–50)"
// @Success 200 {array} db.Reservation
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /v1/me/reservation [get]
func (server *Server) listMyReservation(ctx *gin.Context) {
	var req listMyReservationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListReservationParams{
		UserID: db.ToNull[sql.NullInt32](getUserID(ctx)),
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	reservations, err := server.store.ListReservation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

type createReservationRequest struct {
	SeatID    int32     `json:"seat_id" binding:"required,min=1"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

// createReservation 创建座位预约
// @Summary 创建预约
// @Tags Reservation
// @Accept json
// @Produce json
// @Param data body createReservationRequest true "预约参数"
// @Success 200 {object} db.Reservation
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/reservation [post]
func (server *Server) createReservation(ctx *gin.Context) {
	var req createReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 将时间转换为 UTC
	req.StartTime = req.StartTime.UTC()
	req.EndTime = req.EndTime.UTC()
	now := time.Now().UTC()

	// 检查开始时间
	if req.StartTime.Before(now.Add(server.config.MinReservationAdvanceDuration)) || req.StartTime.After(now.Add(server.config.MaxReservationAdvanceDuration)) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("预约开始时间错误: %s, 预约提前时间范围: %s - %s",
			req.StartTime.Format(time.RFC3339),
			server.config.MinReservationAdvanceDuration.String(),
			server.config.MaxReservationAdvanceDuration.String())))
		return
	}

	// 检查预约时长
	reservationDur := req.EndTime.Sub(req.StartTime)
	if reservationDur < server.config.MinReservationDuration || reservationDur > server.config.MaxReservationDuration {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("预约时长错误: %s, 预约时长范围: %s - %s",
			reservationDur.String(),
			server.config.MinReservationDuration.String(),
			server.config.MaxReservationDuration.String())))
		return
	}

	// 检查座位是否存在且可用
	// 实际数据库会更复杂的判断，包括时间冲突等，这里暂时保留，方便debug
	seat, err := server.store.GetSeat(ctx, req.SeatID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("座位不存在")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !seat.IsAvailable {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("座位不可用")))
		return
	}

	arg := db.CreateReservationParams{
		SeatID:    req.SeatID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID:    getUserID(ctx),
	}

	reservation, err := server.store.CreateReservation(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "预约失败：座位不可用或时间冲突",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	// 提交预约创建任务
	err = server.scheduler.CreateReservationTasks(task.CreateReservationTaskArg{
		ReservationID:        reservation.ID,
		UserID:               reservation.UserID,
		ReservationTime:      reservation.StartTime,
		RemindBeforeDuration: server.config.ReservationRemindBeforeDuration,
		RemindAfterDuration:  server.config.ReservationRemindAfterDuration,
		ViolationDuration:    server.config.ReservationViolationDuration,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

type deleteReservationRequest struct {
	ID string `uri:"id" binding:"required,uuid"` // 预约ID，使用UUID格式
}

// deleteReservation 取消预约
// @Summary 取消指定预约
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path string true "预约ID"
// @Success 200 {string} string "OK"
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/reservation/{id} [delete]
func (server *Server) deleteReservation(ctx *gin.Context) {
	var req deleteReservationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 转换为 UUID
	id := uuid.MustParse(req.ID)
	// 获取预约开始时间
	reservation, err := server.store.GetReservation(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 检查是否可以取消预约
	// query中已有检测逻辑，冗余
	if time.Now().UTC().Add(server.config.CancellableReservationDuration).After(reservation.StartTime) {
		// 如果当前时间已经超过预约开始前的可取消时间，则无法取消
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("预约已开始，无法取消")))
		return
	}

	result, err := server.store.DeleteReservation(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("预约不存在或已取消")))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type checkInRequest struct {
	ID   string `uri:"id" binding:"required,uuid"` // 预约ID，使用UUID格式
	Code string `json:"code" binding:"required"`   // 签到码
}

// checkIn 签到
// @Summary 签到已预约座位
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path string true "预约ID"
// @Param code body string true "签到码"
// @Success 200 {object} db.Reservation
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/reservation/{id}/checkin [post]
func (server *Server) checkIn(ctx *gin.Context) {
	var req checkInRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 查找预约
	rs, err := server.store.GetReservationWithRoomCode(ctx, uuid.MustParse(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("预约不存在")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rs.RoomCode != req.Code { // 检查签到码
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("签到码错误")))
		return
	}
	if rs.UserID != getUserID(ctx) { //检查是否为当前用户的预约
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("无权操作他人预约")))
		return
	}
	if rs.StartTime.After(time.Now().UTC()) { // 检查预约是否已开始
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("预约尚未开始，无法签到")))
		return
	}
	if rs.Status != db.ReservationStatusReserved { // 检查预约状态
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("只能对已预约状态的座位进行签到操作")))
		return
	}

	arg := db.UpdateReservationStatusParams{
		ID:     uuid.MustParse(req.ID),
		Status: db.ReservationStatusCompleted,
	}

	upRs, err := server.store.UpdateReservationStatus(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, upRs)
}

// 格式22:50, 只取时分部分
func timeExtract(t time.Time) time.Duration {
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
}
