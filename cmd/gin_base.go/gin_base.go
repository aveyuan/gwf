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
