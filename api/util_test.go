package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse(t *testing.T) {
	testCases := []struct {
		name         string
		err          error
		expectedBody gin.H
	}{
		{
			name:         "错误响应测试",
			err:          errors.New("测试错误"),
			expectedBody: gin.H{"error": "测试错误"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := errorResponse(tc.err)

			// 验证错误响应结构
			require.Equal(t, tc.err.Error(), response["error"])
		})
	}
}

// Pagination 分页结构体的测试
func TestPagination(t *testing.T) {
	// 创建测试请求
	req, err := http.NewRequest("GET", "/?page=2&page_size=10", nil)
	require.NoError(t, err)

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// 绑定参数
	var pagination Pagination
	err = ctx.ShouldBindQuery(&pagination)
	require.NoError(t, err)

	// 验证分页参数
	require.Equal(t, int32(2), pagination.Page)
	require.Equal(t, int32(10), pagination.PageSize)
}
