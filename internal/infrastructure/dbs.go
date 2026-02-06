package infrastructure

import (
	"github.com/HadesHo3820/ebvn-golang-course/pkg/common"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CreateRedisConn creates a new redis connection
func CreateRedisConn() *redis.Client {
	// Create redis db connection
	redisClient, err := redisPkg.NewClient("")
	common.HandleError(err)
	return redisClient
}

// CreateSQLDBWithMigration creates a new sql db connection and migrate db
func CreateSQLDBWithMigration() *gorm.DB {
	// Create sql db connection
	db, err := sqldb.NewClient("")
	common.HandleError(err)

	err = MigrateDB(db)
	common.HandleError(err)

	return db

}

// MigrateDB executes all pending database migrations from the "./migrations" directory.
// It ensures the database schema is up-to-date by applying all "up" migration files
// found in the configured migration path.
func MigrateDB(sqlDB *gorm.DB) error {
	return sqldb.MigrateSQLDB(sqlDB, "file://./migrations", "up", 0)
}
