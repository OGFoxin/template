package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"template/internal/metric"
	"template/pgk/logger"
	"time"
)

var instance Server
var once sync.Once

type Server interface {
	RunServer(context.Context) error
	createRouters(*gin.Engine) error
	PrintActualStatistic()
	applyConfiguration()
	setupGinInstance() *gin.Engine
}

type server struct {
	logger logger.Logger
	metric metric.Metric
	ticker time.Ticker
	config *Config
}

func ServerInstance(ctx context.Context, configPath string) Server {
	once.Do(func() {
		instance = NewServer(ctx, configPath)
	})

	return instance
}

func NewServer(ctx context.Context, configPath string) Server {
	return &server{
		config: NewConfig(configPath),
		logger: logger.LoggerInstance(),
		metric: metric.NewMetrics(),
	}
}

func (s *server) PrintActualStatistic() {
	var wg = sync.WaitGroup{}
	for {
		select {
		case <-s.ticker.C:
			// if metrics no data in map  don't record's into logs
			s.logger.WriteStatisticToLog(s.metric.GetHttpStats())
			s.logger.WriteCpuInfoToLog(s.metric.GetCpuInfo())
			s.logger.WriteMemoryInfoToLog(s.metric.GetMemoryInfo())

			wg.Add(1)
			go func() {
				if err := s.metric.ResetHttpStat(&wg); err != nil {
					log.Fatal(err)
				}
			}()

			wg.Wait()
		}
	}

}

func (s *server) applyConfiguration() {
	s.ticker = *time.NewTicker(time.Duration(s.config.Server.StatisticRefresh) * time.Second)

	if err := s.logger.SetLogLevel(s.config.Server.LogLevel); err != nil {
		log.Fatal(err)
	}

	s.logger.Write("Configuration updated successfully")
}

func (s *server) setupGinInstance() *gin.Engine {
	// ** GIN configuration ** //
	// set release mode
	gin.SetMode(gin.ReleaseMode)
	// more control for configuration
	ginRouter := gin.New()
	gin.DisableConsoleColor()

	if s.config.Server.UseGin {
		gin.DefaultWriter = s.logger.GetLogFile()
		ginRouter.Use(
			gin.Logger(),
		)
	}

	ginRouter.Use(
		gin.Recovery(),
	)

	return ginRouter
}

func (s *server) RunServer(srvCtx context.Context) error {
	ginRouter := s.setupGinInstance()
	s.ticker = *time.NewTicker(time.Duration(s.config.Server.StatisticRefresh) * time.Second)
	if err := s.logger.SetLogLevel(s.config.Server.LogLevel); err != nil {
		return err
	}

	go func() {
		var err error
		for {
			if s.config, err = s.config.Watchdog("configs/config.yml"); err != nil {
				s.logger.Write("Cant reload config file", err)
			}
			s.applyConfiguration()
		}
	}()

	s.logger.Write("Application waked up")

	go func() {
		if err := s.createRouters(ginRouter); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := ginRouter.Run(s.config.Server.BindPort); err != nil {
			panic(err)
		}
	}()

	go func() {
		s.PrintActualStatistic()
	}()

	<-srvCtx.Done()

	return nil
}
