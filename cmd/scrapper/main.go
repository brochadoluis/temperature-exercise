package main

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/brochadoluis/temperature-exercise/internal/scrapper"
	"github.com/brochadoluis/temperature-exercise/proto"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})

	conn, err := grpc.Dial("server:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	grpcClient := proto.NewTemperatureClient(conn)

	scrapperClient := scrapper.NewClient(grpcClient)

	startGRPCServer(log, scrapperClient)
}

func startGRPCServer(log *logrus.Logger, scrapperClient *scrapper.Client) {
	server := grpc.NewServer()

	serverImpl := &scrapper.Server{
		Client: scrapperClient,
	}

	proto.RegisterTemperatureServer(server, serverImpl)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Info("Scrapper gRPC server started")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
