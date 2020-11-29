package server

import (
	pb "gitlab.com/taskProvider/services/user/proto/user"
)

func userLoginResp(data []byte, res, code string) (*pb.UserLoginResponse, error)  {
	return &pb.UserLoginResponse{Data: data, Message: res, Code: code}, nil
}

func userCheckResp(s bool, res string) (*pb.UserCheckResponse, error)  {
	return &pb.UserCheckResponse{State: s, Message: res}, nil
}

func userCreateResp(r, code string) (*pb.UserCreateResponse, error)  {
	return &pb.UserCreateResponse{Message: r, Code: code}, nil
}

func userRemoveResp(r, code string) (*pb.UserRemoveResponse, error)  {
	return &pb.UserRemoveResponse{Message: r, Code: code}, nil
}

func userGetResp(attr []string, mess, code string) (*pb.UserGetResponse, error)  {
	r := &pb.UserGetResponse{}
	if len(attr) > 0 {
		r.Login = attr[0]
		r.Name  = attr[1]
		r.Email = attr[2]
		r.Id 	= attr[3]
	}
	r.Message = mess
	r.Code    = code

	return r, nil
}

func userListResp(users []*pb.Users, mess, code string) (*pb.UserListResponse, error)  {
	r := &pb.UserListResponse{}

	if (len(users) > 0) && (code == "200") {
		r.Users = users
	}

	r.Message = mess
	r.Code    = code

	return r, nil
}