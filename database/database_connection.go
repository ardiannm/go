package database

import (
	"context"
	"fmt"
	"log"

	"github.com/ardiannm/go/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// client is the global MongoDB client, initialized in init()
var Client *mongo.Client

func init() {
	var err error
	// connect to MongoDB
	Client, err = mongo.Connect(options.Client().ApplyURI(config.MONGODB_URI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	// ping MongoDB to verify connection
	if err := Client.Ping(context.Background(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	// report being connected
	fmt.Println("MongoDB connection established successfully")
}

// return a collection from the global Client
func OpenCollection(collectionName string) *mongo.Collection {
	return Client.Database(config.DATABASE_NAME).Collection(collectionName)
}
