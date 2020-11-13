package server

import (
	"context"

	"gitlab.com/taskProvider/services/broker/proto/user"
	"google.golang.org/grpc"
)

// NewDefUserConn ...
// User protected connection
func newDefUserConn(ctx context.Context, s *initServer) (*grpc.ClientConn, user.UserServiceClient, string, error) {
	conn, c, err := handleConnectUser(s.conf.UserADDR)
	if err != nil {
		return nil, nil, "400", err
	}

	code, err := handleAuthCheck(ctx, c); 
	if err != nil {
		return nil, nil, code, nil
	}

	return conn, c, code, nil
}