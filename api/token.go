package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) refreshAccessToken(ctx *gin.Context) {
	var req refreshAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 验证refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 判断refresh token是否有效
	if err := refreshPayload.Valid(); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 创建新的access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := refreshAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, res)
}
