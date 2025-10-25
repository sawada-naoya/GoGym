package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func WorkoutRoutes(e *echo.Group, w *handler.WorkoutHandler) {
	workout := e.Group("/workouts")

	workout.GET("/records", w.GetWorkoutRecords)
}
