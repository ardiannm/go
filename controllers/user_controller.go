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
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		var mongoCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		count, err := userCollection.CountDocuments(mongoCtx, bson.M{"email": user.Email})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
		}
		if count > 0 {
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		hashed, err := HashPassword(user.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
			return
		}
		user.Password = hashed
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		result, err := userCollection.InsertOne(mongoCtx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var users []models.User
		cursor, err := userCollection.Find(mongoCtx, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer cursor.Close(mongoCtx)
		if err = cursor.All(mongoCtx, &users); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, users)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userLogin models.UserLogin
		if err := ctx.ShouldBindJSON(&userLogin); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var foundUser models.User
		err := userCollection.FindOne(mongoCtx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	}
}
