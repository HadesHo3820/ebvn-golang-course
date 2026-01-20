package infrastructure

import (
	"github.com/HadesHo3820/ebvn-golang-course/internal/model"
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

// MigrateDB migrates db the database according to the User struct.
// It will create the table if it doesn't exist and update the schema if it's outdated.
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
