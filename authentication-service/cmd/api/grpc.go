package main

import (
	"authApp/auth"
	"authApp/cmd/data"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var grpcPort = "50001"

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Model data.Models
}

func (as *AuthServer) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetAuthEntry()

	authEntry := data.Models{
		User: data.User{
			Email:    input.Email,
			Password: input.Password,
		},
	}

	ok, err := as.Model.User.PasswordMatches(authEntry.User.Password)
	if err != nil || ok == false {
		res := &auth.AuthResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &auth.AuthResponse{
		Result: "logged in",
	}

	return res, nil
}

func (app *Config) grpcListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failes to listen for grpc %v", err)
	}

	server := grpc.NewServer()

	auth.RegisterAuthServiceServer(server, &AuthServer{
		Model: app.Models,
	})
	log.Println("grpc server started on port: ", grpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("unable to listen for grpc server %v", err)
	}
}
