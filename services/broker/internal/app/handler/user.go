package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	pb "gitlab.com/taskProvider/services/broker/proto/getway"
	pbUser "gitlab.com/taskProvider/services/broker/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	errAuth = errors.New("Auth error")
)

// InithandleUser ...
type InithandleUser struct {}

// HandleConnectUser ...
func (hl InithandleUser) HandleConnectUser(addr string) (*grpc.ClientConn, pbUser.UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	return conn, pbUser.NewUserServiceClient(conn), nil
}

// HandleLogin ...
func (hl InithandleUser) HandleLogin(in string, c pbUser.UserServiceClient) (*pbUser.UserLoginResponse, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UserLogin(ctx, &pbUser.UserLoginRequest{AuthData: in})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// HandleAuthCheck ...
func (hl InithandleUser) HandleAuthCheck(ctx context.Context, c pbUser.UserServiceClient) (string, error) {
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

// HandleCreateUser ...
func (hl InithandleUser) HandleCreateUser(in *pb.CreateUserRequest, c pbUser.UserServiceClient) (*pbUser.UserCreateResponse, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UserCreate(ctx, &pbUser.UserCreateRequest{
		Login: in.GetLogin(),
		Name: in.GetName(),
		Email: in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// HandleGetUser ...
func (hl InithandleUser) HandleGetUser(ctx context.Context, c pbUser.UserServiceClient) (*pbUser.UserGetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UserGet(ctx, &pbUser.UserGetRequest{Data: strings.Split(md["authorization"][0], " ")[1]})
	if err != nil {
		return nil, err
	}

	return r, nil
}