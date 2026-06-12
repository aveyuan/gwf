package gwf

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)


type HttpServer struct {
	server      *http.Server
	restartChan chan struct{}
	logger      *slog.Logger
	closeChan   chan struct{}
	closeFunc   func()
}

type Option func(*HttpServer)

func WithLogger(logger *slog.Logger) Option {
	return func(srv *HttpServer) {
		srv.logger = logger
	}
}

func WithRestartChan(ch chan struct{}) Option {
	return func(srv *HttpServer) {
		srv.restartChan = ch
	}
}

func WithCloseFunc(closeFunc func()) Option {
	return func(srv *HttpServer) {
		srv.closeFunc = closeFunc
	}
}

func NewHttpServer(server *http.Server, options ...Option) *HttpServer {
	srv := &HttpServer{
		server:    server,
		logger:    slog.Default(),
		closeChan: make(chan struct{}),
	}

	for _, v := range options {
		v(srv)
	}

	return srv
}

func (s *HttpServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			s.logger.Warn("服务关闭完成")
		}
	}()

	s.logger.Warn("服务启动完成", "addr", s.server.Addr)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		s.closeChan <- struct{}{}
	}()

	if s.restartChan != nil {
		// 监听数据
		go func() {
			for range s.restartChan {
				s.logger.Warn("收到重启信号")
				s.stop()
				if err := s.Restart(); err != nil {
					s.logger.Error("重启失败", "error", err.Error())
				}
			}
		}()
	}

	<-ctx.Done()
	s.stop()
}

func (s *HttpServer) stop() {
	s.logger.Warn("关闭服务中...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("正常关闭失败，进行强制关闭", "error", err.Error())
	}

	if s.closeFunc != nil {
		s.logger.Warn("开始关闭资源...")
		s.closeFunc()
	}

}

func (s *HttpServer) Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}

	args := os.Args
	env := os.Environ()
	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err = cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, args, env)
}
