package server

import (
	"context"
	"strings"
	"sync"

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
		Id: in.GetId(),
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
		Id: res.Id,
	}, res.Message, res.Code)
}

// GetUserList  --->  User service
func (s *initServer) GetUserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {	
	return userListTempl(ctx, in, "userList:user", false)
}

// GetUserListFilter  --->  User service
func (s *initServer) GetUserListFilter(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {	
	return userListTempl(ctx, in, "userListFilter:user", true)
}

// GetUserList and GetUserListFilter chunk
func userListTempl(ctx context.Context, in *pb.UserListRequest, id string, f bool)  (*pb.UserListResponse, error) {
	resp := getUserListResp(ctx)

	c := newConn(id)
	u := pbUser.NewUserServiceClient(c)
	defer c.Close()

	code, err := handleAuthCheck(ctx, u); 
	if err != nil {
		return resp(nil, err.Error(), code)
	}

	res, err := u.UserList(ctx, &pbUser.UserListRequest{
		List: in.GetList(),
		Offset: in.GetOffset(),
		Filter: f,
		Value: in.GetValue(),
	})

	if err != nil {
		return resp(nil, err.Error(), "400")
	}

	var users []*pb.Users
	var wg sync.WaitGroup

	for _, el := range res.GetUsers() {
		wg.Add(1)
		go func(item *pbUser.Users) {
			users = append(users, &pb.Users{
				Name:  item.GetName(),
				Email: item.GetEmail(),
				Uuid:  item.GetUuid(),
				Id:    item.GetId(),
			})
			wg.Done()
		}(el)
	}
	wg.Wait()
	
	return resp(users, res.Message, res.Code)
}