package store

import "gitlab.com/taskProvider/services/user/internal/app/model"

// UserRepository ...
// Репозиторий для хранилища пользователя
type UserRepository interface {
	FindByKey(key, value string) (*model.User, error)
	Create(*model.User, string) error
}