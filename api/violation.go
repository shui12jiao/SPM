package api

import (
	"database/sql"
	"man/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// listMyViolation获取我的违规记录
// GET /me/violation

type listMyViolationRequest = Pagination

func (server *Server) listMyViolation(ctx *gin.Context) {
	var req listMyViolationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 获取用户的违规记录
	arg := db.ListViolationParams{
		UserID: db.ToNull[sql.NullInt32](getUserID(ctx)),
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}
	violations, err := server.store.ListViolation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, violations)
}

// updateViolation更新违约记录，填写理由
// PATCH /me/violation/:id

type updateViolationRequest struct {
	// ViolationID 违规记录ID
	ViolationID int32 `uri:"id" binding:"required,min=1"`

	// Reason 违规理由
	Reason string `json:"reason" binding:"required,min=1,max=100"`
}

func (server *Server) updateViolation(ctx *gin.Context) {
	var req updateViolationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 更新违规记录
	arg := db.UpdateViolationParams{
		ID:     req.ViolationID,
		Reason: req.Reason,
	}
	violation, err := server.store.UpdateViolation(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, violation)
}

// listViolation获取违规记录列表
// GET /violation

type listViolationRequest struct {
	Pagination
	// 可选参数
	ReservationID *uuid.UUID `form:"reservation_id"` // 预约ID
	UserID        *int32     `form:"user_id"`        // 用户ID
}

func (server *Server) listViolation(ctx *gin.Context) {
	var req listViolationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 获取违规记录列表
	arg := db.ListViolationParams{
		Limit:         req.PageSize,
		Offset:        (req.Page - 1) * req.PageSize,
		ReservationID: db.ToNull[uuid.NullUUID](req.ReservationID),
		UserID:        db.ToNull[sql.NullInt32](req.UserID),
	}
	violations, err := server.store.ListViolation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, violations)
}

// getViolation获取违规记录详情
// GET /violation/:id

type getViolationRequest struct {
	ViolationID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getViolation(ctx *gin.Context) {
	var req getViolationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 获取违规记录详情
	violation, err := server.store.GetViolation(ctx, req.ViolationID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, violation)
}
