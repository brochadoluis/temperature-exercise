package scrapper

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/brochadoluis/temperature-exercise/proto"
)

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
	// Initialize a logger
	log := logrus.WithContext(ctx).WithField("method", "SaveTemperature")

	log.Infof("Sending SaveTemperature request: %v", req)

	resp, err := c.client.SaveTemperature(ctx, req)
	if err != nil {
		log.Errorf("SaveTemperature request failed: %v", err)
		return nil, err
	}

	log.Infof("SaveTemperature response received: %v", resp)

	return resp, nil
}
