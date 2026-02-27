// Package handler is the Vercel serverless entrypoint.
// It only imports src/app (no internal), so Vercel's build allows it.
package handler

import (
	"net/http"

	"github.com/weeranieb/boonmafarm-backend/src/app"
)

// Handler is the entrypoint Vercel calls for every request.
func Handler(w http.ResponseWriter, r *http.Request) {
	app.Handler(w, r)
}
