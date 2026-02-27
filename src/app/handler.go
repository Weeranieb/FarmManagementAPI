// Package app provides the serverless HTTP handler for the Farm API.
// It lives under src/ so it can import internal packages; the Vercel entrypoint in api/ only imports this package.
package app

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

// Handler is the HTTP handler for each request; Fiber handles routing.
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()
	once.Do(func() {
		fiberApp := buildApp()
		hndlr = adaptor.FiberApp(fiberApp)
	})
	hndlr.ServeHTTP(w, r)
}

func buildApp() *fiber.App {
	conf := loadFn()
	container := di.NewContainer(conf)

	fiberApp := fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024,
	})

	var handlers *apphandler.Handler
	if err := container.Invoke(func(h *apphandler.Handler) {
		handlers = h
	}); err != nil {
		panic("DI: " + err.Error())
	}

	router.SetupRoutes(fiberApp, conf, handlers)
	return fiberApp
}
