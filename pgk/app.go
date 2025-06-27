package pgk

import (
	"context"
	"log"
	"template/internal/app"
	"template/internal/metric"
	"template/pgk/logger"
)

type App struct {
	server app.Server
	logger logger.Logger
	metric metric.Metric
}

func NewApp(configPath string, ctx context.Context) *App {
	return &App{
		server: app.ServerInstance(ctx, configPath),
		logger: logger.LoggerInstance(),
		metric: metric.MetricsInstance(),
	}
}

func (app *App) Start(appCtx context.Context) error {
	go func() {
		if err := app.server.RunServer(appCtx); err != nil {
			if err = app.logger.RenameLog(); err != nil {
				return
			}
			log.Fatal(err)
		}
	}()

	<-appCtx.Done()

	defer func() {
		app.logger.Write("Application successfully stopped")
		if err := app.logger.RenameLog(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
