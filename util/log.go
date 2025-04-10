package util

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(environment string) {
	// 设置日志格式
	zerolog.TimeFieldFormat = time.RFC3339
	// 设置日志级别 environment
	switch environment {
	case "production":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "development":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// 设置日志输出
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
