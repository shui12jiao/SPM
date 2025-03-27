package main

import (
	"man/util"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 加载配置文件
	config := util.LoadConfig()

	// 设置日志输出
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// 初始化数据库, 执行迁移
	db := util.InitDB(config.DBDriver, config.DBSource, config.MigrationURL)

	// 创建存储
	store := db.NewStore(db)
}
