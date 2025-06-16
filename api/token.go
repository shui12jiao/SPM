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

// refreshAccessToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用refresh token刷新access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body refreshAccessTokenRequest true "Refresh Token Request"
// @Success 200 {object} refreshAccessTokenResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /v1/token/refresh [post]
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
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.UserID, refreshPayload.UserRole, server.config.AccessTokenDuration)
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
