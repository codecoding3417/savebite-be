package bootstrap

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"savebite/internal/domain/env"
	"savebite/internal/infra/database"
	"savebite/internal/infra/server"
)

func Init() {
	db := database.NewMySQLConn()

	server := server.NewServer()
	app := server.GetApp()

	app.Get("/metrics", monitor.New())

	server.MountMiddlewares()
	server.MountRoutes(db)
	server.Start(env.AppEnv.AppPort)
}
