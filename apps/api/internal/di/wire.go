//go:build wireinject

package di

import (
	"gogym-api/internal/infra/security"
	"gogym-api/internal/infra/slack"

	"github.com/google/wire"
	"gorm.io/gorm"

	handler "gogym-api/internal/adapter/handler"
	gymrepo "gogym-api/internal/adapter/repository/gym"
	userrepo "gogym-api/internal/adapter/repository/user"
	workoutrepo "gogym-api/internal/adapter/repository/workout"

	contactuc "gogym-api/internal/application/contact"
	gymuc "gogym-api/internal/application/gym"
	sessionuc "gogym-api/internal/application/session"
	useruc "gogym-api/internal/application/user"
	workoutuc "gogym-api/internal/application/workout"
)

type Handlers struct {
	User    *handler.UserHandler
	Session *handler.SessionHandler
	Gym     *handler.GymHandler
	Workout *handler.WorkoutHandler
	Contact *handler.ContactHandler
}

func NewHandlers(
	user *handler.UserHandler,
	session *handler.SessionHandler,
	gym *handler.GymHandler,
	workout *handler.WorkoutHandler,
	contact *handler.ContactHandler,
) *Handlers {
	return &Handlers{
		User:    user,
		Session: session,
		Gym:     gym,
		Workout: workout,
		Contact: contact,
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
	contactuc.NewContactInteractor,
)

var handlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewSessionHandler,
	handler.NewGymHandler,
	handler.NewWorkoutHandler,
	handler.NewContactHandler,
	NewHandlers,
)

// provideSlackGateway converts *slack.Client to contactuc.SlackGateway interface
func provideSlackGateway(client *slack.Client) contactuc.SlackGateway {
	return client
}

var gatewaySet = wire.NewSet(
	provideSlackGateway,
)

func Initialize(db *gorm.DB, slackClient *slack.Client, jwtSecret string) *Handlers {
	wire.Build(
		repositorySet,
		securitySet,
		gatewaySet,
		usecaseSet,
		handlerSet,
	)
	return nil
}
