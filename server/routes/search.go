package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/awendland/whatdidmedicarepay.com/server/repositories"
	"github.com/labstack/echo/v4"
)

// APISearchResponse is the structured response for the endpoint
type APISearchResponse struct {
	Query   string
	Limit   int
	Results []repositories.PUPDEntry
}

// NewAPISearchHandler creates a function for running full-text search queries against the
// Procedure and Provider data set.
func NewAPISearchHandler(config *config.Config) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var query string
		limit := int(50) // default limit is 50
		err = echo.QueryParamsBinder(c).
			String("query", &query).
			Int("limit", &limit).
			BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if limit > 200 {
			return echo.NewHTTPError(http.StatusBadRequest, "`limit` must be less than 100")
		}
		if query == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "`query` must be specified as a URL parameter")
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
		defer cancel()
		results, err := repositories.SearchPUPDEntries(ctx, config, query, limit)
		if err != nil {
			log.Print(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		resp := APISearchResponse{
			Query:   query,
			Limit:   limit,
			Results: results,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
