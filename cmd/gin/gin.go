package main

import (
	"net/http"
	"os"

	"github.com/aveyuan/gwf"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {

	// 定义自定义实现的logger
	logger := log.NewHelper(log.NewStdLogger(os.Stdout))
	// 定义一个重启通道
	var restartChan = make(chan struct{})
	// 定义http停止后需要关闭的资源，例如db等，主动退出
	var closeFunc = func() {
		logger.Info("关闭资源1")
		logger.Info("关闭资源2")
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/restart", func(c *gin.Context) {

		// 向重启通道发送一个信号，通知gwf框架进行重启
		restartChan <- struct{}{}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 使用gwf.NewHttpServer创建一个新的HTTP服务器实例
	// 传入http.Server实例和一些选项
	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	}, gwf.WithCloseFunc(closeFunc), gwf.WithLogger(logger), gwf.WithRestartChan(restartChan))

	// 启动HTTP服务器
	server.Start()
}
