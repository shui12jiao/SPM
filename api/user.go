package api

import (
	"database/sql"
	"man/db"
	"man/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

// loginUser 登录
// @Summary 用户登录
// @Description 使用用户名和密码登录，获取 access 和 refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body loginUserRequest true "登录请求"
// @Success 200 {object} loginUserResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/login [post]
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

// getMe 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取登录用户的基本信息
// @Tags User
// @Produce json
// @Success 200 {object} userResponse
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/me [get]
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

// updateMe 更新当前用户信息
// @Summary 更新当前用户信息
// @Description 仅支持修改密码和邮箱
// @Tags User
// @Accept json
// @Produce json
// @Param data body updateMeRequest true "更新内容"
// @Success 200 {object} userResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /v1/me [patch]
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

// listUser 获取用户列表
// @Summary 获取用户列表
// @Description 支持分页和按角色/部门过滤
// @Tags User
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param role query string false "角色过滤"
// @Param department query string false "部门过滤"
// @Success 200 {array} userResponse
// @Failure 400
// @Failure 500
// @Security BearerAuth
// @Router /v1/user [get]
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

// getUser 获取指定用户信息
// @Summary 获取用户信息
// @Description 管理员通过用户ID获取用户信息
// @Tags User
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} userResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/user/{id} [get]
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

// createUser创建用户信息
// POST /admin/user

type createUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required,password"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}

// createUser 创建用户
// @Summary 创建新用户
// @Description 管理员创建用户账户
// @Tags User
// @Accept json
// @Produce json
// @Param data body createUserRequest true "用户信息"
// @Success 200 {object} userResponse
// @Failure 400
// @Failure 409
// @Failure 500
// @Security BearerAuth
// @Router /admin/user [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 哈希密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:   req.Username,
		Password:   hashedPassword,
		Role:       req.Role,
		Department: req.Department,
		Email:      req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// 如果用户名已存在，返回409错误
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, errorResponse(err))
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

// updateUser 更新用户信息
// @Summary 更新用户
// @Description 管理员更新指定用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param data body updateUserRequest true "更新内容"
// @Success 200 {object} userResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/user/{id} [patch]
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

// deleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除指定用户
// @Tags User
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} nil
// @Failure 400
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /admin/user/{id} [delete]
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
