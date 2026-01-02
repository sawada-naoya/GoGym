package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func ContactRoutes(e *echo.Group, ch *handler.ContactHandler) {
	e.POST("/contact", ch.PostContact)
}
