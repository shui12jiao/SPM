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
