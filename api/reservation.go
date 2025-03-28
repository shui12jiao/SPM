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
