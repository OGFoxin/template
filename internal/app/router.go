package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// create router only when server start
// firstable safe increase counter via mutex
// and sync gorutin with WG
func (s *server) createRouters(router *gin.Engine) error {
	if router == nil {
		return gin.Error{}
	}

	wg := sync.WaitGroup{}
	router.GET("/healthCheck", func(c *gin.Context) {
		wg.Add(1)
		go func() {
			if err := s.metric.IncreaseHttpStat(http.StatusOK, &wg); err != nil {
				return
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": "alive",
		})
		wg.Wait()
	})

	router.GET("/getHttpStat", func(c *gin.Context) {
		ch := make(chan map[int]int)
		go func() {
			ch <- s.metric.GetHttpStats()
		}()

		c.JSON(http.StatusOK, gin.H{
			"httpStats": <-ch,
		})

		close(ch)
	})

	router.GET("/getCpuStat", func(c *gin.Context) {
		ch := make(chan []float64)
		go func() {
			ch <- s.metric.GetCpuInfo()
		}()

		c.JSON(http.StatusOK, gin.H{
			"cpuStats": <-ch,
		})

		close(ch)
	})

	router.GET("/getMemoryStat", func(c *gin.Context) {
		ch := make(chan interface{})

		go func() {
			ch <- s.metric.GetMemoryInfo()
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": <-ch,
		})

		close(ch)
	})

	router.NoRoute(func(c *gin.Context) {
		wg.Add(1)
		go func() {
			if err := s.metric.IncreaseHttpStat(http.StatusNotFound, &wg); err != nil {
				return
			}
		}()

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "resource not found",
		})
		wg.Wait()
	})
	return nil
}
