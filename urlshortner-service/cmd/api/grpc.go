package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"shortnerApp/data"
	"shortnerApp/shortner"
	"time"
)

type ShortnerServer struct {
	shortner.UnimplementedShortnerServiceServer
	Model data.Models
}

func (ss *ShortnerServer) WriteLog(ctx context.Context, req *shortner.ShortnerRequest) (*shortner.ShortnerResponse, error) {
	input := req.GetShortnerEntry()

	shortnerEntry := data.Models{
		ShortnerEntry: data.Request{
			URL:         input.Url,
			CustomShort: input.Custom,
			ExpireTime:  time.Duration(input.ExpireTime),
		},
	}

	p, _ := peer.FromContext(ctx)

	status, res := ShortenURL(shortnerEntry.ShortnerEntry, p.Addr.String())
	if status != 200 {
		res := &shortner.ShortnerResponse{
			Result: "failed",
		}
		return res, errors.New("cannot connect to server")
	}

	result := &shortner.ShortnerResponse{
		Result: fmt.Sprintf("%v", res),
	}
	return result, nil
}

func grpcListen(model data.Models) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failes to listen for grpc %v", err)
	}

	server := grpc.NewServer()

	shortner.RegisterShortnerServiceServer(server, &ShortnerServer{
		Model: model,
	})
	log.Println("grpc server started on port: ", grpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("unable to listen for grpc server %v", err)
	}
}
