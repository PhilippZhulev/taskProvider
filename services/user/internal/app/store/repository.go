package store

import (
	"database/sql"

	"gitlab.com/taskProvider/services/user/internal/app/model"
)

// UserRepository ...
// Репозиторий для хранилища пользователя
type UserRepository interface {
	FindByKey(key, value string) (*model.User, error)
	Create(*model.User, string) error
	Remove(id int) (int64, error)
	GetAllUsers(l, o string) (*sql.Rows, error)
	GetAllUsersAndFiltring(l, o, value string) (*sql.Rows, error)
}