package api

import (
	"fmt"
	"man/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker util.Maker,
	authorizationType string,
	userID int,
	userRole string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(userID, userRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	server := newTestServer(t, nil)

	authPath := "/auth"
	// 用户auth中间件
	server.router.GET(
		authPath,
		authMiddleware(server.tokenMaker),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "正常",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "未授权",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "不支持的授权类型",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupportedType", 1, util.StudentRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "无效格式",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, "", 1, util.StudentRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "过期的token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	server := newTestServer(t, nil)

	adminPath := "/admin"
	// 用户auth中间件
	server.router.GET(
		adminPath,
		adminMiddleware(server.tokenMaker),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "正常",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "非管理员",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, adminPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
