package handler

import (
	"gitlab.com/taskProvider/services/user/internal/app/configurator"
	"gitlab.com/taskProvider/services/user/internal/app/model"
	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
	pb "gitlab.com/taskProvider/services/user/proto/user"
)

// HandleCreate ...
func (init Init) HandleCreate(in *pb.UserCreateRequest, store *sqlstore.Store, config *configurator.Config) error {
	return func() error {

		u := &model.User{
			Login:             in.GetLogin(),
			Name:     		   in.GetName(),
			EncryptedPassword: in.GetPassword(),
			Email:    		   in.GetEmail(), 
		}

		if err := store.User().Create(u, config.Salt); err != nil {
			return err
		}

		return nil
	}()
}