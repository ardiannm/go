package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ardiannm/go/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, Gotrock")
	})

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovie())
	router.POST("/movies", controllers.AddMovie())
	router.DELETE("/movies/:imdb_id", controllers.DeleteMovieByIMDBID())

	router.POST("/users", controllers.RegisterUser())
	router.GET("/users", controllers.GetUsers())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
