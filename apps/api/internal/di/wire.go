//go:build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	handler "gogym-api/internal/adapter/handler"
	gymrepo "gogym-api/internal/adapter/repository/gym"
	userrepo "gogym-api/internal/adapter/repository/user"
	workoutrepo "gogym-api/internal/adapter/repository/workout"
	"gogym-api/internal/infra/security"
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

var repositorySet = wire.NewSet(
	userrepo.NewUserRepository,
	gymrepo.NewGymRepository,
	workoutrepo.NewWorkoutRepository,
	// Bind user repository to interfaces
	wire.Bind(new(useruc.Repository), new(*userrepo.UserRepository)),
	wire.Bind(new(sessionuc.UserRepository), new(*userrepo.UserRepository)),
)

var securitySet = wire.NewSet(
	security.NewBcryptPasswordHasher,
	wire.Bind(new(useruc.PasswordHasher), new(*security.BcryptPasswordHasher)),
	wire.Bind(new(sessionuc.PasswordHasher), new(*security.BcryptPasswordHasher)),
)

var usecaseSet = wire.NewSet(
	useruc.NewUserInteractor,
	sessionuc.NewSessionInteractor,
	gymuc.NewGymInteractor,
	workoutuc.NewWorkoutInteractor,
)

var handlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewSessionHandler,
	handler.NewGymHandler,
	handler.NewWorkoutHandler,
	NewHandlers,
)

func Initialize(db *gorm.DB) *Handlers {
	wire.Build(
		repositorySet,
		securitySet,
		usecaseSet,
		handlerSet,
	)
	return nil
}
