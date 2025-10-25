package routes

import (
	"github.com/ardiannm/go/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.GET("/movies", controllers.GetMovies())
	router.POST("/users", controllers.RegisterUser())
	router.POST("/users/login", controllers.LoginUser())
}
