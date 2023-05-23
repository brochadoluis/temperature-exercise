package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/brochadoluis/temperature-exercise/proto"
)

type Service struct {
	db *mongo.Database
	proto.TemperatureServer
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SaveTemperature(ctx context.Context, req *proto.SaveTemperatureRequest) (*proto.SaveTemperatureResponse, error) {
	// Retrieve the logger from the context
	logger := logrus.WithContext(ctx)

	// Determine the collection based on the conditions
	collectionName := "success"
	if req.GetError() {
		collectionName = "error"
	}

	// Save the temperature data to the appropriate collection
	collection := s.db.Collection(collectionName)
	data := map[string]interface{}{
		"timestamp": time.Now(),
		"request":   req,
	}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		logger.Errorf("Failed to save temperature data to %s collection: %v", collectionName, err)
		return nil, err
	}

	// Save to alerts collection if Alert field is true
	if req.GetAlert() {
		alertsCollection := s.db.Collection("alert")
		_, err := alertsCollection.InsertOne(ctx, data)
		if err != nil {
			logger.Errorf("Failed to save temperature data to alerts collection: %v", err)
			return nil, err
		}
	}

	// Return the response
	return &proto.SaveTemperatureResponse{
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Temperature: req.Temperature,
		Alert:       req.Alert,
		Error:       req.Error,
	}, nil
}
