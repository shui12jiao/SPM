package util

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// 系统环境
	EnvironmentProduction  = "production"  // 生产环境
	EnvironmentDevelopment = "development" // 开发环境
	EnvironmentTest        = "test"        // 测试环境 不手动指定 用于单元测试
)

func InitLogger(environment string) {
	// 设置日志格式
	zerolog.TimeFieldFormat = time.RFC3339
	// 设置日志级别 environment
	switch environment {
	case EnvironmentProduction:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case EnvironmentDevelopment:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		log.Fatal().Msgf("未知环境: %s", environment)
	}
	// 设置日志输出
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
