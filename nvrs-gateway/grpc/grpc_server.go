package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"nvrs-gateway/protos"

	"google.golang.org/grpc"
)

// server is used to implement the AgentService.
type server struct {
	protos.UnimplementedAgentServiceServer
}

// UpdateStatus - gRPC method for updating agent status
func (s *server) UpdateStatus(ctx context.Context, req *protos.StatusRequest) (*protos.StatusResponse, error) {
	log.Printf("Received UpdateStatus request for Agent ID: %d, Status: %s", req.AgentId, req.Status)
	// TODO: Update agent status in your storage or database
	response := &protos.StatusResponse{
		Message: fmt.Sprintf("Status updated for Agent ID %d", req.AgentId),
	}
	return response, nil
}

// SubmitTask - gRPC method for submitting a task
func (s *server) SubmitTask(ctx context.Context, req *protos.TaskRequest) (*protos.TaskResponse, error) {
	log.Printf("Received SubmitTask request for Agent ID: %d, Task: %s", req.AgentId, req.Task)
	// TODO: Store task for the agent in your storage or database
	response := &protos.TaskResponse{
		Message: fmt.Sprintf("Task received for Agent ID %d", req.AgentId),
	}
	return response, nil
}

func main() {
	// Start a listener on port 50051 (or any port you prefer)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the AgentService server
	protos.RegisterAgentServiceServer(grpcServer, &server{})

	log.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
