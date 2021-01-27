package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/awendland/whatdidmedicarepay.com/server/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config.ProvideConfig()

	routes := routes.LoadRoutes(config)

	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Listening on %s...\n", addr)
	err := http.ListenAndServe(addr, routes)
	if err != nil {
		log.Fatal(err)
	}
}
