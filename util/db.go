package util

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import PostgreSQL database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import file source driver
)

func InitDB(dbr, dbd, migrationURL string) *sql.DB {
	db, err := sql.Open(dbr, dbd)
	if err != nil {
		log.Fatal().Err(err).Msg("无法连接数据库")
	}

	// 迁移数据库
	migrateDB(migrationURL, dbd)

	return db
}

func migrateDB(migrationURL string, databaseSource string) {
	m, err := migrate.New(migrationURL, databaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("无法创建迁移实例")
	}

	err = m.Up()
	switch err {
	case nil:
		log.Info().Msg("迁移成功")
	case migrate.ErrNoChange:
		log.Info().Msg("无需迁移")
	default:
		log.Fatal().Err(err).Msg("迁移失败")
	}
}
