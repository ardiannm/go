package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ardiannm/go/database"
	"github.com/ardiannm/go/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var movies []models.Movie
		cursor, err := movieCollection.Find(mongoCtx, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
			return
		}
		defer cursor.Close(mongoCtx)
		if err = cursor.All(mongoCtx, &movies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		movieID := ctx.Param("imdb_id")
		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		ctx.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var movie models.Movie
		if err := ctx.ShouldBindJSON(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		result, err := movieCollection.InsertOne(mongoCtx, movie)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie.", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	}
}

func DeleteMovieByIMDBID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		movieID := ctx.Param("imdb_id")
		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}
		result, err := movieCollection.DeleteOne(ctx, bson.M{"imdb_id": movieID})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.DeletedCount == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"deleted":       false,
				"deleted_count": result.DeletedCount,
				"message":       "No movie found with the provided IMDb ID.",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"deleted":       true,
			"deleted_count": result.DeletedCount,
			"message":       "Movie successfully deleted.",
		})
	}
}
