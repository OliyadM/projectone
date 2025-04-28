package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to the Mongo database using URI and DB name
func ConnectMongo(uri string, dbName string) *mongo.Database {
	if uri == "" {
		log.Fatal("MongoDB URI is not provided") // Removed os.Getenv redundancy
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to create Mongo client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to Mongo:", err)
	}

	return client.Database(dbName)
}
