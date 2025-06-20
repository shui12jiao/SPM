package util

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import PostgreSQL database driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import file source driver
)

const (
	StudentRole = "student"
	AdminRole   = "admin"

	ReservationStatusReserved  = "reserved"
	ReservationStatusCompleted = "completed"
	ReservationStatusCanceled  = "canceled"
	ReservationStatusViolated  = "violated"
)

func InitDB(dbr, dbd, migrationURL string) *sql.DB {
	db, err := sql.Open(dbr, dbd)
	if err != nil {
		log.Fatal().Err(err).Msg("无法连接数据库")
	}

	// 迁移数据库
	migrateDB(migrationURL, dbd)
	// 添加默认管理员用户
	addDefaultAdminUser(db)

	return db
}

func addDefaultAdminUser(db *sql.DB) {
	const defaultAdminUsername = "admin"
	const defaultAdminPassword = "admin123"
	const defaultAdminEmail = "admin@example.com"
	const defaultAdminDepartment = "System"

	// 检查用户是否存在
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM "user" WHERE username = $1`, defaultAdminUsername).Scan(&count)
	if err != nil {
		log.Fatal().Err(err).Msg("查询管理员用户失败")
	}
	if count > 0 {
		log.Info().Msg("默认管理员已存在")
		return
	}

	// 加密密码（使用 bcrypt）
	hashedPwd, err := HashPassword(defaultAdminPassword)
	if err != nil {
		log.Fatal().Err(err).Msg("密码加密失败")
	}

	// 插入默认管理员
	_, err = db.Exec(`
		INSERT INTO "user" (username, password, role, department, email)
		VALUES ($1, $2, $3, $4, $5)
	`, defaultAdminUsername, string(hashedPwd), AdminRole, defaultAdminDepartment, defaultAdminEmail)
	if err != nil {
		log.Fatal().Err(err).Msg("插入默认管理员失败")
	}
	log.Info().Msg("默认管理员用户创建成功")
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
