package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/gotrock/go/controllers"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, Gotrock")
	})
	router.GET("/movies", controller.GetMovies())
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
