package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"man/db"
	"man/db/mockdb"
	"man/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomRoom() db.Room {
	now := time.Now()
	return db.Room{
		ID:         1,
		Name:       util.RandomString(6),
		Department: "IT",
		OpenTime:   now,
		CloseTime:  now.Add(8 * time.Hour),
		Code:       util.RandomString(6),
		QrCode:     util.RandomString(20),
		IsActive:   true,
	}
}

func TestCreateRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name  string
		body  gin.H
		stub  func(store *mockdb.MockStore)
		check func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "正常创建自习室",
			body: gin.H{
				"name":       room.Name,
				"department": room.Department,
				"open_time":  room.OpenTime.Format(time.RFC3339),
				"close_time": room.CloseTime.Format(time.RFC3339),
			},
			stub: func(store *mockdb.MockStore) {
				arg := db.CreateRoomParams{
					Name:       room.Name,
					Department: room.Department,
					OpenTime:   room.OpenTime,
					CloseTime:  room.CloseTime,
				}
				store.EXPECT().
					CreateRoom(mock.Anything, mock.MatchedBy(func(p db.CreateRoomParams) bool {
						// Compare fields, ignoring time zone differences for the check
						return p.Name == arg.Name && p.Department == arg.Department &&
							p.OpenTime.Unix() == arg.OpenTime.Unix() &&
							p.CloseTime.Unix() == arg.CloseTime.Unix()
					})).
					Return(room, nil)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoom(t, recorder.Body, room)
			},
		},
		{
			name: "内部服务器错误",
			body: gin.H{
				"name":       room.Name,
				"department": room.Department,
				"open_time":  room.OpenTime.Format(time.RFC3339),
				"close_time": room.CloseTime.Format(time.RFC3339),
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoom(mock.Anything, mock.Anything).
					Return(db.Room{}, sql.ErrConnDone)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "无效请求",
			body: gin.H{
				// "name" is required
				"department": "IT",
			},
			stub: func(store *mockdb.MockStore) {
				// No DB call should be made
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/admin/room"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// 添加管理员授权
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.check(t, recorder)
		})
	}
}

func TestGetRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name          string
		roomID        int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		stub          func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "正常获取自习室",
			roomID: room.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(mock.Anything, room.ID).
					Return(room, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoom(t, recorder.Body, room)
			},
		},
		{
			name:   "未找到",
			roomID: room.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(mock.Anything, room.ID).
					Return(db.Room{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "内部错误",
			roomID: room.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(mock.Anything, room.ID).
					Return(db.Room{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "无效ID",
			roomID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				// No DB call should be made for invalid ID
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/v1/room/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// 测试更新自习室
func TestUpdateRoom(t *testing.T) {
	oldRoom := randomRoom()
	newRoom := randomRoom()
	newRoom.ID = oldRoom.ID

	testCases := []struct {
		name          string
		roomID        int32
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		stub          func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "正常更新自习室",
			roomID: oldRoom.ID,
			body: gin.H{
				"name":        newRoom.Name,
				"department": newRoom.Department,
				
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(mock.Anything, mock.Anything).
					Times(1).
					Return(newRoom, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRoom(t, recorder.Body, newRoom)
			},
		},
		{
			name:   "内部错误",
			roomID: oldRoom.ID,
			body: gin.H{
				"name":        newRoom.Name,
				"department": newRoom.Department,
				
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(mock.Anything, mock.Anything).
					Times(1).
					Return(db.Room{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/admin/room/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// 测试删除自习室
func TestDeleteRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name          string
		roomID        int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		stub          func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "正常删除自习室",
			roomID: room.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(mock.Anything, room.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "内部错误",
			roomID: room.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(mock.Anything, room.ID).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/admin/room/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// 测试列出自习室
func TestListRoom(t *testing.T) {
	n := 5
	rooms := make([]db.Room, n)
	for i := 0; i < n; i++ {
		rooms[i] = randomRoom()
	}

	testCases := []struct {
		name          string
		page          int32
		pageSize      int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		stub          func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "正常获取自习室列表",
			page:     1,
			pageSize: int32(n),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoom(mock.Anything, mock.Anything).
					Times(1).
					Return(rooms, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRooms(t, recorder.Body, rooms)
			},
		},
		{
			name:     "内部错误",
			page:     1,
			pageSize: int32(n),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.StudentRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoom(mock.Anything, mock.Anything).
					Times(1).
					Return([]db.Room{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/v1/room?page=%d&page_size=%d", tc.page, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// 测试获取自习室使用情况
func TestGetRoomUsage(t *testing.T) {
	mockUsage := db.GetRoomUsageRow{
		ID:             1,
		Name:           "自习室1",
		AvailableSeats: 10,
		OccupiedSeats:  5,
		SocketSeats:    8,
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker util.Maker)
		stub          func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "正常获取自习室使用情况",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoomUsage(mock.Anything).
					Times(1).
					Return([]db.GetRoomUsageRow{mockUsage}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "内部错误",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker util.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, util.AdminRole, time.Minute)
			},
			stub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoomUsage(mock.Anything).
					Times(1).
					Return([]db.GetRoomUsageRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.stub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/admin/room/usage"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchRoom(t *testing.T, body *bytes.Buffer, room db.Room) {
	data, err := json.Marshal(room)
	require.NoError(t, err)

	responseData := body.Bytes()
	// 比较JSON，忽略空格和类似的空白字符差异
	require.JSONEq(t, string(data), string(responseData))
}

func requireBodyMatchRooms(t *testing.T, body *bytes.Buffer, rooms []db.Room) {
	data, err := json.Marshal(rooms)
	require.NoError(t, err)

	responseData := body.Bytes()
	require.JSONEq(t, string(data), string(responseData))
}
