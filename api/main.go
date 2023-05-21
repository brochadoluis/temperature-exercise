package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type TemperatureResponse struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	Alert       bool    `json:"alert"`
	Error       bool    `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	r := gin.Default()

	r.GET("/", Handler)
	fmt.Print("Hello")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello World")
	fmt.Print("Hello World!")
}

