package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
)

type ApiSearchResponse struct {
	Query   string
	Results []string
}

func NewApiSearchHandler(config *config.Config) http.HandlerFunc {
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query()
		sql := p["sql"][0]
		w.Header().Set("content-type", "application/json")
		resp := ApiSearchResponse{
			Query:   sql,
			Results: nil,
		}
		jsonstr, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(jsonstr)
	}
}
