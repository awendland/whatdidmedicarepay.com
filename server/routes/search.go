package routes

import (
	"net/http"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
	"github.com/awendland/whatdidmedicarepay.com/server/repositories"
	"github.com/labstack/echo/v4"
)

type ApiSearchResponse struct {
	Query   string
	Results []repositories.ProcedureProviderEntry
}

func NewApiSearchHandler(config *config.Config) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		sql := c.Param("sql")

		resp := ApiSearchResponse{
			Query:   sql,
			Results: nil,
		}
		return c.JSON(http.StatusOK, resp)
	}
}
