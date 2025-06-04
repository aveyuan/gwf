package main

import (
	"net/http"
	"os"

	"github.com/aveyuan/gwf"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {

	logger := log.NewHelper(log.NewStdLogger(os.Stdout))
	var restartChan = make(chan struct{})
	var closeFunc = func() {
		logger.Info("关闭资源1")
		logger.Info("关闭资源2")
	}
	logger.Error()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/restart", func(c *gin.Context) {

		restartChan <- struct{}{}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	}, gwf.WithCloseFunc(closeFunc), gwf.WithLogger(logger), gwf.WithRestartChan(restartChan))

	server.Start()
}
