package infrastructure

import (
	"github.com/HadesHo3820/ebvn-golang-course/pkg/common"
	redisPkg "github.com/HadesHo3820/ebvn-golang-course/pkg/redis"
	"github.com/HadesHo3820/ebvn-golang-course/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CreateRedisGeneralConn creates a Redis connection for general purposes (DB 0).
// This connection is reserved for non-cache data such as:
//   - Session storage
//   - Rate limiting counters
//   - Temporary application state
//   - Feature flags
//
// Returns a Redis client connected to database 0.
func CreateRedisGeneralConn() *redis.Client {
	redisClient, err := redisPkg.NewClientWithDB("", 0)
	common.HandleError(err)
	return redisClient
}

// CreateRedisCacheConn creates a Redis connection specifically for caching (DB 1).
// This connection should be used for all cache-related operations such as:
//   - Bookmark caching
//   - Query result caching
//   - Any other application-level caching
//
// Returns a Redis client connected to database 1.
func CreateRedisCacheConn() *redis.Client {
	redisClient, err := redisPkg.NewClientWithDB("", 1)
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
