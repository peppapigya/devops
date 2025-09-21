package db

import (
	"k8s-platform-go/internal/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// NewDB 连接数据库
func NewDB() *gorm.DB {
	databaseConfig := config.GetGlobalConfig().DataBase
	var err error
	db, err = gorm.Open(mysql.Open(databaseConfig.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("db init faild: %v", err)
		return nil
	}
	sqlDb, _ := db.DB()
	// 设置最大空闲连接数
	sqlDb.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDb.SetMaxOpenConns(100)
	// 设置每个连接的过期时间
	sqlDb.SetConnMaxLifetime(time.Hour)
	return db
}

func CloseDB() {
	sqlDb, _ := db.DB()
	if err := sqlDb.Close(); err != nil {
		log.Printf("db close faild: %v", err)
		return
	}
}
