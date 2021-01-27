package routes

import (
	"net/http"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
)

// LoadRoutes creates a ServeMux populated with route handlers for
// whatdidmedicarepay.com
func LoadRoutes(config *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", StaticHandler(config))
	mux.HandleFunc("/api/search", NewApiSearchHandler(config))
	return mux
}
