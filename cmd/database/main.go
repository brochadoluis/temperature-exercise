package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/brochadoluis/temperature-exercise/internal/database"
	"github.com/brochadoluis/temperature-exercise/proto"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	// Create a logger with logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create a context with the logger
	ctx := context.WithValue(context.Background(), "logger", logger)

	connectToDB(ctx)

}

func connectToDB(ctx context.Context) {
	// MongoDB connection string
	connectionString := "mongodb://database:27017"
	// Retrieve the logger from the context
	logger := ctx.Value("logger").(*logrus.Logger)

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		logger.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB server: %v", err)
	}
	defer func() {
		// Disconnect from the MongoDB server
		if err := client.Disconnect(ctx); err != nil {
			logger.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Access the database
	db := client.Database("temperatures")

	// Create collections if they don't exist
	collections := []string{"success", "alert", "error"}
	for _, collection := range collections {
		err := createCollection(ctx, db, collection, logger)
		if err != nil {
			logger.Fatalf("Failed to create collection %s: %v", collection, err)
		}
	}

	fmt.Println("Collections created successfully!")

	// Create a new instance of the database service
	dbService := database.NewService(db)

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the database service with the gRPC server
	proto.RegisterTemperatureServer(grpcServer, dbService)

	// Start the gRPC server
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		logger.Fatalf("Failed to start gRPC server: %v", err)
	}
	defer listener.Close()

	logger.Println("Starting gRPC server...")
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatalf("gRPC server stopped: %v", err)
	}
}

func createCollection(ctx context.Context, db *mongo.Database, collectionName string, logger *logrus.Logger) error {
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
