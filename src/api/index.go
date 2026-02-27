// Package handler provides the Vercel serverless entrypoint for the Farm API.
// Vercel invokes Handler(w, r) for each request; the Fiber app handles routing.
package handler

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/di"
	apphandler "github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/router"

	_ "github.com/weeranieb/boonmafarm-backend/docs"
)

var (
	once   sync.Once
	hndlr  http.HandlerFunc
	loadFn = config.LoadConfig
)

// Handler is the entrypoint Vercel calls for every request.
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()
	once.Do(func() {
		app := buildApp()
		hndlr = adaptor.FiberApp(app)
	})
	hndlr.ServeHTTP(w, r)
}

func buildApp() *fiber.App {
	conf := loadFn()
	container := di.NewContainer(conf)

	app := fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024,
	})

	var handlers *apphandler.Handler
	if err := container.Invoke(func(h *apphandler.Handler) {
		handlers = h
	}); err != nil {
		panic("DI: " + err.Error())
	}

	router.SetupRoutes(app, conf, handlers)
	return app
}
