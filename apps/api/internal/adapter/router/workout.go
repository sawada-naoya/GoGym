package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func WorkoutRoutes(e *echo.Group, wh *handler.WorkoutHandler) {
	e.GET("/workouts/records", wh.GetWorkoutRecords)
	e.POST("/workouts/records", wh.CreateWorkoutRecord)
	e.PUT("/workouts/records/:id", wh.UpdateWorkoutRecord)
	e.GET("/workouts/parts", wh.GetWorkoutParts)
	e.POST("/workouts/seed", wh.SeedWorkoutParts)
	e.POST("/workouts/exercises", wh.CreateWorkoutExercise)
	e.DELETE("/workouts/exercises/:id", wh.DeleteWorkoutExercise)
	e.GET("/workouts/exercises/:id/last", wh.GetLastWorkoutRecord)
}
