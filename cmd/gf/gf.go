package main

import (
	"net/http"
	"os"

	"github.com/aveyuan/gwf"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {

	logger := log.NewHelper(log.NewStdLogger(os.Stdout))
	var restartChan = make(chan struct{})
	var closeFunc = func() {
		logger.Info("关闭资源1")
		logger.Info("关闭资源2")
	}
	r := g.Server()
	r.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"message": "pong",
		})
	})
	r.BindHandler("/restart", func(r *ghttp.Request) {
		restartChan <- struct{}{}
		r.Response.WriteJson(g.Map{
			"message": "pong",
		})
	})

	r.Start()

	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	}, gwf.WithCloseFunc(closeFunc), gwf.WithLogger(logger), gwf.WithRestartChan(restartChan))

	server.Start()
}
