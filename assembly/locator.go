package assembly

import (
	"context"

	"courses-service/conf"
	"courses-service/controller"
	"courses-service/repository"
	"courses-service/routes"
	"courses-service/service"

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
	userRepo := repository.NewUser(dbCli)
	authService := service.NewAuth(cfg.Auth, userRepo)
	authCtrl := controller.NewAuth(authService)

	hrouter := routes.Router{
		Auth: authCtrl,
	}
	authMiddleware := routes.NewAuthMiddleware(cfg.Auth.Access.Secret)
	return Config{
		HttpRouter: hrouter.InitRoutes(authMiddleware, endpoint.DefaultWrapper(logger)),
	}, nil
}
