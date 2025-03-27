package api

import (
	"database/sql"
	"man/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 查询自习室详情

type getRoomRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRoom(ctx *gin.Context) {
	var req getRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	room, err := server.store.GetRoom(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// listRoom 获取自习室列表

type listRoomRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
	// 可为空参数
	Department *string `form:"department,omitempty"`
	IsActive   *bool   `form:"is_active,omitempty"`
}

func (server *Server) listRoom(ctx *gin.Context) {
	var req listRoomRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRoomParams{
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
		Department: db.ToNullString(req.Department),
		IsActive:   db.ToNullBool(req.IsActive),
	}

	rooms, err := server.store.ListRoom(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}
