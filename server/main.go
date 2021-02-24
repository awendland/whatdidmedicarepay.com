package main

import (
	"fmt"
	"time"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/awendland/whatdidmedicarepay.com/server/routes"
	"golang.org/x/time/rate"

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
	e.Use(middleware.BodyLimit(config.Env.RequestBodyLimit))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(config.Env.RequestTimeoutSec) * time.Second,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:  rate.Limit(config.Env.RequestRateLimitPerSec),
		Burst: config.Env.RequestBurstLimitPerSec,
	})))
	// TODO e.Use(middleware.Gzip())
	// TODO autotls

	// Routes
	routes.LoadRoutes(e, config)

	// Start server
	addr := fmt.Sprintf(":%d", config.Env.Port)
	e.Logger.Fatal(e.Start(addr))
}
