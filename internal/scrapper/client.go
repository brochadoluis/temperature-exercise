package scrapper

import (
	"context"

	"github.com/brochadoluis/temperature-exercise/proto"
)

type ForecastResponse struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	Alert       bool    `json:"alert"`
	Error       bool    `json:"error"`
}

type Client struct {
	client proto.TemperatureClient
}

func NewClient(client proto.TemperatureClient) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) mustEmbedUnimplementedTemperatureServer() {
	//TODO implement me
}

func (c *Client) SaveTemperature(ctx context.Context, req *proto.SaveTemperatureRequest) (*proto.SaveTemperatureResponse, error) {
	// Placeholder response
	response := &proto.SaveTemperatureResponse{}

	return response, nil
}
