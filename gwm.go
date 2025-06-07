package gwf

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Logger interface {
	Infof(format string, a ...any)
	Errorf(format string, a ...any)
}

func NewVlogger() Logger {
	return &vlogger{}
}

type vlogger struct {
}

func (t *vlogger) Infof(format string, a ...any) {
	log.Printf(format, a...)
}

func (t *vlogger) Errorf(format string, a ...any) {
	log.Printf(format, a...)
}

type HttpServer struct {
	server      *http.Server
	restartChan chan struct{}
	logger      Logger
	closeChan   chan struct{}
	closeFunc   func()
}

type Option func(*HttpServer)

func WithLogger(logger Logger) Option {
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
		logger:    NewVlogger(),
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
			s.logger.Infof("服务关闭完成")
		}
	}()

	s.logger.Infof("服务启动完成")
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
				s.logger.Infof("收到重启信号")
				s.stop()
				if err := s.Restart(); err != nil {
					s.logger.Errorf("重启失败: %s\n", err.Error())
				}
			}
		}()
	}

	<-ctx.Done()
	s.stop()
}

func (s *HttpServer) stop() {
	s.logger.Infof("关闭服务中...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Errorf("正常关闭失败，进行强制关闭: %s\n", err.Error())
	}

	if s.closeFunc != nil {
		s.logger.Infof("开始关闭资源...")
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
