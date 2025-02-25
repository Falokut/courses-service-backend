package assembly

import (
	"context"

	"courses-service/conf"
	"courses-service/controller"
	"courses-service/repository"
	"courses-service/routes"
	"courses-service/service"
	"courses-service/transaction"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/go-kit/validator"
	"github.com/pkg/errors"
)

type Config struct {
	HttpRouter *router.Router
}

// nolint:funlen
func Locator(
	ctx context.Context,
	logger log.Logger,
	dbCli *db.Client,
	imagesCli *client.Client,
	cfg conf.LocalConfig,
) (Config, error) {
	txRunner := transaction.NewManager(dbCli)

	authRepo := repository.NewAuth(dbCli)
	userRepo := repository.NewUser(dbCli)
	roleRepo := repository.NewRole(dbCli)

	authService := service.NewAuth(cfg.Auth, authRepo, userRepo, roleRepo, txRunner)
	if cfg.Auth.InitAdmin != nil {
		err := authService.InitAdmin(ctx, *cfg.Auth.InitAdmin)
		if err != nil {
			return Config{}, errors.WithMessage(err, "init admin")
		}
	}
	auth := controller.NewAuth(authService)

	userService := service.NewUser(authRepo, userRepo)
	user := controller.NewUser(userService)

	roleService := service.NewRole(roleRepo)
	role := controller.NewRole(roleService)
	router := routes.Router{
		Auth: auth,
		User: user,
		Role: role,
	}
	authMiddleware := routes.NewAuthMiddleware(authRepo)
	validator := validator.New(validator.Ru)
	return Config{
		HttpRouter: router.InitRoutes(authMiddleware, endpoint.DefaultWrapper(logger).WithValidator(validator)),
	}, nil
}
