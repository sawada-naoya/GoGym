//go:build wireinject

package di

import (
	"github.com/google/wire"

	"gogym-api/internal/configs"

	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/repository"
	"gogym-api/internal/adapter/router"
	gymuc "gogym-api/internal/usecase/gym"
	reviewuc "gogym-api/internal/usecase/review"
	sessionuc "gogym-api/internal/usecase/session"
	useruc "gogym-api/internal/usecase/user"

	"gogym-api/internal/infra/db"

	"github.com/labstack/echo/v4"
)

var RepositorySet = wire.NewSet(
	repository.NewGymRepository,
	repository.NewReviewRepository,
	repository.NewUserRepository,
	repository.NewTagRepository,
)

var UsecaseSet = wire.NewSet(
	gymuc.NewInteractor,
	reviewuc.NewInteractor,
	sessionuc.NewInteractor,
	useruc.NewInteractor,
)

var HandlerSet = wire.NewSet(
	handler.NewGymHandler,
	handler.NewReviewHandler,
	handler.NewUserHandler,
	handler.NewSessionHandler,
)

var ServerSet = wire.NewSet(
	router.BuildEcho,
)

var InfraSet = wire.NewSet(
	db.NewDB,
)

func provideHTTP(c configs.Config) configs.HTTPConfig { return c.HTTP }

func provideDB(c configs.Config) configs.DatabaseConfig { return c.Database }

func InitEcho(cfg *configs.Config) (*echo.Echo, func(), error) {
	wire.Build(
		provideHTTP,
		provideDB,
		InfraSet,      // DB
		RepositorySet, // Repo
		UsecaseSet,    // Usecase
		HandlerSet,    // Handler
		ServerSet,
	)
	return nil, nil, nil
}
