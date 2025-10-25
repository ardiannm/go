package routes

import (
	"github.com/ardiannm/go/controllers"
	"github.com/ardiannm/go/middleware"
	"github.com/ardiannm/go/models"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/users", controllers.GetUsers())
	router.GET("/movies/:imdb_id", controllers.GetMovie())
	router.POST("/movies", controllers.AddMovie())
	router.DELETE("/movies/:imdb_id", controllers.DeleteMovieByIMDBID())
	router.GET("/movies/recommanded", controllers.GetRecommendedMovies())
	router.PATCH("/movies/review/:imdb_id", middleware.RequireRole(models.ADMIN), controllers.AdminReviewUpdate())
}
