//go:build wireinject

package di

import (
	"gogym-api/internal/configs"

	"github.com/google/wire"

	"gogym-api/internal/adapter/auth"
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
	repository.NewUserRepository, // 具象
	wire.Bind(new(useruc.Repository), new(*repository.UserRepository)),
	wire.Bind(new(sessionuc.UserRepository), new(*repository.UserRepository)),

	repository.NewGymRepository,
	repository.NewReviewRepository,
	repository.NewTagRepository,
)

var internalPlatformSet = wire.NewSet(
	auth.NewBcryptPasswordHasher,
	wire.Bind(new(useruc.PasswordHasher), new(*auth.BcryptPasswordHasher)),
	wire.Bind(new(sessionuc.PasswordHasher), new(*auth.BcryptPasswordHasher)),
)

var UsecaseSet = wire.NewSet(
	internalPlatformSet,
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
	router.RegisterRoutes,
)

var InfraSet = wire.NewSet(
	db.NewDB,
)

func provideHTTP(c *configs.Config) configs.HTTPConfig { return c.HTTP }

func provideDB(c *configs.Config) configs.DatabaseConfig { return c.Database }

func BuildServer(cfg *configs.Config) (*echo.Echo, func(), error) {
	wire.Build(
		provideHTTP, provideDB, // cfg → サブ設定
		InfraSet,      // db
		RepositorySet, // repo
		UsecaseSet,    // usecase（内側でauth/serviceを吸収）
		HandlerSet,    // handler
		ServerSet,     // router(server)
	)
	return nil, nil, nil
}
