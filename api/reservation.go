package api

import (
	"database/sql"
	"man/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 查询预约详情

type getReservationRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

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

// listReservation 获取预约列表
// GET /reservation?start_time=2025-03-01T00:00:00Z&end_time=2025-03-31T23:59:59Z&user_id=123&seat_id=456&status=confirmed&page=1&page_size=10
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

// 获取我的预约列表
// GET /me/reservation?page=1&page_size=10
type listMyReservationRequest = Pagination

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

// 创建预约
// POST /reservation

type createReservationRequest struct {
	SeatID    int32     `json:"seat_id" binding:"required,min=1"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

func (server *Server) createReservation(ctx *gin.Context) {
	var req createReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// 取消预约
// DELETE /reservation/:id
type deleteReservationRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteReservation(ctx *gin.Context) {
	var req deleteReservationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteReservation(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// TODO 二维码签到
// 签到
// POST /reservation/:id/checkin
type checkInRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

func (server *Server) checkIn(ctx *gin.Context) {
	var req checkInRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateReservationStatusParams{
		ID:     req.ID,
		Status: db.ReservationStatusCompleted,
	}

	reservation, err := server.store.UpdateReservationStatus(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}
