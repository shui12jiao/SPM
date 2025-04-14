package main

import (
	"man/api"
	"man/db"
	"man/task"
	"man/util"

	"github.com/rs/zerolog/log"
)

// 添加注释以描述 server 信息
// @title           自习室管理系统 API
// @version         1.0
// @description     该系统用于管理自习室的预约、签到、退座等功能。
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// 加载配置文件
	config := util.LoadConfig()

	// 设置日志输出
	util.InitLogger(config.Environment)

	// 初始化数据库, 执行迁移
	conn := util.InitDB(config.DBDriver, config.DBSource, config.MigrationURL)

	// 创建存储
	store := db.NewStore(conn)

	// 初始化任务调度器
	scheduler := task.NewScheduler(config, store)
	scheduler.Start() // 非阻塞启动调度器

	// HTTP服务器
	runHTTPServer(config, store, scheduler)
}

func runHTTPServer(config util.Config, store db.Store, scheduler task.Scheduler) {
	server, err := api.NewServer(config, store, scheduler)
	if err != nil {
		log.Fatal().Err(err).Msg("无法创建服务器")
		return
	}
	server.Start(config.HTTPServerAddress)
}
