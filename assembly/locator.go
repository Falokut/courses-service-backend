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
	"github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/go-kit/validator"
	"github.com/Falokut/go-kit/worker"
	"github.com/pkg/errors"
)

type Config struct {
	HttpRouter *router.Router
	CleanJob   worker.Job
}

// nolint:funlen
func Locator(
	ctx context.Context,
	logger log.Logger,
	dbCli *db.Client,
	filesCli *client.Client,
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

	userService := service.NewUser(cfg.Auth, authRepo, userRepo)
	user := controller.NewUser(userService)

	roleService := service.NewRole(roleRepo)
	role := controller.NewRole(roleService)

	filesRepo := repository.NewFile(filesCli, cfg.FileStorage.BaseServiceUrl)
	courseRepo := repository.NewCourse(dbCli)
	courseService := service.NewCourse(courseRepo, txRunner, filesRepo)
	course := controller.NewCourse(courseService)

	router := routes.Router{
		Auth:   auth,
		User:   user,
		Role:   role,
		Course: course,
	}

	authMiddleware := routes.NewAuthMiddleware(authRepo)
	validator := validator.New(validator.Ru)
	defaultWrapper := endpoint.DefaultWrapper(logger, endpoint.Log(logger, true, true)).WithValidator(validator)
	wrapperWithoutMaxBodySize := endpoint.DefaultWrapper(logger, nil).WithValidator(validator)
	wrapperWithoutMaxBodySize.Middlewares = []http.Middleware{
		endpoint.RequestId(),
		http.Middleware(endpoint.Log(logger, false, true)),
		endpoint.ErrorHandler(logger),
		endpoint.Recovery(),
	}
	httpRouter := router.InitRoutes(authMiddleware, defaultWrapper, wrapperWithoutMaxBodySize)

	attachmentsCleanerWorker := service.NewAttachmentsCleanerWorker(txRunner, filesRepo)
	attachmentsCleanerWorkerHandler := controller.NewAttachmentsCleanerWorkerHandler(
		attachmentsCleanerWorker, logger)

	return Config{
		HttpRouter: httpRouter,
		CleanJob:   attachmentsCleanerWorkerHandler,
	}, nil
}
