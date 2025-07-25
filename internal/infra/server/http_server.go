package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
	AuthHandler "savebite/internal/app/auth/interface/rest"
	AuthUsecase "savebite/internal/app/auth/usecase"
	AnalysesHandler "savebite/internal/app/ingredient_analyses/interface/rest"
	AnalysesUsecase "savebite/internal/app/ingredient_analyses/usecase"
	"savebite/internal/app/user/repository"
	"savebite/internal/infra/gemini"
	"savebite/internal/middlewares"
	"savebite/pkg/jwt"
	"savebite/pkg/markdown"
	"savebite/pkg/oauth"
)

type HTTPServer interface {
	GetApp() *fiber.App
	Start(socket string)
	MountMiddlewares()
	MountRoutes(db *gorm.DB)
}

type httpServer struct {
	app *fiber.App
}

func NewServer() HTTPServer {
	app := fiber.New()

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(socket string) {
	err := s.app.Listen(socket)

	if err != nil {
		panic(err)
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))
	s.app.Use(helmet.New())
	s.app.Use(logger.New())
	s.app.Use(cache.New())
}

func (s *httpServer) MountRoutes(db *gorm.DB) {
	validator := validator.New()
	oauth := oauth.GoogleOAuth
	jwt := jwt.JWT
	gemini := gemini.Gemini
	md := markdown.Markdown

	app := s.GetApp()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	middleware := middlewares.NewMiddleware(jwt)

	userRepo := repository.NewUserRepo(db)

	authUsecase := AuthUsecase.NewAuthUsecase(userRepo, oauth, jwt)
	analysesUsecase := AnalysesUsecase.NewAnalysesUsecase(gemini, md)

	AuthHandler.NewAuthHandler(v1, validator, authUsecase)
	AnalysesHandler.NewAnalysesHandler(v1, middleware, analysesUsecase)
}
