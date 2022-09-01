package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"
)

type logServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

// WriteLog write log with grpc
func (ls *logServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// inserting data to database by Insert
	// if it fails by any reasons returns failed and exit the function
	err := ls.Model.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged",
	}
	return res, nil
}

// listening for grpc connections
func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failes to listen for grpc %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &logServer{
		Model: app.Models,
	})
	log.Println("grpc server started on port: ", grpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("unable to listen for grpc server %v", err)
	}
}
