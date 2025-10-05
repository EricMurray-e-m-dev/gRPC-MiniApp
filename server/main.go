package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "grpc-proof/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMetricsServiceServer
}

func (s *server) StreamMetrics(stream pb.MetricsService_StreamMetricsServer) error {
	log.Printf("Client connected. Streaming metrics...")

	metricsReceived := 0

	for {
		metric, err := stream.Recv()
		if err == io.EOF {
			// Client finished sending
			log.Printf("Stream ended. Received %d total metrics.\n", metricsReceived)

			return stream.SendAndClose(&pb.Ack{
				Success: true,
				Message: fmt.Sprintf("Received %d metrics.\n", metricsReceived),
			})
		}
		if err != nil {
			log.Printf("Error receiving metrics: %v", err)
			return err
		}

		log.Printf("Received: %s = %.2f (timestamp: %d)", metric.Name, metric.Value, metric.Timestamp)
		metricsReceived++
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetricsServiceServer(grpcServer, &server{})

	log.Println("Server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
