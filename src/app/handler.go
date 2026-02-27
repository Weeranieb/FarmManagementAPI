// Package app provides the serverless HTTP handler for the Farm API.
// It lives under src/ so it can import internal packages; the Vercel entrypoint in api/ only imports this package.
package app

import (
	"log"
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
	log.Printf("[Farm API] %s %s", r.Method, r.URL.Path)
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
