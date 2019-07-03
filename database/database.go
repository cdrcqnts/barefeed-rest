package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect returns the MongoDB collection object which
// is used to execute all database queries.
func Connect() *mongo.Collection {
	fmt.Println("Starting server...")
	// TODO: restructure, use YAML instead of .env
	url := os.Getenv("MONGO_URL")
	if url == "" {
		log.Fatal("$MONGO_URL must be set.")
	}
	db := os.Getenv("MONGO_DB")
	if db == "" {
		log.Fatal("$MONGO_DB must be set.")
	}
	col := os.Getenv("MONGO_COLLECTION")
	if col == "" {
		log.Fatal("$MONGO_COLLECTION must be set.")
	}
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	res := client.Database(db).Collection(col)
	return res
}
