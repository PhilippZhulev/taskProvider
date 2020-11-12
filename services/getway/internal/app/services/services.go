package services

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "gitlab.com/taskProvider/services/getway/proto/getway"
	"google.golang.org/grpc"
)

// Handler ...
type Handler struct {}

// Init ...
func (h Handler) Init(grpcGwMux *runtime.ServeMux) (*grpc.ClientConn, error) {
    grpcUserConn, err := grpc.Dial(
        ":7040",
		grpc.WithInsecure(),
	)
	
    if err != nil{
        log.Fatalln("Failed to connect to User service", err)
	}

    if err = gw.RegisterUserHandler(
        context.Background(),
        grpcGwMux,
        grpcUserConn,
    ); err !=nil {
  		return nil, err
	}

	return grpcUserConn, nil
}