package handler

import (
	"database/sql"
	"fmt"

	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
	pb "gitlab.com/taskProvider/services/user/proto/user"
)

// HandleList ...
func (init Init) HandleList(in *pb.UserListRequest, store *sqlstore.Store) ([]*pb.Users, error) {
	return func() ([]*pb.Users, error) {
		var (
			usersRows *sql.Rows
			err       error
			result []*pb.Users
		)

		if in.GetFilter() {
			usersRows, err = store.User().GetAllUsersAndFiltring(
				in.GetList(), 
				in.GetOffset(), 
				in.GetValue(),
			)
		} else {
			usersRows, err = store.User().GetAllUsers(
				in.GetList(), 
				in.GetOffset(), 
			)
		}

		if err != nil {
			return nil, err
		}
		defer usersRows.Close()

		for usersRows.Next() {
			row := &pb.Users{}
			if err := usersRows.Scan(
				&row.Id,
				&row.Name,
				&row.Email,
				&row.Uuid,
			); err != nil {
				fmt.Println(err.Error())
			}

			result = append(result, row)
		}

		return result, nil
	}()
}