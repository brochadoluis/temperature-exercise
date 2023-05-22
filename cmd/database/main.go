package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	connectToDB()
}

func connectToDB() {
	// MongoDB connection string
	connectionString := "mongodb://localhost:27018"

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		logrus.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB server: %v", err)
	}
	defer func() {
		// Disconnect from the MongoDB server
		if err := client.Disconnect(ctx); err != nil {
			logrus.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Create a logger with logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Add the logger to the context
	ctx = context.WithValue(ctx, "logger", logger)

	// Access the database
	db := client.Database("temperatures")

	// Create collections if they don't exist
	collections := []string{"success", "alert", "error"}
	for _, collection := range collections {
		err := createCollection(ctx, db, collection)
		if err != nil {
			logrus.Fatalf("Failed to create collection %s: %v", collection, err)
		}
	}

	fmt.Println("Collections created successfully!")
}

func createCollection(ctx context.Context, db *mongo.Database, collectionName string) error {
	// Retrieve the logger from the context
	logger := ctx.Value("logger").(*logrus.Logger)

	// Check if the collection already exists
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{"name": collectionName})
	if err != nil {
		logger.Errorf("Failed to list collections: %v", err)
		return err
	}

	// If the collection already exists, return
	if len(collections) > 0 {
		logger.Infof("Collection %s already exists", collectionName)
		return nil
	}

	// Create the collection
	err = db.CreateCollection(ctx, collectionName)
	if err != nil {
		logger.Errorf("Failed to create collection: %v", err)
		return err
	}

	logger.Infof("Collection %s created", collectionName)
	return nil
}
