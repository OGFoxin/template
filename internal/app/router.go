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
	router.GET("/healtCheck", func(c *gin.Context) {
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
