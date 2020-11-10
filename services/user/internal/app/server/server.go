package server

import (
	"context"
	"io/ioutil"
	"net"
	"os"

	"github.com/go-chi/jwtauth"
	"gitlab.com/taskProvider/services/user/internal/app/configurator"
	"gitlab.com/taskProvider/services/user/internal/app/handler"
	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
	pb "gitlab.com/taskProvider/services/user/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)


type initServer struct {
	pb.UnimplementedUserServiceServer
	handlers handler.Init
	store *sqlstore.Store
	config *configurator.Config
	tokenAuth *jwtauth.JWTAuth
}

//-------------------------------------
// GRPC сервер
//-------------------------------------

// NewServer ...
func newServer(config *configurator.Config, store *sqlstore.Store) (grpclog.LoggerV2, error) {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	lis, err := net.Listen("tcp", config.ServiceADDR)
	if err != nil {
		return log, err
	}
	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &initServer{
		store: store, 
		config: config, 
		tokenAuth: jwtauth.New("HS256", []byte(config.Salt), nil),
	})
	
	log.Info("Run user service: ", config.ServiceADDR)
	if err := s.Serve(lis); err != nil {
		return log, err
	}

	return log,nil
}

//-------------------------------------
// GRPС методы
//-------------------------------------

// Login service metod
func (s *initServer) UserLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	attr, err := s.handlers.HandleLogin(in.GetAuthData(), s.store, s.config)
	if err != nil {
		return userLoginResp("", err.Error(), "400")
	}

	token, err := s.handlers.HandleCreateToken(attr, s.tokenAuth)
	if err != nil {
		return userLoginResp("", err.Error(), "400")
	}

	return userLoginResp(token, authSuccess, "200")
}

// User check service metod
func (s *initServer) UserCheck(ctx context.Context, in *pb.UserCheckRequest) (*pb.UserCheckResponse, error) {
	t, err := s.tokenAuth.Decode(in.GetData())
	if err != nil {
		return userCheckResp(false, err.Error())
	}

	c := jwtauth.NewContext(context.Background(), t, err)
	token, _, err := jwtauth.FromContext(c)
	if err != nil {
		return userCheckResp(false, err.Error())
	}

	if token == nil || !token.Valid {
		return userCheckResp(false, tokenNotValid)
	}

	return userCheckResp(true, "")
}

// User create service metod
func (s *initServer) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	err := s.handlers.HandleUserCreate(in, s.store, s.config)
	if err != nil {
		return userCreateResp(err.Error(), "400")
	}

	return userCreateResp(userCreateStatus, "200")
}

// User get service metod
func (s *initServer) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	attr, err := s.handlers.HandleUserGet(in, s.store)
	if err != nil {
		return userGetResp([]string{}, err.Error(), "400")
	}

	return userGetResp(attr, userGetStatus, "200")
}