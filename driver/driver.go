package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectDB() *mongo.Collection {
	fmt.Println("Starting server...")
	url := os.Getenv("MONGO_URL")
	db := os.Getenv("MONGO_DB")
	col := os.Getenv("MONGO_COLLECTION")
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
