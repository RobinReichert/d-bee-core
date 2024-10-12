package internal

import (
	"database/sql"
	"fmt"
	"os"
)

type dbConnection interface {
	Exec(query string, args ...any) error
	Query(query string, args ...any) ([]map[string]any, error)
	QueryRow(query string, args ...any) (map[string]any, error)
}

type database struct {
	connection *sql.DB
}

func Database() database {
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s", user, name, sslmode, password, host)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = connection.Ping()
	if err != nil {
		panic(err)
	}
	return database{connection: connection}
}

func (t *database) Exec(query string, args ...any) error {
	return nil
}

type testDatabase struct {
}

func (t *testDatabase) Exec(query string, args ...any) error {
	return nil
}
