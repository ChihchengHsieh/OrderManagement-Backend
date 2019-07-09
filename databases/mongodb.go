package databases

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB() {
	// Setting autoload in the main funciton to get the environment variables
	clientOptions := options.Client().ApplyURI(os.Getenv("MGDB_APIKEY"))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Checking
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB successfully.")

	DB = client.Database("shoppingAPI")
}
