package main

import (
	"fmt"
	"net/http"

	"github.com/ardiannm/go/controllers"
	"github.com/gin-gonic/gin"
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
	router.POST("/users/login", controllers.LoginUser())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
