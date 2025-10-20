package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotrock/go/database"
	"github.com/gotrock/go/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func GetMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		moviesCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var movies []models.Movie
		cursor, err := movieCollection.Find(moviesCtx, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(moviesCtx)
		if err = cursor.All(moviesCtx, &movies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}
		ctx.JSON(http.StatusOK, movies)
	}
}
