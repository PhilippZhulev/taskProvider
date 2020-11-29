package handler

import (
	"errors"
	"log"
	"strconv"

	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
	pb "gitlab.com/taskProvider/services/user/proto/user"
)

// HandleUserRemove ...
func (init Init) HandleUserRemove(in *pb.UserRemoveRequest, store *sqlstore.Store) error {

	return func() error {
		id, err := strconv.Atoi(in.GetId())
		if err != nil {
			log.Println(err)
			return err;
		}

		
		i, err := store.User().Remove(id); 
		if err != nil {
			log.Println(err)
			return err;
		}

		if int(i) == 0 {
			return errors.New("No such user");
		}

		return nil
	}()
}