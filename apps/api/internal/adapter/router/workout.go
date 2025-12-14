package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/configs"
	"gogym-api/internal/infra/middleware"

	"github.com/labstack/echo/v4"
)

func WorkoutRoutes(e *echo.Group, w *handler.WorkoutHandler, authCfg configs.AuthConfig) {
	workout := e.Group("/workouts")
	workout.Use(middleware.AuthMiddleware(authCfg.JWTSecret))

	workout.GET("/records", w.GetWorkoutRecords)
	workout.POST("/records", w.CreateWorkoutRecord)
	workout.PUT("/records/:id", w.UpdateWorkoutRecord)
	workout.GET("/parts", w.GetWorkoutParts)
	workout.POST("/seed", w.SeedWorkoutParts)
	workout.POST("/exercises/bulk", w.CreateWorkoutExercise)
	workout.DELETE("/exercises/:id", w.DeleteWorkoutExercise)
	workout.GET("/exercises/:id/last", w.GetLastWorkoutRecord)
}
