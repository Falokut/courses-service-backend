package assembly

import (
	"context"
	"fmt"

	"github.com/Falokut/go-kit/http"
	"github.com/Falokut/go-kit/http/client"
	"github.com/Falokut/go-kit/worker"

	"courses-service/conf"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/config"
	"github.com/Falokut/go-kit/healthcheck"
	"github.com/Falokut/go-kit/log"
	"github.com/pkg/errors"
)

type Assembly struct {
	logger                   log.Logger
	db                       *db.Client
	filesCli                 *client.Client
	server                   *http.Server
	attachmentsCleanerWorker *worker.Worker
	healthcheckManager       *healthcheck.Manager
	localCfg                 conf.LocalConfig
}

func New(ctx context.Context, logger log.Logger) (*Assembly, error) {
	localCfg := conf.LocalConfig{}
	err := config.Read(&localCfg)
	if err != nil {
		return nil, errors.WithMessage(err, "read local config")
	}
	dbCli, err := db.NewDB(ctx, localCfg.DB, db.WithMigrationRunner("./migrations", logger))
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}
	server := http.NewServer(logger)
	filesCli := client.Default()

	locatorCfg, err := Locator(ctx, logger, dbCli, filesCli, localCfg)
	if err != nil {
		return nil, errors.WithMessage(err, "locator config")
	}
	server.Upgrade(locatorCfg.HttpRouter)
	attachmentsCleanerWorker := worker.New(locatorCfg.CleanJob)

	healthcheckManager := healthcheck.NewHealthManager(logger, fmt.Sprint(localCfg.HealthcheckPort))
	healthcheckManager.Register("db", dbCli.PingContext)

	return &Assembly{
		logger:                   logger,
		localCfg:                 localCfg,
		db:                       dbCli,
		filesCli:                 filesCli,
		server:                   server,
		healthcheckManager:       &healthcheckManager,
		attachmentsCleanerWorker: attachmentsCleanerWorker,
	}, nil
}

func (a *Assembly) Runners() []app.RunnerFunc {
	return []app.RunnerFunc{
		func(_ context.Context) error {
			return a.server.ListenAndServe(a.localCfg.Listen.GetAddress())
		},
		func(_ context.Context) error {
			return a.healthcheckManager.RunHealthcheckEndpoint()
		},
		func(ctx context.Context) error {
			a.attachmentsCleanerWorker.Run(ctx)
			return nil
		},
	}
}

func (a *Assembly) Closers() []app.CloserFunc {
	return []app.CloserFunc{
		func(_ context.Context) error {
			return a.db.Close()
		},
		func(ctx context.Context) error {
			return a.server.Shutdown(ctx)
		},
	}
}
