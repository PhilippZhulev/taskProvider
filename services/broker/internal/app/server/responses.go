package server

import (
	"context"

	pb "gitlab.com/taskProvider/services/broker/proto/getway"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//-------------------------------------
// Response methods
//-------------------------------------

//login response
func loginResp(ctx context.Context) func(token []byte, res, code string) (*pb.LoginUserResponse, error)  {
	return func(token []byte, res, code string) (*pb.LoginUserResponse, error) {
		grpc.SendHeader(ctx, metadata.Pairs("x-http-code", code))
		return &pb.LoginUserResponse{Token: string(token), Message: res}, nil
	}
}

//create user response
func createUserResp(ctx context.Context) func(r, code string) (*pb.CreateUserResponse, error)  {
	return func(r, code string) (*pb.CreateUserResponse, error) {
		grpc.SendHeader(ctx, metadata.Pairs("x-http-code", code))
		return &pb.CreateUserResponse{Message: r}, nil
	}
}

//get user response
func getUserResp(ctx context.Context) func(d *pb.UserData, mes, code string) (*pb.GetUserResponse, error)  {
	return func(d *pb.UserData, mes, code string) (*pb.GetUserResponse, error) {
		grpc.SendHeader(ctx, metadata.Pairs("x-http-code", code))
		return &pb.GetUserResponse{
			UserData: d,
			Message: mes,
		}, nil
	}
}