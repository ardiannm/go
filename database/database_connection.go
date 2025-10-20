package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client is the global MongoDB client, initialized in init()
var Client *mongo.Client

func init() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, ensure environment variables are set")
	}
	// Get MongoDB URI
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	// Connect to MongoDB
	Client, err = mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	// Ping MongoDB to verify connection
	if err := Client.Ping(context.Background(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	fmt.Println("MongoDB connection established successfully")
}

// OpenCollection returns a collection from the global Client
func OpenCollection(collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME environment variable not set")
	}
	return Client.Database(databaseName).Collection(collectionName)
}
