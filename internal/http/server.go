package http

import (
	"godocapi/internal/config"
	"godocapi/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Server struct {
	app    *fiber.App
	config *config.Config
}

func NewServer(cfg *config.Config, svc *service.DocumentService) *Server {
	app := fiber.New(fiber.Config{
		AppName: "DocAPI v1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	handler := NewDocumentHandler(svc)
	handler.RegisterRoutes(app)

	return &Server{
		app:    app,
		config: cfg,
	}
}

func (s *Server) Run() error {
	return s.app.Listen(s.config.ServerPort)
}
