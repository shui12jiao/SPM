package api

import (
	"errors"
	"fmt"
	"man/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

const (
	authorizationHeaderKey  = "authorization"         // 用于获取token的header key
	authorizationTypeBearer = "bearer"                // token类型, bearer一般用于JWT
	authorizationPayloadKey = "authorization_payload" // 用于存储解析后的token payload的key

	requestIDHeaderKey = "X-Request-ID" // 用于存储请求ID的header key
)

// requestIDMiddleware 是一个中间件，用于生成唯一的请求ID并将其添加到请求头中
// 这里采用xid库
func requestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查请求头中是否已经存在请求ID
		requestID := ctx.GetHeader(requestIDHeaderKey)
		if requestID == "" {
			// 如果不存在，则生成一个新的请求ID
			requestID = xid.New().String()
			// 写入到响应头中！（注意是响应头）
			ctx.Header(requestIDHeaderKey, requestID)
		}
		// 将请求ID存储到上下文中
		ctx.Set(requestIDHeaderKey, requestID)
		ctx.Next()
	}
}

// 用于将gin的日志输出到zerolog中
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 将带有requestID的logger传递给gin的上下文
		logger := log.With().Str("request_id", c.GetString(requestIDHeaderKey)).Logger()
		c.Set("logger", &logger)

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		ginLogger := logger.With().
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("client_ip", clientIP).
			Dur("latency", latency).
			Logger()

		if statusCode >= 500 {
			ginLogger.Error().Msg(errorMessage)
		} else if statusCode >= 400 {
			ginLogger.Warn().Msg(errorMessage)
		} else {
			ginLogger.Info().Msg("请求成功")
		}
	}
}

// 检查用户是否登录的中间件
// 该中间件会检查请求头中是否包含Authorization字段，并验证其有效性
func authMiddleware(tokenMaker util.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("必须提供authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("无效的authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("未支持的authorization类型 %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// adminMiddleware 是一个中间件，用于验证用户是否是管理员
func adminMiddleware(tokenMaker util.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("必须提供authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("无效的authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("未支持的authorization类型 %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 检查用户角色是否是管理员
		if payload.UserRole != util.StudentRole {
			err := errors.New("用户不是管理员")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
