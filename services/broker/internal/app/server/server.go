package server

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"gitlab.com/taskProvider/services/broker/internal/app/configurator"
	pb "gitlab.com/taskProvider/services/broker/proto/getway"
	pbUser "gitlab.com/taskProvider/services/broker/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

var (
	errAuth = errors.New("Auth error")
)


type initServer struct {
	pb.UnimplementedUserServer
	conf *configurator.Config 
}

//-------------------------------------
// GRPC сервер
//-------------------------------------

// NewServer ...
func NewServer(config *configurator.Config) (grpclog.LoggerV2, error) {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	lis, err := net.Listen("tcp", config.ServiceADDR)
	if err != nil {
		return log, err
	}
	s := grpc.NewServer()
	pb.RegisterUserServer(s, &initServer{conf: config})
	
	log.Info("Run broker service: ", config.ServiceADDR)
	if err := s.Serve(lis); err != nil {
		return log, err
	}

	return log, nil
}

// HandleAuthCheck ...
func handleAuthCheck(ctx context.Context, c pbUser.UserServiceClient) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	if len(md["authorization"]) == 0 { 
		return "401", errAuth
	}

	uch, err := c.UserCheck(ctx, &pbUser.UserCheckRequest{
		Data: strings.Split(md["authorization"][0], " ")[1],
	})

	if err != nil {
		return "400", err
	}

	if uch.GetState() == false {
		return "401", errAuth
	}

	return "200", nil
}

func handleConnectUser(addr string) (*grpc.ClientConn, pbUser.UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	return conn, pbUser.NewUserServiceClient(conn), nil
}