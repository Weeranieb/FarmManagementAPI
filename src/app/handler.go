// Package app provides the serverless HTTP handler for the Farm API.
// It lives under src/ so it can import internal packages; the Vercel entrypoint in api/ only imports this package.
package app

import (
	"log"
	"net/http"
	"strings"
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
	// Vercel rewrites send /api/index/:path* so the request path is /api/index/health, /api/index/api/v1/..., etc.
	// Restore the original path so Fiber can route correctly.
	const prefix = "/api/index"
	if path := r.URL.Path; strings.HasPrefix(path, prefix) {
		orig := strings.TrimPrefix(path, prefix)
		if orig == "" {
			orig = "/"
		} else if orig[0] != '/' {
			orig = "/" + orig
		}
		r.URL.Path = orig
		r.RequestURI = orig
		if r.URL.RawQuery != "" {
			r.RequestURI += "?" + r.URL.RawQuery
		}
	} else {
		r.RequestURI = r.URL.String()
	}
	once.Do(func() {
		log.Println("[Farm API] serverless cold start – building app")
		fiberApp := buildApp()
		hndlr = adaptor.FiberApp(fiberApp)
		log.Println("[Farm API] app ready")
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
