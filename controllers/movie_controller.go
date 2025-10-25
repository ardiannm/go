package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ardiannm/go/config"
	"github.com/ardiannm/go/database"
	"github.com/ardiannm/go/models"
	"github.com/ardiannm/go/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies", "reasons": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		movieID := ctx.Param("imdb_id")
		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}
		var movie models.Movie
		err := movieCollection.FindOne(mongoCtx, bson.M{"imdb_id": movieID}).Decode(&movie)
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "reasons": err.Error()})
			return
		}
		result, err := movieCollection.InsertOne(mongoCtx, movie)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie.", "reasons": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	}
}

func DeleteMovieByIMDBID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		movieID := ctx.Param("imdb_id")
		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}
		result, err := movieCollection.DeleteOne(mongoCtx, bson.M{"imdb_id": movieID})
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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, err := utils.GetUserRoleFromContext(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User role not found in this context"})
			return
		}
		if role != "ADMIN" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Not an admin"})
			return
		}
		IMDb_ID := ctx.Param("imdb_id")
		if IMDb_ID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Movie IMDb id is required"})
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var res struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		sentiment, rankingValue, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get review ranking", "reasons": err.Error()})
			return
		}
		filter := bson.M{"imdb_id": IMDb_ID}
		var mongoCtx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		update := bson.M{
			"$set": models.Ranking{
				RankingValue: rankingValue,
				RankingName:  sentiment,
			},
		}
		result, err := movieCollection.UpdateOne(mongoCtx, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
			return
		}
		if result.MatchedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		res.RankingName = sentiment
		res.AdminReview = req.AdminReview
		ctx.JSON(http.StatusOK, res)
	}
}

func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}
	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 0 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ", "
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ", ")
	if config.OPEN_AI_API_KEY == "" {
		return "", 0, errors.New("Could not read OPEN_AI_API_KEY")
	}
	llm, err := openai.New(openai.WithToken(config.OPEN_AI_API_KEY))
	if err != nil {
		return "", 0, err
	}
	BASE_PROMPT := strings.Replace(config.PROMPT_TEMPLATE, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(context.Background(), BASE_PROMPT+adminReview)
	if err != nil {
		return "", 0, err
	}
	rankingValue := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankingValue = ranking.RankingValue
			break
		}
	}
	return response, rankingValue, bson.ErrDecodeToNil
}

func GetRankings() ([]models.Ranking, error) {
	var mongoCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := rankingCollection.Find(mongoCtx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mongoCtx)
	var rankings []models.Ranking
	if err := cursor.All(mongoCtx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := utils.GetUserIDFromContext(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}
		favouriteGenres, err := GetUserFavouriteGenres(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(5)
		filter := bson.M{"genre.genre_name": bson.M{"$in": favouriteGenres}}
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cursor, err := movieCollection.Find(mongoCtx, filter, findOptions)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(mongoCtx)
		var recommendedMovies []models.Movie
		if err := cursor.All(mongoCtx, &recommendedMovies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUserFavouriteGenres(userID string) ([]string, error) {
	mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"user_id": userID}
	projection := bson.M{
		"_id":                         0,
		"favourite_genres.genre_name": 1,
	}
	findOneOptions := options.FindOne().SetProjection(projection)
	var result bson.M
	err := userCollection.FindOne(mongoCtx, filter, findOneOptions).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	favouriteGenres, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("Unable to retrieve favourite genres for user")
	}
	var genreNames []string
	for _, item := range favouriteGenres {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}
	return genreNames, nil
}
