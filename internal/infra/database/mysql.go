package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
)

func NewMySQLConn() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBHost,
		env.AppEnv.DBPort,
		env.AppEnv.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[MySQL][NewMySQLConn] Failed to connect database")
	}

	return db
}
