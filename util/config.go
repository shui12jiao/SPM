package util

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"` // 环境变量，开发环境、生产环境

	DBDriver     string `mapstructure:"DB_DRIVER"`     // 数据库驱动
	DBSource     string `mapstructure:"DB_SOURCE"`     // 数据库源
	MigrationURL string `mapstructure:"MIGRATION_URL"` // 数据库迁移源
	RedisAddress string `mapstructure:"REDIS_ADDRESS"` // Redis地址 暂时不使用

	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`    // HTTP服务器地址
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`    // 对称密钥
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`  // 访问令牌的有效期
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"` // 刷新令牌的有效期

	MaxReservationDuration          time.Duration `mapstructure:"MAX_RESERVATION_DURATION"`           // 预约的最大持续时间
	MaxReservationAdvanceDuration   time.Duration `mapstructure:"MAX_RESERVATION_ADVANCE_DURATION"`   // 预约的最大提前时间
	CancellableReservationDuration  time.Duration `mapstructure:"CANCELLABLE_RESERVATION_DURATION"`   // 预约开始前的可取消时间
	ReservationRemindBeforeDuration time.Duration `mapstructure:"RESERVATION_REMIND_BEFORE_DURATION"` // 预约开始前的提醒时间
	ReservationRemindAfterDuration  time.Duration `mapstructure:"RESERVATION_REMIND_AFTER_DURATION"`  // 预约开始后的提醒时间
	ReservationViolationDuration    time.Duration `mapstructure:"RESERVATION_VIOLATION_DURATION"`     // 预约开始后的违约处理时间

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`     // 邮件发送者名称
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`  // 邮件发送者地址
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"` // 邮件发送者密码
}

func LoadConfig() Config {
	godotenv.Load()

	return Config{
		Environment: MustGetEnvString("ENVIRONMENT"),

		DBDriver:     MustGetEnvString("DB_DRIVER"),
		DBSource:     MustGetEnvString("DB_SOURCE"),
		MigrationURL: MustGetEnvString("MIGRATION_URL"),
		RedisAddress: MustGetEnvString("REDIS_ADDRESS"),

		HTTPServerAddress:    MustGetEnvString("HTTP_SERVER_ADDRESS"),
		TokenSymmetricKey:    MustGetEnvString("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration:  parseDuration(MustGetEnvString("ACCESS_TOKEN_DURATION")),
		RefreshTokenDuration: parseDuration(MustGetEnvString("REFRESH_TOKEN_DURATION")),

		MaxReservationDuration:          parseDuration(MustGetEnvString("MAX_RESERVATION_DURATION")),
		MaxReservationAdvanceDuration:   parseDuration(MustGetEnvString("MAX_RESERVATION_ADVANCE_DURATION")),
		CancellableReservationDuration:  parseDuration(MustGetEnvString("CANCELLABLE_RESERVATION_DURATION")),
		ReservationRemindBeforeDuration: parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),
		ReservationRemindAfterDuration:  parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),
		ReservationViolationDuration:    parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),

		EmailSenderName:     MustGetEnvString("EMAIL_SENDER_NAME"),
		EmailSenderAddress:  MustGetEnvString("EMAIL_SENDER_ADDRESS"),
		EmailSenderPassword: MustGetEnvString("EMAIL_SENDER_PASSWORD"),
	}

}

func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatal().Err(err).Msgf("%s配置格式错误", durationStr)

	}
	return duration
}

func MustGetEnvString(key string) string {
	s := os.Getenv(key)
	if s == "" {
		log.Fatal().Msgf("环境变量 %s 为空", key)
	}
	return s
}
