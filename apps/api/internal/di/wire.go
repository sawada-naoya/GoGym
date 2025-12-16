//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"gogym-api/internal/adapter/handler"
	gymrepo "gogym-api/internal/adapter/repository/gym"
	userrepo "gogym-api/internal/adapter/repository/user"
	workoutrepo "gogym-api/internal/adapter/repository/workout"
	gymuc "gogym-api/internal/usecase/gym"
	sessionuc "gogym-api/internal/usecase/session"
	useruc "gogym-api/internal/usecase/user"
	workoutuc "gogym-api/internal/usecase/workout"
)

type Handlers struct {
	User    *handler.UserHandler
	Session *handler.SessionHandler
	Gym     *handler.GymHandler
	Workout *handler.WorkoutHandler
}

func NewHandlers(
	user *handler.UserHandler,
	session *handler.SessionHandler,
	gym *handler.GymHandler,
	workout *handler.WorkoutHandler,
) *Handlers {
	return &Handlers{
		User:    user,
		Session: session,
		Gym:     gym,
		Workout: workout,
	}
}

func Initialize(e *echo.Echo, db *gorm.DB) *Handlers {
	wire.Build(
		userrepo.NewUserRepository,
		gymrepo.NewGymRepository,
		workoutrepo.NewWorkoutRepository,

		useruc.NewUserInteractor,
		sessionuc.NewSessionInteractor,
		gymuc.NewGymInteractor,
		workoutuc.NewWorkoutInteractor,

		handler.NewUserHandler,
		handler.NewSessionHandler,
		handler.NewGymHandler,
		handler.NewWorkoutHandler,

		NewHandlers,
	)
	return nil
}
