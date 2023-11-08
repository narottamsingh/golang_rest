package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client instance
var DB *mongo.Client = ConnectMongoDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangAPI").Collection(collectionName)
	return collection
}

func ConnectMongoDB() *mongo.Client {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled when the function returns

	// Create a MongoDB client using options
	clientOptions := options.Client().ApplyURI(GetMongoURIFromEnv())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Unable to create DB client ", err)
	}

	// Ping the database to check if the connection is working
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Unable to ping MongoDB: ", err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}
