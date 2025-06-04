package main

import (
	"net/http"
	"os"

	"github.com/aveyuan/gwf"
	"github.com/labstack/echo/v4"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {

	logger := log.NewHelper(log.NewStdLogger(os.Stdout))
	var restartChan = make(chan struct{})
	var closeFunc = func() {
		logger.Info("关闭资源1")
		logger.Info("关闭资源2")
	}
	r := echo.New()
	r.GET("/", func(c echo.Context) error {
		return c.JSON(200, echo.Map{
			"message": "pong",
		})
	})
	r.GET("/restart", func(c echo.Context) error {
		restartChan <- struct{}{}
		return c.JSON(200, echo.Map{
			"message": "pong",
		})
	})

	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	}, gwf.WithCloseFunc(closeFunc), gwf.WithLogger(logger), gwf.WithRestartChan(restartChan))

	server.Start()
}
