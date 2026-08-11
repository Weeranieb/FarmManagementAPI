package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/di"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	appmiddleware "github.com/weeranieb/boonmafarm-backend/src/internal/middleware"
	"github.com/weeranieb/boonmafarm-backend/src/internal/router"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/logging"

	_ "github.com/weeranieb/boonmafarm-backend/docs"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

var (
	app *fiber.App
)

// @title Boonma Farm API
// @version 1.0
// @description A Boonma Farm application with Fiber, GORM, and Dependency Injection
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name jwt_token
// @description JWT token stored in HTTP-only cookie (automatically sent by browser)
func main() {
	conf := config.LoadConfig()
	logging.Init(conf.App.Environment, conf.App.LogLevel)

	// Dependency Injection
	container := di.NewContainer(conf)

	// Wire Fiber + routes; Listen stays in main so the shutdown handler
	// is registered before the blocking call.
	setupServer(conf, container)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			slog.Info("server shutting down")
			shutdownServer()
		}
	}()

	slog.Info("server starting", "address", conf.GetServerAddress())
	if err := app.Listen(conf.GetServerAddress()); err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}

func setupServer(conf *config.Config, container *dig.Container) {
	app = fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024, // 10MB
		ErrorHandler:   appmiddleware.ErrorHandler,
	})

	logging.UseHTTP(app)

	// Construct the Handler using DI container
	var handlers *handler.Handler
	var db *gorm.DB

	err := container.Invoke(func(h *handler.Handler, d *gorm.DB) {
		handlers = h
		db = d
	})
	if err != nil {
		slog.Error("DI error", "err", err)
		os.Exit(1)
	}

	router.SetupRoutes(app, conf, handlers, db)
}

func shutdownServer() {
	slog.Info("fiber shutting down")

	if err := app.Shutdown(); err != nil {
		slog.Error("error shutting down fiber", "err", err)
		os.Exit(1)
	}
	os.Exit(0)
}
