package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/middleware"
	"gogym-api/internal/configs"

	"github.com/labstack/echo/v4"
)

func WorkoutRoutes(e *echo.Group, w *handler.WorkoutHandler, authCfg configs.AuthConfig) {
	workout := e.Group("/workouts")
	workout.Use(middleware.AuthMiddleware(authCfg.JWTSecret))

	workout.GET("/records", w.GetWorkoutRecords)
	workout.POST("/records", w.CreateWorkoutRecord)
	workout.GET("/parts", w.GetWorkoutParts)
}
