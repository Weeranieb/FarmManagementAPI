package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/di"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/router"

	_ "github.com/weeranieb/boonmafarm-backend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/dig"
)

var (
	app *fiber.App
)

var (
	LoadConfigFunc = config.LoadConfig
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
	conf := LoadConfigFunc()

	// Dependency Injection
	container := di.NewContainer(conf)

	// Start Fiber + Router
	setupAndStartServer(conf, container)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			log.Println("Gracefully shutting down...")
			shutdownServer()
		}
	}()

	log.Println("Starting server on " + conf.GetServerAddress())
	if err := app.Listen(conf.GetServerAddress()); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func setupAndStartServer(conf *config.Config, container *dig.Container) {
	app = fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024, // 10MB
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// Construct the Handler using DI container
	var handlers *handler.Handler

	err := container.Invoke(func(h *handler.Handler) {
		handlers = h
	})
	if err != nil {
		log.Fatal("DI error", err)
	}

	router.SetupRoutes(app, conf, handlers)
	app.Listen(conf.GetServerAddress())
}

func shutdownServer() {
	log.Println("Fiber was successfully shut down.")

	if err := app.Shutdown(); err != nil {
		log.Fatal("Error shutting down Fiber", err)
	}
	os.Exit(0)
}
