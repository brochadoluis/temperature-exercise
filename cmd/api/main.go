package main

import (
	"log"
	"net/http"

	"github.com/brochadoluis/temperature-exercise/internal/api"
	"github.com/brochadoluis/temperature-exercise/proto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	scrapperConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to the Scrapper service: %v", err)
	}
	defer scrapperConn.Close()

	apiService := api.NewAPIService(proto.NewTemperatureClient(scrapperConn))

	router := gin.Default()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	router.Use(gin.LoggerWithWriter(logger.Writer()))

	router.GET("/getTemperature", func(c *gin.Context) {
		latitude := c.Query("latitude")
		longitude := c.Query("longitude")

		temperature, err := apiService.GetTemperature(latitude, longitude)
		if err != nil {
			logger.Errorf("Failed to get temperature: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get temperature"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"temperature": temperature})
	})

	// Start the HTTP server
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the HTTP server: %v", err)
	}
}
