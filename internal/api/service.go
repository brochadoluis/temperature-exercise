package api

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/brochadoluis/temperature-exercise/proto"
)

type Service struct {
	client proto.TemperatureClient
}

func NewAPIService(client proto.TemperatureClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetTemperature(latitude, longitude string) (*proto.ListTemperatureResponse, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Errorf("Failed to convert latitude: %v", err)
		return &proto.ListTemperatureResponse{}, err
	}

	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Errorf("Failed to convert longitude: %v", err)
		return &proto.ListTemperatureResponse{}, err
	}
	// Create a gRPC request
	req := &proto.ListTemperatureRequest{
		Latitude:  lat,
		Longitude: lng,
	}

	// Invoke the gRPC method on the client
	resp, err := s.client.ListTemperature(context.Background(), req)
	if err != nil {
		log.Errorf("Failed to call Method: %v", err)
		return &proto.ListTemperatureResponse{}, err
	}
	return resp, nil
}
