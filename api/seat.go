package api

import (
	"database/sql"
	"man/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建座位，单次创建一个座位，用于特殊情况添加
// 一般使用批量创建座位的接口

type createSeatRequest struct {
	RoomID      int32  `json:"room_id" binding:"required,min=1"`
	Number      string `json:"seat_num" binding:"required,min=1"`
	HasSocket   bool   `json:"has_socket" binding:"required"`
	IsAvailable bool   `json:"is_available"`
}

// createSeat 创建单个座位
// @Summary 创建单个座位
// @Tags Seat
// @Accept json
// @Produce json
// @Param data body createSeatRequest true "座位参数"
// @Success 200 {object} db.Seat
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat [post]
func (server *Server) createSeat(ctx *gin.Context) {
	var req createSeatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateSeatParams{
		RoomID:      req.RoomID,
		Number:      req.Number,
		HasSocket:   req.HasSocket,
		IsAvailable: req.IsAvailable,
	}
	seat, err := server.store.CreateSeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seat)
}

type getSeatRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// getSeat 获取座位详情
// @Summary 获取座位详情
// @Tags Seat
// @Accept json
// @Produce json
// @Param id path int true "座位ID"
// @Success 200 {object} db.Seat
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat/{id} [get]
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

type listSeatRequest struct {
	Pagination
	// 可为空参数
	RoomID      *int32 `form:"room_id" binding:"omitempty,min=1"`
	HasSocket   *bool  `form:"has_socket" binding:"omitempty"`
	IsAvailable *bool  `form:"is_available" binding:"omitempty"`
}

// listSeat 获取座位列表
// @Summary 获取座位列表（支持分页和条件过滤）
// @Tags Seat
// @Accept json
// @Produce json
// @Param room_id query int false "自习室ID"
// @Param has_socket query boolean false "是否有插座"
// @Param is_available query boolean false "是否可用"
// @Param page query int true "页码，从1开始"
// @Param page_size query int true "每页数量"
// @Success 200 {array} db.Seat
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat [get]
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

// 更新座位信息，不允许修改room_id，其他参数可选

type updateSeatRequest struct {
	ID          int32   `uri:"id" binding:"required,min=1"`
	Number      *string `json:"seat_num" binding:"omitempty,min=1"`
	HasSocket   *bool   `json:"has_socket" binding:"omitempty"`
	IsAvailable *bool   `json:"is_available"`
}

// updateSeat 更新座位信息
// @Summary 更新座位信息（不可更改RoomID）
// @Tags Seat
// @Accept json
// @Produce json
// @Param id path int true "座位ID"
// @Param data body updateSeatRequest true "可修改字段"
// @Success 200 {object} db.Seat
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat/{id} [patch]
func (server *Server) updateSeat(ctx *gin.Context) {
	var req updateSeatRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateSeatParams{
		ID:          req.ID,
		Number:      db.ToNull[sql.NullString](req.Number),
		HasSocket:   db.ToNull[sql.NullBool](req.HasSocket),
		IsAvailable: db.ToNull[sql.NullBool](req.IsAvailable),
	}
	seat, err := server.store.UpdateSeat(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seat)
}

type deleteSeatRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// deleteSeat 删除座位
// @Summary 删除座位
// @Tags Seat
// @Accept json
// @Produce json
// @Param id path int true "座位ID"
// @Success 200 {string} string "OK"
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat/{id} [delete]
func (server *Server) deleteSeat(ctx *gin.Context) {
	var req deleteSeatRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteSeat(ctx, req.ID)
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

type createSeatsRequest struct {
	RoomID int32 `json:"room_id" binding:"required,min=1"`
	Seats  []struct {
		Number      string `json:"seat_num" binding:"required,min=1"`
		HasSocket   bool   `json:"has_socket" binding:"required"`
		IsAvailable bool   `json:"is_available" binding:"required"`
	} `json:"seats" binding:"required,min=1"`
}

// createSeats 批量创建座位
// @Summary 批量创建座位（用于初始化）
// @Tags Seat
// @Accept json
// @Produce json
// @Param data body createSeatsRequest true "座位数组及所属自习室"
// @Success 200 {array} db.Seat
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat/batch [post]
func (server *Server) createSeats(ctx *gin.Context) {
	var req createSeatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 将Seats中的提取出来
	var numbers []string
	var has_sockets, is_availables []bool
	for _, seat := range req.Seats {
		numbers = append(numbers, seat.Number)
		has_sockets = append(has_sockets, seat.HasSocket)
		is_availables = append(is_availables, seat.IsAvailable)
	}

	arg := db.CreateSeatsParams{
		RoomID:       req.RoomID,
		Numbers:      numbers,
		HasSockets:   has_sockets,
		IsAvailables: is_availables,
	}
	seats, err := server.store.CreateSeats(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seats)
}

// 批量更新座位信息
// 统一将更改座位的插口和可用状态

type updateSeatsRequest struct {
	Seats []struct {
		ID          int32  `json:"ids" binding:"required,min=1"`
		Number      string `json:"seat_num" binding:"required,min=1"`
		HasSocket   bool   `json:"has_socket" binding:"required"`
		IsAvailable bool   `json:"is_available" binding:"required"`
	} `json:"seats" binding:"required,min=1"`
}

// updateSeats 批量更新座位
// @Summary 批量更新座位信息（如状态/插座）
// @Tags Seat
// @Accept json
// @Produce json
// @Param data body updateSeatsRequest true "座位批量更新参数"
// @Success 200 {array} db.Seat
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /admin/seat/batch [patch]
func (server *Server) updateSeats(ctx *gin.Context) {
	var req updateSeatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var ids []int32
	var numbers []string
	var has_sockets, is_availables []bool
	for _, seat := range req.Seats {
		ids = append(ids, seat.ID)
		numbers = append(numbers, seat.Number)
		has_sockets = append(has_sockets, seat.HasSocket)
		is_availables = append(is_availables, seat.IsAvailable)
	}

	arg := db.UpdateSeatsParams{
		Ids:          ids,
		Numbers:      numbers,
		HasSockets:   has_sockets,
		IsAvailables: is_availables,
	}

	seats, err := server.store.UpdateSeats(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, seats)
}
