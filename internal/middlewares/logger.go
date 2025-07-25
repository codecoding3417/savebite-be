package middlewares

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
)

func Logger() fiber.Handler {
	fields := []string{
		"referer",
		"ip",
		"url",
		"latency",
		"status",
		"method",
		"error",
	}

	if env.AppEnv.AppEnv == "development" {
		fields = append(fields, "body")
		fields = append(fields, "reqHeaders")
		fields = append(fields, "resHeaders")
	}

	log := log.GetLogger()
	return fiberzerolog.New(fiberzerolog.Config{
		Logger:          log,
		Fields:          fields,
		FieldsSnakeCase: true,
		Messages: []string{
			"[LoggerMiddleware.LoggerConfig] Server error",
			"[LoggerMiddleware.LoggerConfig] Client error",
			"[LoggerMiddleware.LoggerConfig] Success",
		},
	})

}
