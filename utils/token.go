package utils

import (
	"context"
	"os"
	"time"

	"github.com/ardiannm/go/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserID    string
	jwt.RegisteredClaims
}

var SECRET_ACCESS_KEY = os.Getenv("SECRET_ACCESS_KEY")
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Gotrock",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_ACCESS_KEY))
	if err != nil {
		return "", "", err
	}
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Gotrock",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err
	}
	return signedToken, signedRefreshToken, nil
}

var userCollection *mongo.Collection = database.OpenCollection("users")

func UpdateAllTokens(userId, token, refreshToken string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updatedAt,
		},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updatedData)
	if err != nil {
		return err
	}
	return nil
}
