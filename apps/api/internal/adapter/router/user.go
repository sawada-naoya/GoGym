package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group, uh *handler.UserHandler) {
	user := e.Group("/users")

	user.POST("", uh.SignUp)
}
