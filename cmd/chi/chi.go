package main

import (
	"net/http"
	"os"

	"github.com/aveyuan/gwf"
	"github.com/go-chi/chi/v5"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {

	logger := log.NewHelper(log.NewStdLogger(os.Stdout))
	var restartChan = make(chan struct{})
	var closeFunc = func() {
		logger.Info("关闭资源1")
		logger.Info("关闭资源2")
	}
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/restart", func(w http.ResponseWriter, r *http.Request) {
		restartChan <- struct{}{}
		w.Write([]byte("welcome"))
	})

	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	}, gwf.WithCloseFunc(closeFunc), gwf.WithLogger(logger), gwf.WithRestartChan(restartChan))

	server.Start()
}
