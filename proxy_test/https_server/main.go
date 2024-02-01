package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qppHUST/proxy_test/https_server/cert_file"
)

func main() {
	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok",
		})
	})
	server := http.Server{
		Addr:    ":10001",
		Handler: router,
	}
	if err := server.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil {
		log.Fatal("error", err)
	}
}
