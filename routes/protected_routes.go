package routes

import (
	"github.com/ardiannm/go/controllers"
	"github.com/ardiannm/go/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movies/:imdb_id", controllers.GetMovie())
	router.POST("/movies", controllers.AddMovie())
	router.DELETE("/movies/:imdb_id", controllers.DeleteMovieByIMDBID())
	router.GET("/movies/recommanded", controllers.GetRecommendedMovies())
}
