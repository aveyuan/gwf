<h1 align="center" style="font-size: 40px">GWF (Go Web Fly)</h1>

优雅地启动，关闭，到重启的 `Golang Web` 实现。支持 Gin,GF,Echo,Chi等框架。基于 `http.Server` 统一管理，现在开始 Fly


## GWF 特点
- 能够平滑地关闭 `http.Server`
- 能够平滑地重启 `http.Server`
- 适配多个 `Golang Web` 框架
- 优雅的服务管理最佳实践
- 支持 `logger` 日志接口写入
- `server` 的定义取决你自己的配置
- 支持 `http` 关闭之后，通过 `func` 进行资源的关闭，例如 `db` 连接释放，减少4次握手的时间

## 快速开始

>你可以看到 `cmd/gin` 的例子，直接使用


```go

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

```
