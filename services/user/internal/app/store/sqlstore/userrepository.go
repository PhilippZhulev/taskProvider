package sqlstore

import (
	"database/sql"

	"github.com/google/uuid"
	"gitlab.com/taskProvider/services/user/internal/app/helpers"
	"gitlab.com/taskProvider/services/user/internal/app/model"
	"gitlab.com/taskProvider/services/user/internal/app/store"
)

// UserRepository ...
//Ссылка на хранилище
type UserRepository struct {
	hesh  helpers.Hesh
	store *Store
}

// FindByKey ...
// Поиск пользователя в базе данных
// Поск по логину
func (r *UserRepository) FindByKey(key, value string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE " + key + " = $1",
		value,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
		&u.Email,
		&u.Name,
		&u.UUID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

// Create ...
//Создание пользователя в базе данных
func (r *UserRepository) Create(u *model.User, salt string) error {
	if err := u.Validate(); err != nil {
		return err
	}

	u.UUID = uuid.New().String()
	u.EncryptedPassword = r.hesh.HashPassword(u.EncryptedPassword, salt)
	defer u.Sanitize()

	return r.store.db.QueryRow(
		`
		INSERT INTO users 
		(login_name, encrypted_password, user_name, email, uuid) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
		`,
		u.Login,
		u.EncryptedPassword,
		u.Name,
		u.Email,
		u.UUID,
	).Scan(&u.ID)
}

// Remove ...
// Удаление пользователя
func (r *UserRepository) Remove(id int) (int64, error) {
	s, err := r.store.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)

	if err != nil {
		return int64(0), err
	}

	
	return s.RowsAffected()
}

// GetAllUsers ...
// Получить пользователей
func (r *UserRepository) GetAllUsers(l, o string) (*sql.Rows, error) {
	return r.store.db.Query(`
		SELECT id, user_name, email, uuid
		FROM users  
		LIMIT $1 
		OFFSET $2 * 2
	`, l, o)
}

// GetAllUsersAndFiltring ...
// Получить пользователей Попараметрам филтрации
func (r *UserRepository) GetAllUsersAndFiltring(l, o, value string) (*sql.Rows, error) {
	return r.store.db.Query(`
		SELECT id, user_name, email, uuid
		FROM users  
		WHERE user_name LIKE '`+value+`%'
		LIMIT $1 
		OFFSET $2 * 2
	`, l, o)
}