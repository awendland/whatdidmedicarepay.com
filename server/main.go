package main

import (
	"fmt"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/awendland/whatdidmedicarepay.com/server/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Config
	e := echo.New()
	config := config.ProvideConfig()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))

	// Routes
	routes.LoadRoutes(e, config)

	// Start server
	addr := fmt.Sprintf(":%d", config.Env.Port)
	e.Logger.Fatal(e.Start(addr))
}
