package pgk

import (
	"context"
	"log"
	"template/internal/app"
	"template/internal/metric"
	"template/pgk/logger"
	"time"
)

type App struct {
	server app.Server
	logger logger.Logger
	metric metric.Metric
	config *app.Config
}

func NewApp(configPath string, ctx context.Context) *App {
	return &App{
		config: app.NewConfig(configPath),
		server: app.ServerInstance(ctx),
		logger: logger.LoggerInstance(),
		metric: metric.MetricsInstance(),
	}
}

func (app *App) Start(appCtx context.Context) error {
	if err := app.logger.SetLogLevel(app.config.Server.LogLevel); err != nil {
		return err
	}

	if err := app.server.RunServer(appCtx, app.config.Server.BindPort, time.Duration(app.config.Server.StatisticRefresh)); err != nil {
		if err = app.logger.RenameLog(); err != nil {
			return err
		}
		log.Fatal(err)
	}

	<-appCtx.Done()

	defer func() {
		app.logger.Write("Application successfully stopped")
		if err := app.logger.RenameLog(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
