package server

import (
	"context"
	"strings"

	pb "gitlab.com/taskProvider/services/broker/proto/getway"
	pbUser "gitlab.com/taskProvider/services/broker/proto/user"
	"google.golang.org/grpc/metadata"
)

//-------------------------------------
// GRPС переадресации
//-------------------------------------

// Login  --->  User service
func (s *initServer) Login(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {	
	resp := loginResp(ctx)

	c := newConn("loginUser:user")
	u := pbUser.NewUserServiceClient(c)
	defer c.Close()

	res, err := u.UserLogin(ctx, &pbUser.UserLoginRequest{AuthData: in.GetAuthData().AuthData})


	if err != nil {
		return resp(nil, err.Error(), res.Code)
	}

	return resp(res.GetData(), res.GetMessage(), res.Code)
}

// CreateUser  --->  User service
func (s *initServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {	
	resp := createUserResp(ctx)

	c := newConn("createUser:user")
	u := pbUser.NewUserServiceClient(c)
	defer c.Close()

	code, err := handleAuthCheck(ctx, u); 
	if err != nil {
		return resp(err.Error(), code)
	}

	res, err := u.UserCreate(ctx, &pbUser.UserCreateRequest{
		Login: in.GetLogin(),
		Name: in.GetName(),
		Email: in.GetEmail(),
		Password: in.GetPassword(),
	})

	if err != nil {
		return resp(err.Error(), res.Code)
	}

	return resp(res.GetMessage(), res.Code)
}

// RemoveUser  --->  User service
func (s *initServer) RemoveUser(ctx context.Context, in *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	resp := removeUserResp(ctx)

	c := newConn("removeUser:user")
	u := pbUser.NewUserServiceClient(c)
	defer c.Close()

	code, err := handleAuthCheck(ctx, u); 
	if err != nil {
		return resp(err.Error(), code)
	}

	res, err := u.UserRemove(ctx, &pbUser.UserRemoveRequest{
		Uuid: in.GetId(),
	}) 
	
	if err != nil {
		return resp(err.Error(), res.Code)
	}

	return resp(res.GetMessage(), res.Code)
}

// GetUser  --->  User service
func (s *initServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {	
	md, _ := metadata.FromIncomingContext(ctx)
	resp := getUserResp(ctx)
	
	c := newConn("getUser:user")
	u := pbUser.NewUserServiceClient(c)
	defer c.Close()

	code, err := handleAuthCheck(ctx, u); 
	if err != nil {
		return resp(nil, err.Error(), code)
	}

	res, err := u.UserGet(ctx, &pbUser.UserGetRequest{Data: strings.Split(md["authorization"][0], " ")[1]})
	if err != nil {
		return resp(nil, err.Error(), "400")
	}

	return resp(&pb.UserData{
		Login: res.Login, 
		Name: res.Name, 
		Email: res.Email,
	}, res.Message, res.Code)
}