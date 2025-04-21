package api

import (
	"bytes"
	"encoding/json"
	"io"
	"man/db"
	"man/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"man/db/mockdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user := db.User{
		Username:   util.RandomString(6),
		Password:   hashedPassword,
		Email:      util.RandomEmail(),
		Role:       util.AdminRole,
		Department: util.RandomString(4),
	}
	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var respUser db.User
	err = json.Unmarshal(data, &respUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, respUser.Username)
	require.Equal(t, user.Email, respUser.Email)
	require.Empty(t, respUser.Password)
}

func TestCreateUser(t *testing.T) {
	url := "/admin/user"
	user, password := randomUser(t)

	testCases := []struct {
		name  string
		body  gin.H
		stub  func(store *mockdb.MockStore)
		check func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "正常创建用户",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"email":      user.Email,
				"role":       user.Role,
				"department": user.Department,
			},
			stub: func(store *mockdb.MockStore) {
				// 注意由于哈希加盐的原因，密码不能直接比较
				store.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(arg db.CreateUserParams) bool {
					err := util.CheckPassword(password, arg.Password)
					require.NoError(t, err)
					require.Equal(t, user.Username, arg.Username)
					require.Equal(t, user.Email, arg.Email)
					require.Equal(t, user.Role, arg.Role)
					require.Equal(t, user.Department, arg.Department)
					return true
				})).Return(user, nil).Once()
			},
			check: func(t *testing.T, resp *httptest.ResponseRecorder) {
				require.Equal(t, 200, resp.Code)
				requireBodyMatchUser(t, resp.Body, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 依赖
			store := mockdb.NewMockStore(t)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// 构造存根
			tc.stub(store)

			// 构造请求
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)

			// 发送请求并检查响应
			server.router.ServeHTTP(recorder, request)
			tc.check(t, recorder)
		})
	}
}
