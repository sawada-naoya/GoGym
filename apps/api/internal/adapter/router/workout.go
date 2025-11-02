package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/middleware"

	"github.com/labstack/echo/v4"
)

func WorkoutRoutes(e *echo.Group, w *handler.WorkoutHandler) {
	workout := e.Group("/workouts")
	workout.Use(middleware.AuthMiddleware)

	workout.GET("/records", w.GetWorkoutRecords)
	workout.POST("/records", w.CreateWorkoutRecord)
}
