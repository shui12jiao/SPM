package api

import (
	"fmt"
	"man/db"
	"man/util"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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

	// ==swag==
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==学生端==
	v1Router := router.Group("/v1")
	v1Router.POST("/login", server.loginUser)            // 用户登录
	v1Router.POST("/refresh", server.refreshAccessToken) // 刷新 AccessToken

	authV1Router := v1Router.Group("").Use(authMiddleware(server.tokenMaker))
	// 自习室
	authV1Router.GET("/room", server.listRoom)    // 获取自习室列表
	authV1Router.GET("/room/:id", server.getRoom) // 获取自习室详情
	// 座位
	authV1Router.GET("/seat", server.listSeat)    // 获取座位列表
	authV1Router.GET("/seat/:id", server.getSeat) // 获取座位详情
	// 预约
	authV1Router.POST("/reservation", server.createReservation)       // 创建预约
	authV1Router.GET("/reservation/:id", server.getReservation)       // 获取预约详情
	authV1Router.DELETE("/reservation/:id", server.deleteReservation) // 取消预约
	authV1Router.POST("/reservation/:id/checkin", server.checkIn)     // 签到
	// 用户
	authV1Router.GET("/me", server.getMe)                           // 获取用户信息
	authV1Router.PATCH("/me", server.updateMe)                      // 更新用户信息
	authV1Router.GET("/me/reservation", server.listMyReservation)   // 获我的预约列表
	authV1Router.GET("/me/violation", server.listMyViolation)       // 获取我的违规记录
	authV1Router.PATCH("/me/violation/:id", server.updateViolation) // 更新违约记录，填写理由

	// ==管理端==
	adminRouter := router.Group("/admin").Use(adminMiddleware(server.tokenMaker))
	// 自习室
	adminRouter.GET("/room", server.listRoom)          // 获取自习室列表
	adminRouter.POST("/room", server.createRoom)       // 创建自习室
	adminRouter.GET("/room/:id", server.getRoom)       // 获取自习室详情
	adminRouter.PATCH("/room/:id", server.updateRoom)  // 更新自习室
	adminRouter.DELETE("/room/:id", server.deleteRoom) // 删除自习室
	// 座位
	adminRouter.GET("/seat", server.listSeat)          // 获取座位列表
	adminRouter.POST("/seat", server.createSeat)       // 创建座位
	adminRouter.PATCH("/seat/:id", server.updateSeat)  // 更新座位
	adminRouter.DELETE("/seat/:id", server.deleteSeat) // 删除座位
	adminRouter.POST("/seats", server.createSeats)     // 批量创建座位
	adminRouter.PATCH("/seats", server.updateSeats)    // 批量更新座位
	// 预约
	adminRouter.GET("/reservation", server.listReservation)    // 条件查询所有预约
	adminRouter.GET("/reservation/:id", server.getReservation) // 获取预约详情
	// 用户
	adminRouter.GET("/user", server.listUser)          // 获取用户列表
	adminRouter.GET("/user/:id", server.getUser)       // 获取用户详情
	adminRouter.PATCH("/user/:id", server.updateUser)  // 更新用户信息
	adminRouter.DELETE("/user/:id", server.deleteUser) // 删除用户
	// 违规记录
	adminRouter.GET("/violation", server.listViolation)    // 获取违规记录
	adminRouter.GET("/violation/:id", server.getViolation) // 获取违规详情
	// 资源查看接口
	adminRouter.GET("/room/usage", server.getRoomUsage) // 获取自习室实时使用情况
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
