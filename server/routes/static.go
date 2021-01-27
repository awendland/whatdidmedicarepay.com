package routes

import (
	"log"
	"net/http"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
)

func StaticHandler(config *config.Config) http.Handler {
	log.Printf("Serving static files from %s", config.StaticDir)
	return http.FileServer(http.Dir(config.StaticDir))
}
