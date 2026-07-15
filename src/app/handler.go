// Package app provides the serverless HTTP handler for the Farm API.
// It lives under src/ so it can import internal packages; the Vercel entrypoint in api/ only imports this package.
package app

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/di"
	apphandler "github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	appmiddleware "github.com/weeranieb/boonmafarm-backend/src/internal/middleware"
	"github.com/weeranieb/boonmafarm-backend/src/internal/router"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/logging"

	_ "github.com/weeranieb/boonmafarm-backend/docs"

	"gorm.io/gorm"
)

var (
	once   sync.Once
	hndlr  http.HandlerFunc
	loadFn = config.LoadConfig
)

// Handler is the HTTP handler for each request; Fiber handles routing.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Vercel only invokes this function for /api/index. We pass the original path in __path so
	// the rewrite stays destination=/api/index and we restore the path here for Fiber.
	if orig := r.URL.Query().Get("__path"); orig != "" {
		if orig[0] != '/' {
			orig = "/" + orig
		}
		r.URL.Path = orig
		r.URL.RawPath = ""
		r.RequestURI = orig
		if r.URL.RawQuery != "" {
			q := r.URL.Query()
			q.Del("__path")
			r.URL.RawQuery = q.Encode()
			r.RequestURI = orig
			if r.URL.RawQuery != "" {
				r.RequestURI += "?" + r.URL.RawQuery
			}
		}
	} else {
		r.RequestURI = r.URL.String()
	}
	once.Do(func() {
		fiberApp := buildApp()
		hndlr = adaptor.FiberApp(fiberApp)
		slog.Info("app ready")
	})
	hndlr.ServeHTTP(w, r)
}

func buildApp() *fiber.App {
	conf := loadFn()
	logging.Init(conf.App.Environment, conf.App.LogLevel)
	slog.Info("serverless cold start")

	container := di.NewContainer(conf)

	fiberApp := fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024,
		ErrorHandler:   appmiddleware.ErrorHandler,
	})
	logging.UseHTTP(fiberApp)

	var handlers *apphandler.Handler
	var db *gorm.DB
	if err := container.Invoke(func(h *apphandler.Handler, d *gorm.DB) {
		handlers = h
		db = d
	}); err != nil {
		panic("DI: " + err.Error())
	}

	router.SetupRoutes(fiberApp, conf, handlers, db)
	return fiberApp
}
