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

type ForecastResponse struct {
	Latitude    float64
	Longitude   float64
	Temperature float64
	Alert       bool
	Error       bool
}

type Response struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	CurrentWeather Weather `json:"current_weather"`
}

type Weather struct {
	Temperature float64 `json:"temperature"`
}

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

	forecast.setAlert(ctx)

	log.WithContext(ctx).Infof("Temperature for latitude %f and longitude %f: %f",
		forecast.Latitude,
		forecast.Longitude,
		forecast.Temperature)

	saved, err := s.Client.SaveTemperature(ctx, &proto.SaveTemperatureRequest{
		Latitude:    forecast.Latitude,
		Longitude:   forecast.Longitude,
		Temperature: forecast.Temperature,
		Alert:       forecast.Alert,
		Error:       forecast.Error,
		HttpCode:    int32(resp.StatusCode),
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

func (s *Server) parseResponse(ctx context.Context, data io.ReadCloser) (*Response, error) {
	log.WithContext(ctx).Info("Parsing response object")

	response := &Response{}
	err := json.NewDecoder(data).Decode(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON response")
	}

	err = data.Close()
	if err != nil {
		// Handle error closing the response body
		log.WithContext(ctx).Error("Failed to close response body:", err)
	}

	log.WithContext(ctx).Info("Response object parsed successfully")
	return response, nil
}

func (s *Server) parseTemperature(ctx context.Context, data io.ReadCloser) (*ForecastResponse, error) {
	log.WithContext(ctx).Info("Parsing response object")

	var resp Response
	err := json.NewDecoder(data).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON response")
	}

	err = data.Close()
	if err != nil {
		log.WithContext(ctx).Error("Failed to close response body:", err)
		return nil, errors.Wrap(err, "failed to close response body")
	}

	log.WithContext(ctx).Info("Response object parsed successfully")
	log.WithContext(ctx).Info("Mapping response to forecast object")

	forecast := ForecastResponse{
		Latitude:    resp.Latitude,
		Longitude:   resp.Longitude,
		Temperature: resp.CurrentWeather.Temperature,
	}

	log.WithContext(ctx).Info("Response object parsed successfully")

	return &forecast, nil
}

func (f *ForecastResponse) setAlert(ctx context.Context) {
	log.WithContext(ctx).Infof("Setting alert for temperature: %f", f.Temperature)

	f.Alert = f.Temperature < 10 || f.Temperature > 40

	log.WithContext(ctx).Infof("Alert set to: %v", f.Alert)
}

func (f *ForecastResponse) setError(ctx context.Context, statusCode uint32) {
	log.WithContext(ctx).Infof("Setting error for temperature: %f", f.Temperature)

	if statusCode != http.StatusOK {
		f.Error = true
	}

	log.WithContext(ctx).Infof("Error set to: %v", f.Error)
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
