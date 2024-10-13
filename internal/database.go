package internal

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
)

type Database interface {
	Exec(query string, args ...any) error
	Query(query string, args ...any) ([]map[string]any, error)
	QueryRow(query string, args ...any) (map[string]any, error)
}

type postgresDatabase struct {
	connection *sql.DB
}

func PostgresDatabase() postgresDatabase {
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
	return postgresDatabase{connection: connection}
}

func (t *postgresDatabase) Exec(query string, args ...any) error {
	_, err := t.connection.Exec(query, args...)
	return fmt.Errorf("failed to query without result %w", err)
}

func (t *postgresDatabase) Query(query string, args ...any) ([]map[string]any, error) {
	rows, err := t.connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %w", err)
	}
	result := []map[string]any{}
	for rows.Next() {
		row := make(map[string]any)

		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		for i, col := range columns {
			switch value := values[i].(type) {
			case []byte:
				var arr []string
				if err := pq.Array(&arr).Scan(value); err != nil {
					row[col] = string(value)
				} else {
					row[col] = arr
				}
			default:
				row[col] = value
			}
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}
	return result, nil
}

func (t *postgresDatabase) QueryRow(query string, args ...any) (map[string]any, error) {
	return nil, nil
}

type mockDatabase struct {
}

func MockDatabase() mockDatabase {
	return mockDatabase{}
}

func (t *mockDatabase) Exec(query string, args ...any) error {
	return nil
}

func (t *mockDatabase) Query(query string, args ...any) ([]map[string]any, error) {
	return nil, nil
}

func (t *mockDatabase) QueryRow(query string, args ...any) (map[string]any, error) {
	return nil, nil
}
