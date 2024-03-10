package main

import (
	"boonmafarm/api/pkg/models"

	"boonmafarm/api/routes"
)

// Entrypoint for app.
func main() {
	// Load the routes
	r := routes.SetupRouter()

	// Initialize database
	models.SetupDatabase()

	// Start the HTTP API
	r.Run()
}
