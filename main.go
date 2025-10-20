package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotrock/go/controllers"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, Gotrock")
	})
	router.GET("/movies", controllers.GetMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovie())
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
