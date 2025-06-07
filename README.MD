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

>你可以看到 `cmd` 的例子，参照使用

**基础用法**

```go
package main

import (
	"github.com/aveyuan/gwf"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server := gwf.NewHttpServer(&http.Server{
		Addr:    ":8081",
		Handler: r,
	})

	server.Start()
}

```

**进阶用法**

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


```
