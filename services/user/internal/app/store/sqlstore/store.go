package sqlstore

import (
	"database/sql"

	"gitlab.com/taskProvider/services/user/internal/app/store"
)

// Store ...
// Локализировать сущности
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New ...
// Создать новое хранилище
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
