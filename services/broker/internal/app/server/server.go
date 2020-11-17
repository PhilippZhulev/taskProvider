package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
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

type proxy struct {
	Connect  *grpc.ClientConn
	RespID string
}

var connChanIn = make(chan string)
var connChanOut = make(chan *proxy)

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

	newProxyConnections(map[string]string{
		"user": config.UserADDR,
	});

	log.Info("Run broker service: ", config.ServiceADDR)
	if err := s.Serve(lis); err != nil {
		return log, err
	}

	return log, nil
}

// new connection
func newProxyConnections(addrs map[string]string) {
	go func() {
		for ch := range connChanIn {
			go func(in string) {
				sch := strings.Split(in, ":")
				conn, err := grpc.Dial(addrs[sch[1]], grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					log.Println(err.Error())
				}
				connChanOut <- &proxy{Connect: conn, RespID: sch[0]}
			}(ch)
		}
	}()
}

// new broker connection event
func newConn(ev string) (*grpc.ClientConn) {
	connChanIn <- ev
	var conn *grpc.ClientConn
	for ch := range connChanOut {
		if ch.RespID == strings.Split(ev, ":")[0] {
			conn = ch.Connect
			break
		}

		log.Println("failed:", "connetion id not found")
		break
	}

	return conn
}

// check tocken
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
