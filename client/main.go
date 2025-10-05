package main

import (
	"context"
	"log"
	"time"

	pb "grpc-proof/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMetricsServiceClient(conn)

	stream, err := client.StreamMetrics(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	metrics := []struct {
		name  string
		value float64
	}{
		{"cpu_usage", 45.2},
		{"memory_usage", 40.3},
		{"disk_usage", 120.0},
		{"cpu_usage", 55.5},
		{"memory_usage", 99.0},
	}

	log.Printf("Sending metrics to server...")

	for _, m := range metrics {
		metric := &pb.Metric{
			Name:      m.name,
			Value:     m.value,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(metric); err != nil {
			log.Fatalf("Failed to send metric: %v", err)
		}

		log.Printf("Sent: %s = %.2f", m.name, m.value)
		time.Sleep(500 * time.Millisecond) // Simulate delay
	}

	ack, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed ack: %v", err)
	}

	log.Printf("Server response: %s (success=%v)", ack.Message, ack.Success)
}
