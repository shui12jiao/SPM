package api

import (
	"database/sql"
	"man/db"
	"man/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// db.User => api.userResponse 删去了Password字段
type userResponse struct {
	ID         int32     `json:"id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
}

func newUserResponseFromUser(user db.User) userResponse {
	return userResponse{
		ID:         user.ID,
		Username:   user.Username,
		Role:       user.Role,
		Department: user.Department,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt,
	}
}

// login登录请求
// POST /v1/login
type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,password"`
}

type loginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 通过用户名获取用户
	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 检查密码是否正确
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(int(user.ID), user.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refresh_token, refreshPayload, err := server.tokenMaker.CreateToken(int(user.ID), user.Role, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponseFromUser(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

// getMe获取当前用户信息
// GET /v1/me

func (server *Server) getMe(ctx *gin.Context) {
	id := getUserID(ctx)

	// 通过ID获取用户
	userInfo, err := server.store.GetUser(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponseFromUser(userInfo))
}

// updateMe更新当前用户信息
// 允许更新的字段：password，email
// PATCH /v1/me
type updateMeRequest struct {
	Password string `json:"password" binding:"omitempty,password"`
	Email    string `json:"email" binding:"omitempty,email"`
}

func (server *Server) updateMe(ctx *gin.Context) {
	var req updateMeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := getUserID(ctx)

	arg := db.UpdateUserParams{
		ID:       id,
		Password: db.ToNull[sql.NullString](req.Password),
		Email:    db.ToNull[sql.NullString](req.Email),
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponseFromUser(user))
}

// listUser列出用户

type listUserRequest struct {
	Pagination
	// 可选参数
	Role       *string `form:"role" binding:"omitempty"`
	Department *string `form:"department" binding:"omitempty"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserParams{
		Limit:      req.PageSize,
		Offset:     (req.Page - 1) * req.PageSize,
		Role:       db.ToNull[sql.NullString](req.Role),
		Department: db.ToNull[sql.NullString](req.Department),
	}

	users, err := server.store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var userList []userResponse
	for _, user := range users {
		userList = append(userList, newUserResponseFromUser(user))
	}

	// 返回用户列表
	ctx.JSON(http.StatusOK, userList)
}

// getUser获取用户信息
// GET /admin/user/:id

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponseFromUser(user))
}

// updateUser更新用户信息
// PATCH /admin/user/:id

type updateUserRequest struct {
	ID         int32   `uri:"id" binding:"required,min=1"`
	Username   *string `json:"username" binding:"omitempty"`
	Password   *string `json:"password" binding:"omitempty,password"`
	Role       *string `json:"role" binding:"omitempty"`
	Department *string `json:"department" binding:"omitempty"`
	Email      *string `json:"email" binding:"omitempty,email"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:         req.ID,
		Username:   db.ToNull[sql.NullString](req.Username),
		Password:   db.ToNull[sql.NullString](req.Password),
		Role:       db.ToNull[sql.NullString](req.Role),
		Department: db.ToNull[sql.NullString](req.Department),
		Email:      db.ToNull[sql.NullString](req.Email),
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponseFromUser(user))
}

// deleteUser删除用户信息
// DELETE /admin/user/:id

type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, req.ID)
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
