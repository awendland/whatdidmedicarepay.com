package routes

import (
	"log"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoadRoutes creates a ServeMux populated with route handlers for
// whatdidmedicarepay.com
func LoadRoutes(e *echo.Echo, config *config.Config) {
	// Static handler at root
	log.Printf("Serving static files from %s", config.Env.StaticDir)
	e.Use(middleware.Static(config.Env.StaticDir))

	// API routes
	apiRoutes := e.Group("/api")
	apiRoutes.GET("/search", NewAPISearchHandler(config))
}
