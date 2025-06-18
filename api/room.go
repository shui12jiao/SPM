package api

import (
	"database/sql"
	"man/db"
	"man/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createRoomRequest struct {
	Name       string    `json:"name" binding:"required"`
	Department string    `json:"department" binding:"required"`
	OpenTime   time.Time `json:"open_time" binding:"required"`
	CloseTime  time.Time `json:"close_time" binding:"required"`
}

// createRoom 创建自习室
// @Summary 创建自习室
// @Tags Room
// @Accept json
// @Produce json
// @Param data body createRoomRequest true "自习室参数"
// @Success 200 {object} db.Room
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/room [post]
func (server *Server) createRoom(ctx *gin.Context) {
	var req createRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateRoomParams{
		Name:       req.Name,
		Department: req.Department,
		Code:       util.GenerateSignCode(),
		QrCode:     "暂时不使用二维码",
		OpenTime:   req.OpenTime,
		CloseTime:  req.CloseTime,
	}

	room, err := server.store.CreateRoom(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// 删除自习室

type deleteRoomRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// deleteRoom 删除自习室
// @Summary 删除自习室
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "自习室ID"
// @Success 200 {string} string "OK"
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/room/{id} [delete]
func (server *Server) deleteRoom(ctx *gin.Context) {
	var req deleteRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteRoom(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 返回200
	ctx.JSON(http.StatusOK, nil)
}

// 更新自习室

type updateRoomRequest struct {
	ID         int32      `uri:"id" binding:"required,min=1"`
	Name       *string    `json:"name" binding:"omitempty"`
	Department *string    `json:"department" binding:"omitempty"`
	OpenTime   *time.Time `json:"open_time" binding:"omitempty"`
	CloseTime  *time.Time `json:"close_time" binding:"omitempty"`
	IsActive   *bool      `json:"is_active" binding:"omitempty"`
}

// updateRoom 更新自习室信息
// @Summary 更新自习室
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "自习室ID"
// @Param data body updateRoomRequest true "可更新字段"
// @Success 200 {object} db.Room
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/room/{id} [patch]
func (server *Server) updateRoom(ctx *gin.Context) {
	var req updateRoomRequest
	// 绑定URI参数
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// 绑定JSON参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRoomParams{
		ID:         req.ID,
		Name:       db.ToNull[sql.NullString](req.Name),
		Department: db.ToNull[sql.NullString](req.Department),
		OpenTime:   db.ToNull[sql.NullTime](req.OpenTime),
		CloseTime:  db.ToNull[sql.NullTime](req.CloseTime),
		IsActive:   db.ToNull[sql.NullBool](req.IsActive),
	}

	room, err := server.store.UpdateRoom(ctx, arg)
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

// 查询自习室详情

type getRoomRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// getRoom 获取自习室详情
// @Summary 获取自习室详情
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "自习室ID"
// @Success 200 {object} db.Room
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/room/{id} [get]
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
// @Summary 获取自习室列表（可分页、按学院和状态过滤）
// @Tags Room
// @Accept json
// @Produce json
// @Param page query int true "页码，从1开始"
// @Param page_size query int true "每页数量"
// @Param department query string false "学院"
// @Param is_active query boolean false "是否启用"
// @Success 200 {array} db.Room
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/room [get]
type listRoomRequest struct {
	Pagination
	// 可为空参数
	Department *string `form:"department" binding:"omitempty"`
	IsActive   *bool   `form:"is_active" binding:"omitempty"`
}

func (server *Server) listRoom(ctx *gin.Context) {
	var req listRoomRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRoomParams{
		Limit:      req.PageSize,
		Offset:     (req.Page - 1) * req.PageSize,
		Department: db.ToNull[sql.NullString](req.Department),
		IsActive:   db.ToNull[sql.NullBool](req.IsActive),
	}

	rooms, err := server.store.ListRoom(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}

// getRoomUsage 获取自习室实时使用情况
// @Summary 获取各自习室当前使用情况（预约和签到统计）
// @Tags Room
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500
// @Security BearerAuth
// @Router /admin/room/usage [get]
func (server *Server) getRoomUsage(ctx *gin.Context) {
	usage, err := server.store.GetRoomUsage(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, usage)
}
