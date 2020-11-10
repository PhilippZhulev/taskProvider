package server

import (
	"io/ioutil"
	"net"
	"os"

	pb "gitlab.com/taskProvider/services/broker/proto/getway"

	"gitlab.com/taskProvider/services/broker/internal/app/configurator"

	"gitlab.com/taskProvider/services/broker/internal/app/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type initServer struct {
	pb.UnimplementedUserServer
	user handler.InithandleUser
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
