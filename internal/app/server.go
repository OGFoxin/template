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
	RunServer(context.Context, string, time.Duration) error
	createRouters(*gin.Engine) error
	PrintActualStatistic()
}

type server struct {
	logger logger.Logger
	metric metric.Metric
	ticker time.Ticker
}

func ServerInstance(ctx context.Context) Server {
	once.Do(func() {
		instance = NewServer(ctx)
	})

	return instance
}

func NewServer(ctx context.Context) Server {
	return &server{
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

func (s *server) RunServer(srvCtx context.Context, bindPort string, refreshTicker time.Duration) error {
	// ** GIN configuration ** //
	// set release mode
	gin.SetMode(gin.ReleaseMode)
	// more control for configuration
	ginRouter := gin.New()
	gin.DisableConsoleColor()

	if refreshTicker == 0 {
		s.logger.Write("Timeout not defined, default value 60 seconds")
		refreshTicker = 60
	}

	s.ticker = *time.NewTicker(refreshTicker * time.Second)
	// additional GIN log to FILE
	if s.logger.GetLogLevel() == "debug" {
		gin.DefaultWriter = s.logger.GetLogFile()
		ginRouter.Use(
			gin.Logger(),
		)
	}

	ginRouter.Use(
		gin.Recovery(),
	)

	s.logger.Write("Application waked up")

	go func() {
		if err := s.createRouters(ginRouter); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := ginRouter.Run(bindPort); err != nil {
			panic(err)
		}
	}()

	go func() {
		s.PrintActualStatistic()
	}()

	<-srvCtx.Done()

	return nil
}
