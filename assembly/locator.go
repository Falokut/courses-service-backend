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

	userRepo := repository.NewUser(dbCli)
	authRepo := repository.NewAuth(dbCli)

	authService := service.NewAuth(cfg.Auth, userRepo, txRunner)
	auth := controller.NewAuth(authService)

	userService := service.NewUser(authRepo)
	user := controller.NewUser(userService)
	hrouter := routes.Router{
		Auth: auth,
		User: user,
	}
	authMiddleware := routes.NewAuthMiddleware(authRepo)
	return Config{
		HttpRouter: hrouter.InitRoutes(authMiddleware, endpoint.DefaultWrapper(logger)),
	}, nil
}
