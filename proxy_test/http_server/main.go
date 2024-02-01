package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		url := ctx.Request.URL.String()
		ctx.JSON(200, gin.H{
			"message": "ok",
			"url":     url,
		})
	})
	server := http.Server{
		Addr:    ":10002",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("error", err)
	}
}
