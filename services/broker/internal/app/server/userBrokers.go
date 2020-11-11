package server

import (
	"context"

	pb "gitlab.com/taskProvider/services/broker/proto/getway"
)

//-------------------------------------
// GRPС переадресации
//-------------------------------------

// Login  --->  User service
func (s *initServer) Login(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {	
	resp := loginResp(ctx)
	conn, c, err := s.user.HandleConnectUser(s.conf.UserADDR)
	defer conn.Close()

	res, err := s.user.HandleLogin(in.GetAuthData().AuthData, c)
	if err != nil {
		return resp(nil, err.Error(), "401")
	}

	return resp([]byte(res.GetData()), res.GetMessage(), res.Code)
}

// CreateUser  --->  User service
func (s *initServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {	
	resp := createUserResp(ctx)
	conn, c, err := s.user.HandleConnectUser(s.conf.UserADDR)
	defer conn.Close()

	code, err := s.user.HandleAuthCheck(ctx, c); 
	if err != nil {
		return resp(err.Error(), code)
	}

	res, err := s.user.HandleCreateUser(in, c)
	if err != nil {
		return resp(err.Error(), "400")
	}

	return resp(res.GetMessage(), res.Code)
}

// GetUser  --->  User service
func (s *initServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {	
	resp := getUserResp(ctx)
	conn, c, err := s.user.HandleConnectUser(s.conf.UserADDR)
	defer conn.Close()

	code, err := s.user.HandleAuthCheck(ctx, c); 
	if err != nil {
		return resp(nil, err.Error(), code)
	}

	res, err := s.user.HandleGetUser(ctx, c)
	if err != nil {
		return resp(nil, err.Error(), "400")
	}

	return resp(&pb.UserData{
		Login: res.Login, 
		Name: res.Name, 
		Email: res.Email,
	}, res.Message, res.Code)
}