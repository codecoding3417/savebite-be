package main

import (
	"savebite/internal/domain/entity"
	"savebite/internal/infra/database"
	"savebite/pkg/log"
)

func main() {
	db := database.NewMySQLConn()

	err := db.Migrator().DropTable(&entity.User{}, &entity.Analysis{}, &entity.Ingredient{})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Migrate][main] Failed to drop all database")
	}

	err = db.AutoMigrate(&entity.Ingredient{}, &entity.Analysis{}, &entity.User{})
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[Migrate][main] Failed to migrate database")
	}

	log.Info(nil, "Migration success")
}
