package util

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string // 环境变量，开发环境、生产环境

	DBDriver     string // 数据库驱动
	DBSource     string // 数据库源
	MigrationURL string // 数据库迁移源
	RedisAddress string // Redis地址 暂时不使用

	HTTPServerAddress    string        // HTTP服务器地址
	TokenSymmetricKey    string        // 对称密钥
	AccessTokenDuration  time.Duration // 访问令牌的有效期
	RefreshTokenDuration time.Duration // 刷新令牌的有效期

	BusinessConfig // 业务相关配置

	EmailConfig // 邮件配置
}

// 业务参数配置，可通过api修改
type BusinessConfig struct {
	// 预约相关配置
	MinReservationDuration          time.Duration // 预约的最小持续时间
	MaxReservationDuration          time.Duration // 预约的最大持续时间
	MinReservationAdvanceDuration   time.Duration // 预约的最小提前时间
	MaxReservationAdvanceDuration   time.Duration // 预约的最大提前时间
	CancellableReservationDuration  time.Duration // 预约开始前的可取消时间
	ReservationRemindBeforeDuration time.Duration // 预约开始前的提醒时间
	ReservationRemindAfterDuration  time.Duration // 预约开始后的提醒时间
	ReservationViolationDuration    time.Duration // 预约开始后的违约处理时间
}

func LoadConfig(path string) Config {
	// 读取环境变量配置文件
	// 默认读取工作目录下的.env文件
	// 也可以通过传入参数来指定其他文件
	godotenv.Load(path)

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

		BusinessConfig: BusinessConfig{
			MinReservationDuration:          parseDuration(MustGetEnvString("MIN_RESERVATION_DURATION")),
			MaxReservationDuration:          parseDuration(MustGetEnvString("MAX_RESERVATION_DURATION")),
			MinReservationAdvanceDuration:   parseDuration(MustGetEnvString("MIN_RESERVATION_ADVANCE_DURATION")),
			MaxReservationAdvanceDuration:   parseDuration(MustGetEnvString("MAX_RESERVATION_ADVANCE_DURATION")),
			CancellableReservationDuration:  parseDuration(MustGetEnvString("CANCELLABLE_RESERVATION_DURATION")),
			ReservationRemindBeforeDuration: parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),
			ReservationRemindAfterDuration:  parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),
			ReservationViolationDuration:    parseDuration(MustGetEnvString("RESERVATION_REMIND_BEFORE_DURATION")),
		},

		EmailConfig: EmailConfig{
			SenderName:   MustGetEnvString("EMAIL_SENDER_NAME"),
			SenderEmail:  MustGetEnvString("EMAIL_SENDER_ADDRESS"),
			SMTPPassword: MustGetEnvString("EMAIL_SENDER_PASSWORD"),
			SMTPHost:     MustGetEnvString("EMAIL_SMTP_HOST"),
			SMTPPort:     MustGetEnvInt("EMAIL_SMTP_PORT"),
		},
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

func MustGetEnvInt(key string) int {
	s := MustGetEnvString(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal().Err(err).Msgf("环境变量 %s 转换为整数失败", key)
	}
	return i
}
