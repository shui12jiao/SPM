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

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := server.router

	// ==学生端==
	v1Router := router.Group("/v1")
	v1Router.POST("/login", server.login)     //用户登录
	v1Router.POST("/refresh", server.refresh) //刷新AccessToken

	// 需要认证的路由
	authV1Router := v1Router.Group("").Use(authMiddleware(server.tokenMaker))
	// 自习室
	authV1Router.GET("/room", server.listRoom)                 //获取自习室列表
	authV1Router.GET("/room/:id", server.getRoom)              //获取自习室详情
	authV1Router.POST("/room/:id/reserve", server.reserveRoom) //预约自习室(包括取消预约)
	// 预约
	authV1Router.GET("/reservation/:id", server.getReservation)   //获取预约详情
	authV1Router.POST("/reservation/:id/checkin", server.checkIn) //签到
	// 用户
	authV1Router.GET("/me", server.getMe)                           //获取用户信息
	authV1Router.PUT("/me", server.updateMe)                        //更新用户信息
	authV1Router.GET("/me/reservation", server.listUserReservation) //获取我的预约历史记录

	// ==管理端==
	adminRouter := router.Group("/admin").Use(authMiddleware(server.tokenMaker))
	// 自习室
	adminRouter.GET("/room", server.listRoom)                            //获取自习室列表
	adminRouter.GET("/room/:id", server.getRoom)                         //获取自习室详情
	adminRouter.POST("/room", server.createRoom)                         //创建自习室
	adminRouter.PUT("/room/:id", server.updateRoom)                      //更新自习室
	adminRouter.DELETE("/room/:id", server.deleteRoom)                   //删除自习室
	adminRouter.GET("/room/:id/reservation", server.listRoomReservation) //获取自习室预约历史记录
	// 预约
	adminRouter.GET("/reservation", server.listAllReservation) //获取所有预约历史记录
	adminRouter.GET("/reservation/:id", server.getReservation) //获取预约详情

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
