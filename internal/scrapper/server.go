package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brochadoluis/temperature-exercise/proto"
)

type Server struct {
	proto.UnimplementedTemperatureServer
	Client *Client
}

// Ensure that the Server struct satisfies the temperatureServerType interface
var _ proto.TemperatureServer = (*Server)(nil)

func (s *Server) mustEmbedUnimplementedTemperatureServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListTemperature(ctx context.Context, req *proto.ListTemperatureRequest) (*proto.ListTemperatureResponse, error) {
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()

	err := checkCoordinates(ctx, latitude, longitude)
	if err != nil {
		log.Error("Coordinates are out of range")
		return &proto.ListTemperatureResponse{}, err
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", latitude, longitude)
	log.WithContext(ctx).Infof("Making API call to: %s", url)

	resp, err := s.makeAPICall(ctx, url)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "failed to make API call")
	}
	forecast, err := s.parseTemperature(ctx, resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse temperature")
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	forecast.setAlert(ctx)

	log.WithContext(ctx).Infof("Temperature for latitude %f and longitude %f: %f",
		forecast.Latitude,
		forecast.Longitude,
		forecast.Temperature)

	saved, err := s.client.SaveTemperature(ctx, &proto.SaveTemperatureRequest{
		Latitude:    forecast.Latitude,
		Longitude:   forecast.Longitude,
		Temperature: forecast.Temperature,
		Alert:       forecast.Alert,
		Error:       forecast.Error,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to save temperature")
	}

	return toListTemperatureResponse(saved), nil
}
func checkCoordinates(ctx context.Context, latitude, longitude float64) error {
	if latitude < -100 || latitude > 100 {
		err := errors.New("latitude is out of range")
		log.WithContext(ctx).Error(err)
		return err
	}

	if longitude < -200 || longitude > 200 {
		err := errors.New("longitude is out of range")
		log.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func (s *Server) makeAPICall(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.WithContext(ctx).Error("Failed to create API request", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithContext(ctx).Error("Failed to make API call", err)
		return nil, err
	}

	return resp, nil
}

func (s *Server) parseTemperature(ctx context.Context, body io.Reader) (*ForecastResponse, error) {
	log.WithContext(ctx).Info("Parsing response object")

	var forecast ForecastResponse
	err := json.NewDecoder(body).Decode(&forecast)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON response")
	}

	log.WithContext(ctx).Info("Response object parsed successfully")

	return &forecast, nil
}

func (f *ForecastResponse) setAlert(ctx context.Context) {
	log.WithContext(ctx).Infof("Setting alert for temperature: %f", f.Temperature)

	f.Alert = f.Temperature < 10 || f.Temperature > 40

	log.WithContext(ctx).Infof("Alert set to: %v", f.Alert)
}

func toListTemperatureResponse(s *proto.SaveTemperatureResponse) *proto.ListTemperatureResponse {
	return &proto.ListTemperatureResponse{
		Latitude:    s.GetLatitude(),
		Longitude:   s.GetLongitude(),
		Temperature: s.GetTemperature(),
		Alert:       s.GetAlert(),
		Error:       s.GetError(),
	}
}
