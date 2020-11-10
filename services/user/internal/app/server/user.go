package server

import (
	"database/sql"

	"gitlab.com/taskProvider/services/user/internal/app/configurator"
	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"

	_ "github.com/lib/pq" // ...
)

//-------------------------------------
// APP RUNNER
//-------------------------------------

// Run ...
func Run(config *configurator.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)

	logs, err := newServer(config, store)
	if err != nil {
		logs.Fatalf("failed: %v", err)
	}

	return nil
}

// Инициализировать новую базу данных
func newDB(dbURL string) (*sql.DB, error) {
	// Открыть соединение с db postgres
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
