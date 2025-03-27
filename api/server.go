package api

import (
	"fmt"
	"man/db"
	"man/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker util.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := util.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("无法创建token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	registerValidation()
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := server.router

	// ==学生端==
	v1Router := router.Group("/v1")
	v1Router.POST("/login", server.loginUser)            //用户登录
	v1Router.POST("/refresh", server.refreshAccessToken) //刷新AccessToken

	authV1Router := v1Router.Group("").Use(authMiddleware(server.tokenMaker))
	// 自习室
	authV1Router.GET("/room", server.listRoom)                 //获取自习室列表
	authV1Router.GET("/room/:id", server.getRoom)              //获取自习室详情
	authV1Router.GET("/room/:id/seat", server.listSeat)        //获取自习室座位列表
	authV1Router.POST("/room/:id/reserve", server.reserveRoom) //预约自习室(包括取消预约)
	// 预约
	authV1Router.GET("/reservation/:id", server.getReservation)   //获取预约详情
	authV1Router.POST("/reservation/:id/checkin", server.checkIn) //签到
	// 用户
	authV1Router.GET("/me", server.getMe)                           //获取用户信息
	authV1Router.PUT("/me", server.updateMe)                        //更新用户信息
	authV1Router.GET("/me/reservation", server.listUserReservation) //获取我的预约历史记录

	// ==管理端==
	adminRouter := router.Group("/admin").Use(adminMiddleware(server.tokenMaker))
	// 自习室
	adminRouter.GET("/room", server.listRoom)          //获取自习室列表
	adminRouter.POST("/room", server.createRoom)       //创建自习室
	adminRouter.GET("/room/:id", server.getRoom)       //获取自习室详情
	adminRouter.PUT("/room/:id", server.updateRoom)    //更新自习室
	adminRouter.DELETE("/room/:id", server.deleteRoom) //删除自习室
	// 座位
	adminRouter.GET("/room/:id/seat", server.listSeat)               //获取自习室座位列表
	adminRouter.POST("/room/:id/seat", server.createSeat)            //创建座位
	adminRouter.PUT("/room/:id/seat/:seat_id", server.updateSeat)    //更新座位
	adminRouter.DELETE("/room/:id/seat/:seat_id", server.deleteSeat) //删除座位
	// 预约
	adminRouter.GET("/room/:id/reservation", server.listRoomReservation) //获取自习室预约历史记录
	adminRouter.GET("/reservation", server.listAllReservation)           //获取所有预约历史记录
	adminRouter.GET("/reservation/:id", server.getReservation)           //获取预约详情
	// 用户
	adminRouter.GET("/user", server.listUser)          //获取用户列表
	adminRouter.GET("/user/:id", server.getUser)       //获取用户详情
	adminRouter.PUT("/user/:id", server.updateUser)    //更新用户信息
	adminRouter.DELETE("/user/:id", server.deleteUser) //删除用户

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
