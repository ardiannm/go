package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, Gotrock")
	})

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)

	}
}
