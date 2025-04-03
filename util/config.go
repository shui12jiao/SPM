package util

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`

	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`

	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	MaxReservationDuration         time.Duration `mapstructure:"MAX_RESERVATION_DURATION"`         // 预约的最大持续时间
	MaxReservationAdvanceDuration  time.Duration `mapstructure:"MAX_RESERVATION_ADVANCE_DURATION"` // 预约的最大提前时间
	CancellableReservationDuration time.Duration `mapstructure:"CANCELLABLE_RESERVATION_DURATION"` // 预约开始前的可取消时间

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
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

		MaxReservationDuration:         parseDuration(MustGetEnvString("MAX_RESERVATION_DURATION")),
		MaxReservationAdvanceDuration:  parseDuration(MustGetEnvString("MAX_RESERVATION_ADVANCE_DURATION")), // 预约的最大提前时间
		CancellableReservationDuration: parseDuration(MustGetEnvString("CANCELLABLE_RESERVATION_DURATION")),

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
