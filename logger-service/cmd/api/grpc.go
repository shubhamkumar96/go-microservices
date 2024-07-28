package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/shubhamkumar96/go-microservices/logger-service/data"
	"github.com/shubhamkumar96/go-microservices/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	// below type is required for writting all the grpc server, it is to ensure backward compatibility
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// create the log entry data
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// insert the log entry data to mongoDB
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "Failed"}
		return res, err
	}

	// return the sucessful response
	res := &logs.LogResponse{Result: "Logged Via GRPC!"}
	return res, nil
}

// Define grpc listener
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	// register the grpc methods
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server Stared on port %s", gRpcPort)

	// start the gRPC server
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
