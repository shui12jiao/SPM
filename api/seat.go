package api

import (
	"database/sql"
	"man/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 查询座位详情

type getSeatRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSeat(ctx *gin.Context) {
	var req getSeatRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	seat, err := server.store.GetSeat(ctx, req.ID)
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

// listSeat 获取座位列表
// GET /seats?room_id=1&has_socket=true&is_available=true&page=1&page_size=10
type listSeatRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
	// 可为空参数
	RoomID      *int32 `form:"room_id" binding:"omitempty,min=1"`
	HasSocket   *bool  `form:"has_socket" binding:"omitempty"`
	IsAvailable *bool  `form:"is_available" binding:"omitempty"`
}

func (server *Server) listSeat(ctx *gin.Context) {
	var req listSeatRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListSeatParams{
		Limit:       req.PageSize,
		Offset:      (req.Page - 1) * req.PageSize,
		RoomID:      db.ToNull[sql.NullInt32](req.RoomID),
		HasSocket:   db.ToNull[sql.NullBool](req.HasSocket),
		IsAvailable: db.ToNull[sql.NullBool](req.IsAvailable),
	}

	seats, err := server.store.ListSeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seats)
}
