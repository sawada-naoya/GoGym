package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupHealthRoutes(e *echo.Echo) {
	e.GET("/healthz", HealthCheck)
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"service": "gogym-api",
	})
}